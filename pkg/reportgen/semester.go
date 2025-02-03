package reportgen

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// extractFrontMatterStats 从月报的文档属性中提取查宿和特种工监考次数
func extractFrontMatterStats(content string) (int, int) {
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

	// 生成月报链接列表
	var monthlyLinks strings.Builder
	for _, report := range selectedReports {
		fileName := strings.TrimSuffix(filepath.Base(report.FilePath), ".md")
		monthlyLinks.WriteString(fmt.Sprintf("[[%s]]\n\n", fileName))
	}

	// 合并报告内容并格式化
	content := g.mergeSectionsAndFormat(selectedReports)

	// 统计各项数据
	listeningCount := countListeningClasses(content)
	totalDormCount := 0
	totalExamCount := 0

	// 从每个月报的文档属性中统计查宿和特种工监考次数
	for _, report := range selectedReports {
		dormCount, examCount := extractFrontMatterStats(report.Content)
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

	// 在内容前添加月报链接和文档属性
	content = frontMatter.String() + monthlyLinks.String() + content

	// 生成输出文件名
	outputFile := filepath.Join(g.Config.TargetDir, fmt.Sprintf("%s流水账.md", g.Config.SelectedPeriod))

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
