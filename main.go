package main

import (
	"fmt"
	"protocol_scan/script"
)

func main() {
	result, err := script.Ldap_rootdse_scan("60.204.149.158:389")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(result)
}
