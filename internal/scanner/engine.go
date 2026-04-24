package scanner

import (
	"fmt"
	"port-scanner/internal/cli"
	"port-scanner/internal/protocols"
	"sync"
)

func Run(opts cli.Options) {
	var wg sync.WaitGroup

	port := opts.Port
	IP := opts.IP

	if opts.TCP {

		for i := 0; i < opts.Workers; i++ {
			wg.Add(1)
			go TCPworkers(IP, port, &wg)
		}
		wg.Wait()
		fmt.Println("Work finished.")
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
