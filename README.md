# CGMiner API for Go #

## Installation ##

    # install the library:
    go get github.com/crypt0train/go-cgminer-api

    // Use in your .go code:
    import (
        "github.com/crypt0train/go-cgminer-api"
    )

## API Documentation ##

Not yet.

## Quickstart ##

```go
package main

import (
    "github.com/crypt0train/go-cgminer-api"
    "log"
)

func main() {
    miner := cgminer.New("localhost", 4028, 2)
	stats, err := miner.Stats()
	if err != nil {
		log.Println("Unable to connect to CGMiner: ", err)
		return
	}
	fmt.Printf("%s | %s | temp: %d, %d, %d | GHS avg: %0.2f | fan: %d, %d | %s\n",
		stats.Type, ip, stats.Temp2_6, stats.Temp2_7, stats.Temp2_8, stats.GhsAverage,
		stats.Fan3, stats.Fan6, stats.CompileTime)
}
```