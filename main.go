package main

import (
	"Fofa_scan/convert_model"
	"Fofa_scan/parse_model"
	"Fofa_scan/scan_model"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var jsonDir = "json_report"    // JSON 文件存储路径
var outputFile = "reports.csv" // 转换后的 CSV 文件

// 主函数
func main() {
	if err := os.MkdirAll(jsonDir, 0755); err != nil {
		log.Fatalf("创建 %s 失败: %v\n", jsonDir, err)
	}
	fmt.Printf("成功创建目录: %s\n", jsonDir)

	logo := `
    __  ___                                         _ 
   /  |/  /____ _ _____ ____ _ _____ ____   ____   (_)
  / /|_/ // __ ` + "`" + `// ___// __ ` + "`" + `// ___// __ \ / __ \ / / 
 / /  / // /_/ // /__ / /_/ // /   / /_/ // / / // /  
/_/  /_/ \__,_/ \___/ \__,_//_/    \____//_/ /_//_/   
                                                      
    `
	fmt.Println(logo)

	var choice int
	fmt.Println("请选择查询方式:")
	fmt.Println("1. 单个 IP 查询")
	fmt.Println("2. 批量 IP 查询 (从 .txt 文件读取)")
	fmt.Print("请输入你的选择 (1 或 2): ")
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		handleSingleIPQuery()
	case 2:
		handleBatchIPQuery()
	default:
		fmt.Println("无效输入, 请输入 1 或 2.")
	}
}

// 单个 IP 查询的处理
func handleSingleIPQuery() {
	var ip string
	fmt.Print("请输入 IP 地址: ")
	fmt.Scanln(&ip)

	if err := scan_model.FofaQuery(ip, jsonDir); err != nil {
		fmt.Printf("单个 IP 查询发生错误: %v\n", err)
		return
	}

	if err := convert_model.ConvertJSONToCSV(jsonDir, outputFile); err != nil {
		log.Fatalf("转换失败: %v", err)
	}

	fmt.Println("成功将 JSON 数据转换为 reports.csv 文件.")
}

// 批量 IP 查询的处理
func handleBatchIPQuery() {
	var filePath string
	fmt.Print("请输入 IP 地址文件的路径: ")
	fmt.Scanln(&filePath)

	ipRecords, err := parse_model.ReadIPsFromFile(filePath)
	if err != nil {
		fmt.Printf("读取文件时出错: %v\n", err)
		return
	}

	processIPs(ipRecords)
}

func processIPs(ipRecords []parse_model.IPRecord) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 5) // 控制并发数，例如5个并发请求

	for i := range ipRecords {
		wg.Add(1)
		sem <- struct{}{} // 获取一个信号量

		go func(record *parse_model.IPRecord) {
			defer wg.Done()
			defer func() { <-sem }() // 释放信号量

			time.Sleep(time.Duration(1+rand.Intn(3)) * time.Second) // 随机延时1到3秒
			if err := scan_model.FofaQuery(record.IP, jsonDir); err != nil {
				fmt.Printf("批量查询 %s 时发生错误: %v\n", record.IP, err)
				return
			}
		}(&ipRecords[i])
	}

	wg.Wait() // 等待所有处理完成

	fmt.Println("所有 IP 查询完成...现在开始转换格式")
	if err := convert_model.ConvertJSONToCSV(jsonDir, outputFile); err != nil {
		log.Fatalf("转换失败: %v", err)
	}
	fmt.Println("成功将 JSON 数据转换为 reports.csv 文件.")
}
