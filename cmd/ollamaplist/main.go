/*
ollamaplist
darwin
给通过 Homebrew 安装的 Ollama CLI 工具添加环境变量
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
*/

package main

import (
	"flag"
	"fmt"
	"os"

	"howett.net/plist"
)

type PlistDict map[string]interface{}

func main() {
	// 定义命令行标志
	host := flag.String("s", "0.0.0.0", "OLLAMA_HOST")
	origins := flag.String("o", "*", "OLLAMA_ORIGINS")
	maxLoadedModels := flag.String("m", "2", "OLLAMA_MAX_LOADED_MODELS")
	numParallel := flag.String("p", "4", "OLLAMA_NUM_PARALLEL")
	removeEnv := flag.Bool("r", false, "删除所有环境变量")
	help := flag.Bool("h", false, "显示帮助信息")
	applyDefault := flag.Bool("a", false, "应用默认配置")

	flag.Parse()

	if len(os.Args) == 1 || *help {
		printHelp()
		return
	}

	// 确定 plist 文件路径
	plistPath := "/opt/homebrew/opt/ollama/homebrew.mxcl.ollama.plist"
	outputPath := plistPath

	// 检查 plist 文件是否存在
	if _, err := os.Stat(plistPath); os.IsNotExist(err) {
		fmt.Println("你没有通过 Homebrew 安装 Ollama CLI ，运行下面的命令来安装：")
		fmt.Println(" brew install ollama --formula")
		return
	}

	// 读取 plist 文件
	plistData, err := readPlistFile(plistPath)
	if err != nil {
		fmt.Printf("读取 plist 文件失败: %v\n", err)
		return
	}

	// 删除 EnvironmentVariables 字典
	delete(plistData, "EnvironmentVariables")

	if !*removeEnv {
		// 创建 EnvironmentVariables 字典
		envVars := PlistDict{
			"OLLAMA_HOST":              *host,
			"OLLAMA_ORIGINS":           *origins,
			"OLLAMA_MAX_LOADED_MODELS": *maxLoadedModels,
			"OLLAMA_NUM_PARALLEL":      *numParallel,
		}

		// 检查是否应使用默认值
		if *applyDefault || (flag.NFlag() == 1 && *applyDefault) {
			envVars = PlistDict{
				"OLLAMA_HOST":              "0.0.0.0",
				"OLLAMA_ORIGINS":           "*",
				"OLLAMA_MAX_LOADED_MODELS": "2",
				"OLLAMA_NUM_PARALLEL":      "4",
			}
		}

		// 更新 plist 数据
		plistData["EnvironmentVariables"] = envVars
	}

	// 将 plist 数据写回文件
	err = writePlistFile(outputPath, plistData)
	if err != nil {
		fmt.Printf("写入 plist 文件失败: %v\n", err)
		return
	}

	fmt.Println("Ollama CLI 环境变量已成功更新到 plist 文件，通过以下命令启动服务：")
	fmt.Println(" brew services start ollama")
}

func printHelp() {
	fmt.Println("给通过 Homebrew 安装的 Ollama CLI 工具添加环境变量")
	fmt.Println("用法: ollamaplist [选项]")
	fmt.Println()
	fmt.Println("选项:")
	flag.PrintDefaults()
}

func readPlistFile(path string) (PlistDict, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var plistData PlistDict
	decoder := plist.NewDecoder(file)
	err = decoder.Decode(&plistData)
	if err != nil {
		return nil, err
	}

	return plistData, nil
}

func writePlistFile(path string, data PlistDict) error {
	outputFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	encoder := plist.NewEncoder(outputFile)
	encoder.Indent("\t")
	return encoder.Encode(data)
}
