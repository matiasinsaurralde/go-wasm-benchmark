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

This is what the output looks like 

```
goos: darwin
goarch: amd64
BenchmarkWasmerSum-8            	   10000	    224534 ns/op	     728 B/op	      19 allocs/op
BenchmarkWasmerSumReentrant-8   	  500000	      3506 ns/op	     144 B/op	       5 allocs/op
BenchmarkWagonSum-8             	    5000	    253306 ns/op	 1094268 B/op	     536 allocs/op
BenchmarkWagonSumReentrant-8    	20000000	       118 ns/op	      20 B/op	       2 allocs/op
BenchmarkLifeSum-8              	    2000	    671554 ns/op	 1145509 B/op	     375 allocs/op
PASS
ok  	_/Users/matias/dev/wasm/wasm-benchmark	9.297s
```