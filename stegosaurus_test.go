package stegosaurus_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/gavincabbage/stegosaurus"
	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	var (
		payload  = bytes.NewReader([]byte{0x55}) // 0101 0101
		carrier  = bytes.NewReader(make([]byte, 8))
		result   = bytes.NewBuffer(make([]byte, 0, 8))
		expected = []byte{0x00, 0x80, 0x00, 0x80, 0x00, 0x80, 0x00, 0x80}
	)

	require.NoError(t, stegosaurus.Encode(payload, carrier, result))
	actual, _ := ioutil.ReadAll(result)
	require.Equal(t, expected, actual)
}
