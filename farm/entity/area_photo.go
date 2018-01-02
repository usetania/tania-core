package entity

type AreaPhoto struct {
	UID      string
	Filename string
	MimeType string
	Size     int
}

func CreateAreaPhoto(filename, mimetype string, size int) (AreaPhoto, error) {
	return AreaPhoto{
		Filename: filename,
		MimeType: mimetype,
		Size:     size,
	}, nil
}
