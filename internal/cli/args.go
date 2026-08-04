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
	Install   string
	Uninstall string
	Json      string
}

func ParseArgs() (Options, error) {
	tcp := flag.Bool("tcp", false, "TCP scan")
	udp := flag.Bool("udp", false, "UDP scan")
	ip := flag.String("ip", "", "IP adress")
	port := flag.String("p", "", "Port")
	workers := flag.Int("w", 500, "Workers")
	install := flag.String("install", "", "Install a plugin (lua file)")
	uninstall := flag.String("uninstall", "", "Uninstall a plugin (lua file)")
	json := flag.String("json", "", "JSON output file")

	flag.Parse()

	startPort, endPort, err := portCheck(*port)
	if err != nil && (*install == "" && *uninstall == "") {
		return Options{}, fmt.Errorf("invalid port: %w", err)
	}

	return Options{
		TCP:       *tcp,
		UDP:       *udp,
		IP:        *ip,
		Workers:   *workers,
		StartPort: startPort,
		EndPort:   endPort,
		Install:   *install,
		Uninstall: *uninstall,
		Json:      *json,
	}, nil
}

func portCheck(port string) (int, int, error) {
	if port == "" {
		return 0, 0, nil
	}
	if strings.Contains(port, "-") {
		var port_split = strings.Split(port, "-")

		startPortInt, err := strconv.Atoi(port_split[0])
		if err != nil {
			return 0, 0, fmt.Errorf("'%s' is not a valid port number", port_split[0])
		}

		endPortInt, err := strconv.Atoi(port_split[1])
		if err != nil {
			return 0, 0, fmt.Errorf("'%s' is not a valid port number", port_split[1])
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
		return 0, 0, fmt.Errorf("'%s' is not a valid port number", port)
	}
	if singlePort < 0 || singlePort > 65535 {
		return 0, 0, fmt.Errorf("Port must be between 0-65535")
	}
	return singlePort, singlePort, nil
}
