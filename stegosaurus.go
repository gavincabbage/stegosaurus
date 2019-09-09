package stegosaurus

import (
	"io"
)

// Encoder encodes and decodes secret payloads.
type Encoder interface {
	// Encode the secret payload into the carrier data.
	Encode(payload, carrier io.Reader, data io.Writer) error
	// Decode the secret payload from the carrier data.
	Decode(data io.Reader, payload io.Writer) error
}
