package reader

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
)

type TagsReader struct {
	*CsvReader
}

func NewTagsReader(csvReader *CsvReader) *TagsReader {
	return &TagsReader{
		CsvReader: csvReader,
	}
}

func (tr *TagsReader) Recreate() (*os.File, *TagsReader, error) {
	file, err := os.Open(tr.FilePath)
	if err != nil {
		return nil, nil, err
	}

	reader := csv.NewReader(file)
	reader.LazyQuotes = true

	return file, &TagsReader{&CsvReader{reader, tr.FilePath}}, nil
}

func (tr *TagsReader) GetTagsByPostId(postId int) ([]string, error) {
	var tags []string

	for {
		tag, err := tr.Read()
		if err != nil {
			if err == io.EOF {
				return tags, nil
			}
			return nil, err
		}

		postIdTag, err := strconv.Atoi(tag[1])
		if err != nil {
			return nil, err
		}

		if postIdTag == postId {
			tags = append(tags, tag[2])
		}
	}
}
