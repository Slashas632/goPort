package scanner

import (
	"fmt"
	"port-scanner/internal/cli"
)

func Run(opts cli.Options) {
	fmt.Println(opts.IP)
}
