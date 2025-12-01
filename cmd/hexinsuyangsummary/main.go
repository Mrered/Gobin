/*
hexinsuyangsummary
linux darwin windows
核心素养汇总工具
用法: hexinsuyangsummary [选项]

选项:
  -h    显示帮助信息
  -p string
        输入目录路径
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

func main() {
	pathFlag := flag.String("p", "", "输入目录路径")
	helpFlag := flag.Bool("h", false, "显示帮助信息")

	// Custom usage to avoid default flag output when calling Usage()
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
	// Output file is always summary.xlsx in the input directory
	outputFile := filepath.Join(inputDir, "summary.xlsx")

	// Resolve output file absolute path to avoid processing it
	absOutput, err := filepath.Abs(outputFile)
	if err != nil {
		log.Printf("Warning: Could not resolve absolute path for output file: %v", err)
		absOutput = outputFile // Fallback
	}

	files, err := collectFiles(inputDir, absOutput)
	if err != nil {
		log.Fatalf("Error collecting files: %v", err)
	}

	if len(files) == 0 {
		log.Println("No .xlsx files found in the specified directory.")
		return
	}

	log.Printf("Found %d .xlsx files. Processing...", len(files))

	// Map person name -> PersonData
	data := make(map[string]*PersonData)

	for _, file := range files {
		processFile(file, data)
	}

	if err := writeSummary(outputFile, data); err != nil {
		log.Fatalf("Error writing summary: %v", err)
	}

	log.Printf("Summary successfully written to %s", outputFile)
}

func printHelp() {
	fmt.Println("核心素养汇总工具")
	fmt.Println("用法: hexinsuyangsummary [选项]")
	fmt.Println()
	fmt.Println("选项:")
	flag.PrintDefaults()
}

// collectFiles recursively finds all .xlsx files in the directory.
func collectFiles(dir, absOutput string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".xlsx") {
			// Get absolute path of current file
			absPath, err := filepath.Abs(path)
			if err == nil && absPath == absOutput {
				return nil // Skip output file
			}

			if !strings.HasPrefix(info.Name(), "~$") { // Ignore temp files
				files = append(files, path)
			}
		}
		return nil
	})
	return files, err
}

func processFile(path string, data map[string]*PersonData) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		log.Printf("Warning: Could not open file %s: %v", path, err)
		return
	}
	defer f.Close()

	// Each sheet name is a person's name
	for _, sheetName := range f.GetSheetMap() {
		processSheet(f, sheetName, data)
	}
}

func processSheet(f *excelize.File, personName string, data map[string]*PersonData) {
	rows, err := f.GetRows(personName)
	if err != nil {
		log.Printf("Warning: Could not get rows for sheet %s: %v", personName, err)
		return
	}

	// Get merged cells to handle them correctly
	mergedCells, err := f.GetMergeCells(personName)
	if err != nil {
		log.Printf("Warning: Could not get merged cells for sheet %s: %v", personName, err)
		// Continue without merge info
	}

	// Initialize PersonData if not exists
	if _, ok := data[personName]; !ok {
		data[personName] = &PersonData{
			Dimensions: make(map[string][]float64),
		}
	}

	// Iterate rows. 1-based index for Excel, but rows slice is 0-based.
	// row index i corresponds to Excel row i+1.
	for i, row := range rows {
		rowNum := i + 1

		// Column C is index 2 (dimension name)
		// Column E is index 4 (score)

		// Check if this row is part of a merged cell for Col E and if it's the top-left cell.
		// If it's part of a merged cell but NOT the top-left, we skip it (to avoid double counting).
		eVal, isTopLeftE := getCellValue(f, personName, "E", rowNum, mergedCells, row)
		if !isTopLeftE {
			continue // Skip this row for E-column purposes (it's covered by a previous row's merge)
		}

		// If E value is empty/invalid, we skip
		eVal = strings.TrimSpace(eVal)
		if eVal == "" {
			continue
		}

		score, err := strconv.ParseFloat(eVal, 64)
		if err != nil {
			// Not a number, maybe header or garbage. Skip.
			continue
		}

		// Get corresponding C value (dimension name)
		dimension := getCellStringValue(f, personName, "C", rowNum, mergedCells, row)
		dimension = strings.TrimSpace(dimension)
		if dimension == "" {
			continue // Skip if no dimension name
		}

		// Filter out invalid dimension names (e.g., "姓名：XXX")
		if strings.HasPrefix(dimension, "姓名：") || strings.HasPrefix(dimension, "姓名:") {
			continue // Skip this, it's not a valid dimension
		}

		// Store the score for this person's dimension
		data[personName].Dimensions[dimension] = append(data[personName].Dimensions[dimension], score)
	}
}

// getCellValue returns the value of the cell, and a boolean indicating if this cell should be processed.
// It returns true if the cell is not merged OR if it is the top-left of a merged range.
// It returns false if the cell is part of a merged range but NOT the top-left.
// For the value: if it's a merged cell, it returns the value of the top-left cell.
func getCellValue(f *excelize.File, sheet, colName string, rowNum int, mergedCells []excelize.MergeCell, rowData []string) (string, bool) {
	cellRef := fmt.Sprintf("%s%d", colName, rowNum)

	for _, mc := range mergedCells {
		inRange, err := isCellInRange(cellRef, mc.GetStartAxis(), mc.GetEndAxis())
		if err == nil && inRange {
			// It is in a merged range.
			// Check if it is the top-left.
			if cellRef == mc.GetStartAxis() {
				// It is the top-left. Return value.
				val, _ := f.GetCellValue(sheet, mc.GetStartAxis())
				return val, true
			} else {
				// It is in a merged range but NOT top-left. Skip.
				return "", false
			}
		}
	}

	// Not in any merged range. Return value from rowData if available.
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
	return "", true // Empty, but "valid" in sense of not being skipped due to merge
}

// getCellStringValue gets the value for a cell, handling merges (always returns the value of the merge block if inside one).
// This is for Column C, where we just want the label associated with the current row (or the merge block start).
func getCellStringValue(f *excelize.File, sheet, colName string, rowNum int, mergedCells []excelize.MergeCell, rowData []string) string {
	cellRef := fmt.Sprintf("%s%d", colName, rowNum)

	for _, mc := range mergedCells {
		inRange, err := isCellInRange(cellRef, mc.GetStartAxis(), mc.GetEndAxis())
		if err == nil && inRange {
			// It is in a merged range. Return the value of the top-left.
			val, _ := f.GetCellValue(sheet, mc.GetStartAxis())
			return val
		}
	}

	// Not merged.
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

// isCellInRange checks if cell is within start:end range.
// Simplified check.
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

func writeSummary(outputPath string, data map[string]*PersonData) error {
	f := excelize.NewFile()

	// Create a new sheet or use default "Sheet1"
	sheetName := "Summary"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}
	f.SetActiveSheet(index)
	// Delete default Sheet1 if we created a new one
	if sheetName != "Sheet1" {
		f.DeleteSheet("Sheet1")
	}

	// Collect all dimension names (union of all dimensions across all persons)
	dimensionSet := make(map[string]bool)
	for _, personData := range data {
		for dimension := range personData.Dimensions {
			dimensionSet[dimension] = true
		}
	}

	// Convert to sorted slice for consistent column order
	dimensions := make([]string, 0, len(dimensionSet))
	for dimension := range dimensionSet {
		dimensions = append(dimensions, dimension)
	}
	// Sort dimensions alphabetically (optional, but makes output more predictable)
	// You can customize the sort order if needed
	sort.Strings(dimensions)

	// Collect all person names
	personNames := make([]string, 0, len(data))
	for personName := range data {
		personNames = append(personNames, personName)
	}
	// Sort person names alphabetically (optional)
	sort.Strings(personNames)

	// Write headers: 序号 | 姓名 | dimension1 | dimension2 | ...
	f.SetCellValue(sheetName, "A1", "序号")
	f.SetCellValue(sheetName, "B1", "姓名")
	for i, dimension := range dimensions {
		cell, _ := excelize.CoordinatesToCellName(i+3, 1) // Start from column C (index 3)
		f.SetCellValue(sheetName, cell, dimension)
	}

	// Write data rows
	for idx, personName := range personNames {
		rowNum := idx + 2 // Start from row 2 (row 1 is header)
		personData := data[personName]

		// Column A: 序号 (1-based index)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), idx+1)

		// Column B: 姓名
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), personName)

		// Columns C onwards: dimension scores (average)
		for i, dimension := range dimensions {
			cell, _ := excelize.CoordinatesToCellName(i+3, rowNum)

			if scores, ok := personData.Dimensions[dimension]; ok && len(scores) > 0 {
				// Calculate average
				sum := 0.0
				for _, score := range scores {
					sum += score
				}
				avg := sum / float64(len(scores))

				// Round to 2 decimal places
				avgRounded, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", avg), 64)
				f.SetCellValue(sheetName, cell, avgRounded)
			}
			// If no scores for this dimension, leave cell empty
		}
	}

	return f.SaveAs(outputPath)
}
