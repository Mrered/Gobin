# Gobin

Go 二进制小程序

## 🍺 安装

```sh
brew tap brewforge/chinese
brew install <二进制命令行工具名> --formula
```

## 📋 列表

|                     二进制命令行工具名                     |                        说明                        |
| :--------------------------------------------------------: | :------------------------------------------------: |
| [ollamaplist](https://github.com/Mrered/Gobin#ollamaplist) | 给通过 Homebrew 安装的 Ollama CLI 工具添加环境变量 |
| [makemf](https://github.com/Mrered/Gobin#makemf) | 为 GGUF 文件生成 Makefile |

## 🚀 使用

### makemf

```sh
为 GGUF 文件生成 Makefile
用法: makemf [选项]

选项:
  -a    自动为当前目录下的所有 .gguf 文件生成 Makefile
  -h    显示帮助信息
  -m string
        GGUF 文件名称，包含后缀名
  -n string
        要生成的 Makefile 名称
```

### ollamaplist

```sh
给通过 Homebrew 安装的 Ollama CLI 工具添加环境变量
用法: ollamaplist [选项]

选项:
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

[Homebrew](https://brew.sh) [ChatGPT](https://chatgpt.com)

## 📄 许可

[MIT](https://github.com/Mrered/Gobin/blob/main/LICENSE) © [Mrered](https://github.com/Mrered)
