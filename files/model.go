package files

type File struct {
	Id          int64  `db:"id"`
	Name        string `db:"name"`
	ContentType string `db:"content_type"`
	Data        []byte `db:"data"`
}
