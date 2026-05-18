package TCP

import (
	"context"
	"net"
	"port-scanner/internal/display"
	"port-scanner/internal/output"
	"port-scanner/internal/plugins"
	"port-scanner/internal/ratelimit"
	"strconv"
	"strings"
	"time"
)

const (
	TCPtimeout = time.Second * 1
)

func Tcp(port int, ip string) {
	ratelimit.Limiter.Wait(context.Background())
	fullIp := ip + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout("tcp", fullIp, TCPtimeout)
	if err != nil {
		return
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(TCPtimeout))

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		if _, err := conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n")); err != nil {
			return
		}
		conn.SetReadDeadline(time.Now().Add(TCPtimeout))
		n, err = conn.Read(buf)
		if err != nil {
			return
		}
	}

	banner := strings.SplitN(strings.TrimSpace(string(buf[:n])), "\r\n", 2)[0]

	plugins.RunAll(ip, port, banner)
	display.PrintResult(ip, port, banner)
	if output.JSONPath != "" {
		output.AddResult(ip, port, banner)
	}

}
