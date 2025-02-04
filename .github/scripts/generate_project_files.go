package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Mrered/Gobin/pkg/scripts@test-CI"
)

func main() {
	binDir := "cmd"
	readmeFile := "README.md"
	goreleaserFile := ".goreleaser.yml"

	// 项目名称和描述
	projectName := "Gobin"
	projectDescription := "Go 二进制小程序"

	// 收集帮助信息
	helpTexts := make(map[string]string)
	descriptions := make(map[string]string)
	goosInfo := make(map[string][]string)
	var binaries []string
	err := filepath.Walk(binDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != binDir {
			binaryName := filepath.Base(path)
			binaries = append(binaries, binaryName)
			_, projectDescription, osInfo, helpText, err := scripts.GetHelpTextFromMainGo(filepath.Join(path, "main.go"))
			if err != nil {
				return fmt.Errorf("读取 %s 失败: %v", binaryName, err)
			}
			helpTexts[binaryName] = helpText
			descriptions[binaryName] = projectDescription
			goosInfo[binaryName] = osInfo
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

	readmeContent.WriteString("![Homebrew](https://img.shields.io/badge/-Homebrew-FBB040?labelColor=555555&logoColor=FFFFFF&logo=homebrew) ![CI](https://github.com/Mrered/Gobin/actions/workflows/CI.yml/badge.svg) ![license](https://img.shields.io/github/license/Mrered/Gobin) ![code-size](https://img.shields.io/github/languages/code-size/Mrered/Gobin) ![repo-size](https://img.shields.io/github/repo-size/Mrered/Gobin)\n\n")
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
	readmeContent.WriteString("[Homebrew](https://brew.sh) [ChatGPT](https://chatgpt.com) [Claude](https://claude.ai)\n\n")

	readmeContent.WriteString("## 📄 许可\n\n")
	readmeContent.WriteString("[MIT](https://github.com/Mrered/Gobin/blob/main/LICENSE) © [Mrered](https://github.com/Mrered)\n")

	// 写入 README.md 文件
	err = os.WriteFile(readmeFile, []byte(readmeContent.String()), 0644)
	if err != nil {
		fmt.Println("写入 README.md 文件失败:", err)
		return
	}

	fmt.Println("README.md 文件已生成")

	// 生成 .goreleaser.yml 内容
	var goreleaserContent strings.Builder

	goreleaserContent.WriteString("version: 2\n")
	goreleaserContent.WriteString(fmt.Sprintf("project_name: %s\n\n", projectName))

	goreleaserContent.WriteString("builds:\n")
	for _, binary := range binaries {
		goreleaserContent.WriteString(fmt.Sprintf("  - id: %s\n", binary))
		goreleaserContent.WriteString(fmt.Sprintf("    dir: ./cmd/%s\n", binary))
		goreleaserContent.WriteString(fmt.Sprintf("    binary: %s\n", binary))
		goreleaserContent.WriteString("    goos:\n")
		for _, os := range goosInfo[binary] {
			goreleaserContent.WriteString(fmt.Sprintf("      - %s\n", os))
		}
		goreleaserContent.WriteString("    goarch:\n")
		goreleaserContent.WriteString("      - amd64\n")
		goreleaserContent.WriteString("      - arm64\n")
		goreleaserContent.WriteString("    env:\n")
		goreleaserContent.WriteString("      - CGO_ENABLED=0\n\n")
	}

	goreleaserContent.WriteString("archives:\n")
	for _, binary := range binaries {
		goreleaserContent.WriteString(fmt.Sprintf("  - id: %s\n", binary))
		goreleaserContent.WriteString(fmt.Sprintf("    builds: [%s]\n", binary))
		goreleaserContent.WriteString("    format: tar.gz\n")
		goreleaserContent.WriteString("    name_template: \"{{ .Binary }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}\"\n")
		goreleaserContent.WriteString("    files:\n")
		goreleaserContent.WriteString("      - LICENSE\n")
		goreleaserContent.WriteString("      - README.md\n\n")
	}

	goreleaserContent.WriteString("release:\n")
	goreleaserContent.WriteString("  github:\n")
	goreleaserContent.WriteString("    owner: Mrered\n")
	goreleaserContent.WriteString("    name: Gobin\n")

	// 写入 .goreleaser.yml 文件
	err = os.WriteFile(goreleaserFile, []byte(goreleaserContent.String()), 0644)
	if err != nil {
		fmt.Println("写入 .goreleaser.yml 文件失败:", err)
		return
	}

	fmt.Println(".goreleaser.yml 文件已生成")
}

// getHelpTextFromMainGo 读取 main.go 文件顶部注释内容
func getHelpTextFromMainGo(filePath string) (string, string, []string, string, error) {
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
