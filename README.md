# SocialFinder

A fast and efficient username enumeration tool written in Go that checks for the existence of a username across multiple social media platforms and websites. Inspired by Sherlock but optimized for performance with real-time streaming results.

## Features

- 🚀 **Real-time results** - Output streams immediately as platforms are checked
- 🎨 **Colorful output** - Easy-to-read colored terminal output
- ⚡ **High performance** - Uses httpx for fast HTTP checks
- 📋 **Custom URL lists** - Use default or custom platform lists
- 🔍 **Smart matching** - Handles URL variations and normalization
- 📊 **Result counting** - Shows total active platforms found

## Prerequisites
```
go install github.com/rix4uni/httpx@latest
```

## Installation
```
go install github.com/rix4uni/socialfinder@latest
```

## Download prebuilt binaries
```
wget https://github.com/rix4uni/socialfinder/releases/download/v0.0.2/socialfinder-linux-amd64-0.0.2.tgz
tar -xvzf socialfinder-linux-amd64-0.0.2.tgz
rm -rf socialfinder-linux-amd64-0.0.2.tgz
mv socialfinder ~/go/bin/socialfinder
```
Or download [binary release](https://github.com/rix4uni/socialfinder/releases) for your platform.

## Compile from source
```
git clone --depth 1 https://github.com/rix4uni/socialfinder.git
cd socialfinder; go install
```

## Usage
```
Usage of socialfinder:
  -file string
        Custom URLs file path
  -silent
        silent mode.
  -version
        Print the version of the tool and exit.
```

### Output Example

```yaml
▶ socialfinder rix4uni

[*] Checking username rix4uni on:

[+] Hackthebox: https://app.hackthebox.eu/profile/rix4uni
[+] BugCrowd: https://bugcrowd.com/h/rix4uni
[+] BuyMeACoffee: https://buymeacoff.ee/rix4uni
[+] Intigriti: https://app.intigriti.com/profile/rix4uni
[+] Discord: https://discord.com/users/rix4uni
[+] Twitter: https://x.com/rix4uni
[+] Hackerone: https://hackerone.com/rix4uni
[+] GitHub: https://github.com/rix4uni
[+] LinkedIn: https://www.linkedin.com/in/rix4uni
[+] DockerHub: https://hub.docker.com/u/rix4uni
[+] Medium: https://medium.com/@rix4uni
[+] PayPal: https://www.paypal.com/paypalme/rix4uni
[+] Pypi: https://pypi.org/user/rix4uni
[+] Telegram: https://t.me/rix4uni
[+] Tryhackme: https://tryhackme.com/p/rix4uni
[+] Ko-fi: https://ko-fi.com/rix4uni
[+] Replit: https://replit.com/@rix4uni
[+] Linktr: https://linktr.ee/rix4uni
[+] Twitch: https://www.twitch.tv/rix4uni
[+] Reddit: https://www.reddit.com/user/rix4uni
[+] Asciinema: https://asciinema.org/~rix4uni
[+] Strava: https://www.strava.com/athletes/rix4uni
[+] Archive: https://archive.org/details/@rix4uni
[+] YouTube: https://www.youtube.com/@rix4uni

[*] Search completed with 24 results
```

## Acknowledgments

- Inspired by [Sherlock Project](https://github.com/sherlock-project/sherlock)
- Powered by [httpx](https://github.com/projectdiscovery/httpx) from ProjectDiscovery
- Community contributions for platform lists