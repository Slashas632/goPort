package main

import (
	"fmt"
	"github.com/Slashas632/goPort/internal/cli"
	"github.com/Slashas632/goPort/internal/scanner"
)

func main() {
	opts, err := cli.ParseArgs()
	if err != nil {
		fmt.Println("Error: ", err)
		fmt.Println("\nExample: goPort -tcp -ip 10.0.0.1 -p 0-1024")
		return
	}
	scanner.Run(opts)
}
