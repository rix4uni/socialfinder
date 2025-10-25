package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/rix4uni/socialfinder/banner"
)

const defaultURL = "https://raw.githubusercontent.com/rix4uni/socialfinder/refs/heads/main/urls.txt"

type Platform struct {
	Name string
	URL  string
}

// ANSI color codes - Sherlock style
const (
	ColorReset   = "\033[0m"
	ColorGreen   = "\033[92m"  // Bright Green for [+]
	ColorYellow  = "\033[93m"  // Bright Yellow for [*]
	ColorRed     = "\033[91m"  // Bright Red for errors
	ColorCyan    = "\033[96m"  // Bright Cyan for username
	ColorMagenta = "\033[95m"  // Bright Magenta for platform names
)

func main() {
	fileFlag := flag.String("file", "", "Custom URLs file path")
	silent := flag.Bool("silent", false, "silent mode.")
	version := flag.Bool("version", false, "Print the version of the tool and exit.")
	flag.Parse()

	if *version {
		banner.PrintBanner()
		banner.PrintVersion()
		return
	}

	if !*silent {
		banner.PrintBanner()
	}

	args := flag.Args()
	if len(args) < 1 {
		fmt.Printf("%sUsage: socialfinder [ -file <custom_file> ] <username>%s\n", ColorRed, ColorReset)
		return
	}

	username := args[0]

	fmt.Printf("%s[*]%s Checking username %s%s%s on:\n\n", ColorYellow, ColorReset, ColorCyan, username, ColorReset)

	// Read URLs from file or default URL
	platforms, err := readURLs(*fileFlag)
	if err != nil {
		fmt.Printf("%s[!]%s Error reading URLs: %v\n", ColorRed, ColorReset, err)
		return
	}

	// Check platforms and stream results in real-time
	activeCount := checkPlatformsStream(platforms, username)

	fmt.Printf("\n%s[*]%s Search completed with %s%d%s results\n", ColorYellow, ColorReset, ColorCyan, activeCount, ColorReset)
}

func readURLs(filePath string) ([]Platform, error) {
	var reader io.Reader
	var platforms []Platform

	if filePath != "" {
		// Read from custom file
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %v", filePath, err)
		}
		defer file.Close()
		reader = file
	} else {
		// Read from default URL
		resp, err := http.Get(defaultURL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch default URLs: %v", err)
		}
		defer resp.Body.Close()
		reader = resp.Body
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			platforms = append(platforms, Platform{
				Name: strings.TrimSpace(parts[0]),
				URL:  strings.TrimSpace(parts[1]),
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading URLs: %v", err)
	}

	return platforms, nil
}

// normalizeURL normalizes a URL for better matching
func normalizeURL(url string) string {
	// Remove protocol
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")
	
	// Remove www. prefix
	url = strings.TrimPrefix(url, "www.")
	
	// Remove trailing slashes
	url = strings.TrimSuffix(url, "/")
	
	return strings.ToLower(url)
}

// findPlatformName tries to find the platform name for a given URL
func findPlatformName(url string, urlToName map[string]string, platforms []Platform, username string) string {
	// Try exact match first
	if platformName, exists := urlToName[url]; exists {
		return platformName
	}
	
	// Try normalized version
	normalizedURL := normalizeURL(url)
	for _, platform := range platforms {
		testURL := strings.Replace(platform.URL, "$USERNAME", username, 1)
		normalizedTestURL := normalizeURL(testURL)
		if normalizedURL == normalizedTestURL {
			return platform.Name
		}
	}
	
	// Try to extract from URL path as fallback
	parts := strings.Split(normalizeURL(url), "/")
	if len(parts) > 0 {
		domain := parts[0]
		// Simple heuristic: use domain name as platform name
		dotIndex := strings.Index(domain, ".")
		if dotIndex > 0 {
			return strings.Title(domain[:dotIndex])
		}
		return strings.Title(domain)
	}
	
	return ""
}

func checkPlatformsStream(platforms []Platform, username string) int {
	activeCount := 0

	// Prepare URLs for httpx and create multiple mappings for better matching
	urls := make([]string, len(platforms))
	urlToName := make(map[string]string)
	normalizedUrlToName := make(map[string]string)
	
	for i, platform := range platforms {
		fullURL := strings.Replace(platform.URL, "$USERNAME", username, 1)
		urls[i] = fullURL
		urlToName[fullURL] = platform.Name
		normalizedUrlToName[normalizeURL(fullURL)] = platform.Name
	}

	// Execute httpx command
	cmd := exec.Command("httpx", "-duc", "-silent", "-timeout", "30", "-mc", "200,301,302", "-random-agent")
	
	// Get stdout pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("%s[!]%s Error creating stdout pipe: %v\n", ColorRed, ColorReset, err)
		return 0
	}

	// Pass URLs via stdin
	cmd.Stdin = strings.NewReader(strings.Join(urls, "\n"))
	
	// Start the command
	if err := cmd.Start(); err != nil {
		fmt.Printf("%s[!]%s Error starting httpx: %v\n", ColorRed, ColorReset, err)
		return 0
	}

	// Read output line by line and print immediately
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url != "" {
			// Try to find the platform name using multiple methods
			platformName := findPlatformName(url, urlToName, platforms, username)
			
			if platformName != "" {
				fmt.Printf("%s[+]%s %s%s%s: %s\n", ColorGreen, ColorReset, ColorGreen, platformName, ColorReset, url)
			} else {
				// If we can't find the platform name, show a cleaned up version
				cleanURL := url
				if strings.HasPrefix(cleanURL, "https://") {
					cleanURL = cleanURL[8:]
				} else if strings.HasPrefix(cleanURL, "http://") {
					cleanURL = cleanURL[7:]
				}
				if strings.HasPrefix(cleanURL, "www.") {
					cleanURL = cleanURL[4:]
				}
				// Extract domain name for display
				domainEnd := strings.Index(cleanURL, "/")
				if domainEnd > 0 {
					domain := cleanURL[:domainEnd]
					fmt.Printf("%s[+]%s %s%s%s: %s\n", ColorGreen, ColorReset, ColorGreen, strings.Title(domain), ColorReset, url)
				} else {
					fmt.Printf("%s[+]%s %s\n", ColorGreen, ColorReset, url)
				}
			}
			activeCount++
		}
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		// httpx might return non-zero exit code if no URLs match, which is fine for us
		if _, ok := err.(*exec.ExitError); ok {
			// Command ran but returned non-zero exit status, we can ignore this
		} else {
			fmt.Printf("%s[!]%s Error waiting for httpx: %v\n", ColorRed, ColorReset, err)
		}
	}

	return activeCount
}