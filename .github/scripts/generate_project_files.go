package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mrered/gobin/pkg/scripts"
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
	readmeContent.WriteString("> 请使用简体中文发起工单或拉取请求，谢谢！如果不懂简体中文，请使用 AI 翻译软件。\n\n")
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

	readmeContent.WriteString("## 👍 从本仓库开始\n\n")
	readmeContent.WriteString("本仓库实现了 CI/CD ，只需编写 Go 代码，推送后自动编译发布，自动更新 Homebrew 安装方式。\n\n")
	readmeContent.WriteString("具体功能：\n\n")
	readmeContent.WriteString("- 🌟🌟🌟🌟🌟 **对 `Make` 的支持**：\n")
	readmeContent.WriteString("```sh\n")
	readmeContent.WriteString("make build\n")
	readmeContent.WriteString("```\n")
	readmeContent.WriteString("- 🌟🌟🌟 **对 `GoReleaser` 的支持**：\n")
	readmeContent.WriteString("```yaml\n")
	readmeContent.WriteString("- name: 🚀 发布\n")
	readmeContent.WriteString("  uses: goreleaser/goreleaser-action@v6\n")
	readmeContent.WriteString("  with:\n")
	readmeContent.WriteString("    distribution: goreleaser\n")
	readmeContent.WriteString("    args: release --clean\n")
	readmeContent.WriteString("  env:\n")
	readmeContent.WriteString("    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}\n")
	readmeContent.WriteString("```\n")
	readmeContent.WriteString("- 🌟🌟🌟🌟 **自动生成 `.goreleaser.yml` 和 `README.md`**：\n\n")
	readmeContent.WriteString("    参考 [这个文件](https://github.com/Mrered/Gobin/blob/main/.github/scripts/generate_project_files.go) 和 [这个文件](https://github.com/Mrered/Gobin/blob/main/pkg/scripts/get_info.go) \n\n")
	readmeContent.WriteString("    必要条件：必须在 Go 源码顶端添加如下格式的注释，参考 [这个文件](https://github.com/Mrered/Gobin/blob/main/cmd/reportgen/main.go)\n")
	readmeContent.WriteString("```sh\n")
	readmeContent.WriteString("go run .github/scripts/generate_project_files.go\n")
	readmeContent.WriteString("```\n")
	readmeContent.WriteString("```go\n")
	readmeContent.WriteString("/*\n")
	readmeContent.WriteString("${projectName}\n")
	readmeContent.WriteString("${osInfo}\n")
	readmeContent.WriteString("${projectDescription}\n")
	readmeContent.WriteString("用法: ${projectName} [选项]\n\n")
	readmeContent.WriteString("${helpText.String()}\n")
	readmeContent.WriteString("*/\n")
	readmeContent.WriteString("```\n")
	readmeContent.WriteString("- 🌟🌟🌟 **自动生成 `Homebrew Formula Ruby` 脚本**：\n\n")
	readmeContent.WriteString("    首先使用 [这个文件](https://github.com/Mrered/Gobin/blob/main/.github/scripts/deliver_ruby_config.go) 获取所有命令行工具的信息，格式为 `JSON` ，接着使用 [这个片段](https://github.com/Mrered/Gobin/blob/c63d3021893ba3c12897da15a5f43d005fed43eb/.github/workflows/CI.yml#L97-L124) 中的代码生成 `${name}.rb` 文件\n")
	readmeContent.WriteString("```ruby\n")
	readmeContent.WriteString("class ${capitalized_name} < Formula\n")
	readmeContent.WriteString("  desc \"${desc}\"\n")
	readmeContent.WriteString("  homepage \"https://github.com/Mrered/Gobin\"\n")
	readmeContent.WriteString("  url \"https://github.com/Mrered/Gobin/archive/refs/tags/${VERSION}.tar.gz\"\n")
	readmeContent.WriteString("  sha256 \"${SHA256}\"\n")
	readmeContent.WriteString("  license \"MIT\"\n")
	readmeContent.WriteString("  head \"https://github.com/Mrered/Gobin.git\", branch: \"main\"\n\n")
	readmeContent.WriteString("  depends_on \"go\" => :build\n\n")
	readmeContent.WriteString("  def install\n")
	readmeContent.WriteString("    system \"go\", \"build\", *std_go_args(ldflags: \"-s -w\"), \"./cmd/${name}\"\n")
	readmeContent.WriteString("  end\n\n")
	readmeContent.WriteString("  test do\n")
	readmeContent.WriteString("    system bin/\"${name}\", \"-v\"\n")
	readmeContent.WriteString("  end\n")
	readmeContent.WriteString("end\n")
	readmeContent.WriteString("```\n")

	readmeContent.WriteString("## 🏆 致谢\n\n")
	readmeContent.WriteString("![Homebrew](https://img.shields.io/badge/-Homebrew-FBB040?labelColor=555555&logoColor=FFFFFF&logo=homebrew&link=https%3A%2F%2Fbrew.sh%2F) ![DeepSeek](https://img.shields.io/badge/-DeepSeek-536AF5?labelColor=555555&logoColor=FFFFFF&logo=data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiIHN0YW5kYWxvbmU9Im5vIj8+CjwhLS0gQ3JlYXRlZCB3aXRoIElua3NjYXBlIChodHRwOi8vd3d3Lmlua3NjYXBlLm9yZy8pIC0tPgoKPHN2ZwogICB2ZXJzaW9uPSIxLjEiCiAgIGlkPSJzdmcxIgogICB3aWR0aD0iMjk5LjQwNjAxIgogICBoZWlnaHQ9IjIxOS41OTg3MSIKICAgdmlld0JveD0iMCAwIDI5OS40MDYwMSAyMTkuNTk4NzEiCiAgIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIKICAgeG1sbnM6c3ZnPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CiAgPGRlZnMKICAgICBpZD0iZGVmczEiIC8+CiAgPGcKICAgICBpZD0iZzEiCiAgICAgdHJhbnNmb3JtPSJ0cmFuc2xhdGUoLTAuMjk2OTkwMjIsLTAuMjAwNjQ3NDkpIj4KICAgIDxwYXRoCiAgICAgICBzdHlsZT0iZmlsbDojZmZmZmZmO2ZpbGwtb3BhY2l0eToxO3N0cm9rZTpub25lIgogICAgICAgZD0ibSAxNjAuMDE5MjIsOS41MzM5ODc1IHYgLTIuNjY2NjcgYyAtMTcuNzA1NTIsLTQuMTQ2MDggLTMxLjA4MDI5LDQuMjkwMDc5NSAtNDgsNi45Mjc5Nzk1IC0yNi45Njk3MjMsNC4yMDQ3NSAtNDkuOTI2MDEzLC0zLjQ1NTcxIC03NC42NjY1NDMsMTMuNTEzNzMgLTYwLjMyMjIyNiw0MS4zNzQ3OTggLTM5LjkxNjI3ODIsMTI5LjY4ODUwMyAxMC42NjY1NCwxNjguODU2MjQzIDI4LjU1ODI5LDIyLjExMzUxIDY5LjI2OTAxMywyOS4zNzI4MSAxMDQuMDAwMDAzLDE4Ljk1MzQzIDExLjk4NzA0LC0zLjU5NjE1IDIyLjc1NDc2LC0xMy43NTU0NSAzNC42NjY2NywtMTYuMDU5NzkgMTEuMTc2NDcsLTIuMTYyMDkgNDEuNDcxMzcsMTAuNTcxMjEgNDguMDMyMzUsLTQuNjkyMTQgNC4yNTQyNSwtOS44OTcwMyAtMTkuMzExNTcsLTE1LjIwNTU5IC0yNS4zNjU2OSwtMTYuODMyNzggMTIuMDA4NzMsLTE3LjM2MTIyIDI1Ljg3MzcsLTMxLjcwOTI5IDMzLjE3MzU3LC01MiA0Ljk1OTAyLC0xMy43ODQxMiAzLjY5NTMzLC0zNS4wMjg0MTggMTAuNzk2NTMsLTQ3LjEwNTEwOSA2LjIwNDM4LC0xMC41NTE1MzQgMjcuMzM1NTksLTExLjgzMTA2NSAzNS45NzM3NywtMjMuNTkxNzc2IDMuNTQxNzcsLTQuODIyMDQyIDE2LjY1OTkxLC0zMC40MzE2OTggNi44MDc0OCwtMzQuNDEyODc4IC00LjY4NTg4LC0xLjg5MzQ4IC0xMC45NzEzNiw0LjczMDUxIC0xNC43NTE3Nyw2LjcwNDQyIC0xMi4wNDM1OCw2LjI4ODQ1IC0yNS4xNDU4NCw1LjYyMzA5IC0zNy4zMzI5MSwxMy4wNzIwMDMgbCAtMzIsLTQwLjAwMDAwMjUxIEMgMTg1LjYzMzkxLDkuODg2OTg3NSAyMDUuNTk4OTksNTAuODM1NTI4IDIxNy4xMzI0LDYzLjgzODIzNCBjIDMuODE5NCw0LjMwNTk1OSA4Ljk1NTA0LDE0LjE4MDU1MiAwLjY4NDQxLDE3LjUxMTgxIC04LjUyMjMyLDMuNDMyNjQ4IC0yMC42ODQzOSwtMTAuODE1NDYgLTI1Ljc5NzU5LC0xNS44NDIyNDQgLTE0LjM4NTc2LC0xNC4xNDI2NiAtNTUuNTk5MSwtMzUuODA0NDgzIC0zMiwtNTUuOTczODEyNSBNIDE4MS4zNTI1NSwxOTQuODY3MzIgYyAtMTMuNDc3MTksLTAuMDEyOCAtMjUuNzAzOTgsLTAuOTUxNDQgLTM3LjMzMzMzLC04LjQ5NDE2IC0xMS4zMjM4MSwtNy4zNDQ1NCAtMjQuNTg3ODUsLTIyLjUwNjMxIC0zOC40ODk2MywtMjQuMzc5MDUgLTEyLjc0MDY5MywtMS43MTYzMyAtMC43MzE2NywxNi45MjYyMSAtNC43NTk2LDIyLjIyOTU5IC04LjQ3MDI5MywxMS4xNTI0MyAtMzAuOTU4ODAzLC00LjI2ODkgLTM4LjA2MDI1MywtOS45NDgzOCAtMjAuNjM5MDksLTE2LjUwNjQzIC0zNC40Nzc0NiwtNDAuMzYxOTMgLTQwLjA3NTA3LC02Ni4wNzQ2NyAtMS41MjIzNywtNi45OTMwMSAtNS40MzI4MzcsLTIwLjU0NTA4NCAwLjI4NjIzLC0yNi40NzU5NjIgMTAuNTMzMTEsLTEwLjkyMzIyOCA0MC40NzAyOCwwLjMxMzI4MyA1MS43NjQ5OSw1LjEwOTUyOCAzNi41MDUyOTMsMTUuNTAxNzk0IDUzLjU0MjEyMyw0Ni4yMDE0MDQgNzYuNjUxOTEzLDc2LjAzMzEwNCA5LjAyNTYxLDExLjY1MDg2IDIxLjE0MzAzLDIwLjQxMzEzIDMwLjAxNDc1LDMyIE0gMTYxLjc4NzQ1LDg4LjE0MTE2MSBjIDEzLjI2NDMyLC00LjQ0NzIyNSAyNi4yNzkzNCwxMS43ODMwODYgMzIuMDYwNzksMjEuMzkyMDM5IDMuMjIzMTQsNS4zNTY5OSA3LjMxNjYzLDE0LjEwNjEyIC0wLjg3MTkzLDE3LjE4NTAyIC0xOC4zMzE1OCw2Ljg5MjY5IC0xNy42Nzc2NSwtMTQuMDQ5NjIgLTIzLjI1MzU0LC0yMy4zNTI1NyAtMi42ODA2MywtNC40NzI0NTMgLTguOTAxOTgsLTUuMzIxNjU3IC0xMS4zMTIyNSwtOS42MDkyNTQgLTEuNDk0NzMsLTIuNjU4OTU2IDEuMDM0MDQsLTQuODI5NzEyIDMuMzc2OTMsLTUuNjE1MjM1IG0gLTAuNTgwMTIsMTYuMzA2NDU5IHYgMy4xOTk0MyB6IgogICAgICAgaWQ9InBhdGgxIiAvPgogIDwvZz4KPC9zdmc+Cg==&link=https%3A%2F%2Fwww.deepseek.com%2F) ![ChatGPT](https://img.shields.io/badge/-ChatGPT-412991?labelColor=555555&logoColor=FFFFFF&logo=openai&link=https%3A%2F%2Fchatgpt.com%2F) ![Claude](https://img.shields.io/badge/-Claude-D97757?labelColor=555555&logoColor=FFFFFF&logo=claude&link=https%3A%2F%2Fclaude.ai%2F) ![Trae](https://img.shields.io/badge/-Trae-EC5F4A?labelColor=555555&logoColor=FFFFFF&logo=data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiIHN0YW5kYWxvbmU9Im5vIj8+CjwhLS0gQ3JlYXRlZCB3aXRoIElua3NjYXBlIChodHRwOi8vd3d3Lmlua3NjYXBlLm9yZy8pIC0tPgoKPHN2ZwogICB2ZXJzaW9uPSIxLjEiCiAgIGlkPSJzdmcxIgogICB3aWR0aD0iNjgyLjY2NjY5IgogICBoZWlnaHQ9IjY4Mi42NjY2OSIKICAgdmlld0JveD0iMCAwIDY4Mi42NjY2OSA2ODIuNjY2NjkiCiAgIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIKICAgeG1sbnM6c3ZnPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CiAgPGRlZnMKICAgICBpZD0iZGVmczEiIC8+CiAgPGcKICAgICBpZD0iZzEiPgogICAgPGcKICAgICAgIGlkPSJnNSIKICAgICAgIHRyYW5zZm9ybT0ibWF0cml4KDEuNjA3MTQyOSwwLDAsMS42MDcxNDI5LC0yMDcuMDczODEsLTE5MS45NjA5NykiCiAgICAgICBzdHlsZT0iZmlsbDojZmZmZmZmO2ZpbGwtb3BhY2l0eToxIj4KICAgICAgPHJlY3QKICAgICAgICAgc3R5bGU9ImZpbGw6I2ZmZmZmZjtmaWxsLW9wYWNpdHk6MTtzdHJva2Utd2lkdGg6Mi4yMDc1NSIKICAgICAgICAgaWQ9InJlY3QzIgogICAgICAgICB3aWR0aD0iMzM2IgogICAgICAgICBoZWlnaHQ9IjMzNiIKICAgICAgICAgeD0iMTczLjIzMTExIgogICAgICAgICB5PSIxNjMuODI3NTYiCiAgICAgICAgIHJ4PSI3LjA5OTk5OTkiIC8+CiAgICA8L2c+CiAgICA8ZwogICAgICAgaWQ9Imc2IgogICAgICAgdHJhbnNmb3JtPSJtYXRyaXgoMS42MDcxNDI5LDAsMCwxLjYwNzE0MjksLTIwNy4wNzM4MSwtMTkxLjk2MDk3KSIKICAgICAgIHN0eWxlPSJmaWxsOiM1NTU1NTU7ZmlsbC1vcGFjaXR5OjEiPgogICAgICA8cmVjdAogICAgICAgICBzdHlsZT0iZmlsbDojNTU1NTU1O2ZpbGwtb3BhY2l0eToxO3N0cm9rZTojZmZmZmZmO3N0cm9rZS13aWR0aDoxO3N0cm9rZS1saW5lam9pbjpyb3VuZDtzdHJva2UtZGFzaGFycmF5Om5vbmU7c3Ryb2tlLW9wYWNpdHk6MC42Mjg3NjMiCiAgICAgICAgIGlkPSJyZWN0NSIKICAgICAgICAgd2lkdGg9IjEwOC42MTQ4NCIKICAgICAgICAgaGVpZ2h0PSIzNi45MjEzMDciCiAgICAgICAgIHg9IjMzOS4yNzIxMyIKICAgICAgICAgeT0iNDAxLjMzNzc0IiAvPgogICAgPC9nPgogIDwvZz4KPC9zdmc+Cg==&link=https%3A%2F%2Fwww.trae.ai%2F) \n\n")

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
