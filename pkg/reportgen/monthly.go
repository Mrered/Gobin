package reportgen

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// extractFrontMatterStats 从周报的文档属性中提取查宿和特种工监考次数
func (g *MonthlyGenerator) extractFrontMatterStats(content string) (int, int) {
	var dormCount, examCount int

	// 匹配YAML格式的文档属性
	re := regexp.MustCompile(`(?m)^查宿次数: "(\d+)"$|^特种工监考: "(\d+)"$`)
	matches := re.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if match[1] != "" { // 查宿次数
			if count, err := strconv.Atoi(match[1]); err == nil {
				dormCount = count
			}
		} else if match[2] != "" { // 特种工监考
			if count, err := strconv.Atoi(match[2]); err == nil {
				examCount = count
			}
		}
	}

	return dormCount, examCount
}

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

	// 生成周报链接列表
	var weeklyLinks strings.Builder
	for _, report := range selectedReports {
		fileName := strings.TrimSuffix(filepath.Base(report.FilePath), ".md")
		weeklyLinks.WriteString(fmt.Sprintf("[[%s]]\n\n", fileName))
	}

	// 合并报告内容并格式化
	content := g.mergeSectionsAndFormat(selectedReports)

	// 统计各项数据
	listeningCount := countListeningClasses(content)
	totalDormCount := 0
	totalExamCount := 0

	// 从每个周报的文档属性中统计查宿和特种工监考次数
	for _, report := range selectedReports {
		dormCount, examCount := g.extractFrontMatterStats(report.Content)
		totalDormCount += dormCount
		totalExamCount += examCount
	}

	// 生成文档属性
	var frontMatter strings.Builder
	frontMatter.WriteString("---\n")
	frontMatter.WriteString(fmt.Sprintf("听课次数: \"%d\"\n", listeningCount))
	frontMatter.WriteString(fmt.Sprintf("查宿次数: \"%d\"\n", totalDormCount))
	frontMatter.WriteString(fmt.Sprintf("特种工监考: \"%d\"\n", totalExamCount))
	frontMatter.WriteString("---\n\n")

	// 在内容前添加周报链接和文档属性
	content = frontMatter.String() + weeklyLinks.String() + content

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
