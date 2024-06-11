package main

import (
	"fmt"
	"protocol_scan/lib"
)

func main() {
	fmt.Println(lib.Ldap_rootdse_scan("84.46.250.218"))
}
