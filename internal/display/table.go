package display

import (
	"fmt"
	"strings"
)

const (
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Reset  = "\033[0m"
)

func PrintHeader() {
	fmt.Printf("%-10s %-20s %-8s %s\n", "STATUS", "IP", "PORT", "BANNER")
	fmt.Println(strings.Repeat("─", 80))
}

func PrintResult(ip string, port int, banner string) {
	fmt.Printf("%s%-10s%s %-20s %s%-8d%s %s%s%s\n",
		Green, "[OPEN]", Reset,
		ip,
		Cyan, port, Reset,
		Yellow, banner, Reset,
	)
}
