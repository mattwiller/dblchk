# dblchk

[![GoDoc](https://godoc.org/github.com/mattwiller/dblchk?status.svg)](https://godoc.org/github.com/mattwiller/dblchk)

Zero-dependency implementation of a [two-bit, single-hash function Bloom filter][2bit-bloom]
in Go, achieving significantly better accuracy than a traditional Bloom filter
of the same size with excellent performance.

[2bit-bloom]: https://floedb.ai/blog/two-bits-are-better-than-one-making-bloom-filters-2x-more-accurate

## Installation

```go get -u github.com/mattwiller/dblchk```

## Usage

```go
import "github.com/mattwiller/dblchk"

// Create filter with default size (256 KiB): good for up to
// a few hundred of thousand entries with few false positives (< 5-10%)
filter := dblchk.NewFilter(0) 

// Add an item to the filter
filter.Add([]byte("hello, world!"))

// Check whether an item may be included in the filter
if filter.MayContain(item) {
  // Item is (probably) included in the set approximated by the filter
}
```

## Benchmarks

Tested against popular [Bloom filter][bloom] and [Cuckoo filter][cuckoo] implementations:

```
goos: linux
goarch: amd64
cpu: 12th Gen Intel(R) Core(TM) i5-12400T

              │    dblchk    │                bloom                 │                  cuckoo                   │
              │    sec/op    │    sec/op     vs base                │     sec/op      vs base                   │
AddElement-12   20.45n ± 16%   26.55n ±  0%  +29.80% (p=0.001 n=10)   5761.50n ±  0%  +28073.59% (p=0.000 n=10)
MayContain-12   21.57n ± 10%   22.72n ±  0%   +5.38% (p=0.022 n=10)     11.99n ±  2%     -44.42% (p=0.000 n=10)
```

[bloom]: https://github.com/bits-and-blooms/bloom
[cuckoo]: https://github.com/seiflotfy/cuckoofilter

### False-Positive Rate

| Library | Size | Number of Entries | False Positive Rate |
| ------- | ---- | ----------------- | ------------------- |
| `github.com/mattwiller/dblchk` | 256 KiB | 100k | 1.3% |
| `github.com/bits-and-blooms/bloom` | 256 KiB | 100k | 4.7% |
| `github.com/seiflotfy/cuckoofilter` | 256 KiB | 100k | 1.2% |


## License

Copyright 2026 Matt Willer. Licensed under the [3-Clause BSD License][bsd-3-clause].

[bsd-3-clause]: https://opensource.org/license/bsd-3-clause
