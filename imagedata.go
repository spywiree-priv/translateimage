package translateimage

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
)

type ImageData struct {
	data     []byte
	mimeType string
	img      image.Image
}

func (i *ImageData) ColorModel() color.Model {
	if i.img == nil {
		i.decode()
	}
	return i.img.ColorModel()
}

func (i *ImageData) Bounds() image.Rectangle {
	if i.img == nil {
		i.decode()
	}
	return i.img.Bounds()
}

func (i *ImageData) At(x, y int) color.Color {
	if i.img == nil {
		i.decode()
	}
	return i.img.At(x, y)
}

func (i *ImageData) safeDecode() (image.Image, error) {
	var img image.Image
	var err error

	switch i.mimeType {
	case "image/png":
		img, err = png.Decode(bytes.NewReader(i.data))
	case "image/jpeg":
		img, err = jpeg.Decode(bytes.NewReader(i.data))
	}
	if err != nil {
		return nil, err
	}

	return img, err
}

func (i *ImageData) decode() {
	var err error
	i.img, err = i.safeDecode()
	if err != nil {
		panic(err)
	}
}

var UnknownMimeType = errors.New("unknown mime type")

func (i *ImageData) ConvertTo(mimeType string) ([]byte, error) {
	if i.mimeType == mimeType {
		return i.data, nil
	}

	if i.img == nil {
		i.decode()
	}

	buf := bytes.NewBuffer(nil)
	var err error
	switch mimeType {
	case "image/png":
		err = png.Encode(buf, i.img)
	case "image/jpeg":
		err = jpeg.Encode(buf, i.img, nil)
	default:
		err = UnknownMimeType
	}
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (i *ImageData) Data() []byte {
	return i.data
}

func (i *ImageData) MimeType() string {
	return i.mimeType
}
