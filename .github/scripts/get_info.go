package scripts

import (
	"bufio"
	"bytes"
	"os"
	"strings"
)

// GetHelpTextFromMainGo 读取 main.go 文件顶部注释内容
func GetHelpTextFromMainGo(filePath string) (string, string, []string, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", nil, "", err
	}
	defer file.Close()

	var helpText bytes.Buffer
	var projectDescription string
	var osInfo []string
	scanner := bufio.NewScanner(file)
	inBlockComment := false
	lineNumber := 0

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "/*") {
			inBlockComment = true
			continue
		}
		if strings.HasSuffix(trimmedLine, "*/") {
			inBlockComment = false
			break
		}
		if inBlockComment {
			if lineNumber == 0 {
				// 跳过第一行
				lineNumber++
				continue
			}
			if lineNumber == 1 {
				osInfo = strings.Fields(trimmedLine)
			} else if lineNumber == 2 {
				projectDescription = strings.TrimSpace(line)
			} else {
				helpText.WriteString(line + "\n")
			}
			lineNumber++
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", nil, "", err
	}

	return "", projectDescription, osInfo, helpText.String(), nil
}
