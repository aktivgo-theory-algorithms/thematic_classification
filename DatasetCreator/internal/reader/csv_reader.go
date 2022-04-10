package reader

import (
	"encoding/csv"
	"os"
)

type CsvReader struct {
	*csv.Reader

	FilePath string
}

func NewCsvReader(filePath string) (*os.File, *CsvReader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}

	reader := csv.NewReader(file)
	reader.LazyQuotes = true

	return file, &CsvReader{
		Reader:   reader,
		FilePath: filePath,
	}, nil
}

func (cr *CsvReader) SkipRecord() error {
	_, err := cr.Read()
	if err != nil {
		return nil
	}

	return nil
}
