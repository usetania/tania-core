package imagehelper

import (
	"image"
	"os"
)

func GetImageDimension(srcPath string) (width, height int, err error) {
	src, err := os.Open(srcPath)
	if err != nil {
		return 0, 0, err
	}
	defer src.Close()

	image, _, err := image.DecodeConfig(src)
	if err != nil {
		return 0, 0, err
	}

	return image.Width, image.Height, nil
}
