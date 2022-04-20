package writer

import (
	"encoding/csv"
	"os"
	"sync"
)

type CsvWriter struct {
	*sync.Mutex
	*csv.Writer
}

func NewCsvWriter(file *os.File) (*CsvWriter, error) {
	writer := csv.NewWriter(file)

	return &CsvWriter{
		Mutex:  &sync.Mutex{},
		Writer: writer,
	}, nil
}
