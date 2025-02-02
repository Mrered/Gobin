package reportgen

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	// "strings"
)

// Generate 生成年报
func (g *YearlyGenerator) Generate(sourcePath string, params map[string]string) error {
	// 读取学期报文件
	reports, err := g.readFiles(sourcePath)
	if err != nil {
		return fmt.Errorf("读取学期报文件失败：%v", err)
	}

	// 筛选指定年份的报告
	var selectedReports []Report
	for _, report := range reports {
		// 从文件名中提取年份
		base := filepath.Base(report.FilePath)
		yearRegex := regexp.MustCompile(`^(\d{4})`)
		match := yearRegex.FindString(base)
		if match == "" {
			continue
		}

		if match == g.Config.SelectedPeriod {
			selectedReports = append(selectedReports, report)
		}
	}

	if len(selectedReports) == 0 {
		return fmt.Errorf("未找到 %s 年的学期报", g.Config.SelectedPeriod)
	}

	// 合并报告内容并格式化
	content := g.mergeSectionsAndFormat(selectedReports)

	// 生成输出文件名
	outputFile := filepath.Join(g.Config.TargetDir, fmt.Sprintf("%s - %s 学年.md", g.Config.SelectedPeriod, string(g.Config.SelectedPeriod[0:4])))

	// 写入文件
	return os.WriteFile(outputFile, []byte(content), 0644)
}

// GetAvailablePeriods 获取可用的年份
func (g *YearlyGenerator) GetAvailablePeriods(sourcePath string) ([]string, error) {
	reports, err := g.readFiles(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("读取学期报文件失败：%v", err)
	}

	yearMap := make(map[string]bool)
	for _, report := range reports {
		// 从文件名中提取年份
		base := filepath.Base(report.FilePath)
		yearRegex := regexp.MustCompile(`^(\d{4})`)
		match := yearRegex.FindString(base)
		if match != "" {
			yearMap[match] = true
		}
	}

	var years []string
	for year := range yearMap {
		years = append(years, year)
	}

	// 对年份进行排序
	sort.Strings(years)
	return years, nil
}
