package parse_model

import (
	"io/ioutil"
	"strings"
)

type IPRecord struct {
	IP string
}

// ReadIPsFromFile 从文件中读取 IP 地址
func ReadIPsFromFile(filepath string) ([]IPRecord, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var ipRecords []IPRecord
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	for _, line := range lines {
		ip := strings.TrimSpace(line)
		if ip != "" {
			ipRecords = append(ipRecords, IPRecord{IP: ip})
		}
	}
	return ipRecords, nil
}
