package reader

import (
	"io"
	"log"
	"strconv"
)

type TagsReader struct {
	*CsvReader
	lastTag []string
}

func NewTagsReader(csvReader *CsvReader) *TagsReader {
	return &TagsReader{
		CsvReader: csvReader,
	}
}

func (tr *TagsReader) GetTagsByPostId(postId int) ([]string, error) {
	var tags []string

	if tr.lastTag != nil {
		postIdTag, err := strconv.Atoi(tr.lastTag[1])
		if err != nil {
			return nil, err
		}

		if postIdTag == postId {
			tags = append(tags, tr.lastTag[2])
		} else {
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
					tags = append(tags, tr.lastTag[2])
					break
				}
			}
		}
	}

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

		if postIdTag != postId {
			tr.lastTag = tag
			log.Println(postId, tag)
			return tags, nil
		}

		tags = append(tags, tag[2])
	}
}
