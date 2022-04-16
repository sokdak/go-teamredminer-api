package cgminer

import (
	"context"
	"errors"
	"fmt"
	"strconv"
)

// Version returns version information
//
// For context-based requests use `VersionContext()`
func (c *CGMiner) Version() (*Version, error) {
	return c.VersionContext(context.Background())
}

// VersionContext returns version information using provided context
func (c *CGMiner) VersionContext(ctx context.Context) (*Version, error) {
	resp := new(VersionResponse)
	if err := c.CallContext(ctx, NewCommandWithoutParameter("version"), resp); err != nil {
		return nil, err
	}

	if len(resp.Version) < 1 {
		return nil, errors.New("no version in JSON response")
	}
	if len(resp.Version) > 1 {
		return nil, errors.New("too many versions in JSON response")
	}
	return &resp.Version[0], nil
}

// Devs returns basic information on the miner. See the Devs struct.
func (c *CGMiner) Devs() (*[]Devs, error) {
	return c.DevsContext(context.Background())
}

// DevsContext returns basic information on the miner. See the Devs struct.
func (c *CGMiner) DevsContext(ctx context.Context) (*[]Devs, error) {
	resp := new(devsResponse)
	if err := c.CallContext(ctx, NewCommandWithoutParameter("devs"), resp); err != nil {
		return nil, err
	}

	return &resp.Devs, nil
}

// Summary returns basic information on the miner. See the Summary struct.
func (c *CGMiner) Summary() (*Summary, error) {
	return c.SummaryContext(context.Background())
}

// SummaryContext returns basic information on the miner. See the Summary struct.
func (c *CGMiner) SummaryContext(ctx context.Context) (*Summary, error) {
	resp := new(summaryResponse)
	if err := c.CallContext(ctx, NewCommandWithoutParameter("summary"), resp); err != nil {
		return nil, err
	}

	if len(resp.Summary) > 1 {
		return nil, errors.New("Received multiple Summary objects")
	}
	if len(resp.Summary) < 1 {
		return nil, errors.New("No summary info received")
	}
	return &resp.Summary[0], nil
}

// Stats returns basic information on the miner. See the Stats struct.
func (c *CGMiner) Stats() (Stats, error) {
	return c.StatsContext(context.Background())
}

// StatsContext returns basic information on the miner. See the Stats struct.
func (c *CGMiner) StatsContext(ctx context.Context) (Stats, error) {
	resp := new(statsResponse)
	if err := c.CallContext(ctx, NewCommandWithoutParameter("stats"), resp); err != nil {
		return nil, err
	}

	if len(resp.Stats) < 1 {
		return nil, errors.New("no stats in JSON response")
	}
	// we don't need to check stats has over 2 slices
	//if len(resp.Stats) > 1 {
	//	return nil, errors.New("too many stats in JSON response")
	//}
	return &resp.Stats[0], nil
}

// PoolsContext returns a slice of Pool structs, one per pool.
func (c *CGMiner) PoolsContext(ctx context.Context) ([]Pool, error) {
	resp := new(poolsResponse)
	if err := c.CallContext(ctx, NewCommandWithoutParameter("pools"), resp); err != nil {
		return nil, err
	}

	return resp.Pools, nil
}

// Pools returns a slice of Pool structs, one per pool.
func (c *CGMiner) Pools() ([]Pool, error) {
	return c.PoolsContext(context.Background())
}

// AddPool adds the given URL/username/password combination to the miner's
// pool list.
func (c *CGMiner) AddPool(url, username, password string) error {
	return c.AddPoolContext(context.Background(), url, username, password)
}

// AddPoolContext adds the given URL/username/password combination to the miner's
// pool list with provided context.
func (c *CGMiner) AddPoolContext(ctx context.Context, url, username, password string) error {
	// TODO: Don't allow adding a pool that's already in the pool list
	// TODO: Escape commas in the URL, username, and password
	resp := new(GenericResponse)
	parameter := fmt.Sprintf("%s,%s,%s", url, username, password)
	return c.CallContext(ctx, NewCommand("pools", parameter), resp)
}

func (c *CGMiner) EnablePool(pool *Pool) error {
	return c.EnablePoolContext(context.Background(), pool)
}

func (c *CGMiner) DisablePool(pool *Pool) error {
	return c.DisablePoolContext(context.Background(), pool)
}

func (c *CGMiner) EnablePoolContext(ctx context.Context, pool *Pool) error {
	return c.CallContext(ctx, NewCommand("enablepool", strconv.FormatInt(pool.Pool, 10)), nil)
}

func (c *CGMiner) DisablePoolContext(ctx context.Context, pool *Pool) error {
	return c.CallContext(ctx, NewCommand("disablepool", strconv.FormatInt(pool.Pool, 10)), nil)
}

func (c *CGMiner) RemovePool(pool *Pool) error {
	return c.Call(NewCommand("removepool", strconv.FormatInt(pool.Pool, 10)), nil)
}

func (c *CGMiner) SwitchPool(pool *Pool) error {
	return c.Call(NewCommand("switchpool", strconv.FormatInt(pool.Pool, 10)), nil)
}

// DevDetailContext returns a slice of DeviceDetail structs.
func (c *CGMiner) DevDetailContext(ctx context.Context) ([]DeviceDetail, error) {
	resp := new(deviceDetailResponse)
	if err := c.CallContext(ctx, NewCommandWithoutParameter("devdetails"), resp); err != nil {
		return nil, err
	}

	return resp.DevDetails, nil
}

// DevDetails returns a slice of DeviceDetail structs, one per pool.
func (c *CGMiner) DevDetails() ([]DeviceDetail, error) {
	return c.DevDetailContext(context.Background())
}

func (c *CGMiner) Restart() error {
	return c.Call(NewCommandWithoutParameter("restart"), nil)
}

func (c *CGMiner) Quit() error {
	return c.Call(NewCommandWithoutParameter("quit"), nil)
}

// CheckAvailableCommands - check all commands, that supported by device
// func (miner *CGMiner) CheckAvailableCommands() {
// 	// TODO: add all commands, please note: your ip need to be in "api-allow" range
// 	commandsList := []string{"version"}
// 	for _, cmd := range commandsList {
// 		result, err := miner.runCommand("check", cmd)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		log.Println(result)
// 	}
// }
