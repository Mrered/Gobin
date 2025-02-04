# Gobin

Go 二进制小程序

![Homebrew](https://img.shields.io/badge/-Homebrew-FBB040?labelColor=555555&logoColor=FFFFFF&logo=homebrew) ![CI](https://github.com/Mrered/Gobin/actions/workflows/CI.yml/badge.svg) ![license](https://img.shields.io/github/license/Mrered/Gobin) ![code-size](https://img.shields.io/github/languages/code-size/Mrered/Gobin) ![repo-size](https://img.shields.io/github/repo-size/Mrered/Gobin)

## 🍺 安装

```sh
brew tap brewforge/chinese
brew install <二进制命令行工具名> --formula
```

## 📋 列表

|                     二进制命令行工具名                     |                        说明                        |
| :--------------------------------------------------------: | :------------------------------------------------: |
| [makemf](https://github.com/Mrered/Gobin#makemf) | 为 GGUF 文件生成 Makefile |
| [ollamaplist](https://github.com/Mrered/Gobin#ollamaplist) | 给通过 Homebrew 安装的 Ollama CLI 工具添加环境变量 |
| [reportgen](https://github.com/Mrered/Gobin#reportgen) | 生成报告 |

## 🚀 使用

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

## 🏆 致谢

[Homebrew](https://brew.sh) [ChatGPT](https://chatgpt.com) [Claude](https://claude.ai)

## 📄 许可

[MIT](https://github.com/Mrered/Gobin/blob/main/LICENSE) © [Mrered](https://github.com/Mrered)
