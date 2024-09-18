package main

import (
	"Fofa_scan/convert_model"
	"Fofa_scan/parse_model"
	"Fofa_scan/scan_model"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var jsonDir = "json_report"       // json文件存储的路径
var outputFile = "reports.xlsx"   // 转换后的xlsx文件
var progressFile = "progress.txt" // 记录已处理IP的文件

// 加载已处理的IP地址
func loadProgress() ([]string, error) {
	if _, err := os.Stat(progressFile); os.IsNotExist(err) {
		return []string{}, nil
	}
	data, err := ioutil.ReadFile(progressFile)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(data)), "\n"), nil
}

// 保存已处理的IP地址到文件
func saveProgress(ips []string) error {
	data := strings.Join(ips, "\n")
	return ioutil.WriteFile(progressFile, []byte(data), 0644)
}

// 主函数
func main() {
	if err := os.MkdirAll(jsonDir, 0755); err != nil {
		log.Fatalf("创建 %s 失败: %v", jsonDir, err)
	}

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
	fmt.Println("1. 单个IP查询")
	fmt.Println("2. 批量IP查询(从.txt文件读取)")
	fmt.Print("请输入你的选择(1或2):")
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		var ip string
		fmt.Print("请输入IP地址:")
		fmt.Scanln(&ip)

		if err := scan_model.FofaQuery(ip, jsonDir); err != nil {
			fmt.Printf("单个IP查询发生错误:%v\n", err)
		}
		time.Sleep(2 * time.Second)
		if err := convert_model.ConvertJSONToXLSX(jsonDir, outputFile); err != nil {
			log.Fatalf("转换失败:%v", err)
		}

		fmt.Println("成功将JSON数据转换为reports.xlsx文件.")

	case 2:
		var filePath string
		fmt.Print("请输入IP地址文件的路径:")
		fmt.Scanln(&filePath)

		var ips []string
		var err error
		var continueProgress bool

		// 检查进度文件
		processedIPs, err := loadProgress()
		if err != nil {
			fmt.Printf("加载进度时出错:%v\n", err)
			return
		}

		if len(processedIPs) > 0 {
			fmt.Printf("发现已处理的IP，共 %d 个，是否要继续上次进度? (y/n):", len(processedIPs))
			var answer string
			fmt.Scanln(&answer)
			if strings.ToLower(answer) == "y" {
				// 继续上次进度
				continueProgress = true
			}
		}

		// 从文件中读取IP地址
		ips, err = parse_model.ReadIPsFromFile(filePath)
		if err != nil {
			fmt.Printf("读取文件时出错:%v\n", err)
			return
		}

		// 只保留未处理的IP
		if continueProgress {
			var unprocessedIPs []string
			for _, ip := range ips {
				if !stringInSlice(ip, processedIPs) {
					unprocessedIPs = append(unprocessedIPs, ip) // 添加未处理的IP
				}
			}
			ips = unprocessedIPs // 更新待处理的IP列表
		}

		// 对每个IP执行查询
		for _, ip := range ips {
			time.Sleep(1 * time.Second) // 一秒查询一次
			if err := scan_model.FofaQuery(ip, jsonDir); err != nil {
				fmt.Printf("批量查询%s时发生错误:%v\n", ip, err)
			} else {
				processedIPs = append(processedIPs, ip) // 添加成功处理的IP
				if err := saveProgress(processedIPs); err != nil {
					fmt.Printf("保存进度时出错:%v\n", err)
				}
			}
		}

		// 判断是否所有IP已处理
		if len(processedIPs) == len(ips) {
			fmt.Println("所有IP查询完成...现在开始转换格式")
			time.Sleep(5 * time.Second)
			if err := convert_model.ConvertJSONToXLSX(jsonDir, outputFile); err != nil {
				log.Fatalf("转换失败:%v", err)
			}
			fmt.Println("成功将JSON数据转换为reports.xlsx文件.")
		} else {
			fmt.Println("仍有未处理的IP，请检查进度文件。")
		}

	default:
		fmt.Println("无效输入,请输入1或2.")
	}
}

// 判断字符串是否在切片中
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if a == b {
			return true
		}
	}
	return false
}
