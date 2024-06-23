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
	help := flag.Bool("h", false, "显示帮助信息")

	flag.Parse()

	if *help {
		fmt.Println("用法: ollama_plist [选项]")
		fmt.Println("给通过 Homebrew 安装的 Ollama CLI 工具添加环境变量")
		fmt.Println()
		fmt.Println("选项:")
		flag.PrintDefaults()
		return
	}

	// 确定 plist 文件路径
	plistPath := "/opt/homebrew/opt/ollama/homebrew.mxcl.ollama.plist"
	outputPath := plistPath

	// 检查 plist 文件是否存在
	if _, err := os.Stat(plistPath); os.IsNotExist(err) {
		fmt.Println("你没有通过Homebrew安装Ollama")
		return
	}

	// 读取 plist 文件
	file, err := os.Open(plistPath)
	if err != nil {
		fmt.Printf("打开 plist 文件失败: %v\n", err)
		return
	}
	defer file.Close()

	var plistData PlistDict
	decoder := plist.NewDecoder(file)
	err = decoder.Decode(&plistData)
	if err != nil {
		fmt.Printf("解析 plist 文件失败: %v\n", err)
		return
	}

	fmt.Println("plist 文件解析成功")

	// 创建 EnvironmentVariables 字典
	envVars := PlistDict{
		"OLLAMA_HOST":              *host,
		"OLLAMA_ORIGINS":           *origins,
		"OLLAMA_MAX_LOADED_MODELS": *maxLoadedModels,
		"OLLAMA_NUM_PARALLEL":      *numParallel,
	}

	// 检查是否应使用默认值
	if flag.NFlag() == 0 {
		envVars = PlistDict{
			"OLLAMA_HOST":              "0.0.0.0",
			"OLLAMA_ORIGINS":           "*",
			"OLLAMA_MAX_LOADED_MODELS": "2",
			"OLLAMA_NUM_PARALLEL":      "4",
		}
	}

	// 更新 plist 数据
	plistData["EnvironmentVariables"] = envVars

	// 将 plist 数据写回文件
	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("创建输出 plist 文件失败: %v\n", err)
		return
	}
	defer outputFile.Close()

	encoder := plist.NewEncoder(outputFile)
	encoder.Indent("\t")
	err = encoder.Encode(plistData)
	if err != nil {
		fmt.Printf("编码 plist 文件失败: %v\n", err)
		return
	}

	fmt.Println("plist 文件已修改并保存到", outputPath)
}
