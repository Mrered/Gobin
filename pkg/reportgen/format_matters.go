package reportgen

import (
	"bufio"
	"fmt"
	"strings"
)

// MiscellaneousFormatter 杂事部分的格式化器
type MiscellaneousFormatter struct{}

// Format 实现了 SectionFormatter 接口
func (f *MiscellaneousFormatter) Format(content string) string {
	var nonNumberLines []string
	var importantNumberLines []string

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		// 检查是否以数字开头
		if len(trimmedLine) > 0 && (trimmedLine[0] >= '0' && trimmedLine[0] <= '9') {
			// 检查是否包含关键词
			if strings.Contains(trimmedLine, "查宿") || strings.Contains(trimmedLine, "监考") {
				importantNumberLines = append(importantNumberLines, trimmedLine)
			}
			// 跳过其他数字开头的行
			continue
		}
		// 保留非数字开头的行
		nonNumberLines = append(nonNumberLines, trimmedLine)
	}

	// 组合最终内容
	var result strings.Builder

	// 检查是否只有一个"无"
	if len(nonNumberLines) == 1 && strings.TrimSpace(nonNumberLines[0]) == "无" {
		return "无"
	}

	// 为非数字开头的行添加序号
	for i, line := range nonNumberLines {
		if i > 0 {
			result.WriteString("\n")
		}
		// 跳过空行
		if strings.TrimSpace(line) == "" {
			result.WriteString(line)
			continue
		}
		// 添加序号
		result.WriteString(fmt.Sprintf("%d. %s", i+1, line))
	}

	// 如果有重要的数字开头行，添加到底部
	if len(importantNumberLines) > 0 {
		// 确保在添加重要行之前有一个空行
		if result.Len() > 0 {
			result.WriteString("\n\n")
		}
		// 继续使用前面的序号
		startNum := len(nonNumberLines) + 1
		for i, line := range importantNumberLines {
			if i > 0 {
				result.WriteString("\n")
			}
			result.WriteString(fmt.Sprintf("%d. %s", startNum+i, line))
		}
	}

	return result.String()
}
