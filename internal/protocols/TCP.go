package protocols

import (
	"fmt"
	"net"
	"strconv"
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

	fmt.Printf("Port Open %s\n", full_ip)
}
