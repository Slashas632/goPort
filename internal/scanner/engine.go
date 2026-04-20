package scanner

import (
	"port-scanner/internal/cli"
	"port-scanner/internal/protocols"
)

func Run(opts cli.Options) {
	if opts.IPS {
		protocols.Ipscanner(opts.IP)
	}
}
