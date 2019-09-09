package stegosaurus

import (
	"fmt"
	"strings"
)

// Algorithm represents tuple of byte selection and bit embedding algorithms.
//
// The selection algorithm determines how bytes in the carrier data are chosen
// for embedding. The bit embedding algorithm determines how payload bits are
// embedded into selected carrier bytes.
//
// An Algorithm is a string of the form `selection/embedding` where selection
// represents the selection algorithm and embedding the bit embedding algorithm,
// delimited by a forward slash.
//
// Available algorithms are determined by the Encoder implementation.
type Algorithm string

// NewAlgorithm formats a selection and embedding algorithm string.
func NewAlgorithm(selection, embedding string) Algorithm {
	if selection == "" || embedding == "" {
		return ""
	}

	return Algorithm(fmt.Sprintf("%s/%s", selection, embedding))
}

// Selection returns the Algorithm's selection algorithm.
func (a Algorithm) Selection() string {
	return a.parsed(0)
}

// Embedding returns the Algorithm's embedding algorithm.
func (a Algorithm) Embedding() string {
	return a.parsed(1)
}

func (a Algorithm) parsed(i int) string {
	s := strings.Split(string(a), "/")
	if len(s) != 2 {
		return ""
	}

	return s[i]
}
