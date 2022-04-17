package writer

import (
	"encoding/csv"
	"os"
	"sync"
)

type CsvWriter struct {
	*sync.Mutex
	*csv.Writer

	FilePath string
}

func NewCsvWriter(filePath string) (*os.File, *CsvWriter, error) {
	file, err := os.Open(filePath)
	if err != nil {
		file, err = os.Create(filePath)
		if err != nil {
			return nil, nil, err
		}
	}

	writer := csv.NewWriter(file)

	return file, &CsvWriter{
		Mutex:    &sync.Mutex{},
		Writer:   writer,
		FilePath: filePath,
	}, nil
}
