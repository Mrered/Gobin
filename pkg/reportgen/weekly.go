package reportgen

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// countListeningClasses 统计听课次数
func countListeningClasses(content string) int {
	count := 0
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#### ") && strings.Contains(content, "## 听课") {
			count++
		}
	}
	return count
}

// countKeywordInSection 统计指定部分中关键词出现的次数
func countKeywordInSection(content, section, keyword string) int {
	count := 0
	scanner := bufio.NewScanner(strings.NewReader(content))
	inSection := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "## ") {
			inSection = strings.TrimPrefix(line, "## ") == section
		} else if inSection && strings.Contains(line, keyword) {
			count++
		}
	}
	return count
}

// removeKeywordLines 删除包含指定关键词的行
func removeKeywordLines(content, section, keyword string) string {
	var result strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(content))
	inSection := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "## ") {
			inSection = strings.TrimPrefix(line, "## ") == section
		}

		if !inSection || !strings.Contains(line, keyword) {
			result.WriteString(line + "\n")
		}
	}
	return result.String()
}

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

	// 生成日报链接列表
	var dailyLinks strings.Builder
	for _, report := range selectedReports {
		fileName := strings.TrimSuffix(filepath.Base(report.FilePath), ".md")
		dailyLinks.WriteString(fmt.Sprintf("[[%s]]\n\n", fileName))
	}

	// 合并报告内容并格式化
	content := g.mergeSectionsAndFormat(selectedReports)

	// 统计各项数据
	listeningCount := countListeningClasses(content)
	dormCount := countKeywordInSection(content, "杂事", "查宿")
	examCount := countKeywordInSection(content, "杂事", "特种工监考")

	// 删除统计过的行
	content = removeKeywordLines(content, "杂事", "查宿")
	content = removeKeywordLines(content, "杂事", "特种工监考")

	// 生成文档属性
	var frontMatter strings.Builder
	frontMatter.WriteString("---\n")
	frontMatter.WriteString(fmt.Sprintf("周: \"%s\"\n", g.Config.SelectedPeriod))
	frontMatter.WriteString(fmt.Sprintf("听课次数: \"%d\"\n", listeningCount))
	frontMatter.WriteString(fmt.Sprintf("查宿次数: \"%d\"\n", dormCount))
	frontMatter.WriteString(fmt.Sprintf("特种工监考: \"%d\"\n", examCount))
	frontMatter.WriteString("---\n\n")

	// 在内容前添加日报链接和文档属性
	content = frontMatter.String() + dailyLinks.String() + content

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
