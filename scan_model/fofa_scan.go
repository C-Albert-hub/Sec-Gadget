package scan_model

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	apiUrl = "https://fofa.info/api/v1/search/all"
	apiKey = ""
)

// 利用fofa API 执行 IP 查询并保存响应为单独的 JSON 文件
func FofaQuery(ip string, jsonDir string) error {
	// 构建查询字符串，确保格式为 ip=x.x.x.x
	queryString := fmt.Sprintf("ip=%s", ip)

	// 将查询字符串进行 base64 编码
	qbase64 := base64.StdEncoding.EncodeToString([]byte(queryString))

	// 构建查询参数
	values := url.Values{}
	values.Set("key", apiKey)
	values.Set("qbase64", qbase64)
	values.Set("size", "100")
	values.Set("fields", "ip,domain,port,title,lastupdatetime")

	// 构建完整的请求 URL
	requestURL := fmt.Sprintf("%s?%s", apiUrl, values.Encode())

	// 创建 HTTP GET 请求
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36 Edg/128.0.0.0")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// 输出响应状态
	fmt.Println("Response Status:", resp.Status)

	// 读取响应内容
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	/*	// 打印原始响应体
		fmt.Println("Raw Response Body:")
		fmt.Println(string(responseBody))*/

	// 格式化 JSON 响应
	var formattedJson interface{}
	if err := json.Unmarshal(responseBody, &formattedJson); err != nil {
		return fmt.Errorf("error unmarshaling response: %w", err)
	}
	formattedJsonStr, err := json.MarshalIndent(formattedJson, "", "  ")
	if err != nil {
		return fmt.Errorf("error formatting JSON: %w", err)
	}

	/*	// 打印格式化后的 JSON
		fmt.Println("Formatted Response Body:")
		fmt.Println(string(formattedJsonStr))*/

	// 生成文件名，使用 IP 地址作为文件名的一部分
	sanitizedIP := strings.ReplaceAll(ip, "/", "_") // 将 / 替换为 _
	fileName := fmt.Sprintf("%s.json", sanitizedIP)
	filePath := filepath.Join(jsonDir, fileName) // 使用 json_report 目录

	// 将格式化后的 JSON 数据写入文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件时出错: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(formattedJsonStr); err != nil {
		return fmt.Errorf("写入文件时出错: %w", err)
	}

	fmt.Printf("格式化后的响应已保存到 %s\n", filePath)
	return nil
}
