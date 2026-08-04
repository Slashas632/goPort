package UDP

import (
	"context"
	"net"
	"github.com/Slashas632/goPort/internal/display"
	"github.com/Slashas632/goPort/internal/output"
	"github.com/Slashas632/goPort/internal/plugins"
	"github.com/Slashas632/goPort/internal/ratelimit"
	"strconv"
	"strings"
	"time"
)

var knownServices = map[int]string{
	53:    "DNS",
	111:   "RPC",
	123:   "NTP",
	137:   "NetBIOS",
	161:   "SNMP",
	500:   "IKE",
	514:   "Syslog",
	1900:  "SSDP",
	5353:  "mDNS",
	11211: "Memcached",
	51820: "WireGuard",
}

const (
	UDPtimeout = 300 * time.Millisecond
)

func Udp(port int, ip string) {
	fullIp := ip + ":" + strconv.Itoa(port)

	probes, ok := PortProbes[port]
	if !ok {
		probes = GenericProbes
	}

	for _, probe := range probes {
		if banner, ok := tryProbe(fullIp, probe); ok {
			if !isPrintable(banner) {
				if service, ok := knownServices[port]; ok {
					banner = service
				}
			}
			plugins.RunAll(ip, port, banner)
			display.PrintResult(ip, port, banner)
			if output.JSONPath != "" {
				output.AddResult(ip, port, banner)
			}
			return
		}
	}
}

func isPrintable(s string) bool {
	for _, r := range s {
		if r < 32 || r > 126 {
			return false
		}
	}
	return true
}

func tryProbe(addr string, probe Probe) (string, bool) {
	for attempt := 0; attempt < 2; attempt++ {
		ratelimit.Limiter.Wait(context.Background())
		result, ok := func() (string, bool) {
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
		}()
		if ok {
			return result, true
		}
	}
	return "", false

}
