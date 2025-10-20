package file

import (
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/golog"
)

type Service struct {
	bucketStorages map[Bucket]Storage
	repository     Repository
}

func NewService(bucketStorages map[Bucket]Storage, repository Repository) *Service {
	return &Service{
		bucketStorages: bucketStorages,
		repository:     repository,
	}
}

func (s *Service) Upload(ctx goctx.Context, log golog.Logger, fileHeader *multipart.FileHeader) (File, error) {
	bucket := bucketFromFileName(fileHeader.Filename)

	storage, ok := s.bucketStorages[bucket]
	if !ok {
		return File{}, fmt.Errorf("bucket %s does not exist", bucket)
	}

	payload, err := fileHeader.Open()
	if err != nil {
		return File{}, fmt.Errorf("open file %s failed: %w", fileHeader.Filename, err)
	}
	defer func() {
		if closeErr := payload.Close(); closeErr != nil {
			log.Errorf("close file %s failed: %v", fileHeader.Filename, closeErr)
		}
	}()

	url, err := storage.Add(ctx, fileHeader.Filename, payload, fileHeader.Size)
	if err != nil {
		return File{}, fmt.Errorf("upload file %s failed: %w", fileHeader.Filename, err)
	}

	result, err := s.repository.Add(ctx, fileHeader, bucket, url)
	if err != nil {
		return File{}, fmt.Errorf("save file %s failed: %w", fileHeader.Filename, err)
	}

	return result, nil
}

func bucketFromFileName(name string) Bucket {
	switch filepath.Ext(name) {
	case ".jpg", ".jpeg", ".png":
		return BucketImages
	default:
		return BucketDocuments
	}
}

func (s *Service) GetByID(ctx goctx.Context, id int) (File, error) {
	f, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return File{}, fmt.Errorf("get file by id %d failed: %w", id, err)
	}

	return f, nil
}

func (s *Service) GetByIDs(ctx goctx.Context, ids []int) ([]File, error) {
	files, err := s.repository.GetByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("get files by ids from db: %w", err)
	}

	return files, nil
}
