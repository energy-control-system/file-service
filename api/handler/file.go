package handler

import (
	"file-service/service/file"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/sunshineOfficial/golib/gohttp/gorouter"
	"github.com/sunshineOfficial/golib/pagination"
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

		response = withPublicStorageURL(c.Request(), response)

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

		response = withPublicStorageURL(c.Request(), response)

		return c.WriteJson(http.StatusOK, response)
	}
}

type fileIDsVars struct {
	IDs    []int `query:"id"`
	Limit  int   `query:"limit"`
	Offset int   `query:"offset"`
}

// GetFilesByIDs godoc
// @Summary Get files by IDs
// @Description Returns metadata for several files.
// @Tags files
// @Produce json
// @Param id query []int true "File IDs" collectionFormat(multi)
// @Param limit query int false "Maximum number of items to return; 0 means no limit"
// @Param offset query int false "Number of items to skip"
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

		response, err := s.GetByIDs(c.Ctx(), vars.IDs, pagination.Pagination{Limit: vars.Limit, Offset: vars.Offset})
		if err != nil {
			return fmt.Errorf("get files by ids: %w", err)
		}

		response = withPublicStorageURLs(c.Request(), response)

		return c.WriteJson(http.StatusOK, response)
	}
}

func withPublicStorageURLs(r *http.Request, files []file.File) []file.File {
	result := make([]file.File, 0, len(files))
	for _, f := range files {
		result = append(result, withPublicStorageURL(r, f))
	}

	return result
}

func withPublicStorageURL(r *http.Request, f file.File) file.File {
	parsed, err := url.Parse(f.URL)
	if err != nil || parsed.Path == "" {
		return f
	}

	host := firstHeaderValue(r.Header.Get("X-Forwarded-Host"))
	if host == "" {
		host = r.Host
	}
	if host == "" {
		return f
	}

	scheme := firstHeaderValue(r.Header.Get("X-Forwarded-Proto"))
	if scheme == "" {
		scheme = "http"
		if r.TLS != nil {
			scheme = "https"
		}
	}

	parsed.Scheme = scheme
	parsed.Host = host
	f.URL = parsed.String()

	return f
}

func firstHeaderValue(value string) string {
	value, _, _ = strings.Cut(value, ",")
	return strings.TrimSpace(value)
}
