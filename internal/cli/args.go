package cli

import (
	"fmt"
	"os"
)

func Arguments() {
	protocol := os.Args[1]
	fmt.Println(protocol)
}
