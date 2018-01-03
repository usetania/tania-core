package imagehelper

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func GetImageDimension(srcPath string) (width int, height int, err error) {
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
