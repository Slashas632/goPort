package scanner

import (
	"fmt"
	"port-scanner/internal/cli"
	tcp "port-scanner/internal/protocols/TCP"
	udp "port-scanner/internal/protocols/UDP"
	"sync"
)

const (
	workerChannelMultiplier = 10
)

func Run(opts cli.Options) {
	var wg sync.WaitGroup

	if opts.TCP {
		TCPports := make(chan int, opts.Workers*workerChannelMultiplier)
		for i := 0; i < opts.Workers; i++ {
			wg.Add(1)
			go TCPworkers(TCPports, opts.IP, &wg)
		}
		go func() {
			for port := opts.StartPort; port <= opts.EndPort; port++ {
				TCPports <- port
			}
			close(TCPports)
		}()
	}

	if opts.UDP {
		UDPports := make(chan int, opts.Workers*workerChannelMultiplier)
		for i := 0; i < opts.Workers; i++ {
			wg.Add(1)
			go UDPworkers(UDPports, opts.IP, &wg)
		}
		go func() {
			for port := opts.StartPort; port <= opts.EndPort; port++ {
				UDPports <- port
			}
			close(UDPports)
		}()
	}
	wg.Wait()
	fmt.Println("Work finished.")
}

func TCPworkers(ports <-chan int, ip string, wg *sync.WaitGroup) {
	defer wg.Done()
	for port := range ports {
		tcp.Tcp(port, ip)
	}
}

func UDPworkers(ports <-chan int, ip string, wg *sync.WaitGroup) {
	defer wg.Done()
	for port := range ports {
		udp.Udp(port, ip)
	}
}
