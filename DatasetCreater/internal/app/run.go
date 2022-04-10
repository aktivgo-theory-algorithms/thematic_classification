package app

import (
	"PostsLengthGetter/internal/excel"
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
	sheetName     = "posts_length"
	xlsxFilePath  = "data/dataset.xlsx"
	titles        = []string{"id", "post_id", "title", "text"}
	limitLength   = 7438
)

func Run() error {
	// Creating csv reader
	file, csvReader, err := reader.NewCsvReader(postsFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Creating worksheet
	workSheet, err := excel.NewWorkSheet(sheetName)
	if err != nil {
		return err
	}

	if err := workSheet.SetTitles(titles); err != nil {
		return err
	}

	// Process records from csv file
	totalProcessedRecords, err := processRecords(csvReader, wb, sh)
	if err != nil {
		return err
	}

	log.Printf("Program finished, total records: %d", totalProcessedRecords)

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
		textLength := len(record[3])

		// Skip empty text posts
		if textLength <= 2 {
			continue
		}

		myRecord := model.NewRecord(id, postId, textLength)

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
	cell.SetInt(myRecord.TextLength)

	return nil
}
