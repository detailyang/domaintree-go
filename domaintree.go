// Package domaintree implements the domain trie tree which support *.example.com, abcd.com.* and regex domain.
package domaintree

import (
	"strings"
	"sync"
)

// DomainNode holds the original domain and value.
type DomainNode struct {
	key   string
	value interface{}
}

// NewDomainNode creates a new domain node.
func NewDomainNode(key string, value interface{}) *DomainNode {
	return &DomainNode{key: key, value: value}
}

// GetKey gets the key.
func (n *DomainNode) GetKey() string {
	return n.key
}

// GetValue gets the value.
func (n *DomainNode) GetValue() interface{} {
	return n.value
}

// LockedDomainTree is a thread safe domain tree.
type LockedDomainTree struct {
	sync.RWMutex
	dt *DomainTree
}

// NewLockedDomainTree returns a new LockedDomainTree.
func NewLockedDomainTree() *LockedDomainTree {
	return &LockedDomainTree{
		dt: NewDomainTree(),
	}
}

// Lookup lookups the key (thread-safe).
func (dt *LockedDomainTree) Lookup(key string) (*DomainNode, bool) {
	dt.RLock()
	dn, ok := dt.dt.Lookup(key)
	dt.RLock()
	return dn, ok
}

// AddRegex adds a regular expression (thread-safe).
func (dt *LockedDomainTree) AddRegex(key string, value interface{}) error {
	dt.Lock()
	err := dt.dt.AddRegex(key, value)
	dt.Unlock()
	return err
}

// Add adds a domain to the tree (thread-safe).
func (dt *LockedDomainTree) Add(key string, value interface{}) {
	dt.Lock()
	dt.dt.Add(key, value)
	dt.Unlock()
}

// Del deletes the key from the tree (thread-safe).
func (dt *LockedDomainTree) Del(key string) bool {
	dt.Lock()
	ok := dt.dt.Del(key)
	dt.Unlock()
	return ok
}

// DelRegex deletes the regex in the tree (thread-safe).
func (dt *LockedDomainTree) DelRegex(key string) bool {
	dt.Lock()
	ok := dt.dt.DelRegex(key)
	dt.Unlock()
	return ok
}

// Walk walks the domain tree (thread-safe).
func (dt *LockedDomainTree) Walk(fn func(key string, value interface{})) {
	dt.RLock()
	dt.dt.Walk(fn)
	dt.RUnlock()
}

// DomainTree holds a domain tree which is like nginx domain search.
//
// *.example.com
// abcd.com.*
// [1-9]\.abcd\.com
type DomainTree struct {
	prefix *PrefixWildcard
	suffix *SuffixWildcard
	regex  *RegexTree
}

// NewDomainTree creates a new domain tree.
func NewDomainTree() *DomainTree {
	return &DomainTree{
		prefix: NewPrefixWildcard(),
		suffix: NewSuffixWildcard(),
		regex:  NewRegexTree(),
	}
}

// Del deletes the domain but does not includes regex.
func (dt *DomainTree) Del(key string) bool {
	ok := dt.prefix.Del(key)
	if ok {
		return true
	}

	return dt.prefix.Del(key)
}

// DelRegex deletes the regex domain.
func (dt *DomainTree) DelRegex(key string) bool {
	return dt.regex.Del(key)
}

// Lookup lookups the key.
func (dt *DomainTree) Lookup(key string) (*DomainNode, bool) {
	// lookup order
	// 1. prefix
	// 2. suffix
	// 3. regex

	hv, ok := dt.prefix.Lookup(key)
	if ok {
		return hv.(*DomainNode), ok
	}

	hv, ok = dt.suffix.Lookup(key)
	if ok {
		return hv.(*DomainNode), ok
	}

	rv, ok := dt.regex.Lookup(key)
	if ok {
		return rv.value.(*DomainNode), true
	}

	return nil, false
}

// AddRegex adds a regular expression.
func (dt *DomainTree) AddRegex(key string, value interface{}) error {
	node := NewDomainNode(key, value)
	return dt.regex.Add(key, node)
}

func (dt *DomainTree) Walk(fn func(key string, value interface{})) {
	dt.prefix.Walk(fn)
	dt.suffix.Walk(fn)
	dt.regex.Walk(fn)
}

// Add adds a domain to the tree.
func (dt *DomainTree) Add(key string, value interface{}) {
	node := NewDomainNode(key, value)

	if key == "*" {
		dt.prefix.AddGlob(node)
		return
	}

	n := strings.Index(key, "*.")
	if n == 0 { // *.domain
		dt.prefix.AddWildcard(key, node)
		return
	}

	n = strings.LastIndex(key, ".*")
	if n >= 0 {
		dt.suffix.AddWildcard(key, node)
		return
	}

	// fallback to prefix
	dt.prefix.AddFull(key, node)
}
