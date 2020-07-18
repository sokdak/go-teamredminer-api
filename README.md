# CGMiner API for Go #

This package is a fork or [crypt0train/go-cgminer-api](https://github.com/crypt0train/go-cgminer-api)
which wasn't updated since 2018.

This repo contains some fixes and improvements like:

* Go module support
* [JSON number literal fix](https://github.com/golang/go/issues/34472) for Go 1.14
* Support for alternative request formats (cgminer has also text-based format)
* Context-based requests (like `CGMiner.SummaryContext()`, etc)
* and more!
 
 
## Installation ##

    # install the library:
    go get github.com/x1unix/go-cgminer-api

    // Use in your .go code:
    import (
        cgminer "github.com/x1unix/go-cgminer-api"
    )

## API Documentation ##

I started to completely rewrite forked code.

At this moment fully tested commands: `Version(), Summary(), Stats(), GetPools()`.
You can use `Call()` or `CallContext()` to run any command, that you want.
Version and Summary sections have the same structure over the all devices, but Stats - it is something like hell.

Cgminer/Bmminer - is a best example of very bad JSON api and very bad code at all(it even has `sprintf` buffer overflow in "check" for a years).
When you try get stats from cgminer via api, for a first time you can't parse the answer, because it's not valid and some portion of response of `version` command is mixed to the answer. Yep, fix this, but json response still invalid, float and integer values some times are presented as strings, sometimes you got "" instead of null value. 

Command "stats" returned different output for each device, that i have tested S9, D3, L3+, T9+, even R1-LTC router. 
I wrote test for all of this, you can find textures in **testdata** folder.

My API method Stats() will return generic structure, that have all fields from all devices.
I've created helpers to find data that depends on model, see example below.

Test coverage now: 66.4% of statements

## Quickstart ##

```go
package main

import (
    cgminer "github.com/x1unix/go-cgminer-api"
    "time"
    "log"
    "fmt"
)

func main() {
    miner := cgminer.NewCGMiner("localhost", 4028, 2 * time.Second)
	stats, err := miner.Stats()
	if err != nil {
		log.Println(err)
	}
    fmt.Printf("%s | GHS avg: %0.2f\n", stats.Type, stats.GhsAverage)
    // Get generic stats
    genericStats := stats.Generic()
 
    // When you connected to Antminer S9
    statsS9, _ := stats.S9()
	// When you connected to Antminer D3
    statsD3, _ := stats.D3()
	// When you connected to Antminer L3+
    statsL3, _ := stats.L3()
	// When you connected to Antminer T9+
    statsT9, _ := stats.T9()
	
}
```