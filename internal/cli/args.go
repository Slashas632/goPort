package cli

import (
	"flag"
)

type Options struct {
	TCP  bool
	UDP  bool
	IPS  bool
	IP   string
	Port int
}

func ParseArgs() Options {
	tcp := flag.Bool("tcp", false, "TCP scan")
	udp := flag.Bool("udp", false, "UDP scan")
	ips := flag.Bool("ips", false, "IP scan")
	ip := flag.String("ip", "127.0.0.1", "IP adress")
	port := flag.Int("p", 80, "port")

	flag.Parse()

	return Options{
		TCP:  *tcp,
		UDP:  *udp,
		Port: *port,
		IPS:  *ips,
		IP:   *ip,
	}
}
