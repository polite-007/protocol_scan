package main

import (
	"fmt"
	"protocol_scan/script"
)

func main() {
	result, err := script.Smb_protocol_scan("39.175.75.67:445")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(result)
}
