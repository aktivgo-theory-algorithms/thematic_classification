package app

import (
	"datasetcreator/internal/config"
	"datasetcreator/internal/reader"
	"datasetcreator/internal/worker"
	"datasetcreator/internal/writer"
	"log"
	"os"
)

func Run() error {
	inputFile, err := os.Open(config.InputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	csvReader, err := reader.NewCsvReader(inputFile)
	if err != nil {
		return err
	}

	outputFile, err := os.Open(config.OutputPath)
	if err != nil {
		outputFile, err = os.Create(config.OutputPath)
		if err != nil {
			return err
		}
	}
	defer outputFile.Close()

	csvWriter, err := writer.NewCsvWriter(outputFile)
	if err != nil {
		return err
	}

	w := worker.NewWorker(csvReader, csvWriter)

	totalProcessed, err := w.Process()
	if err != nil {
		return err
	}

	log.Printf("Program finished, total processed: %d", totalProcessed)

	return nil
}
