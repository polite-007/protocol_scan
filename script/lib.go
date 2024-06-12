package script

import (
	"fmt"
	"io"
	"net"
	"time"
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

func readData(conn net.Conn) ([]byte, error) {
	//读取数据
	var buf []byte              // big buffer
	var tmp = make([]byte, 256) // using small tmp buffer for demonstrating
	//设置读取超时Deadline
	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 3))
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
		if len(buf) > 4096 {
			break
		}
		_ = conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	}
	return buf, nil
}

func safeSlice(data []byte, lengthType int) []byte {
	if lengthType >= 0 && lengthType < len(data) {
		return data[lengthType:]
	} else {
		// 指定返回值，比如返回一个空切片或者错误信息
		return []byte{} // 或者 return nil, errors.New("index out of range")
	}
}
