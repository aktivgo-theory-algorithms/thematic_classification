package reader

import (
	"encoding/csv"
	"os"
)

type CsvReader struct {
	*csv.Reader
}

func NewCsvReader(filePath string) (*os.File, *csv.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}

	reader := csv.NewReader(file)

	return file, reader, nil
}
