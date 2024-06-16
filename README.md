### protocol_scan
使用go语言手搓nmap各种协议扫描脚本，不使用nmap接口/Using the go language to implement nmap various protocols scanning scripts

### now achieved
* ldap-rootdse [nmap脚本地址](https://nmap.org/nsedoc/scripts/ldap-rootdse.html)
* smb-protocols 正在编写中

### example
```
package main

import (
	"fmt"
	"protocol_scan/script"
)

func main() {
	result, err := script.Ldap_rootdse_scan("69.39.49.200:389")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(result)
}
```
