### protocol_scan
使用go语言实现nmap各种协议扫描脚本/Using the go language to implement nmap various protocols scanning scripts

### now achieved
ldap-rootdse [nmap脚本地址](https://nmap.org/nsedoc/scripts/ldap-rootdse.html)

### example
```
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
```
