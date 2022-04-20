package worker

import (
	"datasetcreator/internal/config"
	"datasetcreator/internal/consts"
	"datasetcreator/internal/helpers"
	"datasetcreator/internal/reader"
	"datasetcreator/internal/writer"
	"log"
	"strings"
	"sync"
)

type Worker struct {
	CsvReader *reader.CsvReader
	CsvWriter *writer.CsvWriter
}

func NewWorker(csvReader *reader.CsvReader, csvWriter *writer.CsvWriter) *Worker {
	return &Worker{
		CsvReader: csvReader,
		CsvWriter: csvWriter,
	}
}

func (w *Worker) Process() (int, error) {
	totalProcessed := 0

	// Skip first row, because it is titles
	err := w.CsvReader.Skip()
	if err != nil {
		return totalProcessed, nil
	}

	var wg sync.WaitGroup
	wg.Add(2)

	recordChannel := make(chan []string, config.GoroutinesCount)

	// Read rows
	go func() {
		defer func() {
			wg.Done()
			close(recordChannel)
		}()

		var wgInside sync.WaitGroup
		wgInside.Add(config.GoroutinesCount / 2)

		for i := 0; i < config.GoroutinesCount/2; i++ {
			go func() {
				defer wgInside.Done()

				for {
					row, err := w.CsvReader.GetNext()
					if err != nil {
						log.Fatal(err)
					}

					if row == nil {
						return
					}

					text := strings.TrimSpace(row[2])
					tag := strings.TrimSpace(row[4])

					if helpers.Contains(consts.SkipTagNames, tag) {
						continue
					}

					recordChannel <- []string{text, tag}
				}
			}()
		}

		wgInside.Wait()
	}()

	// Write records
	go func() {
		defer wg.Done()

		var wgInside sync.WaitGroup
		wgInside.Add(config.GoroutinesCount / 2)

		for i := 0; i < config.GoroutinesCount/2; i++ {
			go func() {
				defer wgInside.Done()

				for {
					record, ok := <-recordChannel
					if !ok {
						return
					}

					if err = w.write(record); err != nil {
						log.Fatal(err)
					}

					log.Printf("text=%s tag=%s successfully writed", record[0], record[1])

					totalProcessed++
				}
			}()
		}

		wgInside.Wait()
	}()

	wg.Wait()

	return totalProcessed, nil
}
func (w *Worker) write(record []string) error {
	w.CsvWriter.Lock()
	err := w.CsvWriter.Write(record)
	w.CsvWriter.Unlock()
	return err
}
