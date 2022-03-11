package myImageLib

import (
	"image"
	"image/jpeg"
	"io"
	"os"
)

func CreateGrayScaleImage(orgImageReader *io.Reader, newFilePath string) string {
	orgImg, _, _ := image.Decode(*orgImageReader)
	originBounds := orgImg.Bounds()
	newGray := image.NewGray(originBounds)

	for y := originBounds.Min.Y; y < originBounds.Max.Y; y++ {
		for x := originBounds.Min.X; x < originBounds.Max.X; x++ {
			newGray.Set(x, y, orgImg.At(x, y))
		}
	}
	newFile, _ := os.Create(newFilePath)
	defer newFile.Close()
	jpeg.Encode(newFile, newGray, &jpeg.Options{Quality: 75})
	return newFilePath
}
