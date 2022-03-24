package server

import (
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"strings"
)

// File used to handle file path and file operation.
// We use interface so we can swap it to other file storage easily
type File interface {
	GetFile(src string) ([]byte, error)
	Upload(file *multipart.FileHeader, destPath string) error
}

type LocalFile struct{}

func (LocalFile) GetFile(srcPath string) ([]byte, error) {
	file, err := ioutil.ReadFile(srcPath)

	return file, err
}

// Upload saves uploaded file to the destined path
func (LocalFile) Upload(file *multipart.FileHeader, destPath string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create all directory if not exists
	s := strings.Split(destPath, "/")
	s = s[:len(s)-1]
	sJoin := strings.Join(s, "/")

	if _, err := os.Stat(sJoin); os.IsNotExist(err) {
		log.Print("Upload folder is missing. Creating folder...")
		os.MkdirAll(sJoin, os.ModePerm)
		log.Print("Folder created in ", sJoin)
	}

	// Destination
	dst, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}
