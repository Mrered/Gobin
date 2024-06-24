package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	binDir := "cmd"
	readmeFile := "README.md"

	// 项目名称和描述
	projectName := "Gobin"
	projectDescription := "Go 二进制小程序"

	// 收集帮助信息
	helpTexts := make(map[string]string)
	descriptions := make(map[string]string)
	err := filepath.Walk(binDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != binDir {
			binaryName := filepath.Base(path)
			_, projectDescription, helpText, err := getHelpTextFromMainGo(filepath.Join(path, "main.go"))
			if err != nil {
				return fmt.Errorf("读取 %s 失败: %v", binaryName, err)
			}
			helpTexts[binaryName] = helpText
			descriptions[binaryName] = projectDescription
		}
		return nil
	})
	if err != nil {
		fmt.Println("获取帮助信息失败:", err)
		return
	}

	// 生成 README.md 内容
	var readmeContent strings.Builder

	readmeContent.WriteString(fmt.Sprintf("# %s\n\n", projectName))
	readmeContent.WriteString(fmt.Sprintf("%s\n\n", projectDescription))

	readmeContent.WriteString("## 🍺 安装\n\n")
	readmeContent.WriteString("```sh\n")
	readmeContent.WriteString("brew tap brewforge/chinese\n")
	readmeContent.WriteString("brew install <二进制命令行工具名> --formula\n")
	readmeContent.WriteString("```\n\n")

	readmeContent.WriteString("## 📋 列表\n\n")
	readmeContent.WriteString("|                     二进制命令行工具名                     |                        说明                        |\n")
	readmeContent.WriteString("| :--------------------------------------------------------: | :------------------------------------------------: |\n")
	for bin, desc := range descriptions {
		readmeContent.WriteString(fmt.Sprintf("| [%s](https://github.com/Mrered/Gobin#%s) | %s |\n", bin, bin, desc))
	}
	readmeContent.WriteString("\n")

	readmeContent.WriteString("## 🚀 使用\n\n")
	for bin, helpText := range helpTexts {
		readmeContent.WriteString(fmt.Sprintf("### %s\n\n", bin))
		readmeContent.WriteString("```sh\n")
		readmeContent.WriteString(helpText)
		readmeContent.WriteString("```\n\n")
	}

	readmeContent.WriteString("## ⚙️ 构建\n\n")
	readmeContent.WriteString("```sh\n")
	readmeContent.WriteString("# 构建所有二进制文件\n")
	readmeContent.WriteString("make build\n\n")
	readmeContent.WriteString("# 清理生成的文件\n")
	readmeContent.WriteString("make clean\n\n")
	readmeContent.WriteString("# 更新依赖\n")
	readmeContent.WriteString("make tidy\n\n")
	readmeContent.WriteString("# 显示帮助信息\n")
	readmeContent.WriteString("make help\n")
	readmeContent.WriteString("```\n\n")

	readmeContent.WriteString("## 🏆 致谢\n\n")
	readmeContent.WriteString("[Homebrew](https://brew.sh) [ChatGPT](https://chatgpt.com)\n\n")

	readmeContent.WriteString("## 📄 许可\n\n")
	readmeContent.WriteString("[MIT](https://github.com/Mrered/Gobin/blob/main/LICENSE) © [Mrered](https://github.com/Mrered)\n")

	// 写入 README.md 文件
	err = os.WriteFile(readmeFile, []byte(readmeContent.String()), 0644)
	if err != nil {
		fmt.Println("写入 README.md 文件失败:", err)
		return
	}

	fmt.Println("README.md 文件已生成")
}

// getHelpTextFromMainGo 读取 main.go 文件顶部注释内容
func getHelpTextFromMainGo(filePath string) (string, string, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", "", err
	}
	defer file.Close()

	var helpText bytes.Buffer
	var projectDescription string
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
				projectDescription = strings.TrimSpace(line)
			}
			helpText.WriteString(line + "\n")
			lineNumber++
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", "", err
	}

	return "", projectDescription, helpText.String(), nil
}
