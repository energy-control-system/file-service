package file

type File struct {
	ID       int    `db:"id"`
	FileName string `db:"file_name"`
	FileSize int64  `db:"file_size"`
	Bucket   string `db:"bucket"`
	URL      string `db:"url"`
}
