package reader

import (
	"datasetcreator/internal/config"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
)

type TagsReader struct {
	*sync.Mutex
	*CsvReader
}

type Tags struct {
	*sync.Mutex

	Tags []string
}

func NewTagsReader(csvReader *CsvReader) *TagsReader {
	return &TagsReader{&sync.Mutex{}, csvReader}
}

func (tr *TagsReader) Recreate() (*os.File, *TagsReader, error) {
	file, err := os.Open(tr.FilePath)
	if err != nil {
		return nil, nil, err
	}

	reader := csv.NewReader(file)
	reader.LazyQuotes = true

	return file, &TagsReader{&sync.Mutex{}, &CsvReader{&sync.Mutex{}, reader, tr.FilePath}}, nil
}

func (tr *TagsReader) GetTagsByPostId(postId int) ([]string, error) {
	tags := &Tags{&sync.Mutex{}, nil}

	var wg sync.WaitGroup
	wg.Add(config.GoroutinesCount / 4)

	for i := 0; i < config.GoroutinesCount/4; i++ {
		go func() {
			defer wg.Done()

			for {
				tr.Lock()
				tag, err := tr.Read()
				tr.Unlock()
				if err != nil {
					if err == io.EOF {
						return
					}
					log.Fatal(err)
				}

				postIdTag, err := strconv.Atoi(tag[1])
				if err != nil {
					log.Fatal(err)
				}

				if postIdTag == postId {
					tags.Lock()
					tags.Tags = append(tags.Tags, tag[2])
					tags.Unlock()
				}
			}
		}()
	}

	wg.Wait()

	return tags.Tags, nil
}
