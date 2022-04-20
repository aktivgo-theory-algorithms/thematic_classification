package reader

import (
	"io"
	"sync"
)

type PostReader struct {
	*sync.Mutex
	*CsvReader
}

func NewPostsReader(csvReader *CsvReader) *PostReader {
	return &PostReader{&sync.Mutex{}, csvReader}
}

func (pr *PostReader) GetNext() ([]string, error) {
	pr.Lock()
	post, err := pr.Read()
	pr.Unlock()
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	return post, nil
}
