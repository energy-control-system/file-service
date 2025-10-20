package file

import (
	"context"
	_ "embed"
	"file-service/service/file"
	"fmt"
	"mime/multipart"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

//go:embed sql/add.sql
var addSQL string

func (r *Repository) Add(ctx context.Context, fileHeader *multipart.FileHeader, bucket file.Bucket, url string) (file.File, error) {
	var f File
	err := r.db.GetContext(ctx, &f, addSQL, fileHeader.Filename, fileHeader.Size, bucket, url)
	if err != nil {
		return file.File{}, fmt.Errorf("r.db.GetContext: %w", err)
	}

	return MapFromDB(f), nil
}

//go:embed sql/get_by_id.sql
var getByIDSQL string

func (r *Repository) GetByID(ctx context.Context, id int) (file.File, error) {
	var f File
	err := r.db.GetContext(ctx, &f, getByIDSQL, id)
	if err != nil {
		return file.File{}, fmt.Errorf("r.db.GetContext: %w", err)
	}

	return MapFromDB(f), nil
}

//go:embed sql/get_by_ids.sql
var getByIDsSQL string

func (r *Repository) GetByIDs(ctx context.Context, ids []int) ([]file.File, error) {
	var files []File
	err := r.db.SelectContext(ctx, &files, getByIDsSQL, ids)
	if err != nil {
		return nil, fmt.Errorf("r.db.SelectContext: %w", err)
	}

	return MapSliceFromDB(files), nil
}
