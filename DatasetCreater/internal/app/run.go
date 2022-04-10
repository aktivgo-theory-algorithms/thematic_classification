package app

import (
	"DatasetCreater/internal/reader"
	"DatasetCreater/internal/worker"
	"DatasetCreater/internal/xlsx"
	"log"
)

var (
	postsFilePath = "../Dataset/posts.csv"
	sheetName     = "posts_length"
	xlsxFilePath  = "data/posts_length_test.xlsx"
	titles        = []string{"id", "post_id", "title", "text"}
)

func Run() error {
	// Creating csv reader
	file, csvReader, err := reader.NewCsvReader(postsFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Creating worksheet
	workSheet, err := xlsx.NewWorkSheet(sheetName, xlsxFilePath)
	if err != nil {
		return err
	}

	if err := workSheet.SetTitles(titles); err != nil {
		return err
	}

	w := worker.NewWorker(csvReader, workSheet)

	totalProcessedRecords, err := w.Process()
	if err != nil {
		return err
	}

	w.WorkSheet.CloseSheet()

	log.Printf("Program finished, total records: %d", totalProcessedRecords)

	return nil
}
