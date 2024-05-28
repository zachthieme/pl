This is a small POC of using go with paralellization to run generate png labels. On my tiny single board computer it tool 11 ms to generate 5 labels.

1. install go: https://go.dev/doc/install
1. clone the repo locally
  `git clone https://github.com/zachthieme/pl`
2. install the required graphics library
  `go get -u github.com/fogleman/gg`
3. compile the binary
  `go build -o pl`
4. run the binary
  `./pl order.json`
