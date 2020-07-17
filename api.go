package cgminer

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"
)

// StatusGetter is generic command response which provides execution status
type StatusGetter interface {
	// GetStatus returns status
	GetStatus() []Status
}

func (miner *CGMiner) checkStatus(statuses []Status) error {
	for _, status := range statuses {
		switch status.Status {
		case "E":
			return fmt.Errorf("API returned error: Code: %d, Msg: '%s', Description: '%s'", status.Code, status.Msg, status.Description)
		case "F":
			return fmt.Errorf("API returned FATAL error: Code: %d, Msg: '%s', Description: '%s'", status.Code, status.Msg, status.Description)
		}
	}
	return nil
}

func (miner *CGMiner) connect() (net.Conn, error) {
	return net.DialTimeout("tcp", miner.server, miner.timeout)
}

func (miner *CGMiner) connectContext(ctx context.Context) (net.Conn, error) {
	dial := net.Dialer{Timeout: miner.timeout}
	return dial.DialContext(ctx, "tcp", miner.server)
}

func (miner *CGMiner) sendCommand(conn net.Conn, command, argument string) ([]byte, error) {
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(miner.timeout))
	request := &commandRequest{
		Command:   command,
		Parameter: argument,
	}

	requestBody, _ := json.Marshal(request)
	_, err := conn.Write(requestBody)
	if err != nil {
		return nil, err
	}
	result, err := bufio.NewReader(conn).ReadBytes(0x00)
	if err != nil {
		return nil, err
	}
	return bytes.TrimRight(result, "\x00"), nil
}

// Call calls API function and writes result to output or returns error
//
// Output should implement StatusGetter interface to validate command Status (see GenericResponse)
func (miner *CGMiner) Call(ctx context.Context, command, argument string, out StatusGetter) error {
	result, err := miner.commandCtx(ctx, command, argument)
	if err != nil {
		return fmt.Errorf("cgminer: %q - %w", command, err)
	}

	if err = json.Unmarshal(result, out); err != nil {
		return fmt.Errorf("cgminer: %q - %w", command, err)
	}

	if err = miner.checkStatus(out.GetStatus()); err != nil {
		return fmt.Errorf("cgminer: %q - %w", command, err)
	}

	return nil
}

func (miner *CGMiner) commandCtx(ctx context.Context, command, argument string) ([]byte, error) {
	conn, err := miner.connectContext(ctx)
	if err != nil {
		return nil, err
	}
	return miner.sendCommand(conn, command, argument)
}

func (miner *CGMiner) runCommand(command, argument string) ([]byte, error) {
	conn, err := miner.connect()
	if err != nil {
		return nil, err
	}
	return miner.sendCommand(conn, command, argument)
}

// Devs returns basic information on the miner. See the Devs struct.
func (miner *CGMiner) Devs() (*[]Devs, error) {
	result, err := miner.runCommand("devs", "")
	if err != nil {
		return nil, err
	}
	var resp devsResponse
	err = json.Unmarshal(result, &resp)
	if err != nil {
		return nil, err
	}
	err = miner.checkStatus(resp.Status)
	if err != nil {
		return nil, err
	}
	return &resp.Devs, err
}

// Summary returns basic information on the miner. See the Summary struct.
func (miner *CGMiner) Summary() (*Summary, error) {
	return miner.SummaryContext(context.Background())
}

// SummaryContext returns basic information on the miner. See the Summary struct.
func (miner *CGMiner) SummaryContext(ctx context.Context) (*Summary, error) {
	result, err := miner.commandCtx(ctx, "summary", "")
	if err != nil {
		return nil, err
	}
	var resp summaryResponse
	err = json.Unmarshal(result, &resp)
	if err != nil {
		return nil, err
	}
	err = miner.checkStatus(resp.Status)
	if err != nil {
		return nil, err
	}
	if len(resp.Summary) > 1 {
		return nil, errors.New("Received multiple Summary objects")
	}
	if len(resp.Summary) < 1 {
		return nil, errors.New("No summary info received")
	}
	return &resp.Summary[0], err
}

