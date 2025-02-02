package reportgen

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Generate 生成周报
func (g *WeeklyGenerator) Generate(sourcePath string, params map[string]string) error {
	// 读取日报文件
	reports, err := g.readFiles(sourcePath)
	if err != nil {
		return fmt.Errorf("读取日报文件失败：%v", err)
	}

	// 筛选指定周数的报告
	var selectedReports []Report
	for _, report := range reports {
		if week := extractWeekFromContent(report.Content); week == g.Config.SelectedPeriod {
			selectedReports = append(selectedReports, report)
		}
	}

	if len(selectedReports) == 0 {
		return fmt.Errorf("未找到第 %s 周的日报", g.Config.SelectedPeriod)
	}

	// 合并报告内容并格式化
	content := g.mergeSectionsAndFormat(selectedReports)

	// 生成输出文件名
	firstFile := filepath.Base(selectedReports[0].FilePath)
	lastFile := filepath.Base(selectedReports[len(selectedReports)-1].FilePath)
	firstDate := strings.TrimSuffix(firstFile, ".md")
	lastDate := strings.TrimSuffix(lastFile, ".md")
	outputFile := filepath.Join(g.Config.TargetDir, fmt.Sprintf("%s - %s.md", firstDate, lastDate))

	// 写入文件
	return os.WriteFile(outputFile, []byte(content), 0644)
}

// GetAvailablePeriods 获取可用的周数
func (g *WeeklyGenerator) GetAvailablePeriods(sourcePath string) ([]string, error) {
	reports, err := g.readFiles(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("读取日报文件失败：%v", err)
	}

	weekMap := make(map[string]bool)
	for _, report := range reports {
		if week := extractWeekFromContent(report.Content); week != "" {
			weekMap[week] = true
		}
	}

	var weeks []string
	for week := range weekMap {
		weeks = append(weeks, week)
	}
	return weeks, nil
}

// extractWeekFromContent 从文件内容中提取周数
func extractWeekFromContent(content string) string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	inFrontMatter := false

	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" {
			if !inFrontMatter {
				inFrontMatter = true
				continue
			}
			break
		}

		if inFrontMatter && strings.HasPrefix(line, "周:") {
			return strings.Trim(strings.TrimPrefix(line, "周:"), "\" ")
		}
	}

	return ""
}
