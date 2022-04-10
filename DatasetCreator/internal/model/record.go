package model

type Record struct {
	ID     int
	PostID int
	Title  string
	Text   string
	Tags   []string
}

func NewRecord(id int, postId int, title string, text string) *Record {
	return &Record{
		ID:     id,
		PostID: postId,
		Title:  title,
		Text:   text,
	}
}

func (r *Record) AddTag(tag string) {
	r.Tags = append(r.Tags, tag)
}

func (r *Record) AddTags(tags []string) {
	r.Tags = append(r.Tags, tags...)
}
