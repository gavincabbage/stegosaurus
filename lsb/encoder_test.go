package lsb_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/gavincabbage/stegosaurus/lsb"
	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	var cases = []struct {
		name     string
		payload  io.Reader
		carrier  io.Reader
		expected []byte
		err      string
	}{
		{
			name:     "happy path",
			payload:  bytes.NewReader([]byte{0b11011110}),
			carrier:  bytes.NewReader(make([]byte, 8)),
			expected: []byte{0x01, 0x01, 0x00, 0x01, 0x01, 0x01, 0x01, 0x00},
		},
		{
			name:    "happy path two bytes",
			payload: bytes.NewReader([]byte{0b11011110, 0b01010101}),
			carrier: bytes.NewReader(make([]byte, 16)),
			expected: []byte{
				0x01, 0x01, 0x00, 0x01, 0x01, 0x01, 0x01, 0x00,
				0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01,
			},
		},
		{
			name:     "happy path non-zero carrier",
			payload:  bytes.NewReader([]byte{0b11011110}),
			carrier:  bytes.NewReader([]byte{0x41, 0xff, 0xff, 0x01, 0x01, 0x01, 0x01, 0x00}),
			expected: []byte{0x41, 0xff, 0xfe, 0x01, 0x01, 0x01, 0x01, 0x00},
		},
		{
			name:    "happy path extra carrier",
			payload: bytes.NewReader([]byte{0b11011110}),
			carrier: bytes.NewReader([]byte{
				0x01, 0x01, 0x00, 0x01, 0x01, 0x01, 0x01, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x23,
			}),
			expected: []byte{
				0x01, 0x01, 0x00, 0x01, 0x01, 0x01, 0x01, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x23,
			},
		},
		{
			name:     "error payload too large for carrier",
			payload:  bytes.NewReader([]byte{0b11011110, 0xdd}),
			carrier:  bytes.NewReader(make([]byte, 8)),
			expected: []byte{0x01, 0x01, 0x00, 0x01, 0x01, 0x01, 0x01, 0x00},
			err:      "payload too large for carrier",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			var (
				result  = bytes.NewBuffer([]byte{})
				subject = lsb.Encoder{}
			)
			err := subject.Encode(test.payload, test.carrier, result)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
			}
			actual, _ := ioutil.ReadAll(result)
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
			expected: []byte{0b11011110},
		},
		//{
		//	name: "happy path extra carrier",
		//	data: bytes.NewReader([]byte{
		//		0x01, 0x01, 0x00, 0x01, 0x01, 0x01, 0x01, 0x00,
		//		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x23,
		//	}),
		//	expected: []byte{0b11011110, 0x03},
		//},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			var (
				result  = bytes.NewBuffer([]byte{})
				subject = lsb.Encoder{}
			)
			require.NoError(t, subject.Decode(test.data, result))
			actual, _ := ioutil.ReadAll(result)
			require.Equal(t, test.expected, actual)
		})
	}
}

func TestEncodeDecode(t *testing.T) {
	var (
		subject = lsb.Encoder{}
		result  = bytes.NewBuffer([]byte{})
		payload = "abcdefgh"
		carrier = "1234567812345678123456781234567812345678123456781234567812345678"
		decoded = bytes.NewBuffer([]byte{})
	)
	require.NoError(t, subject.Encode(strings.NewReader(payload), strings.NewReader(carrier), result))
	require.NoError(t, subject.Decode(result, decoded))
	require.Equal(t, payload, decoded.String())
}
