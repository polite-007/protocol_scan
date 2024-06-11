package main

import (
	"fmt"
	"protocol_scan/lib"
)

func main() {
	fmt.Println(lib.Ldap_rootdse_scan("201.76.172.51"))
}
