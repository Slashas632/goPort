package protocols

import (
	"fmt"
	"strconv"
	"strings"
)

func Ipscanner(Adress string) {
	fmt.Println(Adress)
	ip := strings.Split(Adress, ".")

	var start int = 0
	var end int = 21

	for start < end {
		prefix := ip[0] + "." + ip[1] + "." + ip[2] + "."

		pilnas_ip := prefix + strconv.Itoa(start)

		start++

		fmt.Println(pilnas_ip)
	}
}
