package entity

type AreaPhoto struct {
	UID      string
	Filename string
	MimeType string
	Size     int

	Area Area
}

func CreateAreaPhoto(area Area, filename, mimetype string, size int) (AreaPhoto, error) {
	return AreaPhoto{
		Area:     area,
		Filename: filename,
		MimeType: mimetype,
		Size:     size,
	}, nil
}
