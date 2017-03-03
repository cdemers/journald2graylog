package blacklist

import (
	"regexp"
	"strings"
)

type Blacklist struct {
	regexp []*regexp.Regexp
}

//IsBlacklisted Return true if any regex in the regex array match this line
func (b *Blacklist) IsBlacklisted(line []byte) bool {
	for _, r := range b.regexp {
		if r.Match(line) {
			return true
		}
	}
	return false
}

//PrepareBlacklist Parse the string using ; separator and return the Blacklist struct
func PrepareBlacklist(blacklist *string) Blacklist {

	b := Blacklist{}

	split := strings.Split(*blacklist, ";")
	for _, r := range split {
		if len(r) > 0 {
			rexp := regexp.MustCompile(r)
			b.regexp = append(b.regexp, rexp)
		}
	}
	return b
}
