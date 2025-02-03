package reportgen

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// MonthlyTeachingFormatter 实现了月报教学部分的格式化
type MonthlyTeachingFormatter struct{}

// Format 格式化教学部分的内容
func (f *MonthlyTeachingFormatter) Format(content string) string {
	if strings.TrimSpace(content) == "" {
		return "无"
	}

	lines := strings.Split(content, "\n")
	blocks := make([]*Block, 0)
	var currentBlock *Block

	// 正则表达式匹配三级标题行
	re := regexp.MustCompile(`^### \[\[(.*?)\]\](.*?)$`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		match := re.FindStringSubmatch(line)
		if len(match) > 0 {
			// 发现新的三级标题行
			title := match[1]
			subtitle := strings.TrimSpace(match[2])
			currentBlock = &Block{Title: title, Content: []string{subtitle}}
			blocks = append(blocks, currentBlock)
		} else if currentBlock != nil {
			// 当前行是正文内容
			currentBlock.Content = append(currentBlock.Content, line)
		}
	}

	// 按照三级标题排序
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].Title < blocks[j].Title
	})

	var result []string
	var lastTitle string

	for _, block := range blocks {
		if block.Title != lastTitle {
			result = append(result, fmt.Sprintf("### [[%s]]", block.Title))
			lastTitle = block.Title
		}

		for _, contentLine := range block.Content {
			result = append(result, contentLine)
		}
	}

	output := strings.Join(result, "\n")
	if output == "" {
		return "无"
	}

	return output
}

// MonthlyListeningFormatter 实现了月报听课部分的格式化
type MonthlyListeningFormatter struct{}

// Format 格式化听课部分的内容
func (f *MonthlyListeningFormatter) Format(content string) string {
	return (&MonthlyTeachingFormatter{}).Format(content)
}

// MonthlyTrainingFormatter 实现了月报培训学习部分的格式化
type MonthlyTrainingFormatter struct{}

// Format 格式化培训学习部分的内容
func (f *MonthlyTrainingFormatter) Format(content string) string {
	if strings.TrimSpace(content) == "" {
		return "无"
	}

	lines := strings.Split(content, "\n")
	var result []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}

	if len(result) == 0 {
		return "无"
	}

	return strings.Join(result, "\n")
}

// MonthlyMattersFormatter 实现了月报杂事部分的格式化
type MonthlyMattersFormatter struct{}

// Format 格式化杂事部分的内容
func (f *MonthlyMattersFormatter) Format(content string) string {
	if strings.TrimSpace(content) == "" {
		return "无"
	}

	lines := strings.Split(content, "\n")
	var result []string
	var counter int

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && line != "无" {
			// 移除原有的序号（如果有的话）
			line = strings.TrimLeft(line, "0123456789. ")
			counter++
			result = append(result, fmt.Sprintf("%d. %s", counter, line))
		}
	}

	if len(result) == 0 {
		return "无"
	}

	return strings.Join(result, "\n")
}
