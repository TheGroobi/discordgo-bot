package helper

import (
	"image"
	_ "image/png"
	"io"
	"strconv"

	"github.com/cenkalti/dominantcolor"
)

func FindDominantColor(file io.Reader) (int, error) {

	img, _, err := image.Decode(file)
	if err != nil {
		return 0, err
	}

	hex := dominantcolor.Hex(dominantcolor.Find(img))
	hex = hex[1:]

	color, err := strconv.ParseInt(hex, 16, 32)
	
	if err != nil {
		return 0, err
	}

	return int(color), nil
}
