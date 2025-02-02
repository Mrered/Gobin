package reportgen

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// BaseGenerator 提供基本的报告生成功能
type BaseGenerator struct {
	Config *Config
}

// WeeklyGenerator 周报生成器
type WeeklyGenerator struct {
	BaseGenerator
}

// MonthlyGenerator 月报生成器
type MonthlyGenerator struct {
	BaseGenerator
}

// SemesterGenerator 学期报生成器
type SemesterGenerator struct {
	BaseGenerator
}

// YearlyGenerator 年报生成器
type YearlyGenerator struct {
	BaseGenerator
}

// NewGenerator 创建对应类型的报告生成器
func NewGenerator(config *Config) (ReportGenerator, error) {
	switch config.ReportType {
	case "w":
		return &WeeklyGenerator{BaseGenerator{config}}, nil
	case "m":
		return &MonthlyGenerator{BaseGenerator{config}}, nil
	case "s":
		return &SemesterGenerator{BaseGenerator{config}}, nil
	case "y":
		return &YearlyGenerator{BaseGenerator{config}}, nil
	default:
		return nil, fmt.Errorf("不支持的报告类型：%s", config.ReportType)
	}
}

// readFiles 读取指定目录下的所有 Markdown 文件
func (g *BaseGenerator) readFiles(sourcePath string) ([]Report, error) {
	var reports []Report
	err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".md") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			report := Report{
				FilePath: path,
				Content:  string(content),
			}
			reports = append(reports, report)
		}
		return nil
	})
	return reports, err
}

// extractSections 从报告内容中提取各个部分
func (g *BaseGenerator) extractSections(content string) map[string][]string {
	sections := make(map[string][]string)
	currentSection := ""
	var currentContent []string

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "## ") {
			if currentSection != "" && len(currentContent) > 0 {
				sections[currentSection] = append(sections[currentSection], strings.Join(currentContent, "\n"))
			}
			currentSection = strings.TrimSpace(strings.TrimPrefix(line, "## "))
			currentContent = nil
		} else if currentSection != "" {
			currentContent = append(currentContent, line)
		}
	}

	// 处理最后一个部分
	if currentSection != "" && len(currentContent) > 0 {
		sections[currentSection] = append(sections[currentSection], strings.Join(currentContent, "\n"))
	}

	return sections
}

// mergeSections 合并多个报告的相同部分
func (g *BaseGenerator) mergeSections(reports []Report) map[string][]string {
	merged := make(map[string][]string)
	for _, report := range reports {
		sections := g.extractSections(report.Content)
		for section, content := range sections {
			merged[section] = append(merged[section], content...)
		}
	}
	return merged
}

// formatReport 格式化最终报告
func (g *BaseGenerator) formatReport(sections map[string][]string) string {
	var result strings.Builder
	sectionOrder := []string{TeachingSection, ListeningSection, TrainingSection, MiscellaneousSection}

	for _, section := range sectionOrder {
		if content, ok := sections[section]; ok && len(content) > 0 {
			result.WriteString(fmt.Sprintf("## %s\n\n", section))
			formattedContent := FormatMarkdown(strings.Join(content, "\n"))
			result.WriteString(formattedContent)
			result.WriteString("\n\n")
		}
	}

	return result.String()
}
