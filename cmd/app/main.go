package main

import (
	"fmt"
	"port-scanner/internal/cli"
	"port-scanner/internal/scanner"
)

func main() {
	opts, err := cli.ParseArgs()
	if err != nil {
		fmt.Println("Erro: ", err)
		return
	}
	scanner.Run(opts)
}
