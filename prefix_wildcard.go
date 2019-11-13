package domaintree

import "strings"

type PrefixWildcard struct {
	glob interface{}
	wh   *WildcardHash
}

func PrefixIndexer(s, substr string) (left, right string, ok bool) {
	n := strings.LastIndex(s, substr)
	if n >= 0 {
		return s[n+1:], s[:n], true
	}
	return s, s, false
}

func NewPrefixWildcard() *PrefixWildcard {
	return &PrefixWildcard{
		wh: NewWildcardHash(PrefixIndexer),
	}
}

func (wc *PrefixWildcard) Walk(fn func(key string, value interface{})) {
	wc.wh.Walk(fn)
	if wc.glob != nil {
		fn("*", wc.glob)
	}
}

func (wc *PrefixWildcard) DelFull(key string) bool {
	return wc.wh.DelFull(key)
}

func (wc *PrefixWildcard) Lookup(key string) (interface{}, bool) {
	hv, typ := wc.wh.Lookup(key)
	if typ > NodeHashValueType {
		if typ == FullHashValueType {
			return hv.fullvalue, true
		}
		return hv.wildcardvalue, true
	}

	if wc.glob != nil {
		return wc.glob, true
	}

	return nil, false
}

func (wc *PrefixWildcard) Del(key string) bool {
	if key == "*" {
		wc.DelGlob()
		return true
	}

	first := strings.Index(key, "*")
	if first == -1 {
		return wc.DelFull(key)
	}

	return wc.DelWildcard(key)
}

func (wc *PrefixWildcard) DelGlob() {
	wc.glob = nil
}

func (wc *PrefixWildcard) AddGlob(value interface{}) {
	wc.glob = value
}

// Add adds the key to the trie tree.
func (wc *PrefixWildcard) Add(key string, value interface{}) {
	if key == "*" {
		wc.AddGlob(value)
		return
	}

	first := strings.Index(key, "*")
	if first == -1 {
		wc.AddFull(key, value)
		return
	}

	wc.AddWildcard(key, value)
}

// AddFull adds the key to the trie tree.
func (wc *PrefixWildcard) AddFull(key string, value interface{}) {
	wc.wh.add(key, value, FullHashValueType)
}

// AddWildcard adds the suffix match like "*.abcd.com".
func (wc *PrefixWildcard) AddWildcard(key string, value interface{}) {
	n := strings.Index(key, "*.")
	if n >= 0 {
		key = key[n+2:]
		wc.wh.add(key, value, WildcardHashValueType)
		return
	}
	wc.AddFull(key, value)
}

// DelWildcard deletes the wildcard match.
func (wc *PrefixWildcard) DelWildcard(key string) bool {
	n := strings.Index(key, "*.")
	if n >= 0 {
		key = key[n+2:]
	}
	return wc.wh.delWildcard(key)
}

func (wc *PrefixWildcard) String() string {
	return wc.wh.String()
}
