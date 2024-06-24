/*
makemf
linux darwin windows
为 GGUF 文件生成 Makefile
用法: makemf [选项]

选项:
  -a    自动为当前目录下的所有 .gguf 文件生成 Makefile
  -h    显示帮助信息
  -m string
        GGUF 文件名称，包含后缀名
  -n string
        要生成的 Makefile 名称
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	// 定义命令行标志
	name := flag.String("n", "", "要生成的 Makefile 名称")
	modelFile := flag.String("m", "", "GGUF 文件名称，包含后缀名")
	autoMode := flag.Bool("a", false, "自动为当前目录下的所有 .gguf 文件生成 Makefile")
	help := flag.Bool("h", false, "显示帮助信息")

	flag.Parse()

	if len(os.Args) == 1 || *help {
		printHelp()
		return
	}

	// 获取当前目录下的所有 .gguf 文件
	ggufFiles, err := getGGUFFiles()
	if err != nil {
		log.Fatalf("读取目录失败: %v", err)
	}

	// 处理自动模式
	if *autoMode {
		handleAutoMode(ggufFiles)
		return
	}

	if *name == "" || *modelFile == "" {
		log.Fatalf("必须同时提供 -n 和 -m 参数，或者使用 -a 参数")
	}

	// 处理手动模式
	handleManualMode(*name, *modelFile, ggufFiles)
}

// printResult 打印生成结果
func printResult(outputFile string) {
	name := strings.TrimSuffix(outputFile, ".mf")
	fmt.Printf("已生成：%s\n", outputFile)
	fmt.Println("")
	fmt.Println("生成模型文件：")
	fmt.Println("")
	fmt.Printf(" ollama create %s -f ./%s\n", name, outputFile)
	fmt.Println("")
	fmt.Println("运行大模型：")
	fmt.Println("")
	fmt.Printf(" ollama run %s\n", name)
	fmt.Println("")
}

func getGGUFFiles() ([]string, error) {
	var ggufFiles []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".gguf" {
			ggufFiles = append(ggufFiles, path)
		}
		return nil
	})
	return ggufFiles, err
}

func handleAutoMode(ggufFiles []string) {
	for _, file := range ggufFiles {
		outputFile := strings.TrimSuffix(file, ".gguf") + ".mf"
		// 生成 Makefile
		if err := generateMakefile(strings.TrimSuffix(file, ".gguf"), file); err != nil {
			log.Fatalf("写入文件失败: %v", err)
		}
		printResult(outputFile)
	}
}

func handleManualMode(name, modelFile string, ggufFiles []string) {
	// 生成 Makefile
	if err := generateMakefile(name, modelFile); err != nil {
		log.Fatalf("写入文件失败: %v", err)
	}
	printResult(name + ".mf")
	fmt.Println("当前目录下的 .gguf 文件：")
	for i, file := range ggufFiles {
		fmt.Printf("%d. %s\n", i+1, file)
	}
	fmt.Println("请输入文件序号（多个序号用空格隔开，或使用 1-3 表示连续序号）：")
	var userInput string
	fmt.Scanln(&userInput)
	// 解析用户输入的序号
	selectedFiles := parseUserInput(userInput, ggufFiles)
	for _, file := range selectedFiles {
		outputFile := strings.TrimSuffix(file, ".gguf") + ".mf"
		// 生成 Makefile
		if err := generateMakefile(strings.TrimSuffix(file, ".gguf"), file); err != nil {
			log.Fatalf("写入文件失败: %v", err)
		}
		printResult(outputFile)
	}
}

func parseUserInput(input string, ggufFiles []string) []string {
	var selectedFiles []string
	inputs := strings.Fields(input)
	for _, input := range inputs {
		if strings.Contains(input, "-") {
			rangeParts := strings.Split(input, "-")
			// 将字符串解析为整数索引
			startIndex := parseIndex(rangeParts[0])
			endIndex := parseIndex(rangeParts[1])
			for i := startIndex; i <= endIndex; i++ {
				selectedFiles = append(selectedFiles, ggufFiles[i-1])
			}
		} else {
			index := parseIndex(input)
			selectedFiles = append(selectedFiles, ggufFiles[index-1])
		}
	}
	return selectedFiles
}

func generateMakefile(name, modelFile string) error {
	outputFile := name + ".mf"
	content := fmt.Sprintf("FROM ./%s", modelFile)
	return os.WriteFile(outputFile, []byte(content), 0644)
}

func parseIndex(input string) int {
	index, err := strconv.Atoi(input)
	if err != nil {
		log.Fatalf("无效的输入: %v", input)
	}
	return index
}

func printHelp() {
	fmt.Println("为 GGUF 文件生成 Makefile")
	fmt.Println("用法: makemf [选项]")
	fmt.Println()
	fmt.Println("选项:")
	flag.PrintDefaults()
}
