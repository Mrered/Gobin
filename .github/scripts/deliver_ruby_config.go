package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mrered/gobin/pkg/scripts"
)

func main() {
	// 获取cmd目录的绝对路径
	cmdDir := "./cmd"
	abs, err := filepath.Abs(cmdDir)
	if err != nil {
		fmt.Printf("获取绝对路径失败: %v\n", err)
		os.Exit(1)
	}

	// 创建一个切片来存储所有项目信息
	var projects []map[string]string

	// 遍历cmd目录
	err = filepath.Walk(abs, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查是否为main.go文件
		if !info.IsDir() && info.Name() == "main.go" {
			// 读取文件的帮助信息
			_, description, _, _, err := scripts.GetHelpTextFromMainGo(path)
			if err != nil {
				fmt.Printf("读取文件失败 %s: %v\n", path, err)
				return nil
			}

			// 获取项目名称（目录名）
			projectName := filepath.Base(filepath.Dir(path))

			// 创建项目信息并添加到切片中
			projectInfo := map[string]string{
				"project":     projectName,
				"description": description,
			}
			projects = append(projects, projectInfo)
		}
		return nil
	})

	// 将所有项目信息转换为JSON并输出
	jsonData, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		fmt.Printf("JSON 转换失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonData))

	if err != nil {
		fmt.Printf("遍历目录失败: %v\n", err)
		os.Exit(1)
	}
}
