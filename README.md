# clickjackingXD
Detect clickjacking vulnerablities from list of urls. It simply check the reponse if it has the ```X-Frame-Options``` in the header. If not present, it will output that url to stdout and you can take a screenshot of that to see if it is worth reporting. It takes input from stdin

# How to install
```go get github.com/noobexploiter/clickjackingXD```
