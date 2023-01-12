package qrcode

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "github.com/chai2010/webp"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/pkg/errors"
)

func Decode(img image.Image) (string, error) {
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", err
	}

	r, err := qrcode.NewQRCodeReader().Decode(bmp, nil)
	if err != nil {
		return "", errors.Wrap(err, "decode failed")
	}

	return r.String(), nil
}
