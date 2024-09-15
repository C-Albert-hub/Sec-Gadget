package convert_model

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tealeg/xlsx"
)

// ConvertJSONToXLSX 将指定目录中的 JSON 文件中的 results 数组转换为 XLSX 文件
func ConvertJSONToXLSX(jsonDir string, outputFile string) error {
	xlFile := xlsx.NewFile()
	sheet, err := xlFile.AddSheet("Reports")
	if err != nil {
		return fmt.Errorf("创建工作表时出错: %w", err)
	}

	// 添加表头,后续根据fofa的field字段对应添加
	headerRow := sheet.AddRow()
	headerRow.AddCell().Value = "IP"
	headerRow.AddCell().Value = "域名"
	headerRow.AddCell().Value = "端口"
	headerRow.AddCell().Value = "标题"
	headerRow.AddCell().Value = "最近更新时间"

	err = filepath.Walk(jsonDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".json" {
			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("读取文件 %s 时出错: %w", path, err)
			}

			var jsonObject map[string]interface{}

			if err := json.Unmarshal(data, &jsonObject); err != nil {
				return fmt.Errorf("解析 JSON 文件 %s 时出错: %w", path, err)
			}

			results, ok := jsonObject["results"].([]interface{})
			if !ok {
				return fmt.Errorf("文件 %s 中找不到 results 字段", path)
			}

			// 将 results 数据写入 XLSX 文件
			for _, item := range results {
				resultItem, ok := item.([]interface{})
				row := sheet.AddRow()

				if !ok {
					// 如果是非数组类型，继续处理其他项
					continue
				}

				// 检查 IP 和 域名 的索引，并写入单元格
				if len(resultItem) >= 2 {
					row.AddCell().Value = fmt.Sprintf("%v", resultItem[0]) // IP
					row.AddCell().Value = fmt.Sprintf("%v", resultItem[1]) // 域名
					row.AddCell().Value = fmt.Sprintf("%v", resultItem[2]) // 端口
					row.AddCell().Value = fmt.Sprintf("%v", resultItem[3]) // 标题
					row.AddCell().Value = fmt.Sprintf("%v", resultItem[4]) //最近更新时间
				}
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("遍历目录时出错: %w", err)
	}

	if err := xlFile.Save(outputFile); err != nil {
		return fmt.Errorf("保存 XLSX 文件时出错: %w", err)
	}

	return nil
}
