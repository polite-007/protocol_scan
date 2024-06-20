package script

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func readDataSmb(conn net.Conn) ([]byte, error) {
	var bufAll []byte
	var smbFirst = make([]byte, 4)
	_, err := io.ReadFull(conn, smbFirst)
	if err == io.EOF {
		return nil, fmt.Errorf("no data on tcp or premature EOF")
	} else if err != nil {
		return nil, err
	}
	bufAll = append(bufAll, smbFirst...)
	smbTwo := make([]byte, bytesToInt(smbFirst))
	for {
		n, err := io.ReadFull(conn, smbTwo)
		if n > 0 {
			bufAll = append(bufAll, smbTwo[:n]...)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error reading from connection: %v", err)
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}
	return bufAll, nil
}

func readDataLdap(conn net.Conn) ([]byte, error) {
	var bufAll []byte
	var ldapFirst = make([]byte, 2)
	_, err := io.ReadFull(conn, ldapFirst)
	if err == io.EOF {
		return nil, fmt.Errorf("no data on tcp or premature EOF")
	} else if err != nil {
		return nil, err
	}
	ldapNumber := int(ldapFirst[1]) - 48
	var ldapTwoLen int
	if ldapNumber >= 81 && ldapNumber <= 89 {
		var ldapTwo = make([]byte, ldapNumber-80)
		_, err = io.ReadFull(conn, ldapTwo)
		if err != nil {
			return nil, err
		}
		ldapTwoLen = int(bytesToInt(ldapTwo))
		bufAll = append(bufAll, ldapFirst...)
		bufAll = append(bufAll, ldapTwo...)
	} else {
		ldapTwoLen = int(bytesToInt(ldapFirst[1:]))
		bufAll = append(bufAll, ldapFirst...)
	}
	var ldapThree = make([]byte, ldapTwoLen)
	for {
		n, err := io.ReadFull(conn, ldapThree)
		if n > 0 {
			bufAll = append(bufAll, ldapThree[:n]...)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error reading from connection: %v", err)
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}
	return bufAll, nil
}

// 判断是否为可打印字符
func isPrintableInfo(bytes []byte) string {
	str := ""
	for _, b := range bytes {
		if b >= 32 && b <= 126 {
			str += fmt.Sprintf("%c", b)
		} else {
			str += fmt.Sprintf("\\x%02X", b)
		}
	}
	return str
}

// 将字节数组转换为整数
func bytesToInt(b []byte) uint64 {
	var result uint64
	for _, byteVal := range b {
		result = (result << 8) | uint64(byteVal)
	}
	return result
}
