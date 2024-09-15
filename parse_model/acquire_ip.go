package parse_model

import (
	"bufio"
	"os"
)

// ReadIPsFromFile 从指定文件读取 IP 地址列表
func ReadIPsFromFile(path string) ([]string, error) {
	var ips []string
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ip := scanner.Text()
		if ip != "" {
			ips = append(ips, ip)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ips, nil
}
