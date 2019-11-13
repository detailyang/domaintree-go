package domaintree

import (
	"bytes"
	"fmt"
	"io"
)

type StringIndexer func(s, substr string) (left, right string, ok bool)

// HashValueType represents the type of the hashvalue.
type HashValueType uint8

var (
	NodeHashValueType     HashValueType = 0x00
	FullHashValueType     HashValueType = 0x01
	WildcardHashValueType HashValueType = 0x02
)

func (hvt HashValueType) String() string {
	switch hvt {
	case NodeHashValueType:
		return "."
	case FullHashValueType:
		return "="
	case WildcardHashValueType:
		return "*"
	case FullHashValueType | WildcardHashValueType:
		return "=*"
	}
	return "unknown"
}

// HashValue returns the hashvalue.
type HashValue struct {
	typ           HashValueType
	fullvalue     interface{}
	wildcardvalue interface{}
	hash          *WildcardHash
}

// GetValue returns the value.
func (hv *HashValue) GetValue() interface{} {
	if hv.typ&FullHashValueType == FullHashValueType {
		return hv.fullvalue
	}

	if hv.typ&WildcardHashValueType == WildcardHashValueType {
		return hv.wildcardvalue
	}

	return nil
}

// GetFullValue gets the full value.
func (hv *HashValue) GetFullValue() interface{} { return hv.fullvalue }

// GetWildcardValue gets the wildcard value.
func (hv *HashValue) GetWildcardValue() interface{} { return hv.wildcardvalue }

// GetType returns the type.
func (hv *HashValue) GetType() HashValueType { return hv.typ }

// String returns the string representation.
func (hv *HashValue) String() string {
	switch hv.typ {
	case NodeHashValueType:
		return hv.typ.String()
	case FullHashValueType:
		return fmt.Sprintf("%s[%+v]", hv.typ, hv.fullvalue)
	case WildcardHashValueType:
		return fmt.Sprintf("%s[%+v]", hv.typ, hv.wildcardvalue)
	case FullHashValueType | WildcardHashValueType:
		return fmt.Sprintf("%s[%+v-%v]", hv.typ, hv.fullvalue, hv.wildcardvalue)
	}
	return "unknown"
}

// WildcardHash represents the trie tree which support prefix wildcard
//
// *.example.com
// example.com
// abcd.example.com
type WildcardHash struct {
	glob    interface{}
	indexer StringIndexer
	hash    map[string]*HashValue
}

// NewWildcardHash returns a new WildcardHash.
func NewWildcardHash(indexer StringIndexer) *WildcardHash {
	return &WildcardHash{
		indexer: indexer,
		hash:    make(map[string]*HashValue, 4),
	}
}

func (wc *WildcardHash) delWildcard(key string) bool {
	sub, remaining, success := wc.indexer(key, ".")

	hv, ok := wc.hash[sub]
	if ok {
		if success {
			hv.hash.delWildcard(remaining)
			if hv.hash.Len() == 0 {
				// TODO(detailyang): cleanup self node
			}
			return true
		}

		if hv.typ&WildcardHashValueType == WildcardHashValueType {
			hv.wildcardvalue = nil
			hv.typ ^= WildcardHashValueType
		}

		if hv.typ&FullHashValueType == FullHashValueType {
			hv.fullvalue = nil
		}

		if hv.typ == NodeHashValueType {
			delete(wc.hash, sub)
		}
		return true
	}

	return false
}

// DelFull deletes the full match.
func (wc *WildcardHash) DelFull(key string) bool {
	sub, remaining, success := wc.indexer(key, ".")

	hv, ok := wc.hash[sub]
	if ok {
		if success {
			hv.hash.DelFull(remaining)
			if hv.hash.Len() == 0 {
				// TODO(detailyang): cleanup self node
			}
			return true
		}

		if hv.typ&FullHashValueType == FullHashValueType {
			hv.fullvalue = nil
			hv.typ ^= FullHashValueType
		}

		if hv.typ&WildcardHashValueType == WildcardHashValueType {
			hv.wildcardvalue = nil
		}

		if hv.typ == NodeHashValueType {
			delete(wc.hash, sub)
		}
		return true
	}

	return false
}

