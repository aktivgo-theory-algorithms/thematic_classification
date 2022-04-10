package app

import (
	"DatasetCreater/internal/config"
	"DatasetCreater/internal/reader"
	"DatasetCreater/internal/worker"
	"DatasetCreater/internal/xlsx"
	"log"
)

func Run() error {
	filePosts, csvReaderPosts, err := reader.NewCsvReader(config.PostsFilePath)
	if err != nil {
		return err
	}
	defer filePosts.Close()

	fileTags, csvReaderTags, err := reader.NewCsvReader(config.TagsFilePath)
	if err != nil {
		return err
	}
	defer fileTags.Close()

	postsReader := reader.NewPostsReader(csvReaderPosts)
	tagsReader := reader.NewTagsReader(csvReaderTags)

	workSheet, err := xlsx.NewWorkSheet(config.SheetName, config.XlsxFilePath)
	if err != nil {
		return err
	}

	if err := workSheet.SetTitles(config.Titles); err != nil {
		return err
	}

	w := worker.NewWorker(postsReader, tagsReader, workSheet)

	totalProcessedRecords, err := w.Process()
	if err != nil {
		return err
	}

	w.WorkSheet.CloseSheet()

	log.Printf("Program finished, total records: %d", totalProcessedRecords)

	return nil
}
