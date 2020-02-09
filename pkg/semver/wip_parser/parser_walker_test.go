package wip_parser

import (
	assert "github.com/stretchr/testify/require"
	"testing"
)

func TestStringCandidateWalker(t *testing.T) {
	newWalker := func() iWalker {
		return stringCandidateWalkerFactory([]string{"alpha", "beta", "rc", "release"})
	}

	t.Run("alpha", func(t *testing.T) {
		walker := newWalker()

		assert.True(t, walker('a'))
		assert.True(t, walker('l'))
		assert.True(t, walker('p'))
		assert.True(t, walker('h'))
		assert.True(t, walker('a'))
		assert.False(t, walker(' '))
	})

	t.Run("beta", func(t *testing.T) {
		walker := newWalker()

		assert.True(t, walker('b'))
		assert.True(t, walker('e'))
		assert.True(t, walker('t'))
		assert.True(t, walker('a'))
		assert.False(t, walker(' '))
	})

	t.Run("rc", func(t *testing.T) {
		walker := newWalker()

		assert.True(t, walker('r'))
		assert.True(t, walker('c'))
		assert.False(t, walker(' '))
	})

	t.Run("release", func(t *testing.T) {
		walker := newWalker()

		assert.True(t, walker('r'))
		assert.True(t, walker('e'))
		assert.True(t, walker('l'))
		assert.True(t, walker('e'))
		assert.True(t, walker('a'))
		assert.True(t, walker('s'))
		assert.True(t, walker('e'))
		assert.False(t, walker(' '))
	})

	t.Run("not_valid", func(t *testing.T) {
		walker := newWalker()

		assert.True(t, walker('a'))
		//assert.False(t, walker('c'))
		//assert.False(t, walker(' '))
	})
}
