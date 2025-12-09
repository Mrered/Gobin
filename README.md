# Gobin

Go 二进制小程序

![Homebrew](https://img.shields.io/badge/-Homebrew-FBB040?labelColor=555555&logoColor=FFFFFF&logo=homebrew) ![CI](https://github.com/Mrered/Gobin/actions/workflows/CI.yml/badge.svg) ![license](https://img.shields.io/github/license/Mrered/Gobin) ![code-size](https://img.shields.io/github/languages/code-size/Mrered/Gobin) ![repo-size](https://img.shields.io/github/repo-size/Mrered/Gobin)

> 请使用简体中文发起工单或拉取请求，谢谢！如果不懂简体中文，请使用 AI 翻译软件。

## 🍺 安装

```sh
brew tap brewforge/chinese
brew install <二进制命令行工具名> --formula
```

## 📋 列表

|                     二进制命令行工具名                     |                        说明                        |
| :--------------------------------------------------------: | :------------------------------------------------: |
| [reportgen](https://github.com/Mrered/Gobin#reportgen) | 生成报告 |
| [shicaojiaoan](https://github.com/Mrered/Gobin#shicaojiaoan) | 实操教案格式化生成器 |
| [hexinsuyangsummary](https://github.com/Mrered/Gobin#hexinsuyangsummary) | 核心素养汇总工具 |
| [makemf](https://github.com/Mrered/Gobin#makemf) | 为 GGUF 文件生成 Makefile |
| [ollamaplist](https://github.com/Mrered/Gobin#ollamaplist) | 给通过 Homebrew 安装的 Ollama CLI 工具添加环境变量 |

## 🚀 使用

### shicaojiaoan

```sh
用法: shicaojiaoan [选项] [输入文件]

选项:
  -h    显示帮助信息
  -p    生成 PDF 文件（需要安装 typst）
  -t    生成空白模板文件 template.md
  -v    显示详细输出信息
```

### hexinsuyangsummary

```sh
用法: hexinsuyangsummary [选项]

选项:
  -h    显示帮助信息
  -p string
        输入目录路径
  -c string
        指定一个模板 Excel 文件 (用于读取 H3/H4/H5)
  -m    开启修改模式 (将模板数据写入目标文件)
  -a    全量输出（database-style 详细记录模式）
```

### makemf

```sh
用法: makemf [选项]

选项:
  -a    自动为当前目录下的所有 .gguf 文件生成 Makefile
  -h    显示帮助信息
  -m string
        GGUF 文件名称，包含后缀名
  -n string
        要生成的 Makefile 名称
  -v    显示版本号
```

### ollamaplist

```sh
用法: ollamaplist [选项]

选项:
  -a    应用默认配置
  -h    显示帮助信息
  -m string
        OLLAMA_MAX_LOADED_MODELS (default "2")
  -o string
        OLLAMA_ORIGINS (default "*")
  -p string
        OLLAMA_NUM_PARALLEL (default "4")
  -r    删除所有环境变量
  -s string
        OLLAMA_HOST (default "0.0.0.0")
  -v    显示版本号
```

### reportgen

```sh
用法: reportgen [选项]

选项:
  -d string
        指定工作目录
  -f    是否格式化内容
  -h    显示帮助信息
  -m string
        指定月份 (格式: YYYYMM)
  -s string
        指定学期 (格式: YYYY - YYYY 春/秋)
  -t string
        指定报告类型 (w: 周报, m: 月报, s: 学期报, y: 年报)
  -v    显示版本号
  -w string
        指定周数
  -y string
        指定年份 (格式: YYYY)
```

## ⚙️ 构建

```sh
# 构建所有二进制文件
make build

# 清理生成的文件
make clean

# 更新依赖
make tidy

# 显示帮助信息
make help
```

## 👍 从本仓库开始

本仓库实现了 CI/CD ，只需编写 Go 代码，推送后自动编译发布，自动更新 Homebrew 安装方式。

具体功能：

- 🌟🌟🌟🌟🌟 **对 `Make` 的支持**：
```sh
make build
```
- 🌟🌟🌟 **对 `GoReleaser` 的支持**：
```yaml
- name: 🚀 发布
  uses: goreleaser/goreleaser-action@v6
  with:
    distribution: goreleaser
    args: release --clean
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```
- 🌟🌟🌟🌟 **自动生成 `.goreleaser.yml` 和 `README.md`**：

    参考 [这个文件](https://github.com/Mrered/Gobin/blob/main/.github/scripts/generate_project_files.go) 和 [这个文件](https://github.com/Mrered/Gobin/blob/main/pkg/scripts/get_info.go) 

    必要条件：必须在 Go 源码顶端添加如下格式的注释，参考 [这个文件](https://github.com/Mrered/Gobin/blob/main/cmd/reportgen/main.go)
```sh
go run .github/scripts/generate_project_files.go
```
```go
/*
${projectName}
${osInfo}
${projectDescription}
用法: ${projectName} [选项]

${helpText.String()}
*/
```
- 🌟🌟🌟 **自动生成 `Homebrew Formula Ruby` 脚本**：

    首先使用 [这个文件](https://github.com/Mrered/Gobin/blob/main/.github/scripts/deliver_ruby_config.go) 获取所有命令行工具的信息，格式为 `JSON` ，接着使用 [这个片段](https://github.com/Mrered/Gobin/blob/c63d3021893ba3c12897da15a5f43d005fed43eb/.github/workflows/CI.yml#L97-L124) 中的代码生成 `${name}.rb` 文件
```ruby
class ${capitalized_name} < Formula
  desc "${desc}"
  homepage "https://github.com/Mrered/Gobin"
  url "https://github.com/Mrered/Gobin/archive/refs/tags/${VERSION}.tar.gz"
  sha256 "${SHA256}"
  license "MIT"
  head "https://github.com/Mrered/Gobin.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/${name}"
  end

  test do
    system bin/"${name}", "-v"
  end
end
```
## 🏆 致谢

![Homebrew](https://img.shields.io/badge/-Homebrew-FBB040?labelColor=555555&logoColor=FFFFFF&logo=homebrew&link=https%3A%2F%2Fbrew.sh%2F) ![DeepSeek](https://img.shields.io/badge/-DeepSeek-536AF5?labelColor=555555&logoColor=FFFFFF&logo=data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiIHN0YW5kYWxvbmU9Im5vIj8+CjwhLS0gQ3JlYXRlZCB3aXRoIElua3NjYXBlIChodHRwOi8vd3d3Lmlua3NjYXBlLm9yZy8pIC0tPgoKPHN2ZwogICB2ZXJzaW9uPSIxLjEiCiAgIGlkPSJzdmcxIgogICB3aWR0aD0iMjk5LjQwNjAxIgogICBoZWlnaHQ9IjIxOS41OTg3MSIKICAgdmlld0JveD0iMCAwIDI5OS40MDYwMSAyMTkuNTk4NzEiCiAgIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIKICAgeG1sbnM6c3ZnPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CiAgPGRlZnMKICAgICBpZD0iZGVmczEiIC8+CiAgPGcKICAgICBpZD0iZzEiCiAgICAgdHJhbnNmb3JtPSJ0cmFuc2xhdGUoLTAuMjk2OTkwMjIsLTAuMjAwNjQ3NDkpIj4KICAgIDxwYXRoCiAgICAgICBzdHlsZT0iZmlsbDojZmZmZmZmO2ZpbGwtb3BhY2l0eToxO3N0cm9rZTpub25lIgogICAgICAgZD0ibSAxNjAuMDE5MjIsOS41MzM5ODc1IHYgLTIuNjY2NjcgYyAtMTcuNzA1NTIsLTQuMTQ2MDggLTMxLjA4MDI5LDQuMjkwMDc5NSAtNDgsNi45Mjc5Nzk1IC0yNi45Njk3MjMsNC4yMDQ3NSAtNDkuOTI2MDEzLC0zLjQ1NTcxIC03NC42NjY1NDMsMTMuNTEzNzMgLTYwLjMyMjIyNiw0MS4zNzQ3OTggLTM5LjkxNjI3ODIsMTI5LjY4ODUwMyAxMC42NjY1NCwxNjguODU2MjQzIDI4LjU1ODI5LDIyLjExMzUxIDY5LjI2OTAxMywyOS4zNzI4MSAxMDQuMDAwMDAzLDE4Ljk1MzQzIDExLjk4NzA0LC0zLjU5NjE1IDIyLjc1NDc2LC0xMy43NTU0NSAzNC42NjY2NywtMTYuMDU5NzkgMTEuMTc2NDcsLTIuMTYyMDkgNDEuNDcxMzcsMTAuNTcxMjEgNDguMDMyMzUsLTQuNjkyMTQgNC4yNTQyNSwtOS44OTcwMyAtMTkuMzExNTcsLTE1LjIwNTU5IC0yNS4zNjU2OSwtMTYuODMyNzggMTIuMDA4NzMsLTE3LjM2MTIyIDI1Ljg3MzcsLTMxLjcwOTI5IDMzLjE3MzU3LC01MiA0Ljk1OTAyLC0xMy43ODQxMiAzLjY5NTMzLC0zNS4wMjg0MTggMTAuNzk2NTMsLTQ3LjEwNTEwOSA2LjIwNDM4LC0xMC41NTE1MzQgMjcuMzM1NTksLTExLjgzMTA2NSAzNS45NzM3NywtMjMuNTkxNzc2IDMuNTQxNzcsLTQuODIyMDQyIDE2LjY1OTkxLC0zMC40MzE2OTggNi44MDc0OCwtMzQuNDEyODc4IC00LjY4NTg4LC0xLjg5MzQ4IC0xMC45NzEzNiw0LjczMDUxIC0xNC43NTE3Nyw2LjcwNDQyIC0xMi4wNDM1OCw2LjI4ODQ1IC0yNS4xNDU4NCw1LjYyMzA5IC0zNy4zMzI5MSwxMy4wNzIwMDMgbCAtMzIsLTQwLjAwMDAwMjUxIEMgMTg1LjYzMzkxLDkuODg2OTg3NSAyMDUuNTk4OTksNTAuODM1NTI4IDIxNy4xMzI0LDYzLjgzODIzNCBjIDMuODE5NCw0LjMwNTk1OSA4Ljk1NTA0LDE0LjE4MDU1MiAwLjY4NDQxLDE3LjUxMTgxIC04LjUyMjMyLDMuNDMyNjQ4IC0yMC42ODQzOSwtMTAuODE1NDYgLTI1Ljc5NzU5LC0xNS44NDIyNDQgLTE0LjM4NTc2LC0xNC4xNDI2NiAtNTUuNTk5MSwtMzUuODA0NDgzIC0zMiwtNTUuOTczODEyNSBNIDE4MS4zNTI1NSwxOTQuODY3MzIgYyAtMTMuNDc3MTksLTAuMDEyOCAtMjUuNzAzOTgsLTAuOTUxNDQgLTM3LjMzMzMzLC04LjQ5NDE2IC0xMS4zMjM4MSwtNy4zNDQ1NCAtMjQuNTg3ODUsLTIyLjUwNjMxIC0zOC40ODk2MywtMjQuMzc5MDUgLTEyLjc0MDY5MywtMS43MTYzMyAtMC43MzE2NywxNi45MjYyMSAtNC43NTk2LDIyLjIyOTU5IC04LjQ3MDI5MywxMS4xNTI0MyAtMzAuOTU4ODAzLC00LjI2ODkgLTM4LjA2MDI1MywtOS45NDgzOCAtMjAuNjM5MDksLTE2LjUwNjQzIC0zNC40Nzc0NiwtNDAuMzYxOTMgLTQwLjA3NTA3LC02Ni4wNzQ2NyAtMS41MjIzNywtNi45OTMwMSAtNS40MzI4MzcsLTIwLjU0NTA4NCAwLjI4NjIzLC0yNi40NzU5NjIgMTAuNTMzMTEsLTEwLjkyMzIyOCA0MC40NzAyOCwwLjMxMzI4MyA1MS43NjQ5OSw1LjEwOTUyOCAzNi41MDUyOTMsMTUuNTAxNzk0IDUzLjU0MjEyMyw0Ni4yMDE0MDQgNzYuNjUxOTEzLDc2LjAzMzEwNCA5LjAyNTYxLDExLjY1MDg2IDIxLjE0MzAzLDIwLjQxMzEzIDMwLjAxNDc1LDMyIE0gMTYxLjc4NzQ1LDg4LjE0MTE2MSBjIDEzLjI2NDMyLC00LjQ0NzIyNSAyNi4yNzkzNCwxMS43ODMwODYgMzIuMDYwNzksMjEuMzkyMDM5IDMuMjIzMTQsNS4zNTY5OSA3LjMxNjYzLDE0LjEwNjEyIC0wLjg3MTkzLDE3LjE4NTAyIC0xOC4zMzE1OCw2Ljg5MjY5IC0xNy42Nzc2NSwtMTQuMDQ5NjIgLTIzLjI1MzU0LC0yMy4zNTI1NyAtMi42ODA2MywtNC40NzI0NTMgLTguOTAxOTgsLTUuMzIxNjU3IC0xMS4zMTIyNSwtOS42MDkyNTQgLTEuNDk0NzMsLTIuNjU4OTU2IDEuMDM0MDQsLTQuODI5NzEyIDMuMzc2OTMsLTUuNjE1MjM1IG0gLTAuNTgwMTIsMTYuMzA2NDU5IHYgMy4xOTk0MyB6IgogICAgICAgaWQ9InBhdGgxIiAvPgogIDwvZz4KPC9zdmc+Cg==&link=https%3A%2F%2Fwww.deepseek.com%2F) ![ChatGPT](https://img.shields.io/badge/-ChatGPT-412991?labelColor=555555&logoColor=FFFFFF&logo=openai&link=https%3A%2F%2Fchatgpt.com%2F) ![Claude](https://img.shields.io/badge/-Claude-D97757?labelColor=555555&logoColor=FFFFFF&logo=claude&link=https%3A%2F%2Fclaude.ai%2F) ![Trae](https://img.shields.io/badge/-Trae-EC5F4A?labelColor=555555&logoColor=FFFFFF&logo=data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiIHN0YW5kYWxvbmU9Im5vIj8+CjwhLS0gQ3JlYXRlZCB3aXRoIElua3NjYXBlIChodHRwOi8vd3d3Lmlua3NjYXBlLm9yZy8pIC0tPgoKPHN2ZwogICB2ZXJzaW9uPSIxLjEiCiAgIGlkPSJzdmcxIgogICB3aWR0aD0iNjgyLjY2NjY5IgogICBoZWlnaHQ9IjY4Mi42NjY2OSIKICAgdmlld0JveD0iMCAwIDY4Mi42NjY2OSA2ODIuNjY2NjkiCiAgIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIKICAgeG1sbnM6c3ZnPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CiAgPGRlZnMKICAgICBpZD0iZGVmczEiIC8+CiAgPGcKICAgICBpZD0iZzEiPgogICAgPGcKICAgICAgIGlkPSJnNSIKICAgICAgIHRyYW5zZm9ybT0ibWF0cml4KDEuNjA3MTQyOSwwLDAsMS42MDcxNDI5LC0yMDcuMDczODEsLTE5MS45NjA5NykiCiAgICAgICBzdHlsZT0iZmlsbDojZmZmZmZmO2ZpbGwtb3BhY2l0eToxIj4KICAgICAgPHJlY3QKICAgICAgICAgc3R5bGU9ImZpbGw6I2ZmZmZmZjtmaWxsLW9wYWNpdHk6MTtzdHJva2Utd2lkdGg6Mi4yMDc1NSIKICAgICAgICAgaWQ9InJlY3QzIgogICAgICAgICB3aWR0aD0iMzM2IgogICAgICAgICBoZWlnaHQ9IjMzNiIKICAgICAgICAgeD0iMTczLjIzMTExIgogICAgICAgICB5PSIxNjMuODI3NTYiCiAgICAgICAgIHJ4PSI3LjA5OTk5OTkiIC8+CiAgICA8L2c+CiAgICA8ZwogICAgICAgaWQ9Imc2IgogICAgICAgdHJhbnNmb3JtPSJtYXRyaXgoMS42MDcxNDI5LDAsMCwxLjYwNzE0MjksLTIwNy4wNzM4MSwtMTkxLjk2MDk3KSIKICAgICAgIHN0eWxlPSJmaWxsOiM1NTU1NTU7ZmlsbC1vcGFjaXR5OjEiPgogICAgICA8cmVjdAogICAgICAgICBzdHlsZT0iZmlsbDojNTU1NTU1O2ZpbGwtb3BhY2l0eToxO3N0cm9rZTojZmZmZmZmO3N0cm9rZS13aWR0aDoxO3N0cm9rZS1saW5lam9pbjpyb3VuZDtzdHJva2UtZGFzaGFycmF5Om5vbmU7c3Ryb2tlLW9wYWNpdHk6MC42Mjg3NjMiCiAgICAgICAgIGlkPSJyZWN0NSIKICAgICAgICAgd2lkdGg9IjEwOC42MTQ4NCIKICAgICAgICAgaGVpZ2h0PSIzNi45MjEzMDciCiAgICAgICAgIHg9IjMzOS4yNzIxMyIKICAgICAgICAgeT0iNDAxLjMzNzc0IiAvPgogICAgPC9nPgogIDwvZz4KPC9zdmc+Cg==&link=https%3A%2F%2Fwww.trae.ai%2F) 

