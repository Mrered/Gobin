/*
reportgen
linux darwin windows
生成报告
用法: reportgen [选项]

选项:
  -d string
        指定工作目录
  -h    显示帮助信息
  -t string
        指定报告类型 (w: 周报, m: 月报, s: 学期报, y: 年报)
  -w string
        指定周数
  -m string
        指定月份 (格式: YYYYMM)
  -s string
        指定学期 (格式: YYYY - YYYY 春/秋)
  -y string
        指定年份 (格式: YYYY)
  -v    显示版本号
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mrered/gobin/pkg/reportgen"
)

// 版本号，默认值 "dev"，在编译时通过 -ldflags 动态设置
var version = "dev"

func main() {
	// 定义命令行参数
	dirPath := flag.String("d", "", "指定工作目录")
	reportType := flag.String("t", "", "指定报告类型 (w: 周报, m: 月报, s: 学期报, y: 年报)")
	formatting := flag.Bool("f", false, "是否格式化内容")
	week := flag.String("w", "", "指定周数")
	month := flag.String("m", "", "指定月份 (格式: YYYYMM)")
	semester := flag.String("s", "", "指定学期 (格式: YYYY - YYYY 春/秋)")
	year := flag.String("y", "", "指定年份 (格式: YYYY)")
	help := flag.Bool("h", false, "显示帮助信息")
	showVersion := flag.Bool("v", false, "显示版本号")

	flag.Parse()

	// 显示版本号
	if *showVersion {
		fmt.Println("reportgen 版本:", version)
		return
	}

	// 显示帮助信息
	if len(os.Args) == 1 || *help {
		printHelp()
		return
	}

	// 检查工作目录
	if *dirPath == "" {
		log.Fatal("错误：必须指定工作目录 (-d)")
	}

	// 验证工作目录结构
	if err := reportgen.ValidateWorkingDir(*dirPath); err != nil {
		log.Fatal(err)
	}

	// 如果未指定报告类型，提供选择
	if *reportType == "" {
		selected, err := selectReportType()
		if err != nil {
			log.Fatal(err)
		}
		*reportType = selected
	}

	// 创建配置
	config := &reportgen.Config{
		ReportType: *reportType,
		Formatting: *formatting,
	}

	// 创建生成器
	generator, err := reportgen.NewGenerator(config)
	if err != nil {
		log.Fatal(err)
	}

	// 根据报告类型设置源目录和目标目录
	switch *reportType {
	case "w":
		config.SourceDir = filepath.Join(*dirPath, "日报")
		config.TargetDir = filepath.Join(*dirPath, "周报")
		if *week == "" {
			selected, err := selectPeriod(config)
			if err != nil {
				log.Fatal(err)
			}
			// 对每个选中的时间段生成报告
			for _, period := range selected {
				config.SelectedPeriod = period
				if err := generator.Generate(config.SourceDir, nil); err != nil {
					log.Fatal(err)
				}
			}
			return
		}
		config.SelectedPeriod = *week

	case "m":
		config.SourceDir = filepath.Join(*dirPath, "周报")
		config.TargetDir = filepath.Join(*dirPath, "月报")
		if *month == "" {
			selected, err := selectPeriod(config)
			if err != nil {
				log.Fatal(err)
			}
			// 对每个选中的时间段生成报告
			for _, period := range selected {
				config.SelectedPeriod = period
				if err := generator.Generate(config.SourceDir, nil); err != nil {
					log.Fatal(err)
				}
			}
			return
		}
		config.SelectedPeriod = *month

	case "s":
		config.SourceDir = filepath.Join(*dirPath, "月报")
		config.TargetDir = filepath.Join(*dirPath, "学期报")
		if *semester == "" {
			selected, err := selectPeriod(config)
			if err != nil {
				log.Fatal(err)
			}
			// 对每个选中的时间段生成报告
			for _, period := range selected {
				config.SelectedPeriod = period
				if err := generator.Generate(config.SourceDir, nil); err != nil {
					log.Fatal(err)
				}
			}
			return
		}
		config.SelectedPeriod = *semester

	case "y":
		config.SourceDir = filepath.Join(*dirPath, "学期报")
		config.TargetDir = filepath.Join(*dirPath, "年报")
		if *year == "" {
			selected, err := selectPeriod(config)
			if err != nil {
				log.Fatal(err)
			}
			// 对每个选中的时间段生成报告
			for _, period := range selected {
				config.SelectedPeriod = period
				if err := generator.Generate(config.SourceDir, nil); err != nil {
					log.Fatal(err)
				}
			}
			return
		}
		config.SelectedPeriod = *year

	default:
		log.Fatal("不支持的报告类型：", *reportType)
	}

	// 创建生成器
	generator, err = reportgen.NewGenerator(config)
	if err != nil {
		log.Fatal(err)
	}

	// 生成报告
	if err := generator.Generate(config.SourceDir, nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println("报告生成完成")
}

func selectReportType() (string, error) {
	reportTypes := []string{
		"归纳周报 (w)",
		"归纳月报 (m)",
		"归纳学期报 (s)",
		"归纳年报 (y)",
	}

	var selected string
	prompt := &survey.Select{
		Message: "请选择要生成的报告类型：",
		Options: reportTypes,
	}

	err := survey.AskOne(prompt, &selected)
	if err != nil {
		return "", fmt.Errorf("选择报告类型失败：%v", err)
	}

	return strings.ToLower(string(selected[len(selected)-2])), nil
}

func selectPeriod(config *reportgen.Config) ([]string, error) {
	// 创建生成器
	generator, err := reportgen.NewGenerator(config)
	if err != nil {
		return nil, err
	}

	// 获取可用的时间段
	periods, err := generator.GetAvailablePeriods(config.SourceDir)
	if err != nil {
		return nil, err
	}

	if len(periods) == 0 {
		return nil, fmt.Errorf("未找到可用的时间段")
	}

	// 对时间段进行排序
	sort.Slice(periods, func(i, j int) bool {
		// 提取数字部分
		ni, _ := strconv.Atoi(strings.TrimLeft(periods[i], "第"))
		nj, _ := strconv.Atoi(strings.TrimLeft(periods[j], "第"))
		return ni < nj
	})

	var selected []string
	prompt := &survey.MultiSelect{
		Message: "请选择时间段（空格键选择，回车键确认）：",
		Options: periods,
	}

	err = survey.AskOne(prompt, &selected)
	if err != nil {
		return nil, fmt.Errorf("选择时间段失败：%v", err)
	}

	if len(selected) == 0 {
		return nil, fmt.Errorf("未选择任何时间段")
	}

	return selected, nil
}

func printHelp() {
	fmt.Println("生成报告")
	fmt.Println("用法: reportgen [选项]")
	fmt.Println()
	fmt.Println("选项:")
	flag.PrintDefaults()
}
