package blacklist

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsBlacklisted(t *testing.T) {

	list := "foo.*"
	b := PrepareBlacklist(&list)

	bytes := []byte{'f', 'o', 'o'}

	ret := b.IsBlacklisted(bytes)

	assert.True(t, ret)
}

func TestIsNotBlacklisted(t *testing.T) {
	list := "foo.*"
	b := PrepareBlacklist(&list)
	bytes := []byte{'b', 'a', 'r'}

	ret := b.IsBlacklisted(bytes)

	assert.False(t, ret)
}

func TestIsNotBlacklistedWhenEmpty(t *testing.T) {
	list := ""

	b := PrepareBlacklist(&list)
	bytes := []byte{'b', 'a', 'r'}

	ret := b.IsBlacklisted(bytes)

	assert.False(t, ret)
}

func TestPrepareBlacklist(t *testing.T) {

	option := "foo.*;bar.*"
	b := PrepareBlacklist(&option)

	bytes := []byte{'f', 'o', 'o'}
	foo := b.IsBlacklisted(bytes)
	assert.True(t, foo)

	bytes = []byte{'b', 'a', 'r'}
	bar := b.IsBlacklisted(bytes)
	assert.True(t, bar)

	bytes = []byte{'d', 'o', 'e'}
	doe := b.IsBlacklisted(bytes)
	assert.False(t, doe)
}
