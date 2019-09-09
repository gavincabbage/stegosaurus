package image

import (
	"errors"
	"image"
	"image/draw"
	"image/png"
	"io"

	_ "image/jpeg"
)

type Image struct {
	img  image.Image
	pix  []uint8
	fmt  string
	max  int
	r, w int
}

func NewImage(r io.Reader) (*Image, error) {
	decoded, format, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	var (
		img    = image.NewNRGBA(image.Rect(0, 0, decoded.Bounds().Dx(), decoded.Bounds().Dy()))
		bounds = img.Bounds()
	)
	draw.Draw(img, bounds, decoded, decoded.Bounds().Min, draw.Src)

	return &Image{
		img: img,
		pix: img.Pix,
		max: bounds.Dy() * bounds.Dx() * 4,
		fmt: format,
	}, nil
}

func (img *Image) Read(p []byte) (int, error) {
	var (
		s = img.r
	)
	for i, _ := range p {
		if img.r >= img.max {
			return img.r - s, io.EOF
		}

		p[i] = img.pix[img.r]

		img.r++
	}

	return img.r - s, nil
}

func (img *Image) Write(p []byte) (int, error) {
	var (
		s = img.w
	)
	for _, v := range p {
		if img.w >= img.max {
			return img.w - s, nil
		}

		img.pix[img.w] = v

		img.w++
	}

	return img.w - s, nil
}

func (img *Image) Encode(w io.Writer) error {
	switch img.fmt {
	case "png", "jpeg":
		return png.Encode(w, img.img)
	default:
		return errors.New("unsupported format")
	}
}