func (wc *WildcardHash) add(key string, value interface{}, typ HashValueType) {
	sub, remaining, success := wc.indexer(key, ".")

	hv, ok := wc.hash[sub]
	if ok {
		if success {
			hv.hash.add(remaining, value, typ)
			return
		}

		hv.typ |= typ
		if typ == WildcardHashValueType {
			hv.wildcardvalue = value
		} else if typ == FullHashValueType {
			hv.fullvalue = value
		}

		return
	}

	wch := NewWildcardHash(wc.indexer)
	nhv := &HashValue{ // intermediate node
		hash: wch,
	}

	wc.hash[sub] = nhv
	if !success {
		nhv.typ |= typ
		if typ == WildcardHashValueType {
			nhv.wildcardvalue = value
		} else if typ == FullHashValueType {
			nhv.fullvalue = value
		}
		return
	}

	wch.add(remaining, value, typ)
}

// Len returns the length of the underlying hash.
func (wc *WildcardHash) Len() int {
	return len(wc.hash)
}

func (wc *WildcardHash) String() string {
	var buf bytes.Buffer
	wc.pretty(&buf, "")
	return buf.String()
}

func (wc *WildcardHash) pretty(w io.Writer, prefix string) {
	if prefix != "" {
		prefix = prefix + "."
	}
	for k := range wc.hash {
		v := wc.hash[k]
		if v.typ > NodeHashValueType {
			if v.hash.Len() == 0 {
				fmt.Fprintf(w, "%s[%s]\n", prefix+k, v)
				continue
			}
			if v.typ > NodeHashValueType {
				fmt.Fprintf(w, "%s[%s]\n", prefix+k, v)
			}
		}
		v.hash.pretty(w, prefix+k)
	}
}

// Walk walks the tree recursively.
func (wc *WildcardHash) Walk(fn func(key string, value interface{})) {
	wc.walk("", fn)
}

func (wc *WildcardHash) walk(prefix string, fn func(key string, value interface{})) {
	if prefix != "" {
		prefix = prefix + "."
	}
	for k := range wc.hash {
		v := wc.hash[k]
		if v.typ > NodeHashValueType {
			if v.hash.Len() == 0 {
				if v.fullvalue != nil {
					fn(prefix+k, v.fullvalue)
				}
				if v.wildcardvalue != nil {
					fn(prefix+k, v.wildcardvalue)
				}
				continue
			}
			if v.typ > NodeHashValueType {
				if v.fullvalue != nil {
					fn(prefix+k, v.fullvalue)
				}
				if v.wildcardvalue != nil {
					fn(prefix+k, v.wildcardvalue)
				}
			}
		}
		v.hash.walk(prefix+k, fn)
	}
}

// Lookup lookups the key in trie tree.
func (wc *WildcardHash) Lookup(key string) (*HashValue, HashValueType) {
	sub, remaining, success := wc.indexer(key, ".")

	hash, ok := wc.hash[sub]
	if !ok {
		return nil, NodeHashValueType
	}

	if hash.typ == NodeHashValueType { // intermediate layer
		if !success {
			return nil, NodeHashValueType
		}
		return hash.hash.Lookup(remaining)
	}

	if hash.typ&FullHashValueType == FullHashValueType {
		if !success {
			return hash, FullHashValueType
		}

		subh, typ := hash.hash.Lookup(remaining)
		if typ > NodeHashValueType {
			return subh, typ
		}

		// continue see if it's wildcard
	}

	if hash.typ&WildcardHashValueType == WildcardHashValueType {
		if !success {
			return hash, WildcardHashValueType
		}
		subh, typ := hash.hash.Lookup(remaining)
		if typ > NodeHashValueType {
			return subh, typ
		}

		return hash, WildcardHashValueType
	}

	return nil, NodeHashValueType
}
