# SocialFinder
 
## Usage `replace rix4uni to your username`
```
for url in $(cat urls.txt); do echo "$url" | sed 's/$/rix4uni/' | httpx -mc 200 -random-agent -silent;done
```
