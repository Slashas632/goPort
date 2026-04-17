package main

import (
	"port-scanner/internal/cli"
	"port-scanner/internal/scanner"
)

func main() {
	opts := cli.ParseArgs()
	scanner.Run(opts)
}