// Stats returns basic information on the miner. See the Stats struct.
func (miner *CGMiner) Stats() (Stats, error) {
	return miner.StatsContext(context.Background())
}

// StatsContext returns basic information on the miner. See the Stats struct.
func (miner *CGMiner) StatsContext(ctx context.Context) (Stats, error) {
	result, err := miner.commandCtx(ctx, "stats", "")
	if err != nil {
		return nil, err
	}
	var resp statsResponse
	// fix incorrect json response from miner "}{"
	fixResponse := bytes.Replace(result, []byte("}{"), []byte(","), 1)
	err = json.Unmarshal(fixResponse, &resp)
	if err != nil {
		return nil, err
	}
	err = miner.checkStatus(resp.Status)
	if err != nil {
		return nil, err
	}
	if len(resp.Stats) < 1 {
		return nil, errors.New("no stats in JSON response")
	}
	if len(resp.Stats) > 1 {
		return nil, errors.New("too many stats in JSON response")
	}
	return &resp.Stats[0], nil
}

// PoolsContext returns a slice of Pool structs, one per pool.
func (miner *CGMiner) PoolsContext(ctx context.Context) ([]Pool, error) {
	result, err := miner.commandCtx(ctx, "pools", "")
	if err != nil {
		return nil, err
	}
	var resp poolsResponse
	err = json.Unmarshal(result, &resp)
	if err != nil {
		return nil, err
	}
	err = miner.checkStatus(resp.Status)
	if err != nil {
		return nil, err
	}
	return resp.Pools, nil
}

// Pools returns a slice of Pool structs, one per pool.
func (miner *CGMiner) Pools() ([]Pool, error) {
	return miner.PoolsContext(context.Background())
}

// AddPoolCtx adds the given URL/username/password combination to the miner's
// pool list.
func (miner *CGMiner) AddPoolCtx(ctx context.Context, url, username, password string) error {
	// TODO: Don't allow adding a pool that's already in the pool list
	// TODO: Escape commas in the URL, username, and password
	parameter := fmt.Sprintf("%s,%s,%s", url, username, password)
	result, err := miner.commandCtx(ctx, "addpool", parameter)
	if err != nil {
		return err
	}
	var resp GenericResponse
	err = json.Unmarshal(result, &resp)
	if err != nil {
		return err
	}
	err = miner.checkStatus(resp.Status)
	if err != nil {
		return err
	}
	return nil
}

// AddPool adds the given URL/username/password combination to the miner's
// pool list.
func (miner *CGMiner) AddPool(url, username, password string) error {
	return miner.AddPoolCtx(context.Background(), url, username, password)
}

func (miner *CGMiner) Enable(pool *Pool) error {
	parameter := fmt.Sprintf("%d", pool.Pool)
	_, err := miner.runCommand("enablepool", parameter)
	return err
}

func (miner *CGMiner) Disable(pool *Pool) error {
	parameter := fmt.Sprintf("%d", pool.Pool)
	_, err := miner.runCommand("disablepool", parameter)
	return err
}

func (miner *CGMiner) Delete(pool *Pool) error {
	parameter := fmt.Sprintf("%d", pool.Pool)
	_, err := miner.runCommand("removepool", parameter)
	return err
}

func (miner *CGMiner) SwitchPool(pool *Pool) error {
	parameter := fmt.Sprintf("%d", pool.Pool)
	_, err := miner.runCommand("switchpool", parameter)
	return err
}

func (miner *CGMiner) Restart() error {
	_, err := miner.runCommand("restart", "")
	return err
}

func (miner *CGMiner) Quit() error {
	_, err := miner.runCommand("quit", "")
	return err
}

// Version - reply section VERSION
func (miner *CGMiner) Version() (*Version, error) {
	response, err := miner.runCommand("version", "")
	if err != nil {
		return nil, err
	}
	resp := &VersionResponse{}
	err = json.Unmarshal(response, resp)
	if err != nil {
		return nil, err
	}
	err = miner.checkStatus(resp.Status)
	if err != nil {
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
