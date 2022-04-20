package worker

import (
	"datasetcreator/internal/config"
	"datasetcreator/internal/reader"
	"datasetcreator/internal/writer"
	"log"
	"strconv"
	"strings"
	"sync"
)

var (
	limitLength = 7438
)

type Worker struct {
	PostsReader *reader.PostReader
	CsvWriter   *writer.CsvWriter
}

func NewWorker(postsReader *reader.PostReader, csvWriter *writer.CsvWriter) *Worker {
	return &Worker{
		PostsReader: postsReader,
		CsvWriter:   csvWriter,
	}
}

func (w *Worker) Process() (int, error) {
	totalProcessedPosts := 0

	// Skip first record, because it is titles
	err := w.PostsReader.SkipRecord()
	if err != nil {
		return totalProcessedPosts, nil
	}

	var wg sync.WaitGroup
	wg.Add(2)

	recordChannel := make(chan []string, config.GoroutinesCount)

	// Read records
	go func() {
		defer func() {
			wg.Done()
			close(recordChannel)
		}()

		var wgInside sync.WaitGroup
		wgInside.Add(config.GoroutinesCount / 6)

		for i := 0; i < config.GoroutinesCount/6; i++ {
			go func() {
				defer wgInside.Done()

				for {
					post, err := w.PostsReader.GetNext()
					if err != nil {
						log.Fatal(err)
					}

					if post == nil {
						return
					}

					textLength := len(post[3])

					// Skip empty and > limitLength text posts
					if textLength <= 2 || textLength > limitLength {
						continue
					}

					postID, err := strconv.Atoi(post[1])
					if err != nil {
						log.Fatal(err)
					}
					title := strings.TrimSpace(post[2])
					text := strings.TrimSpace(post[3])

					fileTags, csvReaderTags, err := reader.NewCsvReader(config.TagsFilePath)
					if err != nil {
						log.Fatalln(err)
					}

					tagsReader := reader.NewTagsReader(csvReaderTags)

					// Skip first record, because it is titles
					err = tagsReader.SkipRecord()
					if err != nil {
						log.Fatal(err)
					}

					tags, err := tagsReader.GetTagsByPostId(postID)
					if err != nil {
						fileTags.Close()
						log.Fatal(err)
					}

					if tags == nil {
						return
					}

					for _, tag := range tags {
						var record []string

						record = append(record, strconv.Itoa(postID), title, text, tag)
						recordChannel <- record
					}

					totalProcessedPosts++
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

					log.Printf("post_id=%s tag=%s successfully writed", record[0], record[3])
				}
			}()
		}
	}()

	wg.Wait()

	return totalProcessedPosts, nil
}
func (w *Worker) write(record []string) error {
	w.CsvWriter.Lock()
	err := w.CsvWriter.Write(record[1:])
	w.CsvWriter.Unlock()
	return err
}
