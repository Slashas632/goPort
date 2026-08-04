package scanner

import (
	"fmt"
	"port-scanner/internal/cli"
	"port-scanner/internal/display"
	"port-scanner/internal/output"
	"port-scanner/internal/plugins"
	tcp "port-scanner/internal/protocols/TCP"
	udp "port-scanner/internal/protocols/UDP"
	"port-scanner/internal/ratelimit"
	"sync"
)

const (
	workerChannelMultiplier = 10
)

func Run(opts cli.Options) {
	if !ErrorHandler(opts) {
		return
	}

	output.Init(opts.Json)

	display.PrintHeader()

	var wg sync.WaitGroup
	ratelimit.Init(opts.Workers)

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
	if opts.Json != "" {
		if err := output.SaveJson(opts.Json); err != nil {
			fmt.Println("Error saving JSON:", err)
		}
	}
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

func ErrorHandler(opts cli.Options) bool {
	if opts.Install != "" {
		if err := plugins.Install(opts.Install); err != nil {
			fmt.Printf("Plugin installation failed: %s\n", err)
		}
		return false
	}
	if opts.Uninstall != "" {
		if err := plugins.Uninstall(opts.Uninstall); err != nil {
			fmt.Printf("Plugin uninstallation failed: %s\n", err)
		}
		return false
	}
	if !opts.TCP && !opts.UDP {
		fmt.Println("Error: specify -tcp and/or -udp")
		return false
	}
	if opts.IP == "" {
		fmt.Println("Error: specify -ip")
		return false
	}
	if opts.StartPort == 0 && opts.EndPort == 0 {
		fmt.Println("Error: specify -p")
		return false
	}
	return true
}
