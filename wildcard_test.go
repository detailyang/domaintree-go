package domaintree

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimpleSuffixWildcard(t *testing.T) {
	wc := NewSuffixWildcard()
	wc.AddFull("a", "a")
	wc.AddFull("b", "b")
	wc.AddFull("c", "c")
	wc.AddFull("a.b.c.com", "a.b.c.com")
	wc.AddFull("a.b.d.com", "a.b.d.com")
	wc.AddWildcard("abcd.com.*", "abcd.com.*")
	wc.AddFull("com", "com")
	wc.AddWildcard("defg.com.*", "defg.com.*")
	// wc.Walk(func(domain string, value *HashValue) {
	// 	fmt.Println(domain, value.typ.String())
	// })
	// fmt.Println(wc.String())

	hv, ok := wc.Lookup("a.b.c.com")
	require.True(t, ok)
	require.Equal(t, "a.b.c.com", hv)

	hv, ok = wc.Lookup("abcd.com.a.b")
	require.True(t, ok)
	require.Equal(t, "abcd.com.*", hv)

	hv, ok = wc.Lookup("abcd.com.a.b.c")
	require.True(t, ok)
	require.Equal(t, "abcd.com.*", hv)

	hv, ok = wc.Lookup("defg.com.a.b.c")
	require.True(t, ok)
	require.Equal(t, "defg.com.*", hv)
}

func TestSimplePrefixWildcard(t *testing.T) {
	wc := NewPrefixWildcard()
	wc.AddFull("a", "a")
	wc.AddFull("b", "b")
	wc.AddFull("c", "c")
	wc.AddFull("a.b.c.com", "a.b.c.com")
	wc.AddFull("a.b.d.com", "a.b.d.com")
	wc.AddWildcard("*.com", "*.com")
	wc.AddFull("com", "com")
	wc.AddWildcard("*.example.com", "*.example.com")
	// wc.Walk(func(domain string, value *HashValue) {
	// 	fmt.Println(domain, value.typ.String())
	// })
	// fmt.Println(wc.String())

	hv, ok := wc.Lookup("a.b.c.com")
	require.True(t, ok)
	require.Equal(t, "a.b.c.com", hv)

	hv, ok = wc.Lookup("a")
	require.True(t, ok)
	require.Equal(t, "a", hv)

	hv, ok = wc.Lookup("z.example.com")
	require.True(t, ok)
	require.Equal(t, "*.example.com", hv)

	hv, ok = wc.Lookup("abcd.com")
	require.True(t, ok)
	require.Equal(t, "*.com", hv)

	hv, ok = wc.Lookup("213797123.abcd-2323.com")
	require.True(t, ok)
	require.Equal(t, "*.com", hv)

	hv, ok = wc.Lookup("a.b.d.com")
	require.True(t, ok)
	require.Equal(t, "a.b.d.com", hv)

	wc.DelFull("a.b.d.com")
	// fmt.Println(wc.String())

	hv, ok = wc.Lookup("a.b.d.com")
	require.True(t, ok)
	require.Equal(t, "*.com", hv)

	wc.DelWildcard("*.com")
	hv, ok = wc.Lookup("a.b.d.com")
	require.False(t, ok)

	hv, ok = wc.Lookup("a.b.c.com")
	require.True(t, ok)
	require.Equal(t, "a.b.c.com", hv)
}
