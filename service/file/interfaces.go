package file

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/pagination"
)

type Storage interface {
	Add(ctx goctx.Context, fileName string, payload io.Reader, payloadLength int64) (string, error)
}

type Repository interface {
	Add(ctx context.Context, fileHeader *multipart.FileHeader, bucket Bucket, url string) (File, error)
	GetByID(ctx context.Context, id int) (File, error)
	GetByIDs(ctx context.Context, ids []int, page pagination.Pagination) ([]File, error)
}
