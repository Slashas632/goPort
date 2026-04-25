package scanner

import (
	"fmt"
	"port-scanner/internal/cli"
	"port-scanner/internal/protocols"
	"strconv"
	"strings"
	"sync"
)

func Run(opts cli.Options) {
	var wg sync.WaitGroup

	port := opts.Port
	IP := opts.IP

	PortCheck(port)

	if opts.TCP {

		for i := 0; i < opts.Workers; i++ {
			wg.Add(1)
			go TCPworkers(IP, port, &wg)
		}
		wg.Wait()
		fmt.Println("Work finished.")
	}
}

func PortCheck(port int) (int, int, error) {

	var port_strings string = strconv.Itoa(port)

	if strings.Contains(port_strings, "-") {
		var port_split = strings.Split(port_strings, "-")

		startPortInt, err := strconv.Atoi(port_split[0])
		if err != nil {
			return 0, 0, fmt.Errorf("neteisingas pradžios portas: %w", err)
		}

		endPortInt, err := strconv.Atoi(port_split[1])
		if err != nil {
			return 0, 0, fmt.Errorf("neteisingas pabaigos portas: %w", err)
		}

		return startPortInt, endPortInt, nil
	}
}

func TCPworkers(ip string, port int, wg *sync.WaitGroup) {
	defer wg.Done()
	var ports int = 0
	for ports < port {
		ports++
		protocols.Tcp(ports, ip)
	}
}
