package files

import (
	"github.com/jmoiron/sqlx"
	"meal-planner/core"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(c core.Context) *Repository {
	return &Repository{db: c.Db()}
}

const filesInsert = `INSERT INTO files (name, content_type, data) VALUES (:name, :content_type, :data)`

func (r *Repository) CreateFile(file *File) error {
	result, err := r.db.NamedExec(filesInsert, file)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	file.Id = id
	return nil
}
func (r *Repository) GetFile(id int64) (*File, error) {
	var file File
	err := r.db.Get(&file, "SELECT * FROM files WHERE id=$1", id)
	return &file, err
}

func (r *Repository) DeleteFile(id int64) error {
	_, err := r.db.Exec("DELETE FROM files WHERE id=$1", id)
	return err
}
