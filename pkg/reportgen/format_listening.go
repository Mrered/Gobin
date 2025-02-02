package reportgen

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// ListeningFormatter 实现了听课部分的格式化
type ListeningFormatter struct{}

// Format 格式化听课部分的内容
func (f *ListeningFormatter) Format(content string) string {
	if strings.TrimSpace(content) == "" {
		return "无" // 如果输入为空，直接返回“无”
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
			subtitle := strings.TrimSpace(match[2]) // 提取标题后的内容
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

		if len(block.Content) == 1 && strings.TrimSpace(block.Content[0]) == "" {
			result = append(result, "#### 无") // 添加“无”
		} else {
			for i, contentLine := range block.Content {
				if i == 0 {
					if strings.HasPrefix(contentLine, "#### ") {
						result = append(result, contentLine)
					} else {
						result = append(result, fmt.Sprintf("#### %s", contentLine))
					}
				} else {
					result = append(result, contentLine)
				}

			}
		}
	}

	output := strings.Join(result, "\n")
	if output == "" {
		return "无" // 如果最终输出为空，则返回“无”
	}

	return output
}

// Block 存储一个三级标题及其对应的内容
type Block struct {
	Title   string
	Content []string
}
