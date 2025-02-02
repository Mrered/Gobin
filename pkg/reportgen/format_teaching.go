package reportgen

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// SectionFormatter 定义了各个部分内容格式化的接口
type SectionFormatter interface {
	Format(content string) string
}

// TeachingFormatter 实现了教学部分的格式化
type TeachingFormatter struct{}

// Format 格式化教学部分的内容
func (f *TeachingFormatter) Format(content string) string {
	// 如果内容为空，直接返回
	if strings.TrimSpace(content) == "" {
		return ""
	}

	// 按行分割内容
	lines := strings.Split(content, "\n")
	var result []string

	// 检查是否所有内容都是"无"
	allNone := true
	for _, line := range lines {
		if strings.TrimSpace(line) != "无" && strings.TrimSpace(line) != "" {
			allNone = false
			break
		}
	}

	// 如果全是"无"，只保留一个"无"
	if allNone {
		return "无"
	}

	// 用于存储相同标题内容的映射
	titleMap := make(map[string]*titleContent)

	var currentTitle string
	var currentContent []string
	var nonCourseContent []string

	// 处理每一行
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		// 跳过空行和单独的"无"
		if line == "" || line == "无" {
			continue
		}

		// 处理三级标题
		if strings.HasPrefix(line, "### ") {
			// 如果是非课程相关的标题（以 ### # 开头），直接添加到非课程内容中
			if strings.HasPrefix(line, "### #") {
				// 如果之前有课程相关的内容，先处理它
				if currentTitle != "" && len(currentContent) > 0 {
					addToTitleMap(titleMap, currentTitle, currentContent)
					currentTitle = ""
					currentContent = nil
				}
				// 添加非课程标题
				nonCourseContent = append(nonCourseContent, line)
				continue
			}

			// 处理课程相关的标题（以 ### [[ 开头）
			if strings.HasPrefix(line, "### [[") {
				// 如果已有累积的内容，先处理之前的内容
				if currentTitle != "" && len(currentContent) > 0 {
					addToTitleMap(titleMap, currentTitle, currentContent)
				}

				// 提取 [[...]] 中的内容作为标题
				titlePattern := regexp.MustCompile(`\[\[([^\]]+)\]\]`)
				match := titlePattern.FindStringSubmatch(line)
				if len(match) > 1 {
					currentTitle = match[1]
					// 提取标签
					tagStart := strings.Index(line, "]]")
					if tagStart != -1 {
						tagStr := strings.TrimSpace(line[tagStart+2:])
						if tagStr != "" {
							// 直接将标题行后面的内容作为标签
							tagPattern := regexp.MustCompile(`#[^\s]+`)
							tags := tagPattern.FindAllString(tagStr, -1)
							for _, tag := range tags {
								if _, exists := titleMap[currentTitle]; !exists {
									titleMap[currentTitle] = &titleContent{
										tags:    make(map[string]bool),
										content: []string{},
									}
								}
								titleMap[currentTitle].tags[tag] = true
							}
						}
					}
				}
				currentContent = nil
				continue
			}

			// 其他三级标题，直接添加到非课程内容中
			if currentTitle != "" && len(currentContent) > 0 {
				addToTitleMap(titleMap, currentTitle, currentContent)
				currentTitle = ""
				currentContent = nil
			}
			nonCourseContent = append(nonCourseContent, line)
			continue
		}

		// 处理内容
		if currentTitle != "" {
			// 处理课程相关内容
			currentContent = append(currentContent, line)
		} else if len(nonCourseContent) > 0 {
			// 处理非课程相关内容
			nonCourseContent = append(nonCourseContent, line)
		}
	}

	// 处理最后一个标题的内容
	if currentTitle != "" && len(currentContent) > 0 {
		addToTitleMap(titleMap, currentTitle, currentContent)
	}

	// 将处理后的课程内容转换为结果
	for title, content := range titleMap {
		result = append(result, formatTitleContent(title, content))
	}

	// 对课程内容结果进行排序
	sort.Strings(result)

	// 添加非课程内容
	if len(nonCourseContent) > 0 {
		result = append(result, strings.Join(nonCourseContent, "\n"))
	}

	return strings.Join(result, "\n\n")
}

// titleContent 用于存储标题相关的内容
type titleContent struct {
	tags    map[string]bool
	content []string
}

// addToTitleMap 将内容添加到标题映射中
func addToTitleMap(titleMap map[string]*titleContent, title string, content []string) {
	if _, exists := titleMap[title]; !exists {
		titleMap[title] = &titleContent{
			tags:    make(map[string]bool),
			content: []string{},
		}
	}

	// 处理每一行内容
	for _, line := range content {
		// 直接添加内容，不处理标签
		if strings.TrimSpace(line) != "" {
			titleMap[title].content = append(titleMap[title].content, line)
		}
	}
}

// formatTitleContent 格式化标题和内容
func formatTitleContent(title string, content *titleContent) string {
	// 构建标题行
	var result strings.Builder
	result.WriteString(fmt.Sprintf("### [[%s]]", title))

	// 添加标签
	var tags []string
	for tag := range content.tags {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	if len(tags) > 0 {
		result.WriteString(" " + strings.Join(tags, " "))
	}

	// 添加内容
	if len(content.content) > 0 {
		// 去重内容
		contentMap := make(map[string]bool)
		var uniqueContent []string
		for _, line := range content.content {
			if !contentMap[line] {
				contentMap[line] = true
				uniqueContent = append(uniqueContent, line)
			}
		}
		result.WriteString("\n" + strings.Join(uniqueContent, "\n"))
	}

	return result.String()
}