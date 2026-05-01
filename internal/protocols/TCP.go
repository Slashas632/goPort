package protocols

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func Tcp(port int, ip string) {
	full_ip := ip + ":" + strconv.Itoa(port)
	timeout := time.Second * 2
	conn, err := net.DialTimeout("tcp", full_ip, timeout)
	if err != nil {
		return
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(timeout))

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
		conn.SetReadDeadline(time.Now().Add(timeout))
		n, err = conn.Read(buf)
		if err != nil {
			return
		}
	}

	banner := strings.TrimSpace(string(buf[:n]))
	firstLine := strings.SplitN(banner, "\n", 2)[0]

	const (
		Green  = "\033[32m"
		Yellow = "\033[33m"
		Cyan   = "\033[36m"
		Reset  = "\033[0m"
	)

	fmt.Printf("%s[OPEN]%s %s | %s%s%s\n", Green, Reset, full_ip, Yellow, firstLine, Reset)

}
