package lsb

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
)

const mask = 0x01

type Encoder struct{}

func (_ Encoder) Encode(payload, carrier io.Reader, result io.Writer) error {
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
				return errors.New("payload too large for carrier")
			} else if err != nil {
				return err
			}

			return errors.New("reading carrier")
		}

		for i, b := 8, p[0]; i > 0; i, b = i-1, b>>1 {
			r[i-1] = b&mask | c[i-1]&^mask
		}

		_, err = result.Write(r)
		if err != nil {
			return fmt.Errorf("writing result: %w", err)
		}
	}

	b, err := ioutil.ReadAll(carrier)
	if err != nil {
		return fmt.Errorf("reading remaining carrier: %w", err)
	}

	_, err = result.Write(b)
	if err != nil {
		return fmt.Errorf("writing remaining carrier: %w", err)
	}

	return nil
}

func (_ Encoder) Decode(data io.Reader, result io.Writer) error {
	var (
		d = make([]byte, 8)
	)
	for {
		n, err := data.Read(d)
		if n < 8 {
			if err != nil && err != io.EOF {
				return fmt.Errorf("reading data: %w", err)
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
			return fmt.Errorf("writing result: %w", err)
		}
	}

	return nil
}
