package file

import "file-service/service/file"

func MapFromDB(f File) file.File {
	return file.File{
		ID:       f.ID,
		FileName: f.FileName,
		FileSize: f.FileSize,
		Bucket:   file.Bucket(f.Bucket),
		URL:      f.URL,
	}
}

func MapSliceFromDB(files []File) []file.File {
	result := make([]file.File, 0, len(files))
	for _, f := range files {
		result = append(result, MapFromDB(f))
	}

	return result
}
