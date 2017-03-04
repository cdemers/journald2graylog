package blacklist

import "testing"

func TestIsBlacklisted(t *testing.T) {

	list := "foo.*"
	b := PrepareBlacklist(&list)

	bytes := []byte{'f', 'o', 'o'}

	ret := b.IsBlacklisted(bytes)

	if ret != true {
		t.Error()
	}
}

func TestIsNotBlacklisted(t *testing.T) {
	list := "foo.*"
	b := PrepareBlacklist(&list)
	bytes := []byte{'b', 'a', 'r'}

	ret := b.IsBlacklisted(bytes)

	if ret != false {
		t.Error()
	}
}

func TestIsNotBlacklistedWhenEmpty(t *testing.T) {
	list := ""

	b := PrepareBlacklist(&list)
	bytes := []byte{'b', 'a', 'r'}

	ret := b.IsBlacklisted(bytes)

	if ret != false {
		t.Error()
	}
}

func TestPrepareBlacklist(t *testing.T) {

	option := "foo.*;bar.*"
	b := PrepareBlacklist(&option)

	bytes := []byte{'f', 'o', 'o'}
	foo := b.IsBlacklisted(bytes)

	if foo != true {
		t.Error()
	}

	bytes = []byte{'b', 'a', 'r'}
	bar := b.IsBlacklisted(bytes)

	if bar != true {
		t.Error()
	}

	bytes = []byte{'d', 'o', 'e'}
	doe := b.IsBlacklisted(bytes)

	if doe != false {
		t.Error()
	}
}
