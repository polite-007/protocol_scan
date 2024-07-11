### protocol_scan
使用go语言手搓nmap各种协议扫描脚本，不使用nmap接口/Using the go language hand rubbed nmap various protocols scanning scripts, do not use the nmap interface
* [nmap脚本地址](https://nmap.org/nsedoc/scripts)
### now supported
* ldap-rootdse
* smb-protocols
* smb-os-discovery 
* smb-ls 正在编写中


### example
```
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

```