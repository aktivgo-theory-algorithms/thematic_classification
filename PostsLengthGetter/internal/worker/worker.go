package worker

import (
	"PostsLengthGetter/internal/model"
	"PostsLengthGetter/internal/reader"
	"PostsLengthGetter/internal/xlsx"
	"io"
	"log"
	"strconv"
	"time"
)

type Worker struct {
	CSVReader *reader.CsvReader
	WorkSheet *xlsx.WorkSheet
}

func NewWorker(csvReader *reader.CsvReader, workSheet *xlsx.WorkSheet) *Worker {
	return &Worker{
		CSVReader: csvReader,
		WorkSheet: workSheet,
	}
}

func (w *Worker) Process() (int, error) {
	totalProcessedRecords := 0

	// Skip first record, because it is titles
	_, err := w.CSVReader.Read()
	if err != nil {
		return totalProcessedRecords, nil
	}

	row := 1
	for {
		// Read next record
		record, err := w.CSVReader.Read()
		if err != nil {
			// If end save file and return
			if err == io.EOF {
				if err = w.WorkSheet.Save(); err != nil {
					return totalProcessedRecords, err
				}
				return totalProcessedRecords, nil
			}

			return totalProcessedRecords, err
		}

		id, err := strconv.Atoi(record[0])
		postId, err := strconv.Atoi(record[1])
		textLength := len(record[3])

		// Skip empty text posts
		if textLength <= 2 {
			continue
		}

		curRecord := model.NewRecord(id, postId, textLength)

		// Fill row
		if err = w.fillRow(row, curRecord); err != nil {
			return totalProcessedRecords, err
		}

		// Save file every 29-31 Seconds
		if time.Now().Second() >= 29 && time.Now().Second() <= 31 {
			if err = w.WorkSheet.Save(); err != nil {
				return totalProcessedRecords, err
			}
		}

		log.Printf("record (id = %d, post_id = %d, text_length = %d) successfully recorded to row %d",
			curRecord.ID, curRecord.PostID, curRecord.TextLength, row)

		totalProcessedRecords++
		row++
	}
}

func (w *Worker) fillRow(row int, myRecord *model.Record) error {
	cell, err := w.WorkSheet.Sheet.Cell(row, 0)
	if err != nil {
		return err
	}
	cell.SetInt(myRecord.ID)

	cell, err = w.WorkSheet.Sheet.Cell(row, 1)
	if err != nil {
		return err
	}
	cell.SetInt(myRecord.PostID)

	cell, err = w.WorkSheet.Sheet.Cell(row, 2)
	if err != nil {
		return err
	}
	cell.SetInt(myRecord.TextLength)

	return nil
}
