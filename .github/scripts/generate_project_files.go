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
	readmeContent.WriteString("![Homebrew](https://img.shields.io/badge/-Homebrew-FBB040?labelColor=555555&logoColor=FFFFFF&logo=homebrew&link=https%3A%2F%2Fbrew.sh%2F) ![DeepSeek](https://img.shields.io/badge/-DeepSeek-536AF5?labelColor=555555&logoColor=FFFFFF&logo=data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiIHN0YW5kYWxvbmU9Im5vIj8+CjxzdmcKICAgd2lkdGg9IjU2LjIwMjQxMiIKICAgaGVpZ2h0PSI0MS4zNTk0NTkiCiAgIHZpZXdCb3g9IjAgMCA1Ni4yMDI0MTIgNDEuMzU5NDg0IgogICBmaWxsPSJub25lIgogICB2ZXJzaW9uPSIxLjEiCiAgIGlkPSJzdmc3IgogICB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciCiAgIHhtbG5zOnN2Zz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgogIDxkZXNjCiAgICAgaWQ9ImRlc2MxIj4KCQkJQ3JlYXRlZCB3aXRoIFBpeHNvLgoJPC9kZXNjPgogIDxkZWZzCiAgICAgaWQ9ImRlZnMxIiAvPgogIDxwYXRoCiAgICAgaWQ9InBhdGg3IgogICAgIGQ9Im0gNTUuNjEyNzc0LDMuNDcxMjIyNSBjIC0wLjU5NTMsLTAuMjkxNzUgLTAuODUxNywwLjI2NDE2IC0xLjE5OTgsMC41NDY2MyAtMC4xMTkxLDAuMDkxMDcgLTAuMjE5OCwwLjIwOTQ3IC0wLjMyMDYsMC4zMTg4NSAtMC44NzAxLDAuOTI5MiAtMS44ODY3LDEuNTM5NzkgLTMuMjE0OCwxLjQ2NjggLTEuOTQxNywtMC4xMDkzOCAtMy41OTk1LDAuNTAxMjIgLTUuMDY1LDEuOTg2MzIgLTAuMzExNCwtMS44MzEyOSAtMS4zNDYzLC0yLjkyNDggLTIuOTIxNywtMy42MjYyMiAtMC44MjQyLC0wLjM2NDUgLTEuNjU3NywtMC43MjkgLTIuMjM0OCwtMS41MjE3MiAtMC40MDI5LC0wLjU2NDcgLTAuNTEyOSwtMS4xOTMzNiAtMC43MTQzLC0xLjgxMjk5MzAyIC0wLjEyODMsLTAuMzczNTM1IC0wLjI1NjUsLTAuNzU2MzQ3MiAtMC42ODcsLTAuODIwMDY3OTQgLTAuNDY3MSwtMC4wNzI3NTM4NiAtMC42NTAzLDAuMzE4ODQ3OTQgLTAuODMzNSwwLjY0Njk3Mjk0IC0wLjczMjcsMS4zMzkzNTgwMiAtMS4wMTY2LDIuODE1NDI4MDIgLTAuOTg5Miw0LjMwOTU2ODAyIDAuMDY0MSwzLjM2MjA2IDEuNDgzOCw2LjA0MDU2OTUgNC4zMDQ4LDcuOTQ0ODY5NSAwLjMyMDYsMC4yMTg3IDAuNDAzLDAuNDM3MiAwLjMwMjMsMC43NTYzIC0wLjE5MjQsMC42NTYgLTAuNDIxNCwxLjI5MzcgLTAuNjIyOCwxLjk0OTcgLTAuMTI4MywwLjQxOTIgLTAuMzIwNywwLjUxMDMgLTAuNzY5NCwwLjMyNzkgLTEuNTQ3OSwtMC42NDY3IC0yLjg4NTIsLTEuNjAzNSAtNC4wNjY3LC0yLjc2MDUgLTIuMDA1OCwtMS45NDA3IC0zLjgxOTMsLTQuMDgxNzg5NSAtNi4wODE1LC01Ljc1ODMwOTUgLTAuNTMxMywtMC4zOTE4NCAtMS4wNjI1LC0wLjc1NjEgLTEuNjEyMSwtMS4xMDI1NCAtMi4zMDgxLC0yLjI0MTIxIDAuMzAyMywtNC4wODE3OCAwLjkwNjgsLTQuMzAwMjkgMC42MzE5LC0wLjIyNzc4IDAuMjE5OCwtMS4wMTE0NyAtMS44MjI3LC0xLjAwMjIgLTIuMDQyNSwwLjAwOTA0IC0zLjkxMDksMC42OTIzOSAtNi4yOTIyLDEuNjAzNTIgLTAuMzQ4LDAuMTM2NzIgLTAuNzE0NSwwLjIzNjgyIC0xLjA5LDAuMzE4ODUgLTIuMTYxNSwtMC40MDk5MiAtNC40MDU1LC0wLjUwMTIyIC02Ljc1MDIsLTAuMjM2ODIgLTQuNDE0Njg5OCwwLjQ5MTk0IC03Ljk0MDgxOTgsMi41NzgzNyAtMTAuNTMyODU5OCw2LjE0MDg3IC0zLjExNDEzNDAzLDQuMjgyMjE5NSAtMy44NDY4MDAwMyw5LjE0NzQxOTUgLTIuOTQ5MjE3MDMsMTQuMjIyNDE5NSAwLjk0MzM1NzAzLDUuMzQ4MSAzLjY3MjcyNzAzLDkuNzc2MSA3Ljg2NzU1NzAzLDEzLjIzODUgNC4zNTA2MTk4LDMuNTg5NiA5LjM2MDYxOTgsNS4zNDgyIDE1LjA3NTgxOTgsNS4wMTEgMy40NzEzLC0wLjIwMDQgNy4zMzY0LC0wLjY2NSAxMS42OTYxLC00LjM1NSAxLjA5OSwwLjU0NjcgMi4yNTMxLDAuNzY1MiA0LjE2NzQsMC45MjkyIDEuNDc0NiwwLjEzNjcgMi44OTQzLC0wLjA3MjcgMy45OTMzLC0wLjMwMDUgMS43MjE5LC0wLjM2NDUgMS42MDI5LC0xLjk1OSAwLjk4MDEsLTIuMjUwNSAtNS4wNDY2LC0yLjM1MDYgLTMuOTM4NSwtMS4zOTQgLTQuOTQ1OSwtMi4xNjg1IDIuNTY0NSwtMy4wMzM5IDYuNDI5NywtNi4xODY1IDcuOTQwOSwtMTYuNDAwMSAwLjExOSwtMC44MTA4IDAuMDE4MywtMS4zMjExIDAsLTEuOTc3MSAtMC4wMDkyLC0wLjQwMDggMC4wODI0LC0wLjU1NTYgMC41NDA0LC0wLjYwMTMgMS4yNjM5LC0wLjE0NTcgMi40OTEyLC0wLjQ5MTkgMy42MTc4LC0xLjExMTUgMy4yNjk4LC0xLjc4NTcgNC41ODg2LC00LjcxOTUyOTUgNC45LC04LjIzNjM3OTUgMC4wNDU5LC0wLjUzNzU5IC0wLjAwOTEsLTEuMDkzNSAtMC41NzcsLTEuMzc1NzMgeiBtIC0yOC40OTM4LDMxLjY1MTgwOTUgYyAtNC44OTA5LC0zLjg0NDcgLTcuMjYzLC01LjExMTMgLTguMjQzMSwtNS4wNTY2IC0wLjkxNTksMC4wNTQ3IC0wLjc1MSwxLjEwMjUgLTAuNTQ5NiwxLjc4NTkgMC4yMTA3LDAuNjc0MSAwLjQ4NTUsMS4xMzg5IDAuODcwMSwxLjczMSAwLjI2NTYsMC4zOTE4IDAuNDQ4OSwwLjk3NDggLTAuMjY1NSwxLjQxMjMgLTEuNTc1NCwwLjk3NDkgLTQuMzE0LC0wLjMyODEgLTQuNDQyMywtMC4zOTE4IC0zLjE4NzIsLTEuODc3IC01Ljg1MjQ4OTgsLTQuMzU1MyAtNy43MzAxNzk4LC03Ljc0NDQgLTEuODEzNDcsLTMuMjYyIC0yLjg2NjcsLTYuNzYwNSAtMy4wNDA3NywtMTAuNDk2MSAtMC4wNDU3NywtMC45MDE5IDAuMjE5ODUsLTEuMjIxIDEuMTE3NDMsLTEuMzg0OCAxLjE4MTUyLC0wLjIxODcgMi4zOTk2NiwtMC4yNjQ0IDMuNTgxMTgsLTAuMDkxMyA0Ljk5MTczOTgsMC43MjkgOS4yNDE0Mzk4LDIuOTYxMiAxMi44MDQzMzk4LDYuNDk2MyAyLjAzMzMsMi4wMTM1IDMuNTcyLDQuNDE5IDUuMTU2Niw2Ljc2OTYgMS42ODUyLDIuNDk2MyAzLjQ5ODcsNC44NzQ1IDUuODA2OCw2LjgyNDIgMC44MTUxLDAuNjgzMyAxLjQ2NTQsMS4yMDI2IDIuMDg4MiwxLjU4NTQgLTEuODc3NSwwLjIwOTUgLTUuMDEsMC4yNTUyIC03LjE1MzIsLTEuNDM5NyB6IG0gMi4zNDQ3LC0xNS4wNzg4IGMgMCwtMC40MDA5IDAuMzIwNiwtMC43MTk3IDAuNzIzNywtMC43MTk3IDAuMDkxNSwwIDAuMTczOSwwLjAxOCAwLjI0NzIsMC4wNDU0IDAuMTAwOCwwLjAzNjYgMC4xOTI0LDAuMDkxMyAwLjI2NTYsMC4xNzMxIDAuMTI4MywwLjEyNzcgMC4yMDE1LDAuMzA5OCAwLjIwMTUsMC41MDEyIDAsMC40MDA5IC0wLjMyMDUsMC43MTk3IC0wLjcyMzUsMC43MTk3IC0wLjQwMzEsMCAtMC43MTQ1LC0wLjMxODggLTAuNzE0NSwtMC43MTk3IHogbSA3LjI4MTUsMy43MzU2IGMgLTAuNDY3MSwwLjE5MTQgLTAuOTM0MiwwLjM1NTIgLTEuMzgzLDAuMzczNSAtMC42OTYxLDAuMDM2NCAtMS40NTYzLC0wLjI0NjEgLTEuODY4NCwtMC41OTIzIC0wLjY0MTEsLTAuNTM3NiAtMS4wOTkxLC0wLjgzODEgLTEuMjkxNSwtMS43NzY2IC0wLjA4MjQsLTAuNDAwOSAtMC4wMzY3LC0xLjAyMDUgMC4wMzY3LC0xLjM3NTcgMC4xNjQ4LC0wLjc2NTQgLTAuMDE4NCwtMS4yNTczIC0wLjU1ODcsLTEuNzAzOSAtMC40Mzk3LC0wLjM2NDUgLTAuOTk4NCwtMC40NjQ2IC0xLjYxMjEsLTAuNDY0NiAtMC4yMjksMCAtMC40Mzk1LC0wLjEwMDMgLTAuNTk1MywtMC4xODIzIC0wLjI1NjUsLTAuMTI3NSAtMC40NjcsLTAuNDQ2MyAtMC4yNjU2LC0wLjgzODIgMC4wNjQxLC0wLjEyNzQgMC4zNzU2LC0wLjQzNzIgMC40NDg4LC0wLjQ5MTkgMC44MzM1LC0wLjQ3MzkgMS43OTUyLC0wLjMxODkgMi42ODM2LDAuMDM2NCAwLjgyNDQsMC4zMzcxIDEuNDQ3MiwwLjk1NjcgMi4zNDQ3LDEuODMxMyAwLjkxNTksMS4wNTY4IDEuMDgwNywxLjM0ODYgMS42MDI4LDIuMTQxMSAwLjQxMjMsMC42MTk2IDAuNzg3OCwxLjI1NzMgMS4wNDQyLDEuOTg2MyAwLjE1NTcsMC40NTU2IC0wLjA0NTgsMC44MjkxIC0wLjU4NjIsMS4wNTY5IHoiCiAgICAgZmlsbC1ydWxlPSJub256ZXJvIgogICAgIGZpbGw9IiM0RDZCRkUiCiAgICAgc3R5bGU9ImZpbGw6I2ZmZmZmZjtmaWxsLW9wYWNpdHk6MSIgLz4KPC9zdmc+Cg==&link=https%3A%2F%2Fwww.deepseek.com%2F) ![ChatGPT](https://img.shields.io/badge/-ChatGPT-412991?labelColor=555555&logoColor=FFFFFF&logo=openai&link=https%3A%2F%2Fchatgpt.com%2F) ![Claude](https://img.shields.io/badge/-Claude-D97757?labelColor=555555&logoColor=FFFFFF&logo=claude&link=https%3A%2F%2Fclaude.ai%2F) ![Trae](https://img.shields.io/badge/-Trae-EC5F4A?labelColor=555555&logoColor=FFFFFF&logo=data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiIHN0YW5kYWxvbmU9Im5vIj8+CjwhLS0gQ3JlYXRlZCB3aXRoIElua3NjYXBlIChodHRwOi8vd3d3Lmlua3NjYXBlLm9yZy8pIC0tPgoKPHN2ZwogICB2ZXJzaW9uPSIxLjEiCiAgIGlkPSJzdmcxIgogICB3aWR0aD0iNjgyLjY2NjY5IgogICBoZWlnaHQ9IjY4Mi42NjY2OSIKICAgdmlld0JveD0iMCAwIDY4Mi42NjY2OSA2ODIuNjY2NjkiCiAgIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIKICAgeG1sbnM6c3ZnPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CiAgPGRlZnMKICAgICBpZD0iZGVmczEiIC8+CiAgPGcKICAgICBpZD0iZzEiPgogICAgPGcKICAgICAgIGlkPSJnNSIKICAgICAgIHRyYW5zZm9ybT0ibWF0cml4KDEuNjA3MTQyOSwwLDAsMS42MDcxNDI5LC0yMDcuMDczODEsLTE5MS45NjA5NykiCiAgICAgICBzdHlsZT0iZmlsbDojZmZmZmZmO2ZpbGwtb3BhY2l0eToxIj4KICAgICAgPHJlY3QKICAgICAgICAgc3R5bGU9ImZpbGw6I2ZmZmZmZjtmaWxsLW9wYWNpdHk6MTtzdHJva2Utd2lkdGg6Mi4yMDc1NSIKICAgICAgICAgaWQ9InJlY3QzIgogICAgICAgICB3aWR0aD0iMzM2IgogICAgICAgICBoZWlnaHQ9IjMzNiIKICAgICAgICAgeD0iMTczLjIzMTExIgogICAgICAgICB5PSIxNjMuODI3NTYiCiAgICAgICAgIHJ4PSI3LjA5OTk5OTkiIC8+CiAgICA8L2c+CiAgICA8ZwogICAgICAgaWQ9Imc2IgogICAgICAgdHJhbnNmb3JtPSJtYXRyaXgoMS42MDcxNDI5LDAsMCwxLjYwNzE0MjksLTIwNy4wNzM4MSwtMTkxLjk2MDk3KSIKICAgICAgIHN0eWxlPSJmaWxsOiM1NTU1NTU7ZmlsbC1vcGFjaXR5OjEiPgogICAgICA8cmVjdAogICAgICAgICBzdHlsZT0iZmlsbDojNTU1NTU1O2ZpbGwtb3BhY2l0eToxO3N0cm9rZTojZmZmZmZmO3N0cm9rZS13aWR0aDoxO3N0cm9rZS1saW5lam9pbjpyb3VuZDtzdHJva2UtZGFzaGFycmF5Om5vbmU7c3Ryb2tlLW9wYWNpdHk6MC42Mjg3NjMiCiAgICAgICAgIGlkPSJyZWN0NSIKICAgICAgICAgd2lkdGg9IjEwOC42MTQ4NCIKICAgICAgICAgaGVpZ2h0PSIzNi45MjEzMDciCiAgICAgICAgIHg9IjMzOS4yNzIxMyIKICAgICAgICAgeT0iNDAxLjMzNzc0IiAvPgogICAgPC9nPgogIDwvZz4KPC9zdmc+Cg==&link=https%3A%2F%2Fwww.trae.ai%2F) \n\n")

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
