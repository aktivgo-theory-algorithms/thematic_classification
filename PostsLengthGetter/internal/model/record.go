package model

type Record struct {
	ID         int
	PostID     int
	TextLength int
}

func NewRecord(id int, postId int, textLength int) *Record {
	return &Record{
		ID:         id,
		PostID:     postId,
		TextLength: textLength,
	}
}
