package model

type Record struct {
	ID         int
	PostID     int
	Title      string
	Text       string
	TextLength int
}

func NewRecord(id int, postId int, title string, text string, textLength int) *Record {
	return &Record{
		ID:         id,
		PostID:     postId,
		Title:      title,
		Text:       text,
		TextLength: textLength,
	}
}
