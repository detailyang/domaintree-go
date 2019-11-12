package domaintree

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegexDomainTree(t *testing.T) {
	dt := NewDomainTree()
	dt.Add("www.example.com", 1)
	dt.Add("abcd.example.com", 2)
	dt.Add("defg.example.com", 3)
	dt.Add("232.example.com", 4)
	dt.Add("*.example.com", 5)
	dt.Add("example.com", 6)
	dt.AddRegex(`[0-9]\.abcd\.com`, 7)

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
			"1.abcd.com",
			`[0-9]\.abcd\.com`,
		},
	} {
		dn, ok := dt.Lookup(tt.input)
		require.True(t, ok)
		require.Equal(t, tt.expect, dn.GetKey(), tt.input)
	}
}

func TestDomainTree(t *testing.T) {
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
		dn, ok := dt.Lookup(tt.input)
		require.True(t, ok)
		require.Equal(t, tt.expect, dn.GetKey(), tt.input)
	}
}
