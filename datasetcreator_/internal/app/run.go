package app

import (
	"datasetcreator/internal/config"
	"datasetcreator/internal/reader"
	"datasetcreator/internal/worker"
	"datasetcreator/internal/writer"
	"log"
)

func Run() error {
	filePosts, csvReaderPosts, err := reader.NewCsvReader(config.PostsFilePath)
	if err != nil {
		return err
	}
	defer filePosts.Close()

	postsReader := reader.NewPostsReader(csvReaderPosts)

	fileResult, csvWriter, err := writer.NewCsvWriter(config.ResultFilePath)
	if err != nil {
		return err
	}
	defer fileResult.Close()

	w := worker.NewWorker(postsReader, csvWriter)

	totalProcessedRecords, err := w.Process()
	if err != nil {
		return err
	}

	log.Printf("Program finished, total records: %d", totalProcessedRecords)

	return nil
}
