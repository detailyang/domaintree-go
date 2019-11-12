<p align="center">
  <b>
    <span style="font-size:larger;">domaintree-go</span>
  </b>
  <br />
   <a href="https://travis-ci.org/detailyang/domaintree-go"><img src="https://travis-ci.org/detailyang/domaintree-go.svg?branch=master" /></a>
   <a href="https://ci.appveyor.com/project/detailyang/domaintree-go"><img src="https://ci.appveyor.com/api/projects/status/hbpj944ankoy9sh5?svg=true" /></a>
   <br />
   <b>domaintree-go is a yet another domain tree which based on the trie tree but support the following features</b>
   <ul>
    <li>*.example.com</li>
    <li>example.com.*</li>
    <li>[0-9]+\.abcd.com</li>
   </ul>
</p>

```bash
go test -v -benchmem -run="^$" github.com/detailyang/domaintree-go -bench Benchmark
goos: darwin
goarch: amd64
pkg: github.com/detailyang/domaintree-go
BenchmarkDomainTree/www.example.com-8         	15066765	        71.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDomainTree/abcd.example.com-8        	14446810	        79.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDomainTree/defg.example.com-8        	13326453	        84.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDomainTree/232.example.com-8         	13465729	        85.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDomainTree/example.com-8             	19866852	        58.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDomainTree/11111111.example.com-8    	10720773	       110 ns/op	       0 B/op	       0 allocs/op
BenchmarkDomainTree/abcd.com-8                	22146808	        50.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkWildcard/full_match-8                	13434924	        85.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkWildcard/suffix_match-8              	15529063	        75.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkWildcard/apex_match-8                	38242728	        29.5 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/detailyang/domaintree-go	12.232s
```
