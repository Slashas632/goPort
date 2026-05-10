package UDP

import (
	"net"
	"port-scanner/internal/display"
	"strconv"
	"strings"
	"time"
)

const (
	UDPtimeout = 500 * time.Millisecond
)

func Udp(port int, ip string) {
	fullIp := ip + ":" + strconv.Itoa(port)

	probes, ok := PortProbes[port]
	if !ok {
		probes = GenericProbes
	}

	for _, probe := range probes {
		if banner, ok := tryProbe(fullIp, probe); ok {
			display.PrintResult(ip, port, banner)
			return
		}
	}
}

func tryProbe(addr string, probe Probe) (string, bool) {
	conn, err := net.DialTimeout("udp", addr, UDPtimeout)
	if err != nil {
		return "", false
	}
	defer conn.Close()

	conn.Write(probe.Payload)
	conn.SetReadDeadline(time.Now().Add(UDPtimeout))

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return "", false
	}
	banner := strings.SplitN(strings.TrimSpace(string(buf[:n])), "\r\n", 2)[0]
	return banner, true
}
