package main

import (
	"flag"
	"fmt"
	"protocol_scan/script"
	"sync"
)

var wg sync.WaitGroup // 用于等待所有goroutine完成

func scanLdap(server string, results chan<- map[string]interface{}, scriptName string) {
	defer wg.Done() // 完成一个goroutine时减计数器
	scriptLists := map[string]func(string) (string, error){
		"ldap_rootdse":  script.Ldap_rootdse_scan,
		"smb_protocols": script.Smb_protocol_scan,
	}
	scriptname := scriptLists[scriptName]
	if scriptname == nil {
		results <- map[string]interface{}{server: fmt.Errorf("脚本不存在")}
	}
	result, err := scriptname(server)
	if err != nil {
		results <- map[string]interface{}{server: err}
	} else {
		results <- map[string]interface{}{server: len(result)}
	}
}

func main() {
	//servers := []string{"45.151.2.178:389", "another.server:389", "yet.another:389"} // LDAP服务器列表
	fmt.Println("现支持脚本名称: ldap_rootdse,smb_protocols")
	host := flag.String("host", "", "<host>:<port>")
	scriptName := flag.String("script", "", "扫描的脚本名称")
	number := flag.Int("number", 30, "测试线程数")
	flag.Parse()

	//判断输入的脚本名称是否存在
	scriptLists := map[string]func(string) (string, error){
		"ldap_rootdse":  script.Ldap_rootdse_scan,
		"smb_protocols": script.Smb_protocol_scan,
	}
	scriptname := scriptLists[*scriptName]
	if scriptname == nil {
		fmt.Println("脚本不存在")
		return
	}

	// 检查输入参数
	if *host == "" || *scriptName == "" {
		fmt.Println("请输入完整的参数")
		return
	}

	results := make(chan map[string]interface{}) // 创建通道用于传递结果

	// 启动指定的goroutine数量进行并发测试扫描
	for i := 1; i <= *number; i++ {
		wg.Add(1) // 对每个goroutine增加计数
		go func(host string, scriptName string) {
			scanLdap(host, results, scriptName)
		}(*host, *scriptName)
	}
	go func() {
		wg.Wait()      // 等待所有goroutine完成
		close(results) // 所有任务完成后关闭通道
	}()

	// 收集并打印结果
	for result := range results {
		for server, res := range result {
			if err, ok := res.(error); ok {
				fmt.Printf("Server %s Error: %v\n", server, err)
			} else {
				fmt.Printf("Server %s Result: %v\n", server, res)
			}
		}
	}
}
