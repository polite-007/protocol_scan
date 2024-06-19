package script

import (
	"fmt"
	"io"
	"net"
	"strconv"
)

func readData(conn net.Conn) ([]byte, error) {
	//读取数据
	var buf []byte
	var tmp = make([]byte, 256)
	//循环读取数据
	for {
		length, err := conn.Read(tmp)
		buf = append(buf, tmp[:length]...)
		if length < len(tmp) {
			break
		}
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
	}
	return buf, nil
}

func readDataSmb(conn net.Conn) ([]byte, error) {
	var tmp = make([]byte, 4)
	var smbContent []byte
	_, err := conn.Read(tmp)
	if err != nil || len(tmp) < 4 {
		return nil, err
	}
	if len(tmp) < 4 {
		return nil, fmt.Errorf("smb length too short")
	}
	lengthSmb := bytesToInt(tmp)
	for {
		var tmpLice = make([]byte, 256)
		length, err := conn.Read(tmpLice)
		smbContent = append(smbContent, tmpLice[:length]...)
		if length < len(tmp) {
			break
		}
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		if len(smbContent) >= int(lengthSmb) {
			break
		}
	}
	return append(tmp, smbContent[:lengthSmb]...), nil
}

func readDataLdap(conn net.Conn) ([]byte, error) {
	var tmp = make([]byte, 2)
	_, err := conn.Read(tmp)
	if err != nil || len(tmp) < 2 {
		return nil, fmt.Errorf("read fail or ldap length too short")
	}

	numberLdap, _ := strconv.Atoi(fmt.Sprintf("%x", tmp[1]))
	if numberLdap > 80 && numberLdap < 90 {
		var bufAll []byte
		var allLength = make([]byte, numberLdap-80)
		tmpLength, err := conn.Read(allLength)
		if err != nil || tmpLength == 0 {
			return nil, err
		}
		contentLength := int(bytesToInt(allLength))
		tmpLengthValue := allLength
		var ldapContent []byte
		for {
			var tmpLice = make([]byte, 256)
			length, err := conn.Read(tmpLice)
			ldapContent = append(ldapContent, tmpLice[:length]...)
			if length < len(tmp) {
				break
			}
			if err != nil {
				if err != io.EOF {
					return nil, err
				}
				break
			}
			if len(ldapContent) >= contentLength {
				break
			}
		}
		bufAll = append(bufAll, tmp...)
		bufAll = append(bufAll, tmpLengthValue...)
		bufAll = append(bufAll, ldapContent[:contentLength]...)
		return bufAll, err
	} else {
		var bufAll = make([]byte, 4096)
		length, err := conn.Read(bufAll)
		if err != nil || length == 0 {
			return nil, err
		}
		return bufAll[:length], nil
	}
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
