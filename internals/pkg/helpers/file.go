package helpers

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func SetFilePermission(realtiveFilepath string) error {
	return os.Chmod(realtiveFilepath, 0644)
}

// func UploadFile(w http.ResponseWriter, userPath, fileName string, file *multipart.File) error {
// 	if err := os.MkdirAll(userPath, os.ModePerm); err != nil {

// 		return err
// 	}
// 	relativePath := filepath.Join(userPath, fileName)
// 	dst, err := os.Create(relativePath)
// 	if err != nil {

// 		return err
// 	}
// 	defer dst.Close()
// 	if err := SetFilePermission(relativePath); err != nil {

// 		return err
// 	}
// 	_, err = io.Copy(dst, *file)
// 	if err != nil {

//			return err
//		}
//		return nil
//	}
func UploadFile(w http.ResponseWriter, userPath string, fileHeader *multipart.FileHeader) (string, int64, error) {
	if err := os.MkdirAll(userPath, os.ModePerm); err != nil {

		return "", 0, err
	}
	file, err := fileHeader.Open()
	if err != nil {
		return "", 0, err
	}
	defer file.Close()
	fileName := fileHeader.Filename
	relativePath := filepath.Join(userPath, fileName)
	dst, err := os.Create(relativePath)
	if err != nil {

		return "", 0, err
	}
	defer dst.Close()
	if err := SetFilePermission(relativePath); err != nil {

		return "", 0, err
	}
	_, err = io.Copy(dst, file)
	if err != nil {

		return "", 0, err
	}
	return fileName, fileHeader.Size, nil
}
