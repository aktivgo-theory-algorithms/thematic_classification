package reader

import (
	"encoding/csv"
	"io"
	"os"
	"sync"
)

type CsvReader struct {
	*sync.Mutex
	*csv.Reader
}

func NewCsvReader(file *os.File) (*CsvReader, error) {
	reader := csv.NewReader(file)

	return &CsvReader{
		Mutex:  &sync.Mutex{},
		Reader: reader,
	}, nil
}

func (cr *CsvReader) Skip() error {
	_, err := cr.Read()
	if err != nil {
		return nil
	}

	return nil
}

func (cr *CsvReader) GetNext() ([]string, error) {
	cr.Lock()
	row, err := cr.Read()
	cr.Unlock()
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	return row, nil
}
