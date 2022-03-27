package app

import (
	"PostsLengthGetter/internal/model"
	"PostsLengthGetter/internal/reader"
	"encoding/csv"
	"github.com/tealeg/xlsx/v3"
	"log"
	"strconv"
	"time"
)

var (
	postsFilePath = "../Dataset/posts.csv"
	xlsxFilePath  = "data/posts_length.xlsx"
)

func Run() error {
	// Creating csv reader
	file, r, err := reader.NewCsvReader(postsFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Creating workbook
	wb, sh, err := InitWorkBook()
	if err != nil {
		return err
	}

	// Prepare table
	if err = setTableTitles(sh); err != nil {
		return err
	}

	// Process records from csv file
	totalProcessedRecords, err := processRecords(r, wb, sh)
	if err != nil {
		return err
	}

	log.Printf("Program finished, total records: %d", totalProcessedRecords)

	return nil
}

func InitWorkBook() (*xlsx.File, *xlsx.Sheet, error) {
	// Creating new xlsx file
	wb := xlsx.NewFile()

	// Add new sheet
	sh, err := wb.AddSheet("posts_length")
	if err != nil {
		return nil, nil, err
	}

	return wb, sh, nil
}

func setTableTitles(sh *xlsx.Sheet) error {
	cell, err := sh.Cell(0, 0)
	if err != nil {
		return err
	}
	cell.SetValue("id")

	cell, err = sh.Cell(0, 1)
	if err != nil {
		return err
	}
	cell.SetValue("post_id")

	cell, err = sh.Cell(0, 2)
	if err != nil {
		return err
	}
	cell.SetValue("title")

	cell, err = sh.Cell(0, 3)
	if err != nil {
		return err
	}
	cell.SetValue("text")

	cell, err = sh.Cell(0, 4)
	if err != nil {
		return err
	}
	cell.SetValue("text_length")

	return nil
}

func processRecords(reader *csv.Reader, wb *xlsx.File, sh *xlsx.Sheet) (int, error) {
	totalProcessedRecords := 0

	// Skip first record
	_, err := reader.Read()
	if err != nil {
		return totalProcessedRecords, nil
	}

	row := 1
	for {
		// Read next record
		record, err := reader.Read()
		if err != nil {
			// If end save file and return
			if err = wb.Save(xlsxFilePath); err != nil {
				return totalProcessedRecords, err
			}
			sh.Close()
			return totalProcessedRecords, nil
		}

		id, err := strconv.Atoi(record[0])
		postId, err := strconv.Atoi(record[1])
		title := record[2]
		text := record[3]
		textLength := len(record[3])

		// Skip empty text posts
		if textLength <= 2 {
			continue
		}

		myRecord := model.NewRecord(id, postId, title, text, textLength)

		// Fill row
		if err = fillRow(sh, row, myRecord); err != nil {
			return totalProcessedRecords, err
		}

		// Save file every 29-31 Seconds
		if time.Now().Second() >= 29 && time.Now().Second() <= 31 {
			if err = wb.Save(xlsxFilePath); err != nil {
				return totalProcessedRecords, err
			}
		}

		log.Printf("record (id = %d, post_id = %d, text_length = %d) successfully recorded to row %d",
			myRecord.ID, myRecord.PostID, myRecord.TextLength, row)

		totalProcessedRecords++
		row++
	}
}

func fillRow(sh *xlsx.Sheet, row int, myRecord *model.Record) error {
	cell, err := sh.Cell(row, 0)
	if err != nil {
		return err
	}
	cell.SetInt(myRecord.ID)

	cell, err = sh.Cell(row, 1)
	if err != nil {
		return err
	}
	cell.SetInt(myRecord.PostID)

	cell, err = sh.Cell(row, 2)
	if err != nil {
		return err
	}
	cell.SetString(myRecord.Title)

	cell, err = sh.Cell(row, 3)
	if err != nil {
		return err
	}
	cell.SetString(myRecord.Text)

	cell, err = sh.Cell(row, 4)
	if err != nil {
		return err
	}
	cell.SetInt(myRecord.TextLength)

	return nil
}
