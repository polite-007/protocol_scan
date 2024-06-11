package script

import (
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"
)

func Smb_protocol_scan(addr string) (string, error) {
	var res []byte

	// 尝试TCP连接
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// 判断LDAP服务是否开启
	payload, _ := hex.DecodeString("000000b4fe534d424000000000000000000000000000000000000000000000000000000000000000000000000000000000000000313233343536373839303132333435362400050001000000000000003132333435363738393031323334353670000000020000000202100200030203110300000200060000000000020002000100000001002c00000000000200020001000100200001000000000000000000000000000000000001000000000000000000000000000000")
	_, err = conn.Write(payload)
	if err != nil {
		return "", err
	}
	res, err = readData(conn)
	if err != nil {
		return "", err
	}
	if !strings.Contains(fmt.Sprintf("%x", res), "fe534d42") {
		return "", fmt.Errorf("smb-protocol scan failed")
	}
	return fmt.Sprintf("%x", res), nil
}
