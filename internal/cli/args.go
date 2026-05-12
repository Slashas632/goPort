package cli

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type Options struct {
	TCP       bool
	UDP       bool
	IP        string
	Workers   int
	StartPort int
	EndPort   int
}

func ParseArgs() (Options, error) {

	tcp := flag.Bool("tcp", false, "TCP scan")
	udp := flag.Bool("udp", false, "UDP scan")
	ip := flag.String("ip", "127.0.0.1", "IP adress")
	port := flag.String("p", "65535", "port")
	workers := flag.Int("w", 500, "workers")

	flag.Parse()

	startPort, endPort, err := portCheck(*port)
	if err != nil {
		return Options{}, fmt.Errorf("invalid port: %w", err)
	}

	return Options{
		TCP:       *tcp,
		UDP:       *udp,
		IP:        *ip,
		Workers:   *workers,
		StartPort: startPort,
		EndPort:   endPort,
	}, nil
}

func portCheck(port string) (int, int, error) {
	if strings.Contains(port, "-") {
		var port_split = strings.Split(port, "-")

		startPortInt, err := strconv.Atoi(port_split[0])
		if err != nil {
			return 0, 0, fmt.Errorf("Bad start port: %w", err)
		}

		endPortInt, err := strconv.Atoi(port_split[1])
		if err != nil {
			return 0, 0, fmt.Errorf("Bad end port: %w", err)
		}
		if startPortInt < 0 || endPortInt > 65535 {
			return 0, 0, fmt.Errorf("Port must be between 0-65535")
		}
		if startPortInt > endPortInt {
			return 0, 0, fmt.Errorf("Start port must be <= end port")
		}
		return startPortInt, endPortInt, nil
	}

	singlePort, err := strconv.Atoi(port)

	if err != nil {
		return 0, 0, fmt.Errorf("Bad port: %w", err)
	}
	return singlePort, singlePort, nil
}
