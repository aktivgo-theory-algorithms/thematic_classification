package xlsx

import "github.com/tealeg/xlsx/v3"

type WorkSheet struct {
	File  *xlsx.File
	Sheet *xlsx.Sheet

	Path string
}

func NewWorkSheet(sheetName string, path string) (*WorkSheet, error) {
	workSheet := &WorkSheet{
		File: xlsx.NewFile(),
		Path: path,
	}

	if err := workSheet.createSheet(sheetName); err != nil {
		return nil, err
	}

	return workSheet, nil
}

func (wb *WorkSheet) createSheet(sheetName string) error {
	sheet, err := wb.File.AddSheet(sheetName)
	if err != nil {
		return err
	}

	wb.Sheet = sheet

	return nil
}

func (wb *WorkSheet) SetTitles(titles []string) error {
	for i, title := range titles {
		cell, err := wb.Sheet.Cell(0, i)
		if err != nil {
			return err
		}

		cell.SetValue(title)
	}

	return nil
}

func (wb *WorkSheet) Save() error {
	if err := wb.File.Save(wb.Path); err != nil {
		return err
	}

	return nil
}

func (wb *WorkSheet) CloseSheet() {
	wb.Sheet.Close()
}
