go-wasm-benchmark
==

A very simple benchmark to try out WASM VMs in Go. Only a simple `sum` function is used for now. I might extend this in the future.

To run:

```
$ go test -bench=. -benchmem 
```

## Motivation and limitations

I've been doing a few experiments with WASM and Go, the technology has great potential but features like WASI support aren't available everywhere (yet).

Some of the limitations I've found:
- Not all VMs are reentrant (see `life` for example).
- I've been working with Go interop projects for three years and `cgo` is a big limitation when looking at performance. From my perspective we should focus on enhancing the performance of pure-Go VMs. `wagon` and `life` are the best candidates here.

## Results

*Hardware:* Macbook Pro Mid 2014, i7, 16 GB RAM.

This is what the output looks like 

```
goos: darwin
goarch: amd64
pkg: github.com/matiasinsaurralde/go-wasm-benchmark
BenchmarkWasmerSum-8            	   10000	    213710 ns/op	     727 B/op	      19 allocs/op
BenchmarkWasmerSumReentrant-8   	  500000	      3101 ns/op	     144 B/op	       5 allocs/op
BenchmarkWagonSum-8             	   10000	    432710 ns/op	 1093881 B/op	     535 allocs/op
BenchmarkWagonSumReentrant-8    	 3000000	       538 ns/op	     152 B/op	      10 allocs/op
BenchmarkLifeSum-8              	    2000	    616651 ns/op	 1145511 B/op	     375 allocs/op
BenchmarkWASM3Sum-8             	    2000	    981291 ns/op	     351 B/op	      14 allocs/op
BenchmarkWASM3SumReentrant-8    	 2000000	       891 ns/op	     128 B/op	       7 allocs/op
PASS
ok  	github.com/matiasinsaurralde/go-wasm-benchmark	17.329s
```