package image

import (
	"io"

	"github.com/gavincabbage/stegosaurus/lsb"

	"github.com/gavincabbage/stegosaurus"
)

type Encoder struct {
	alg stegosaurus.Algorithm
	key []byte
}

func NewEncoder(algorithm stegosaurus.Algorithm, key []byte) Encoder {
	return Encoder{
		alg: algorithm,
		key: key,
	}
}

func (e Encoder) Encode(payload, carrier io.Reader, data io.Writer) error {
	img, err := NewImage(carrier)
	if err != nil {
		return err
	}
	encoder := lsb.Encoder{}

	if err := encoder.Encode(payload, img, img); err != nil {
		return err
	}

	return img.Encode(data)
}

func (e Encoder) Decode(data io.Reader, payload io.Writer) error {
	img, err := NewImage(data)
	if err != nil {
		return err
	}
	encoder := lsb.Encoder{}

	return encoder.Decode(img, payload)
}
