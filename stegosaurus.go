package stegosaurus

import (
	"errors"
	"io"
)

func Encode(payload, carrier io.Reader, result io.Writer) error {
	var (
		p = make([]byte, 1)
		c = make([]byte, 8)
	)

	for {
		n, err := payload.Read(p)
		if err == io.EOF {
			break
		}

		if n != len(p) {
			return errors.New("reading payload")
		} else if err != nil {
			return err
		}

		n, err = carrier.Read(c)
		if err == io.EOF {
			return errors.New("payload too large for carrier")
		}

		if n != len(p) {
			return errors.New("reading carrier")
		} else if err != nil {
			return err
		}

		for i, b := 0, p[0]; i < 8; i, b = i+1, b<<1 {
			n, err := result.Write([]byte{b | c[i]})
			if n != 1 {
				return errors.New("writing result")
			} else if err != nil {
				return err
			}
		}

	}

	return nil
}

func Decode(data io.Reader, result io.Writer) error {

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
