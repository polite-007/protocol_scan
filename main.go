package main

import (
	"fmt"
	"github.com/polite-007/protocol_scan/script"
)

func main() {
	result, err := script.Smb_protocol_scan("183.244.234.105:445")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(result)
}
