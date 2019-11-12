package domaintree

import (
	"errors"
	"regexp"
)

type regexValue struct {
	key   string
	regex *regexp.Regexp
	value interface{}
}

// RegexTree represents a regular expression tree.
type RegexTree struct {
	regex []*regexValue
}

// NewRegexTree creates a new regex tree.
func NewRegexTree() *RegexTree {
	return &RegexTree{}
}

// Walk walks the regex tree.
func (rt *RegexTree) Walk(fn func(key string, value interface{})) {
	for i := range rt.regex {
		fn(rt.regex[i].key, rt.regex[i].value)
	}
}

// Lookup lookups the key in the regex tree.
func (rt *RegexTree) Lookup(key string) (*regexValue, bool) {
	for i := range rt.regex {
		if rt.regex[i].regex.MatchString(key) {
			return rt.regex[i], true
		}
	}
	return nil, false
}

// Add adds a regular expression.
func (rt *RegexTree) Add(key string, value interface{}) error {
	// Compile(expr string) (*Regexp, error)
	for i := range rt.regex {
		if rt.regex[i].key == key {
			return errors.New("duplicated key")
		}
	}

	rex, err := regexp.Compile(key)
	if err != nil {
		return err
	}

	rt.regex = append(rt.regex, &regexValue{
		key:   key,
		regex: rex,
		value: value,
	})

	return nil
}
