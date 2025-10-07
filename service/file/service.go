package file

import (
	"file-service/database/file"
	"file-service/database/object"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/golog"
)

type Service struct {
	bucketStorages map[Bucket]*object.Minio
	fileRepository *file.Postgres
}

func NewService(bucketStorages map[Bucket]*object.Minio, fileRepository *file.Postgres) *Service {
	return &Service{
		bucketStorages: bucketStorages,
		fileRepository: fileRepository,
	}
}

func (s *Service) Upload(ctx goctx.Context, log golog.Logger, fileHeader *multipart.FileHeader) (File, error) {
	result := File{
		FileName: fileHeader.Filename,
		FileSize: fileHeader.Size,
		Bucket:   bucketFromFileName(fileHeader.Filename),
	}

	storage, ok := s.bucketStorages[result.Bucket]
	if !ok {
		return File{}, fmt.Errorf("bucket %s does not exist", result.Bucket)
	}

	payload, err := fileHeader.Open()
	if err != nil {
		return File{}, fmt.Errorf("open file %s failed: %w", result.FileName, err)
	}
	defer func() {
		if closeErr := payload.Close(); closeErr != nil {
			log.Errorf("close file %s failed: %v", result.FileName, closeErr)
		}
	}()

	result.URL, err = storage.Add(ctx, result.FileName, payload, result.FileSize)
	if err != nil {
		return File{}, fmt.Errorf("upload file %s failed: %w", result.FileName, err)
	}

	result.ID, err = s.fileRepository.Add(ctx, file.File{
		FileName: result.FileName,
		FileSize: result.FileSize,
		Bucket:   string(result.Bucket),
		URL:      result.URL,
	})
	if err != nil {
		return File{}, fmt.Errorf("save file %s failed: %w", result.FileName, err)
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
	f, err := s.fileRepository.GetByID(ctx, id)
	if err != nil {
		return File{}, fmt.Errorf("get file by id %d failed: %w", id, err)
	}

	return File{
		ID:       f.ID,
		FileName: f.FileName,
		FileSize: f.FileSize,
		Bucket:   Bucket(f.Bucket),
		URL:      f.URL,
	}, nil
}
