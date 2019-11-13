package domaintree

import "strings"

type SuffixWildcard struct {
	wh *WildcardHash
}

func SuffixIndexer(s, substr string) (left, right string, ok bool) {
	n := strings.Index(s, substr)
	if n >= 0 {
		return s[:n], s[n+1:], true
	}
	return s, s, false
}

func NewSuffixWildcard() *SuffixWildcard {
	return &SuffixWildcard{
		wh: NewWildcardHash(SuffixIndexer),
	}
}

func (wc *SuffixWildcard) Del(key string) bool {
	first := strings.Index(key, "*")
	if first == -1 {
		return wc.DelFull(key)
	}

	return wc.DelWildcard(key)
}

func (wc *SuffixWildcard) DelFull(key string) bool {
	return wc.wh.DelFull(key)
}

// Walk walks the suffix tree.
func (wc *SuffixWildcard) Walk(fn func(key string, value interface{})) {
	wc.wh.Walk(fn)
}

func (wc *SuffixWildcard) Lookup(key string) (interface{}, bool) {
	hv, typ := wc.wh.Lookup(key)
	if typ > NodeHashValueType {
		if typ == FullHashValueType {
			return hv.fullvalue, true
		}
		return hv.wildcardvalue, true
	}

	return nil, false
}

// Add adds the key to the trie tree.
func (wc *SuffixWildcard) Add(key string, value interface{}) {
	first := strings.LastIndex(key, "*")
	if first == 0 {
		wc.AddFull(key, value)
		return
	}

	wc.AddWildcard(key, value)
}

// AddFull adds the key to the trie tree.
func (wc *SuffixWildcard) AddFull(key string, value interface{}) {
	wc.wh.add(key, value, FullHashValueType)
}

// AddWildcard adds the suffix match like "*.abcd.com".
func (wc *SuffixWildcard) AddWildcard(key string, value interface{}) {
	n := strings.LastIndex(key, ".*")
	if n >= 0 {
		key = key[:n]
		wc.wh.add(key, value, WildcardHashValueType)
		return
	}
	wc.AddFull(key, value)
}

// DelWildcard deletes the wildcard match.
func (wc *SuffixWildcard) DelWildcard(key string) bool {
	n := strings.LastIndex(key, ".*")
	if n >= 0 {
		key = key[:n]
	}
	return wc.wh.delWildcard(key)
}

func (wc *SuffixWildcard) String() string {
	return wc.wh.String()
}
