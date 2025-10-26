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

// Hardcoded list of NSFW websites
var nsfwWebsites = []string{
	"https://onlyfans.com/$USERNAME",
	"https://chaturbate.com/$USERNAME",
	"https://bongacams.com/profile/$USERNAME",
	"https://pornhub.com/users/$USERNAME",
	"https://xvideos.com/profiles/$USERNAME",
	"https://xhamster.com/users/$USERNAME",
	"https://redtube.com/users/$USERNAME",
	"https://youporn.com/uservids/$USERNAME",
	"https://tnaflix.com/profile/$USERNAME",
	"https://rockettube.com/$USERNAME",
	"https://motherless.com/m/$USERNAME",
	"https://erome.com/$USERNAME",
	"https://heavy-r.com/user/$USERNAME",
	"https://imagefap.com/profile/$USERNAME",
	"https://apclips.com/$USERNAME",
	"https://admireme.vip/$USERNAME",
	"https://pocketstars.com/$USERNAME",
	"https://modelhub.com/$USERNAME/videos",
	"https://royalcams.com/profile/$USERNAME",
	"https://livejasmin.com/$USERNAME",
	"https://cams.com/$USERNAME",
	"https://stripchat.com/$USERNAME",
	"https://myfreecams.com/$USERNAME",
	"https://cam4.com/$USERNAME",
	"https://flirt4free.com/$USERNAME",
	"https://jasmin.com/$USERNAME",
	"https://imlive.com/$USERNAME",
	"https://streamate.com/$USERNAME",
	"https://xhamsterlive.com/$USERNAME",
	"https://camsoda.com/$USERNAME",
	"https://rabbitscams.com/$USERNAME",
	"https://wowcams.com/$USERNAME",
	"https://joingy.com/$USERNAME",
	"https://shagle.com/$USERNAME",
	"https://coomeet.com/$USERNAME",
	"https://camfuze.com/$USERNAME",
	"https://camslurp.com/$USERNAME",
	"https://camsfind.com/$USERNAME",
	"https://fapchat.com/$USERNAME",
	"https://dirtyroulette.com/$USERNAME",
	"https://chatroulette.com/$USERNAME",
	"https://omegle.tv/$USERNAME",
	"https://flingster.com/$USERNAME",
	"https://adultchat.net/$USERNAME",
	"https://321sexchat.com/$USERNAME",
	"https://chat-avenue.com/$USERNAME",
	"https://freechatnow.com/$USERNAME",
	"https://chatzy.com/$USERNAME",
	"https://wireclub.com/$USERNAME",
	"https://chatib.us/$USERNAME",
	"https://chatblink.com/$USERNAME",
	"https://y99.in/$USERNAME",
	"https://chatfriends.com/$USERNAME",
	"https://talkwithstranger.com/$USERNAME",
	"https://strangercam.com/$USERNAME",
	"https://camgo.com/$USERNAME",
	"https://shychat.com/$USERNAME",
	"https://chatki.com/$USERNAME",
	"https://chatspin.com/$USERNAME",
	"https://fruzo.com/$USERNAME",
	"https://yesichat.com/$USERNAME",
	"https://chat-random.com/$USERNAME",
	"https://cameraroll.com/$USERNAME",
	"https://coco.chat/$USERNAME",
	"https://holla.world/$USERNAME",
	"https://azarcams.com/$USERNAME",
	"https://lushstories.com/profile/$USERNAME",
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
	silent := flag.Bool("silent", false, "Silent mode.")
	version := flag.Bool("version", false, "Print the version of the tool and exit.")
	nsfw := flag.Bool("nsfw", false, "Include NSFW websites in the check.")
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
		fmt.Printf("%sUsage: socialfinder [ -file <custom_file> ] [ -nsfw ] <username>%s\n", ColorRed, ColorReset)
		return
	}

	username := args[0]

	fmt.Printf("%s[*]%s Checking username %s%s%s on:\n\n", ColorYellow, ColorReset, ColorCyan, username, ColorReset)

	// Read URLs from file or default URL
	platforms, err := readURLs(*fileFlag, *nsfw, username)
	if err != nil {
		fmt.Printf("%s[!]%s Error reading URLs: %v\n", ColorRed, ColorReset, err)
		return
	}

	// Check platforms and stream results in real-time
	activeCount := checkPlatformsStream(platforms, username)

	fmt.Printf("\n%s[*]%s Search completed with %s%d%s results\n", ColorYellow, ColorReset, ColorCyan, activeCount, ColorReset)
}

func readURLs(filePath string, includeNSFW bool, username string) ([]Platform, error) {
	var platforms []Platform

	// Read from custom file or default URL
	var reader io.Reader
	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %v", filePath, err)
		}
		defer file.Close()
		reader = file
	} else {
		resp, err := http.Get(defaultURL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch default URLs: %v", err)
		}
		defer resp.Body.Close()
		reader = resp.Body
	}

	// Normalize NSFW URLs for comparison
	nsfwURLs := make(map[string]bool)
	for _, url := range nsfwWebsites {
		normalized := normalizeURL(strings.Replace(url, "$USERNAME", username, 1))
		nsfwURLs[normalized] = true
	}

	// Read URLs from the source
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			platformURL := strings.TrimSpace(parts[1])
			normalizedURL := normalizeURL(strings.Replace(platformURL, "$USERNAME", username, 1))

			// Include platform if it's not NSFW or if NSFW is allowed
			if includeNSFW || !nsfwURLs[normalizedURL] {
				platforms = append(platforms, Platform{
					Name: strings.TrimSpace(parts[0]),
					URL:  platformURL,
				})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading URLs: %v", err)
	}

	// Add NSFW platforms if -nsfw flag is set
	if includeNSFW {
		for _, url := range nsfwWebsites {
			// Extract platform name from URL (e.g., "onlyfans.com" -> "OnlyFans")
			domain := strings.Split(strings.Replace(url, "$USERNAME", username, 1), "/")[2]
			name := strings.Title(strings.Split(domain, ".")[0])
			platforms = append(platforms, Platform{
				Name: name,
				URL:  url,
			})
		}
	}

	return platforms, nil
}

// normalizeURL normalizes a URL for better matching
func normalizeURL(url string) string {
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "www.")
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