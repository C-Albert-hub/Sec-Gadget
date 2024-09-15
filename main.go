// main.go
package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"tracer_scanner/convert_model"
	"tracer_scanner/parse_model"
	"tracer_scanner/scan_model"
)

var jsonDir = "json_report"     //json文件存储的路径
var outputFile = "reports.xlsx" //转换后的xlsx文件

func main() {
	// 检查 jsonDir 是否存在，如果不存在则创建
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
	fmt.Println("1.单个IP查询")
	fmt.Println("2.批量IP查询(从.txt文件读取)")
	fmt.Print("请输入你的选择(1或2):")
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		var ip string
		fmt.Print("请输入IP地址:")
		fmt.Scanln(&ip)

		//调用单个IP查询并保存结果
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

		//从指定文件读取IP地址
		ips, err := parse_model.ReadIPsFromFile(filePath)
		if err != nil {
			fmt.Printf("读取文件时出错:%v\n", err)
			return
		}

		//对每个IP执行查询
		for _, ip := range ips {
			time.Sleep(1 * time.Second) //一秒种查询一次
			if err := scan_model.FofaQuery(ip, jsonDir); err != nil {
				fmt.Printf("批量查询%s时发生错误:%v\n", ip, err)
			}
		}
		fmt.Println("批量IP查询完成...现在开始转换格式")
		time.Sleep(5 * time.Second)
		if err := convert_model.ConvertJSONToXLSX(jsonDir, outputFile); err != nil {
			log.Fatalf("转换失败:%v", err)
		}

		fmt.Println("成功将JSON数据转换为reports.xlsx文件.")

	default:
		fmt.Println("无效输入,请输入1或2.")
	}
}
