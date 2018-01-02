package server

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strings"

	"github.com/Tanibox/tania-server/helper/stringhelper"
)

// File used to handle file path and file operation.
// We use interface so we can swap it to other file storage easily
type File interface {
	GetFile(src string) ([]byte, error)
	Upload(file *multipart.FileHeader, destPath string) error

	GetAreaFilepath() string
}

type LocalFile struct {
	AreaFilepath string
}

func InitLocalFile() File {
	// Set custom env variable
	// TODO: Move the setter to something like `.env` file
	os.Setenv("TANIA_AREA_FILE_PATH", "/Users/user/Code/golang/src/github.com/Tanibox/tania-server/uploads")

	areaFilepath := stringhelper.Join(os.Getenv("TANIA_AREA_FILE_PATH"), "/areas/")

	return LocalFile{
		AreaFilepath: areaFilepath,
	}
}

func (f LocalFile) GetFile(srcPath string) ([]byte, error) {
	file, err := ioutil.ReadFile(srcPath)

	return file, err
}

// Upload saves uploaded file to the destined path
func (f LocalFile) Upload(file *multipart.FileHeader, destPath string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create all directory if not exists
	s := strings.Split(destPath, "/")
	s = s[:len(s)-1]
	sJoin := strings.Join(s, "/")
	os.MkdirAll(sJoin, os.ModePerm)

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

func (f LocalFile) GetAreaFilepath() string {
	return f.AreaFilepath
}
