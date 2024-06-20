package script

import (
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"
)

func Smb_os_discovery_scan(addr string) (string, error) {
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// Negotiate Protocol Request/判断是否有smb服务
	negotiateRequest, _ := hex.DecodeString("00000031ff534d4272000000001845680000000000000000000000000000461500000100000e00024e54204c4d20302e3132000200")
	_, err = conn.Write(negotiateRequest)
	if err != nil {
		return "", err
	}
	res, err := readDataSmb(conn)
	if err != nil {
		return "", err
	}
	if !strings.Contains(fmt.Sprintf("%x", res[4:8]), "fe534d42") && !strings.Contains(fmt.Sprintf("%x", res[4:8]), "ff534d42") {
		fmt.Println(fmt.Sprintf("%x", res))
		return "", fmt.Errorf("no smb service")
	}

	// Session Setup And Xact Secondary Request/获取操作系统信息
	sessionRequest, _ := hex.DecodeString("00000091ff534d4273000000001845680000d3cf5cf2d0e5359600000000bd07000001000cff009100ffff0100010000000000420000000000500000805600604006062b0601050502a0363034a00e300c060a2b06010401823702020aa22204204e544c4d535350000100000015820800000000000000000000000000000000004e6d6170004e6174697665204c616e6d616e0000")
	_, err = conn.Write(sessionRequest)
	if err != nil {
		return "", err
	}
	res, err = readDataSmb(conn)
	if err != nil {
		return "", err
	}

	sessionResponseContent := res[36:]
	if len(sessionResponseContent) < 4 {
		return isPrintableInfo(res), nil
	}
	securityBlobLength := bytesToInt(append([]byte{}, sessionResponseContent[8], sessionResponseContent[7]))

	//securityBlobContent := res[47 : 47+securityBlobLength]

	res = res[47+securityBlobLength:]
	var nativeOs string
	var nativeLanMan string
	for i, _ := range res {
		if res[i] == 0x00 {
			nativeOs = fmt.Sprintf("%s", res[:i])
			nativeLanMan = fmt.Sprintf("%s", res[i+1:])
			break
		}
	}
	return "OS: " + nativeOs + "\n" + "Software: " + nativeLanMan, nil
}
