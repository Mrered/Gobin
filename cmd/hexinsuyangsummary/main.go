/*
hexinsuyangsummary
linux darwin windows
核心素养汇总工具
用法: hexinsuyangsummary [选项]

选项:
  -h    显示帮助信息
  -p string
        输入目录路径
  -a    全量输出（database-style 详细记录模式）
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

	"github.com/xuri/excelize/v2"
)

// PersonData stores the scores for a specific person across all dimensions.
type PersonData struct {
	Dimensions map[string][]float64 // dimension name -> list of scores
}

// Record represents a single row in the database-style output
type Record struct {
	SourceFile string   // 源文件名
	SheetName  string   // sheet 名称
	CValue     string   // C 列的内容
	EValue     float64  // E 列的数值
}

func main() {
	// Command line flags
	pathFlag := flag.String("p", "", "输入目录路径")
	allFlag := flag.Bool("a", false, "全量输出（database-style 详细记录模式）")
	helpFlag := flag.Bool("h", false, "显示帮助信息")

	// Custom usage
	flag.Usage = printHelp

	flag.Parse()

	if len(os.Args) == 1 || *helpFlag {
		printHelp()
		return
	}

	if *pathFlag == "" {
		fmt.Println("错误: 必须提供 -p 参数指定路径")
		printHelp()
		os.Exit(1)
	}

	inputDir := *pathFlag

	// Collect files
	files, err := collectFiles(inputDir)
	if err != nil {
		log.Fatalf("Error collecting files: %v", err)
	}

	if len(files) == 0 {
		log.Println("No .xlsx files found in the specified directory.")
		return
	}

	log.Printf("Found %d .xlsx files. Processing...", len(files))

	if *allFlag {
		// Database-style detailed output
		var records []Record
		for _, file := range files {
			fileRecords := processFileForDetail(file)
			records = append(records, fileRecords...)
		}

		if len(records) > 0 {
			outputFile := filepath.Join(inputDir, "detail.xlsx")
			if err := writeDetail(outputFile, records); err != nil {
				log.Fatalf("Error writing detail file: %v", err)
			}
			log.Printf("Detail records written to %s", outputFile)
		} else {
			log.Println("No valid records found.")
		}
	} else {
		// Summary mode (average scores)
		data := make(map[string]*PersonData)
		for _, file := range files {
			processFileForSummary(file, data)
		}

		if len(data) > 0 {
			outputFile := filepath.Join(inputDir, "summary.xlsx")
			if err := writeSummary(outputFile, data); err != nil {
				log.Fatalf("Error writing summary file: %v", err)
			}
			log.Printf("Summary written to %s", outputFile)
		} else {
			log.Println("No valid data found.")
		}
	}
}

func printHelp() {
	fmt.Println("核心素养汇总工具")
	fmt.Println("用法: hexinsuyangsummary [选项]")
	fmt.Println()
	fmt.Println("选项:")
	flag.PrintDefaults()
}

// collectFiles finds all .xlsx files in the directory
func collectFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".xlsx") {
			if !strings.HasPrefix(info.Name(), "~$") { // Ignore temp files
				files = append(files, path)
			}
		}
		return nil
	})
	return files, err
}

// processFileForSummary processes a file for summary mode (averages)
func processFileForSummary(path string, data map[string]*PersonData) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		log.Printf("Warning: Could not open file %s: %v", path, err)
		return
	}
	defer f.Close()

	for _, sheetName := range f.GetSheetMap() {
		processSheetForSummary(f, sheetName, data)
	}
}

// processSheetForSummary processes a sheet for summary mode
func processSheetForSummary(f *excelize.File, personName string, data map[string]*PersonData) {
	rows, err := f.GetRows(personName)
	if err != nil {
		log.Printf("Warning: Could not get rows for sheet %s: %v", personName, err)
		return
	}

	mergedCells, err := f.GetMergeCells(personName)
	if err != nil {
		log.Printf("Warning: Could not get merged cells for sheet %s: %v", personName, err)
	}

	if _, ok := data[personName]; !ok {
		data[personName] = &PersonData{
			Dimensions: make(map[string][]float64),
		}
	}

	for i, row := range rows {
		rowNum := i + 1

		eVal, isTopLeftE := getCellValue(f, personName, "E", rowNum, mergedCells, row)
		if !isTopLeftE {
			continue
		}

		eVal = strings.TrimSpace(eVal)
		if eVal == "" {
			continue
		}

		score, err := strconv.ParseFloat(eVal, 64)
		if err != nil {
			continue
		}

		dimension := getCellStringValue(f, personName, "C", rowNum, mergedCells, row)
		dimension = strings.TrimSpace(dimension)
		if dimension == "" {
			continue
		}

		if strings.HasPrefix(dimension, "姓名：") || strings.HasPrefix(dimension, "姓名:") {
			continue
		}

		data[personName].Dimensions[dimension] = append(data[personName].Dimensions[dimension], score)
	}
}

// processFileForDetail processes a file for detail mode
func processFileForDetail(path string) []Record {
	f, err := excelize.OpenFile(path)
	if err != nil {
		log.Printf("Warning: Could not open file %s: %v", path, err)
		return nil
	}
	defer f.Close()

	var records []Record
	for _, sheetName := range f.GetSheetMap() {
		sheetRecords := processSheetForDetail(f, path, sheetName)
		records = append(records, sheetRecords...)
	}
	return records
}

// processSheetForDetail processes a sheet for detail mode
func processSheetForDetail(f *excelize.File, filePath, sheetName string) []Record {
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Printf("Warning: Could not get rows for sheet %s: %v", sheetName, err)
		return nil
	}

	mergedCells, err := f.GetMergeCells(sheetName)
	if err != nil {
		log.Printf("Warning: Could not get merged cells for sheet %s: %v", sheetName, err)
	}

	var records []Record
	for i, row := range rows {
		rowNum := i + 1

		cVal := getCellStringValue(f, sheetName, "C", rowNum, mergedCells, row)
		eVal, isTopLeftE := getCellValue(f, sheetName, "E", rowNum, mergedCells, row)

		if !isTopLeftE {
			continue
		}

		cVal = strings.TrimSpace(cVal)
		eVal = strings.TrimSpace(eVal)

		if cVal == "" || eVal == "" {
			continue
		}

		eValue, err := strconv.ParseFloat(eVal, 64)
		if err != nil {
			continue
		}

		record := Record{
			SourceFile: filepath.Base(filePath),
			SheetName:  sheetName,
			CValue:     cVal,
			EValue:     eValue,
		}
		records = append(records, record)
	}

	return records
}

// getCellValue returns the value of the cell, and a boolean indicating if this cell should be processed
func getCellValue(f *excelize.File, sheet, colName string, rowNum int, mergedCells []excelize.MergeCell, rowData []string) (string, bool) {
	cellRef := fmt.Sprintf("%s%d", colName, rowNum)

	for _, mc := range mergedCells {
		inRange, err := isCellInRange(cellRef, mc.GetStartAxis(), mc.GetEndAxis())
		if err == nil && inRange {
			if cellRef == mc.GetStartAxis() {
				val, _ := f.GetCellValue(sheet, mc.GetStartAxis())
				return val, true
			} else {
				return "", false
			}
		}
	}

	colIdx := -1
	switch colName {
	case "C":
		colIdx = 2
	case "E":
		colIdx = 4
	}

	if colIdx >= 0 && colIdx < len(rowData) {
		return rowData[colIdx], true
	}
	return "", true
}

// getCellStringValue gets the value for a cell, handling merges
func getCellStringValue(f *excelize.File, sheet, colName string, rowNum int, mergedCells []excelize.MergeCell, rowData []string) string {
	cellRef := fmt.Sprintf("%s%d", colName, rowNum)

	for _, mc := range mergedCells {
		inRange, err := isCellInRange(cellRef, mc.GetStartAxis(), mc.GetEndAxis())
		if err == nil && inRange {
			val, _ := f.GetCellValue(sheet, mc.GetStartAxis())
			return val
		}
	}

	colIdx := -1
	switch colName {
	case "C":
		colIdx = 2
	case "E":
		colIdx = 4
	}

	if colIdx >= 0 && colIdx < len(rowData) {
		return rowData[colIdx]
	}
	return ""
}

// isCellInRange checks if cell is within start:end range
func isCellInRange(cell, start, end string) (bool, error) {
	col, row, err := excelize.CellNameToCoordinates(cell)
	if err != nil {
		return false, err
	}

	c1, r1, err := excelize.CellNameToCoordinates(start)
	if err != nil {
		return false, err
	}

	c2, r2, err := excelize.CellNameToCoordinates(end)
	if err != nil {
		return false, err
	}

	return col >= c1 && col <= c2 && row >= r1 && row <= r2, nil
}

// writeSummary writes the summary output (averages)
func writeSummary(outputPath string, data map[string]*PersonData) error {
	f := excelize.NewFile()

	sheetName := "Summary"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}
	f.SetActiveSheet(index)
	if sheetName != "Sheet1" {
		f.DeleteSheet("Sheet1")
	}

	// Collect all dimension names
	dimensionSet := make(map[string]bool)
	for _, personData := range data {
		for dimension := range personData.Dimensions {
			dimensionSet[dimension] = true
		}
	}

	dimensions := make([]string, 0, len(dimensionSet))
	for dimension := range dimensionSet {
		dimensions = append(dimensions, dimension)
	}
	sort.Strings(dimensions)

	personNames := make([]string, 0, len(data))
	for personName := range data {
		personNames = append(personNames, personName)
	}
	sort.Strings(personNames)

	// Write headers
	f.SetCellValue(sheetName, "A1", "序号")
	f.SetCellValue(sheetName, "B1", "姓名")
	for i, dimension := range dimensions {
		cell, _ := excelize.CoordinatesToCellName(i+3, 1)
		f.SetCellValue(sheetName, cell, dimension)
	}

	// Write data
	for idx, personName := range personNames {
		rowNum := idx + 2
		personData := data[personName]

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), idx+1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), personName)

		for i, dimension := range dimensions {
			cell, _ := excelize.CoordinatesToCellName(i+3, rowNum)
			if scores, ok := personData.Dimensions[dimension]; ok && len(scores) > 0 {
				sum := 0.0
				for _, score := range scores {
					sum += score
				}
				avg := sum / float64(len(scores))
				avgRounded, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", avg), 64)
				f.SetCellValue(sheetName, cell, avgRounded)
			}
		}
	}

	return f.SaveAs(outputPath)
}

// writeDetail writes the detail output (database-style)
func writeDetail(outputPath string, records []Record) error {
	f := excelize.NewFile()

	sheetName := "Detail"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}
	f.SetActiveSheet(index)
	if sheetName != "Sheet1" {
		f.DeleteSheet("Sheet1")
	}

	// Write headers
	headers := []string{"SourceFile", "SheetName", "CValue", "EValue"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Write data
	for idx, record := range records {
		rowNum := idx + 2

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), record.SourceFile)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), record.SheetName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), record.CValue)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), record.EValue)
	}

	return f.SaveAs(outputPath)
}
