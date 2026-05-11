package main

import (
	"fmt"
	"port-scanner/internal/cli"
	"port-scanner/internal/display"
	"port-scanner/internal/scanner"
)

func main() {
	opts, err := cli.ParseArgs()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	display.PrintHeader()
	scanner.Run(opts)
}
