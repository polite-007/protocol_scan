package script

import (
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"
)

func Smb_protocol_scan(addr string) (string, error) {
	var versionList string
	var payloadListMap = map[string]string{
		"2.0.2": "00000066fe534d424000000000000000000000000000000000000000000000000000000000000000000000000000000000000000313233343536373839303132333435362400010001000000000000003132333435363738393031323334353600000000000000000202",
		"2.1.0": "00000066fe534d424000000000000000000000000000000000000000000000000000000000000000000000000000000000000000313233343536373839303132333435362400010001000000000000003132333435363738393031323334353600000000000000001002",
		"3.0.0": "00000066fe534d424000000000000000000000000000000000000000000000000000000000000000000000000000000000000000313233343536373839303132333435362400010001000000000000003132333435363738393031323334353600000000000000000003",
		"3.0.2": "00000066fe534d424000000000000000000000000000000000000000000000000000000000000000000000000000000000000000313233343536373839303132333435362400010001000000000000003132333435363738393031323334353600000000000000000203",
		"3.1.1": "00000066fe534d424000000000000000000000000000000000000000000000000000000000000000000000000000000000000000313233343536373839303132333435362400010001000000000000003132333435363738393031323334353600000000000000001103",
		"all":   "000000b4fe534d424000000000000000000000000000000000000000000000000000000000000000000000000000000000000000313233343536373839303132333435362400050001000000000000003132333435363738393031323334353670000000020000000202100200030203110300000200060000000000020002000100000001002c00000000000200020001000100200001000000000000000000000000000000000001000000000000000000000000000000",
	}
	var versionLists = map[string]string{
		"1103": "3.1.1",
		"0203": "3.0.2",
		"0003": "3.0.0",
		"1002": "2.1.0",
		"0202": "2.0.2",
	}
	var versionListArray = []string{"2.0.2", "2.1.0", "3.0.0", "3.0.2", "3.1.1"}
	// 尝试TCP连接
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// 判断SMB服务是否开启和版本
	payloadAll, _ := hex.DecodeString(payloadListMap["all"])
	_, err = conn.Write(payloadAll)
	if err != nil {
		return "", err
	}
	res, err := readData(conn)
	if err != nil {
		return "", err
	}
	if !strings.Contains(fmt.Sprintf("%x", res), "fe534d42") {
		return "", fmt.Errorf("no smb service")
	}
	version := versionLists[fmt.Sprintf("%x", res[72:74])]
	if version == "" {
		return "", fmt.Errorf("smb version contain fail")
	}

	for _, i := range versionListArray {
		if version == i {
			break
		}
		payload, _ := hex.DecodeString(payloadListMap[i])

		conn, err = net.DialTimeout("tcp", addr, 5*time.Second)
		if err != nil {
			return "", err
		}
		_, err = conn.Write(payload)
		if err != nil {
			return "", err
		}
		res, err = readData(conn)
		if err != nil {
			return "", err
		}
		if !strings.Contains(fmt.Sprintf("%x", res), "fe534d42") || versionLists[fmt.Sprintf("%x", res[72:74])] == "" {
			continue
		}
		versionList += " " + versionLists[fmt.Sprintf("%x", res[72:74])] + "\n"
	}
	return "NT LM 0.12 (SMBv1)\n" + versionList + " " + version, err
}
