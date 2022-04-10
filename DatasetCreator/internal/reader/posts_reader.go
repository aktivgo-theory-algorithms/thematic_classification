package reader

import "io"

type PostReader struct {
	*CsvReader
}

func NewPostsReader(csvReader *CsvReader) *PostReader {
	return &PostReader{csvReader}
}

func (pr *PostReader) GetNext() ([]string, error) {
	post, err := pr.Read()
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	return post, nil
}
