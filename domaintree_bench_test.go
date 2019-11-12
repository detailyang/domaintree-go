package domaintree

import (
	"log"
	"testing"
)

func BenchmarkDomainTree(b *testing.B) {
	dt := NewDomainTree()
	dt.Add("www.example.com", 1)
	dt.Add("abcd.example.com", 2)
	dt.Add("defg.example.com", 3)
	dt.Add("232.example.com", 4)
	dt.Add("*.example.com", 5)
	dt.Add("example.com", 6)
	dt.Add("*", 7)

	for _, tt := range []struct {
		input  string
		expect string
	}{
		{
			"www.example.com",
			"www.example.com",
		},
		{
			"abcd.example.com",
			"abcd.example.com",
		},
		{
			"defg.example.com",
			"defg.example.com",
		},
		{
			"232.example.com",
			"232.example.com",
		},
		{
			"example.com",
			"example.com",
		},
		{
			"11111111.example.com",
			"*.example.com",
		},
		{
			"abcd.com",
			"*",
		},
	} {
		b.Run(tt.input, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, ok := dt.Lookup(tt.input)
				if !ok {
					log.Fatal("failed to lookup")
				}
			}
		})
	}
}
