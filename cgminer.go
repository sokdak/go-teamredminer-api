package cgminer

import (
	"fmt"
	"time"
)

// NewCGMiner returns a CGMiner pointer, which is used to communicate with a running
// cgminer instance. Note that New does not attempt to connect to the miner.
func NewCGMiner(hostname string, port int64, timeout int) *CGMiner {
	miner := new(CGMiner)
	miner.server = fmt.Sprintf("%s:%d", hostname, port)
	miner.timeout = time.Second * time.Duration(timeout)
	return miner
}
