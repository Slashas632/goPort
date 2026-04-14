package main

import (
	"fmt"
	"strconv"
	"strings"
)

func ip_scanner() {

	// IP

	var ip string
	fmt.Println("Enter IP which you want to scan: ")
	fmt.Scan(&ip)
	ip_parts := strings.Split(ip, ".")

	// Diapazonas

	var dp string
	fmt.Println("Pasirinkite diapazoną kurį norite nuskanuoti")
	fmt.Scan(&dp)

	diapazonas := strings.Split(dp, "-")
	start, _ := strconv.Atoi(diapazonas[0])
	end, _ := strconv.Atoi(diapazonas[1])

	for start < end {
		prefix := ip_parts[0] + "." + ip_parts[1] + "." + ip_parts[2] + "."

		pilnas_ip := prefix + "." + strconv.Itoa(start)

		start++
	}

	fmt.Println("Your ip is: ", ip_parts[3])
	fmt.Println(prefix)
}

func main() {

	var i int

	fmt.Print("Select your software: ")
	fmt.Scan(&i)

	if i == 1 {
		fmt.Println("Your software is ip scanner")
		ip_scanner()
	} else {
		fmt.Println("This program doesn't exist")
	}

}
