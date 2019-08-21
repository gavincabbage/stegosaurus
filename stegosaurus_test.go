package stegosaurus_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/gavincabbage/stegosaurus"
	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	var cases = []struct {
		name     string
		payload  io.Reader
		carrier  io.Reader
		result   io.ReadWriter
		expected []byte
	}{
		{
			name:     "happy path",
			payload:  bytes.NewReader([]byte{0xde}),
			carrier:  bytes.NewReader(make([]byte, 8)),
			result:   bytes.NewBuffer(make([]byte, 0, 8)),
			expected: []byte{0x01, 0x01, 0x00, 0x01, 0x01, 0x01, 0x01, 0x00},
		},
		{
			name:     "happy path non-zero carrier",
			payload:  bytes.NewReader([]byte{0xde}),
			carrier:  bytes.NewReader([]byte{0x41, 0xff, 0x00, 0x01, 0x01, 0x01, 0x01, 0x00}),
			result:   bytes.NewBuffer(make([]byte, 0, 8)),
			expected: []byte{0x41, 0xff, 0x00, 0x01, 0x01, 0x01, 0x01, 0x00},
		},
		{
			name:    "happy path extra carrier",
			payload: bytes.NewReader([]byte{0xde}),
			carrier: bytes.NewReader([]byte{
				0x01, 0x01, 0x00, 0x01, 0x01, 0x01, 0x01, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x23,
			}),
			result: bytes.NewBuffer(make([]byte, 0, 16)),
			expected: []byte{
				0x01, 0x01, 0x00, 0x01, 0x01, 0x01, 0x01, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x23,
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			require.NoError(t, stegosaurus.Encode(test.payload, test.carrier, test.result))
			actual, _ := ioutil.ReadAll(test.result)
			require.Equal(t, test.expected, actual)
		})
	}
}

func TestDecode(t *testing.T) {
	var cases = []struct {
		name     string
		data     io.Reader
		result   io.ReadWriter
		expected []byte
	}{
		{
			name:     "happy path",
			data:     bytes.NewReader([]byte{0x01, 0x01, 0x00, 0x01, 0x01, 0x01, 0x01, 0x00}),
			result:   bytes.NewBuffer(make([]byte, 0, 8)),
			expected: []byte{0xde},
		},
		{
			name: "happy path extra carrier",
			data: bytes.NewReader([]byte{
				0x01, 0x01, 0x00, 0x01, 0x01, 0x01, 0x01, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x23,
			}),
			result:   bytes.NewBuffer(make([]byte, 0, 16)),
			expected: []byte{0xde, 0x03},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			require.NoError(t, stegosaurus.Decode(test.data, test.result))
			actual, _ := ioutil.ReadAll(test.result)
			require.Equal(t, test.expected, actual)
		})
	}
}
