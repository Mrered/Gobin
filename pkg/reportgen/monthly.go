package reportgen

import (
	"fmt"
	"os"
	"path/filepath"
)

// Generate 生成月报
func (g *MonthlyGenerator) Generate(sourcePath string, params map[string]string) error {
	// 读取周报文件
	reports, err := g.readFiles(sourcePath)
	if err != nil {
		return fmt.Errorf("读取周报文件失败：%v", err)
	}

	// 筛选指定月份的报告
	var selectedReports []Report
	for _, report := range reports {
		date, err := ExtractDateFromFilename(report.FilePath)
		if err != nil {
			continue
		}

		if date.Format("200601") == g.Config.SelectedPeriod {
			selectedReports = append(selectedReports, report)
		}
	}

	if len(selectedReports) == 0 {
		return fmt.Errorf("未找到 %s 月的周报", g.Config.SelectedPeriod)
	}

	// 合并报告内容并格式化
	content := g.mergeSectionsAndFormat(selectedReports)

	// 生成输出文件名
	outputFile := filepath.Join(g.Config.TargetDir, fmt.Sprintf("%s.md", g.Config.SelectedPeriod))

	// 写入文件
	return os.WriteFile(outputFile, []byte(content), 0644)
}

// GetAvailablePeriods 获取可用的月份
func (g *MonthlyGenerator) GetAvailablePeriods(sourcePath string) ([]string, error) {
	reports, err := g.readFiles(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("读取周报文件失败：%v", err)
	}

	monthMap := make(map[string]bool)
	for _, report := range reports {
		date, err := ExtractDateFromFilename(report.FilePath)
		if err != nil {
			continue
		}

		monthMap[date.Format("200601")] = true
	}

	var months []string
	for month := range monthMap {
		months = append(months, month)
	}
	return months, nil
}
