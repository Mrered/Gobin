package reportgen

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Generate 生成学期报
func (g *SemesterGenerator) Generate(sourcePath string, params map[string]string) error {
	// 读取月报文件
	reports, err := g.readFiles(sourcePath)
	if err != nil {
		return fmt.Errorf("读取月报文件失败：%v", err)
	}

	// 筛选指定学期的报告
	var selectedReports []Report
	for _, report := range reports {
		date, err := ExtractMonthFromFilename(report.FilePath)
		if err != nil {
			continue
		}

		if GetSemesterPeriod(date) == g.Config.SelectedPeriod {
			selectedReports = append(selectedReports, report)
		}
	}

	if len(selectedReports) == 0 {
		return fmt.Errorf("未找到 %s 学期的月报", g.Config.SelectedPeriod)
	}

	// 合并报告内容并格式化
	content := g.mergeSectionsAndFormat(selectedReports)

	// 生成输出文件名
	outputFile := filepath.Join(g.Config.TargetDir, fmt.Sprintf("%s.md", g.Config.SelectedPeriod))

	// 写入文件
	return os.WriteFile(outputFile, []byte(content), 0644)
}

// GetAvailablePeriods 获取可用的学期
func (g *SemesterGenerator) GetAvailablePeriods(sourcePath string) ([]string, error) {
	reports, err := g.readFiles(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("读取月报文件失败：%v", err)
	}

	semesterMap := make(map[string]bool)
	for _, report := range reports {
		date, err := ExtractMonthFromFilename(report.FilePath)
		if err != nil {
			continue
		}

		semesterMap[GetSemesterPeriod(date)] = true
	}

	var semesters []string
	for semester := range semesterMap {
		semesters = append(semesters, semester)
	}

	// 对学期进行排序
	sort.Strings(semesters)
	return semesters, nil
}
