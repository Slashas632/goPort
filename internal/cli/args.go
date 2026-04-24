package cli

import (
	"flag"
)

type Options struct {
	TCP     bool
	UDP     bool
	IP      string
	Port    int
	Workers int
}

func ParseArgs() Options {
	tcp := flag.Bool("tcp", false, "TCP scan")
	udp := flag.Bool("udp", false, "UDP scan")
	ip := flag.String("ip", "127.0.0.1", "IP adress")
	port := flag.Int("p", 80, "port")
	workers := flag.Int("w", 100, "workers")

	flag.Parse()

	return Options{
		TCP:     *tcp,
		UDP:     *udp,
		Port:    *port,
		IP:      *ip,
		Workers: *workers,
	}
}
