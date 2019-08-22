package stegosaurus

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
)

const mask = 0x01

func Encode(payload, carrier io.Reader, result io.Writer) error {
	var (
		p = make([]byte, 1)
		c = make([]byte, 8)
		r = make([]byte, 8)
	)
	for {
		n, err := payload.Read(p)
		if n < 1 {
			if err != io.EOF {
				return err
			}

			break
		}

		n, err = carrier.Read(c)
		if n < 8 {
			if err == io.EOF {
				return errors.New("stegosaurus: encode: payload too large for carrier")
			} else if err != nil {
				return err
			}

			return errors.New("stegosaurus: encode: reading carrier")
		}

		for i, b := 8, p[0]; i > 0; i, b = i-1, b>>1 {
			r[i-1] = b&mask | c[i-1]
		}

		_, err = result.Write(r)
		if err != nil {
			return fmt.Errorf("stegosaurus: encode: writing result: %w", err)
		}
	}

	b, err := ioutil.ReadAll(carrier)
	if err != nil {
		return fmt.Errorf("stegosaurus: encode: reading remaining carrier: %w", err)
	}

	_, err = result.Write(b)
	if err != nil {
		return fmt.Errorf("stegosaurus: encode: writing remaining carrier: %w", err)
	}

	return nil
}

func Decode(data io.Reader, result io.Writer) error {
	var (
		d = make([]byte, 8)
	)
	for {
		n, err := data.Read(d)
		if n < 8 {
			if err != io.EOF {
				return fmt.Errorf("stegosaurus: decode: reading data: %w", err)
			}

			break
		}

		var b byte
		for i := 0; i < 8; i++ {
			b = d[i]&mask | b
			if i < 7 {
				b = b << 1
			}
		}

		_, err = result.Write([]byte{b})
		if err != nil {
			return fmt.Errorf("stegosaurus: decode: writing result: %w", err)
		}
	}

	return nil
}

/*

  LSBSteg.py encode -i <input> -o <output> -f <file>
  LSBSteg.py decode -i <input> -o <output>

--

The payload is the data covertly communicated. T

he carrier is the signal, stream, or data file that hides the payload, which differs from
the channel, which typically means the type of input, such as a JPEG image.

The resulting signal, stream, or data file with the encoded payload is sometimes called the package, stego file, or covert message.

The proportion of bytes, samples, or other signal elements modified to encode the payload is called
the encoding density and is typically expressed as a number between 0 and 1.

*/
