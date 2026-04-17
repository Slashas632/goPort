package scanner

import (
	"port-scanner/internal/cli"
)

func Run(opts cli.Options) {

}

/*
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
*/
