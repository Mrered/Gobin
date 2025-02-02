package reportgen

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// ValidateWorkingDir 验证工作目录是否包含所需的子目录
func ValidateWorkingDir(dirPath string) error {
	requiredDirs := []string{"日报", "周报", "月报", "学期报", "年报"}
	for _, dir := range requiredDirs {
		if _, err := os.Stat(filepath.Join(dirPath, dir)); err != nil {
			return fmt.Errorf("当前目录不完整，无法归纳总结：%s 目录不存在", dir)
		}
	}
	return nil
}

// ExtractDateFromFilename 从文件名中提取日期
func ExtractDateFromFilename(filename string) (time.Time, error) {
	// 移除文件扩展名
	base := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))

	// 尝试匹配日期格式（YYYYMMDD）
	dateRegex := regexp.MustCompile(`^(\d{8})`)
	match := dateRegex.FindString(base)
	if match == "" {
		return time.Time{}, fmt.Errorf("无法从文件名 %s 中提取日期", filename)
	}

	return time.Parse("20060102", match)
}

// ExtractMonthFromFilename 从月报文件名中提取日期
func ExtractMonthFromFilename(filename string) (time.Time, error) {
	// 移除文件扩展名
	base := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))

	// 尝试匹配月份格式（YYYYMM）
	monthRegex := regexp.MustCompile(`^(\d{6})$`)
	match := monthRegex.FindString(base)
	if match == "" {
		return time.Time{}, fmt.Errorf("无法从月报文件名 %s 中提取日期", filename)
	}

	// 解析日期，将月份设置为该月的第一天
	return time.Parse("200601", match)
}

// GetSemesterPeriod 根据日期获取学期信息
func GetSemesterPeriod(date time.Time) string {
	year := date.Year()
	month := date.Month()

	if month >= 2 && month <= 7 {
		return fmt.Sprintf("%d - %d 春", year-1, year)
	} else if month == 1 {
		return fmt.Sprintf("%d - %d 秋", year-1, year)
	}
	return fmt.Sprintf("%d - %d 秋", year, year+1)
}

// FormatMarkdown 格式化 Markdown 内容
func FormatMarkdown(content string) string {
	// 移除空行
	lines := strings.Split(content, "\n")
	var result []string

	for _, line := range lines {
		// 跳过以 > 或 [[ 开头的行
		if strings.HasPrefix(strings.TrimSpace(line), ">") ||
			strings.HasPrefix(strings.TrimSpace(line), "[[") {
			continue
		}

		// 保留非空行
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}
