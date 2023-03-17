package excel

import (
	"fmt"
	"github.com/extrame/xls"
	"github.com/xuri/excelize/v2"
)

type ExcelCloStyle struct {
	Width float64
}

func GetSheetData(fileName, sheetName string) ([][]string, error) {
	file, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, err
	}
	rows, err := file.GetRows(sheetName)

	if err != nil {
		return nil, err
	}
	return rows, nil
}

func XlsGetSheetData(fileName string, sheet int) ([][]string, error) {
	file, err := xls.Open(fileName, "utf-8")
	if err != nil {
		return nil, err
	}
	sheetData := file.GetSheet(sheet)
	maxRow := int(sheetData.MaxRow)
	rows := make([][]string, 0)
	for i := 0; i <= maxRow; i++ {
		xlsRow := sheetData.Row(i)
		row := make([]string, 0)
		for j := 0; j <= xlsRow.LastCol(); j++ {
			row = append(row, xlsRow.Col(j))
		}
		rows = append(rows, row)
	}

	if err != nil {
		return nil, err
	}
	return rows, nil
}

func WriteSheet(fileName, sheetName string, rows [][]string, closStyles ...ExcelCloStyle) error {
	file := excelize.NewFile()
	// 创建一个名字为stock的sheet页
	file.NewSheet(sheetName)
	headerStyleId, _ := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	defaultStyleId, _ := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: false, Size: 12},
		Alignment: &excelize.Alignment{Horizontal: "center", WrapText: true, Vertical: "center"},
	})
	closStyleMap := map[int]ExcelCloStyle{}
	if len(closStyles) > 0 {
		for k, v := range closStyles {
			closStyleMap[k] = v
		}
	}
	for i, row := range rows {
		for j, v := range row {
			index := getExcelIndex(j)
			cloIndex := fmt.Sprintf("%s%d", index, i+1)
			_ = file.SetCellStr(sheetName, cloIndex, v)
			if i == 0 {
				_ = file.SetCellStyle(sheetName, cloIndex, cloIndex, headerStyleId)
				if s, ok := closStyleMap[j]; ok {
					if s.Width > 0 {
						_ = file.SetColWidth(sheetName, index, index, s.Width)
					}
				}
			} else {
				_ = file.SetCellStyle(sheetName, cloIndex, cloIndex, defaultStyleId)
			}
		}
	}
	defer file.Close()
	err := file.SaveAs(fileName)
	return err
}

func GetSheetFileFlow(sheetName string, rows [][]string, closStyles ...ExcelCloStyle) *excelize.File {
	file := excelize.NewFile()
	// 创建一个名字为stock的sheet页
	file.NewSheet(sheetName)
	headerStyleId, _ := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	defaultStyleId, _ := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: false, Size: 12},
		Alignment: &excelize.Alignment{Horizontal: "center", WrapText: true, Vertical: "center"},
	})
	closStyleMap := map[int]ExcelCloStyle{}
	if len(closStyles) > 0 {
		for k, v := range closStyles {
			closStyleMap[k] = v
		}
	}
	for i, row := range rows {
		for j, v := range row {
			index := getExcelIndex(j)
			cloIndex := fmt.Sprintf("%s%d", index, i+1)
			_ = file.SetCellStr(sheetName, cloIndex, v)
			if i == 0 {
				_ = file.SetCellStyle(sheetName, cloIndex, cloIndex, headerStyleId)
				if s, ok := closStyleMap[j]; ok {
					if s.Width > 0 {
						_ = file.SetColWidth(sheetName, index, index, s.Width)
					}
				} else {
					_ = file.SetColWidth(sheetName, index, index, 15)
				}
			} else {
				_ = file.SetCellStyle(sheetName, cloIndex, cloIndex, defaultStyleId)
			}
		}
	}
	return file
}

func getExcelIndex(index int) string {
	if index > 26*26-1 {
		return ""
	}
	indexStr := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	modIndex := index % 26
	num := index / 26
	if num == 0 {
		return string(indexStr[modIndex])
	}
	indexBytes := make([]byte, 0)
	indexBytes = append(indexBytes, indexStr[num-1])
	indexBytes = append(indexBytes, indexStr[modIndex])
	return string(indexBytes)
}
