package handler

import (
	"file-service/service/file"
	"fmt"
	"net/http"

	"github.com/sunshineOfficial/golib/gohttp/gorouter"
)

// UploadFile godoc
// @Summary Upload file
// @Description Uploads one file to object storage and stores its metadata.
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param File formData file true "File to upload"
// @Success 200 {object} file.File
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /files [post]
func UploadFile(s *file.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		files, err := c.FormFiles("File")
		if err != nil {
			return fmt.Errorf("parse form files: %w", err)
		}
		if len(files) != 1 {
			return fmt.Errorf("parse form files: got %d files, expected 1", len(files))
		}

		response, err := s.Upload(c.Ctx(), c.Log(), files[0])
		if err != nil {
			return fmt.Errorf("upload file: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}

type fileIDVars struct {
	ID int `path:"id"`
}

// GetFileByID godoc
// @Summary Get file by ID
// @Description Returns file metadata and a URL by file identifier.
// @Tags files
// @Produce json
// @Param id path int true "File ID"
// @Success 200 {object} file.File
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 404 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /files/{id} [get]
func GetFileByID(s *file.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars fileIDVars
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("get file id: %w", err)
		}

		response, err := s.GetByID(c.Ctx(), vars.ID)
		if err != nil {
			return fmt.Errorf("get file by id: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}

type fileIDsVars struct {
	IDs []int `query:"id"`
}

// GetFilesByIDs godoc
// @Summary Get files by IDs
// @Description Returns metadata for several files.
// @Tags files
// @Produce json
// @Param id query []int true "File IDs" collectionFormat(multi)
// @Success 200 {array} file.File
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /files [get]
func GetFilesByIDs(s *file.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars fileIDsVars
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("get files ids: %w", err)
		}

		response, err := s.GetByIDs(c.Ctx(), vars.IDs)
		if err != nil {
			return fmt.Errorf("get files by ids: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}
