package multireader

import (
	"errors"
	"io"
)

type MultiReader struct {
	readers []reader
}

type reader struct {
	io.Reader
	n int
}

func New() *MultiReader {
	return &MultiReader{}
}

func (m *MultiReader) Add(r io.Reader, n int) {
	m.readers = append(m.readers, reader{r, n})
}

func (m *MultiReader) Read() ([][]byte, error) {
	var (
		read [][]byte
		errs []error
	)
	for _, r := range m.readers {
		s := make([])
	}
}
