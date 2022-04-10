package worker

import (
	"DatasetCreater/internal/config"
	"DatasetCreater/internal/model"
	"DatasetCreater/internal/reader"
	"DatasetCreater/internal/xlsx"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	limitLength = 7438
)

type Worker struct {
	CSVReaderPosts *reader.CsvReader
	CSVReaderTags  *reader.CsvReader
	WorkSheet      *xlsx.WorkSheet
}

func NewWorker(csvReader *reader.CsvReader, workSheet *xlsx.WorkSheet) *Worker {
	return &Worker{
		CSVReaderPosts: csvReader,
		WorkSheet:      workSheet,
	}
}

func (w *Worker) Process() (int, error) {
	totalProcessedRecords := 0

	// Skip first record, because it is titles
	_, err := w.CSVReaderPosts.Read()
	if err != nil {
		return totalProcessedRecords, nil
	}

	row := 1
	for {
		// Read next post
		post, err := w.CSVReaderPosts.Read()
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

		textLength := len(post[3])

		// Skip empty and > limitLength text posts
		if textLength <= 2 || textLength > limitLength {
			continue
		}

		id, err := strconv.Atoi(post[0])
		postId, err := strconv.Atoi(post[1])
		title := strings.TrimSpace(post[2])
		text := strings.TrimSpace(post[3])

		curRecord := model.NewRecord(id, postId, title, text)

		// Get tags

		fileTags, csvTagsReader, err := reader.NewCsvReader(config.TagsFilePath)
		if err != nil {
			return totalProcessedRecords, err
		}
		w.CSVReaderTags = csvTagsReader

		// Skip first record, because it is titles
		_, err = w.CSVReaderTags.Read()
		if err != nil {
			fileTags.Close()
			return totalProcessedRecords, err
		}

		for {
			// Read next tag
			tag, err := w.CSVReaderTags.Read()
			if err != nil {
				// If end save file and return
				if err == io.EOF {
					if err := w.WorkSheet.Save(); err != nil {
						fileTags.Close()
						return totalProcessedRecords, err
					}

					break
				}
				return totalProcessedRecords, err
			}

			postIdTag, err := strconv.Atoi(tag[1])
			if postIdTag != postId {
				continue
			}

			curRecord.AddTag(tag[2])
		}

		// Fill row
		if err = w.fillRow(row, curRecord); err != nil {
			return totalProcessedRecords, err
		}

		// Save file every 29-31 Seconds
		if time.Now().Second() >= 29 && time.Now().Second() <= 31 {
			if err = w.WorkSheet.Save(); err != nil {
				fileTags.Close()
				return totalProcessedRecords, err
			}
		}

		log.Printf("post (id = %d, post_id = %d) successfully recorded to row %d",
			curRecord.ID, curRecord.PostID, row)

		totalProcessedRecords++
		row++
	}

	return totalProcessedRecords, nil
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
	cell.SetString(myRecord.Title)

	cell, err = w.WorkSheet.Sheet.Cell(row, 3)
	if err != nil {
		return err
	}
	cell.SetString(myRecord.Text)

	cell, err = w.WorkSheet.Sheet.Cell(row, 4)
	if err != nil {
		return err
	}
	cell.SetString(strings.Join(myRecord.Tags, ","))

	return nil
}
