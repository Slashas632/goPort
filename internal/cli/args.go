package cli

import (
	"flag"
)

type Options struct {
	TCP  bool
	UDP  bool
	Port int
}

func ParseArgs() Options {
	tcp := flag.Bool("tcp", false, "TCP scan")
	udp := flag.Bool("udp", false, "UDP scan")
	port := flag.Int("p", 80, "port")

	flag.Parse()

	return Options{
		TCP:  *tcp,
		UDP:  *udp,
		Port: *port,
	}
}
