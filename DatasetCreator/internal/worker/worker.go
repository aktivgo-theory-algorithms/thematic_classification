package worker

import (
	"DatasetCreater/internal/model"
	"DatasetCreater/internal/reader"
	"DatasetCreater/internal/xlsx"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	limitLength = 7438
)

type Worker struct {
	PostsReader *reader.PostReader
	TagsReader  *reader.TagsReader
	WorkSheet   *xlsx.WorkSheet
}

func NewWorker(postsReader *reader.PostReader, tagsReader *reader.TagsReader, workSheet *xlsx.WorkSheet) *Worker {
	return &Worker{
		PostsReader: postsReader,
		TagsReader:  tagsReader,
		WorkSheet:   workSheet,
	}
}

func (w *Worker) Process() (int, error) {
	defer w.WorkSheet.Save()

	totalProcessedRecords := 0

	// Skip first record, because it is titles
	err := w.PostsReader.SkipRecord()
	if err != nil {
		return totalProcessedRecords, nil
	}

	row := 1
	for {
		post, err := w.PostsReader.GetNext()
		if err != nil {
			return totalProcessedRecords, err
		}

		if post == nil {
			return totalProcessedRecords, nil
		}

		textLength := len(post[3])

		// Skip empty and > limitLength text posts
		if textLength <= 2 || textLength > limitLength {
			continue
		}

		id, err := strconv.Atoi(post[0])
		if err != nil {
			return totalProcessedRecords, err
		}
		postId, err := strconv.Atoi(post[1])
		if err != nil {
			return totalProcessedRecords, err
		}
		title := strings.TrimSpace(post[2])
		text := strings.TrimSpace(post[3])

		curRecord := model.NewRecord(id, postId, title, text)

		tagsFile, tagsReader, err := w.TagsReader.Recreate()
		if err != nil {
			return totalProcessedRecords, err
		}

		w.TagsReader = tagsReader

		// Skip first record, because it is titles
		err = w.TagsReader.SkipRecord()
		if err != nil {
			return totalProcessedRecords, err
		}

		tags, err := w.TagsReader.GetTagsByPostId(postId)
		if err != nil {
			tagsFile.Close()
			return totalProcessedRecords, err
		}

		if tags != nil {
			curRecord.AddTags(tags)
		}

		// Fill row
		if err = w.fillRow(row, curRecord); err != nil {
			tagsFile.Close()
			return totalProcessedRecords, err
		}

		// Save file every 29-31 Seconds
		if time.Now().Second() >= 29 && time.Now().Second() <= 31 {
			if err = w.WorkSheet.Save(); err != nil {
				tagsFile.Close()
				return totalProcessedRecords, err
			}
			log.Println("File saved")
			time.Sleep(5 * time.Second)
		}

		log.Printf("post (id = %d, post_id = %d) successfully recorded to row %d",
			curRecord.ID, curRecord.PostID, row)

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
