package giftImageLib

import (
	"github.com/disintegration/gift"
	"image"
	"image/jpeg"
	"io"
	"os"
)

func CreateInvertImage(orgImageReader *io.Reader, newFilePath string) string {
	orgImg, _, _ := image.Decode(*orgImageReader)

	// 1. Create a new filter list and add some filters.
	g := gift.New(
		gift.Invert(),
	)

	dst := image.NewRGBA(g.Bounds(orgImg.Bounds()))
	g.Draw(dst, orgImg)

	newFile, _ := os.Create(newFilePath)
	defer newFile.Close()
	jpeg.Encode(newFile, dst, &jpeg.Options{Quality: 75})
	return newFilePath
}

func CreateColorizeImage(orgImageReader *io.Reader, newFilePath string) string {
	orgImg, _, _ := image.Decode(*orgImageReader)

	// 1. Create a new filter list and add some filters.
	g := gift.New(
		gift.Colorize(2.0, 2.0, 63),
	)

	dst := image.NewRGBA(g.Bounds(orgImg.Bounds()))
	g.Draw(dst, orgImg)

	newFile, _ := os.Create(newFilePath)
	defer newFile.Close()
	jpeg.Encode(newFile, dst, &jpeg.Options{Quality: 75})
	return newFilePath
}
