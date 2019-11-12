package domaintree

import (
	"testing"
)

func BenchmarkWildcard(b *testing.B) {
	wc := NewPrefixWildcard()
	wc.AddFull("a", "a")
	wc.AddFull("b", "b")
	wc.AddFull("c", "c")
	wc.AddFull("a.b.c.com", "a.b.c.com")
	wc.AddFull("a.b.d.com", "a.b.d.com")
	wc.AddWildcard("*.com", "*.com")
	wc.AddFull("com", "com")
	wc.AddWildcard("*.example.com", "*.example.com")

	b.Run("full match", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, ok := wc.Lookup("a.b.c.com")
			if !ok {
				b.Fatal("failed")
			}
		}
	})

	b.Run("suffix match", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, ok := wc.Lookup("z.example.com")
			if !ok {
				b.Fatal("failed")
			}
		}
	})

	b.Run("apex match", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, ok := wc.Lookup("com")
			if !ok {
				b.Fatal("failed")
			}
		}
	})
}
