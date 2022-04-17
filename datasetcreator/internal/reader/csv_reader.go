package reader

import (
	"encoding/csv"
	"os"
	"sync"
)

type CsvReader struct {
	*sync.Mutex
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
		Mutex:    &sync.Mutex{},
		Reader:   reader,
		FilePath: filePath,
	}, nil
}

func (cr *CsvReader) SkipRecord() error {
	cr.Lock()
	_, err := cr.Read()
	cr.Unlock()
	if err != nil {
		return nil
	}

	return nil
}
