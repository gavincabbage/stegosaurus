package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

type options struct {
	name   string
	fmt    string
	width  int
	height int
	img    draw.Image
	color  color.Color
}

func main() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	todo := []options{
		{
			name:   "uniform blue NRGBA",
			fmt:    "png",
			width:  256,
			height: 256,
			img:    image.NewNRGBA(image.Rect(0, 0, 256, 256)),
			color:  color.NRGBA{0, 0, 255, 255},
		},
		{
			name:   "random NRGBA",
			fmt:    "png",
			width:  256,
			height: 256,
			img:    image.NewNRGBA(image.Rect(0, 0, 256, 256)),
		},
		{
			name:   "uniform blue RGBA",
			fmt:    "png",
			width:  256,
			height: 256,
			img:    image.NewRGBA(image.Rect(0, 0, 256, 256)),
			color:  color.RGBA{0, 0, 255, 255},
		},
		{
			name:   "random RGBA",
			fmt:    "png",
			width:  256,
			height: 256,
			img:    image.NewRGBA(image.Rect(0, 0, 256, 256)),
		},
		{
			name:   "uniform blue CMYK",
			fmt:    "jpeg",
			width:  256,
			height: 256,
			img:    image.NewCMYK(image.Rect(0, 0, 256, 256)),
			color:  color.NRGBA{0, 0, 255, 255},
		},
		{
			name:   "random CMYK",
			fmt:    "jpeg",
			width:  256,
			height: 256,
			img:    image.NewCMYK(image.Rect(0, 0, 256, 256)),
		},
		{
			name:   "gray",
			fmt:    "png",
			width:  256,
			height: 256,
			img:    image.NewGray(image.Rect(0, 0, 256, 256)),
			color:  color.Gray{Y: 23},
		},
	}

	var g errgroup.Group
	for _, t := range todo {
		t := t
		g.Go(func() (err error) {
			t.name = strings.Replace(t.name, " ", "_", -1)
			filename := fmt.Sprintf("%s/%s.%s", path, t.name, t.fmt)
			file, err := os.Create(filename)
			if err != nil {
				return fmt.Errorf("opening %s for creation: %w", filename, err)
			}
			defer func() {
				if e := file.Close(); e != nil && err == nil {
					err = e
				}
			}()

			var src image.Image
			if t.color == nil {
				src = randomImage(t.width, t.height)
			} else {
				src = &image.Uniform{t.color}
			}

			draw.Draw(t.img, t.img.Bounds(), src, src.Bounds().Min, draw.Src)

			switch t.fmt {
			case "png":
				err = png.Encode(file, t.img)
			case "jpeg":
				err = jpeg.Encode(file, t.img, &jpeg.Options{100})
			default:
				err = fmt.Errorf("unsupported format: %s", t.fmt)
			}

			return
		})

		if err := g.Wait(); err != nil {
			log.Fatal(err)
		}

	}
}

func randomImage(w, h int) *image.NRGBA {
	var (
		img    = image.NewNRGBA(image.Rect(0, 0, w, h))
		bounds = img.Bounds()
	)

	rand.Seed(time.Now().Unix())
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Set(x, y, color.NRGBA{
				R: uint8(rand.Uint32()),
				G: uint8(rand.Uint32()),
				B: uint8(rand.Uint32()),
				A: uint8(rand.Uint32()),
			})
		}
	}

	return img
}
