package stegosaurus

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
)

type Image struct {
	img           *image.NRGBA
	fmt           string
	height, width int
	read, written int
}

type imageSetter interface {
	image.Image
	Set(int, int, color.Color)
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
	fmt.Printf("%d %d %d %d\n", bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y)

	return &Image{
		img:    img,
		height: bounds.Dy(),
		width:  bounds.Dx(),
		fmt:    format,
	}, nil
}

func (img *Image) Read(p []byte) (int, error) {
	var (
		s = img.read
	)
	//fmt.Printf("len Pix = %d\n", len(img.img.Pix))
	for i, _ := range p {
		if img.read >= img.height*img.width*3 {
			return img.read - s, io.EOF
		}

		if (img.read+1)%4 == 0 {
			img.read++
			s++
		}

		p[i] = img.img.Pix[img.read]

		img.read++
	}

	return img.read - s, nil
}

// RGBA is an in-memory image whose At method returns color.RGBA values.
type myRGBA struct {
	// Pix holds the image's pixels, in R, G, B, A order. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect int // Rectangle type
}

func (img *Image) Write(p []byte) (int, error) {
	var (
		s = img.written
	)
	for _, v := range p {
		if img.written >= img.height*img.width*3 {
			return img.written - s, nil
		}

		if (img.written+1)%4 == 0 {
			img.written++
			s++
		}

		img.img.Pix[img.written] = v

		img.written++
	}

	return img.written - s, nil
}

func (img *Image) Encode(w io.Writer) error {
	switch img.fmt {
	case "png":
		return png.Encode(w, img.img)
	case "jpeg":
		return jpeg.Encode(w, img.img, &jpeg.Options{100})
	default:
		return errors.New("unsupported format")
	}
}
