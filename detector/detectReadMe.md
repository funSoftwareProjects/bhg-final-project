This program can be used for payload detection in a png.
It searches for executable headers. Future iterations could validate by calling the validate function to pattern match the data for code once it has been converted to a string.
Usage:
- go run detector.go    --simply runs the file once. Follow the prompts to use
- go build detector.go
  .\detector.go         --creates an executable to run multiple times