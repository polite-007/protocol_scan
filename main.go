package main

import (
	"fmt"
	"protocol_scan/script"
)

func main() {
	result, err := script.Smb_os_discovery_scan("183.244.234.105:445")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(result)
}
