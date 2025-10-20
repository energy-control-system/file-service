package handler

import (
	"file-service/service/file"
	"fmt"
	"net/http"

	"github.com/sunshineOfficial/golib/gohttp/gorouter"
)

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
