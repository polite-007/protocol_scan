package script

import (
	"fmt"
	"net"
)

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

// 从TCP中读取数据
func readData(conn net.Conn) ([]byte, error) {
	bufferBind := make([]byte, 4096)
	r, err := conn.Read(bufferBind)
	if err != nil {
		panic(err)
	}
	return bufferBind[:r], nil
}
