package lib

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func Ldap_rootdse_scan(addr string) string {
	var res []byte
	var result string

	// 尝试TCP连接
	conn, err := net.DialTimeout("tcp", addr+":389", 5*time.Second)
	if err != nil {
		return "TCP Connection Failed"
	}
	defer conn.Close()

	// 判断LDAP服务是否开启
	_, err = conn.Write([]byte("\x30\x0c\x02\x01\x01\x60\x07\x02\x01\x03\x04\x00\x80\x00"))
	if err != nil {
		return "Bind Request Failed"
	}
	res, err = readData(conn)
	if err != nil {
		return "read error"
	}
	if !strings.Contains(fmt.Sprintf("%x", res), "070a010004000400") && !strings.Contains(fmt.Sprintf("%x", res), "616e6f6e796d6f75732062696e6420646973616c6c6f776564") {
		return ""
	}

	// 读取LDAP的数据1
	_, err = conn.Write([]byte("0%\x02\x01\x02c \x04\x00\x0a\x01\x00\x0a\x01\x00\x02\x01\x00\x02\x01\x00\x01\x01\x00\x87\x0bobjectclass0\x00"))
	if err != nil {
		return "Have LDAP Server, But No Data"
	}
	res, err = readData(conn)
	if err != nil {
		return "read error"
	}
	if fmt.Sprintf("%x", res) != "300c02010265070a010004000400" && !strings.Contains(fmt.Sprintf("%x", res), "746f70040f4f70656e4c444150726f6f74445345") {
		result, err = searchResEntryParse(res)
		if err != nil {
			fmt.Println(err)
		}
		return result
	}

	// 读取LDAP的数据2
	_, err = conn.Write([]byte("0\x82\x02\x1a\x02\x01\x03c\x82\x02\x13\x04\x00\x0a\x01\x00\x0a\x01\x00\x02\x01\x00\x02\x01\x00\x01\x01\x00\x87\x0bobjectclass0\x82\x01\xf1\x04\x1e_domainControllerFunctionality\x04\x1aconfigurationNamingContext\x04\x0bcurrentTime\x04\x14defaultNamingContext\x04\x0bdnsHostName\x04\x13domainFunctionality\x04\x0ddsServiceName\x04\x13forestFunctionality\x04\x13highestCommittedUSN\x04\x14isGlobalCatalogReady\x04\x0eisSynchronized\x04\x13ldap-get-baseobject\x04\x0fldapServiceName\x04\x0enamingContexts\x04\x17rootDomainNamingContext\x04\x13schemaNamingContext\x04\x0aserverName\x04\x11subschemaSubentry\x04\x15supportedCapabilities\x04\x10supportedControl\x04\x15supportedLDAPPolicies\x04\x14supportedLDAPVersion\x04\x17supportedSASLMechanisms\x04\x09altServer\x04\x12supportedExtension"))
	if err != nil {
		return "Have LDAP Server, But No Data"
	}

	for fmt.Sprintf("%x", res) == "300c02010265070a010004000400" || strings.Contains(fmt.Sprintf("%x", res), "746f70040f4f70656e4c444150726f6f74445345") {
		res, err = readData(conn)
		if err != nil {
			return "read error"
		}
		if len(res) == 0 {
			return "Have LDAP Server, But No Data"
		}
	}
	result, err = searchResEntryParse(res)
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func searchResEntryParse(data []byte) (string, error) {
	var searchResEntry []byte
	var contentAll string
	var content string

	//替换掉可能出现的searchResDne数据
	if strings.Contains(fmt.Sprintf("%x", data), "070a010004000400") {
		for i := len(data); ; i-- {
			if data[i-1] == 0x30 {
				data = data[:i-1]
				break
			}
		}
	}

	if len(data) == 0 {
		return "Have LDAP Server, But No Data", nil
	}

	// 获取searchResEntry数据
	numberOne, _ := strconv.Atoi(fmt.Sprintf("%x", data[1]))
	if numberOne < 81 || numberOne > 89 {
		return "", fmt.Errorf("searchResEntryParse: %s", "data is not searchResEntry")
	} else {
		searchResEntry = data[numberOne-75:]
	}

	// 解析出LDAP的attributes数据
	cname, str := obtainObject(searchResEntry)
	if cname == "" {
		fmt.Println("objectname:null")
	} else {
		fmt.Println("objectname:" + cname)
	}

	// 解析出LDAP的type数据和vals
	for len(str) != 0 {
		content, str = obtainAttribute(str)
		contentAll += content + "\n"
	}
	return contentAll, nil
}

// 解析出LDAP的type数据和vals
func obtainAttribute(data []byte) (string, []byte) {
	var content string

	numberOne, _ := strconv.Atoi(fmt.Sprintf("%x", data[1]))
	if numberOne < 81 || numberOne > 89 {
		data = data[3:]
	} else {
		data = data[numberOne-77:]
	}

	lengthType, _ := strconv.ParseInt(fmt.Sprintf("%x", data[0]), 16, 64)
	valueType := fmt.Sprintf("%s", data[1:lengthType+1])
	data = data[lengthType+1:]
	content = valueType + ":\n "
	numberOne, _ = strconv.Atoi(fmt.Sprintf("%x", data[1]))
	if numberOne < 81 || numberOne > 89 {
		data = data[2:]
	} else {
		data = data[numberOne-78:]
	}

	for {
		if len(data) == 0 {
			break
		}
		if data[0] != 0x04 {
			break
		}
		numberOne, _ = strconv.Atoi(fmt.Sprintf("%x", data[1]))
		if numberOne >= 81 && numberOne <= 89 {
			numberTwo := int(bytesToInt(data[2 : numberOne-78]))
			Value := isPrintableInfo(data[numberOne-78 : numberTwo+numberOne-78])
			data = data[numberTwo+numberOne-78:]
			content += Value + "\n "
		} else {
			lengthValue, _ := strconv.ParseInt(fmt.Sprintf("%x", data[1]), 16, 64)
			Value := fmt.Sprintf("%s", data[2:lengthValue+2])
			data = data[lengthValue+2:]
			content += Value + "\n "
		}
	}
	return strings.TrimRight(content, "\n "), data
}

// 解析出LDAP的attributes数据
func obtainObject(data []byte) (string, []byte) {
	var cname string
	numberOne, _ := strconv.Atoi(fmt.Sprintf("%x", data[1]))
	if numberOne < 81 || numberOne > 89 {
		return "", data[4:]
	} else {
		data = data[numberOne-77:]
		length, _ := strconv.ParseInt(fmt.Sprintf("%x", data[0]), 16, 64)
		cname = fmt.Sprintf("%s", data[1:length+1])
		data = data[length+1:]
	}
	numberOne, _ = strconv.Atoi(fmt.Sprintf("%x", data[1]))
	if numberOne < 81 || numberOne > 89 {
		return "", data[4:]
	} else {
		return cname, data[numberOne-78:]
	}
}

func readData(conn net.Conn) ([]byte, error) {
	bufferBind := make([]byte, 4096)
	r, err := conn.Read(bufferBind)
	if err != nil {
		panic(err)
	}
	return bufferBind[:r], nil
}

func bytesToInt(b []byte) uint64 {
	var result uint64
	for _, byteVal := range b {
		result = (result << 8) | uint64(byteVal)
	}
	return result
}

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
