package scanner

import (
	"fmt"
	"port-scanner/internal/cli"
	"port-scanner/internal/protocols"
	"sync"
)

func Run(opts cli.Options) {
	const buffer = 10
	ports := make(chan int, opts.Workers*buffer)
	var wg sync.WaitGroup

	if opts.TCP {
		for i := 0; i < opts.Workers; i++ {
			wg.Add(1)
			go TCPworkers(ports, opts.IP, &wg)
		}
		for port := opts.StartPort; port <= opts.EndPort; port++ {
			ports <- port
		}
		close(ports)
		wg.Wait()
		fmt.Println("Work finished.")
	}
}

func TCPworkers(ports <-chan int, ip string, wg *sync.WaitGroup) {
	defer wg.Done()
	for port := range ports {
		protocols.Tcp(port, ip)
	}
}
