package convert_model

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ConvertJSONToCSV 将指定目录中的 JSON 文件中的 results 数组转换为 CSV 文件
func ConvertJSONToCSV(jsonDir string, outputFile string) error {
	// 打开 CSV 文件
	csvFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("创建 CSV 文件时出错: %w", err)
	}
	defer csvFile.Close()

	// 添加 UTF-8 BOM
	_, err = csvFile.Write([]byte{0xEF, 0xBB, 0xBF}) // 添加 BOM
	if err != nil {
		return fmt.Errorf("写入 BOM 时出错: %w", err)
	}

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	// 添加表头
	header := []string{"IP", "域名", "端口", "标题", "最近更新时间", "服务器"}
	if err := csvWriter.Write(header); err != nil {
		return fmt.Errorf("写入表头时出错: %w", err)
	}

	err = filepath.Walk(jsonDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理 JSON 文件
		if filepath.Ext(path) == ".json" {
			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("读取文件 %s 时出错: %w", path, err)
			}

			// 解析 JSON 文件
			var jsonObject struct {
				Results [][]interface{} `json:"results"`
			}

			if err := json.Unmarshal(data, &jsonObject); err != nil {
				return fmt.Errorf("解析 JSON 文件 %s 时出错: %w", path, err)
			}

			// 将 results 数据写入 CSV 文件
			for _, resultItem := range jsonObject.Results {
				if len(resultItem) >= 5 {
					// 创建一行数据
					record := []string{
						fmt.Sprintf("%v", resultItem[0]), // IP
						fmt.Sprintf("%v", resultItem[1]), // 域名
						fmt.Sprintf("%v", resultItem[2]), // 端口
						fmt.Sprintf("%v", resultItem[3]), // 标题
						fmt.Sprintf("%v", resultItem[4]), // 最近更新时间
						fmt.Sprintf("%v", resultItem[5]), //服务器
					}
					if err := csvWriter.Write(record); err != nil {
						return fmt.Errorf("写入数据行时出错: %w", err)
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("遍历目录时出错: %w", err)
	}

	return nil
}
