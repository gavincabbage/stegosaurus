package stegosaurus_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gavincabbage/stegosaurus"
)

func TestAlgorithm(t *testing.T) {
	normal := stegosaurus.NewAlgorithm("selection", "embedding")
	assert.Equal(t, "selection/embedding", normal)
	assert.Equal(t, "selection", normal.Selection())
	assert.Equal(t, "embedding", normal.Embedding())

	empty := stegosaurus.NewAlgorithm("", "foo")
	assert.Equal(t, "", empty)
	assert.Equal(t, "", empty.Selection())
	assert.Equal(t, "", empty.Embedding())
}
