package cgminer

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type CGMiner struct {
	server  string
	timeout time.Duration
}

type status struct {
	Code        int
	Description string
	Status      string `json:"STATUS"`
	When        int64
}

// CGfloat32 is a float32 type with own unmarshaller to fix empty JSON string values
type CGfloat32 float32

// UnmarshalJSON for Antpool json timestamp format: 2006-01-02 15:04:05
func (a *CGfloat32) UnmarshalJSON(b []byte) (err error) {
	s := string(b)
	s = s[1 : len(s)-1]

	if len(s) == 0 {
		*a = 0
		return
	}
	value, err := strconv.ParseFloat(s, 32)
	if err != nil {
		msg := fmt.Sprintf("json: failed to unmarshal \"%s\" into float32", s)
		return errors.New(msg)
	}
	*a = CGfloat32(value)
	return
}

// Stats - get stats from antminer S9
type Stats struct {
	// "miner_id": "81680d4162b51111",
	MinerID string `json:"miner_id"`
	// "Type": "Antminer S9",
	Type string `json:"Type"`
	// "CompileTime": "Tue Aug 15 11:37:49 CST 2017",
	CompileTime string `json:"CompileTime"`
	// "Miner": "16.8.1.3",
	Miner string `json:"Miner"`
	// "BMMiner": "2.0.0"
	BMMiner string `json:"BMMiner"`
	// "miner_version": "16.8.1.3",
	MinerVersion string `json:"miner_version"`
	// "miner_count": 3,
	MinerCount int16 `json:"miner_count"`

	// "Elapsed": 1434,
	Elapsed int64 `json:"Elapsed"`
	// "Wait": 0.0,
	// "Device Hardware%": 0.0,
	// "STATS": 0,
	// "Max": 0.0,
	// "no_matching_work": 2,
	// "ID": "BC50",
	// "Calls": 0,
	// "Min": 99999999.0,

	// "total_acn": 189,
	TotalAcn int16 `json:"total_acn"`
	// "total_rate": 13709.27,
	TotalRate float32 `json:"total_rate"`
	// "total_rateideal": 13501.4,
	TotalRateIdeal float32 `json:"total_rateideal"`
	// "total_freqavg": 633.38,
	TotalFrequencyAvg float32 `json:"total_freqavg"`

	// "frequency": "643",
	Frequency CGfloat32 `json:"frequency"`
	// "freq_avg1": 0.0,
	// "freq_avg2": 0.0,
	// "freq_avg3": 0.0,
	// "freq_avg4": 0.0,
	// "freq_avg5": 0.0,
	// "freq_avg6": 633.47,
	FrequencyAvg6 CGfloat32 `json:"freq_avg6"`
	// "freq_avg7": 632.53,
	FrequencyAvg7 CGfloat32 `json:"freq_avg7"`
	// "freq_avg8": 634.14,
	FrequencyAvg8 CGfloat32 `json:"freq_avg8"`
	// "freq_avg9": 0.0,
	// "freq_avg12": 0.0,
	// "freq_avg13": 0.0,
	// "freq_avg10": 0.0,
	// "freq_avg11": 0.0,
	// "freq_avg16": 0.0,
	// "freq_avg14": 0.0,
	// "freq_avg15": 0.0,

	// "fan_num": 2,
	FunNum int16 `json:"fun_num"`
	// "fan1": 0,
	// "fan3": 4200,
	Fan3 int16 `json:"fan3"`
	// "fan2": 0,
	// "fan5": 0,
	// "fan4": 0,
	// "fan7": 0,
	// "fan6": 6000,
	Fan6 int16 `json:"fan6"`
	// "fan8": 0,

	// "temp_max": 68,
	TempMax int16 `json:"temp_max"`
	// "temp_num": 3,
	TempNum int16 `json:"temp_num"`
	// "temp1": 0,
	// "temp2": 0,
	// "temp3": 0,
	// "temp4": 0,
	// "temp5": 0,
	// "temp6": 68,
	Temp6 int16 `json:"temp6,omitempty"`
	// "temp7": 62,
	Temp7 int16 `json:"temp7,omitempty"`
	// "temp8": 61,
	Temp8 int16 `json:"temp8,omitempty"`
	// "temp9": 0,
	// "temp10": 0,
	// "temp11": 0,
	// "temp12": 0,
	// "temp13": 0,
	// "temp14": 0,
	// "temp15": 0,
	// "temp16": 0,
	// "temp2_1": 0,
	// "temp2_2": 0,
	// "temp2_3": 0,
	// "temp2_4": 0,
	// "temp2_5": 0,
	// "temp2_6": 83,
	Temp2_6 int16 `json:"temp2_6,omitempty"`
	// "temp2_7": 77,
	Temp2_7 int16 `json:"temp2_7,omitempty"`
	// "temp2_8": 76,
	Temp2_8 int16 `json:"temp2_8,omitempty"`
	// "temp2_9": 0,
	// "temp2_10": 0,
	// "temp2_11": 0,
	// "temp2_12": 0,
	// "temp2_13": 0,
	// "temp2_14": 0,
	// "temp2_15": 0,
	// "temp2_16": 0,
	// "temp3_1": 0,
	// "temp3_2": 0,
	// "temp3_3": 0,
	// "temp3_4": 0,
	// "temp3_5": 0,
	// "temp3_6": 0,
	// "temp3_7": 0,
	// "temp3_9": 0,
	// "temp3_8": 0,
	// "temp3_10": 0,
	// "temp3_11": 0,
	// "temp3_12": 0,
	// "temp3_13": 0,
	// "temp3_14": 0,
	// "temp3_15": 0,
	// "temp3_16": 0,

	// "GHS 5s": "13709.27",
	Ghs5s CGfloat32 `json:"GHS 5s"`
	// "GHS av": 13681.36,
	GhsAverage CGfloat32 `json:"GHS av"`

	// "chain_hw1": 0,
	// "chain_hw2": 0,
	// "chain_hw3": 0,
	// "chain_hw4": 0,
	// "chain_hw5": 0,
	// "chain_hw6": 0,
	// "chain_hw7": 2,
	// "chain_hw8": 0,
	// "chain_hw9": 0,
	// "chain_hw10": 0,
	// "chain_hw11": 0,
	// "chain_hw12": 0,
	// "chain_hw13": 0,
	// "chain_hw14": 0,
	// "chain_hw15": 0,
	// "chain_hw16": 0,

	// "chain_acs1": "",
	// "chain_acs2": "",
	// "chain_acs3": "",
	// "chain_acs4": "",
	// "chain_acs5": "",
	// "chain_acs6": " oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo ooooooo",
	ChainAcs6 string `json:"chain_acs6"`
	// "chain_acs7": " oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo ooooooo",
	ChainAcs7 string `json:"chain_acs7"`
	// "chain_acs8": " oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo ooooooo",
	ChainAcs8 string `json:"chain_acs8"`
	// "chain_acs9": "",
	// "chain_acs10": "",
	// "chain_acs11": "",
	// "chain_acs12": "",
	// "chain_acs13": "",
	// "chain_acs14": "",
	// "chain_acs15": "",
	// "chain_acs16": "",

	// "chain_acn1": 0,
	// "chain_acn2": 0,
	// "chain_acn3": 0,
	// "chain_acn4": 0,
	// "chain_acn6": 63,
	// "chain_acn5": 0,
	// "chain_acn7": 63,
	// "chain_acn8": 63,
	// "chain_acn9": 0,
	// "chain_acn10": 0,
	// "chain_acn11": 0,
	// "chain_acn12": 0,
	// "chain_acn13": 0,
	// "chain_acn14": 0,
	// "chain_acn15": 0,
	// "chain_acn16": 0,

	// "chain_rate1": "",
	// "chain_rate2": "",
	// "chain_rate3": "",
	// "chain_rate4": "",
	// "chain_rate5": "",
	// "chain_rate6": "4554.34",
	ChainRate6 CGfloat32 `json:"chain_rate6,omitempty"`
	// "chain_rate7": "4573.79",
	ChainRate7 CGfloat32 `json:"chain_rate7,omitempty"`
	// "chain_rate8": "4581.14",
	ChainRate8 CGfloat32 `json:"chain_rate8,omitempty"`
	// "chain_rate9": "",
	// "chain_rate10": "",
	// "chain_rate11": "",
	// "chain_rate12": "",
	// "chain_rate13": "",
	// "chain_rate14": "",
	// "chain_rate15": "",
	// "chain_rate16": "",

	// "chain_rateideal1": 0.0,
	// "chain_rateideal2": 0.0,
	// "chain_rateideal3": 0.0,
	// "chain_rateideal4": 0.0,
	// "chain_rateideal5": 0.0,
	// "chain_rateideal6": 4500.16,
	// "chain_rateideal7": 4500.72,
	// "chain_rateideal8": 4500.5,
	// "chain_rateideal9": 0.0,
	// "chain_rateideal10": 0.0,
	// "chain_rateideal11": 0.0,
	// "chain_rateideal12": 0.0,
	// "chain_rateideal13": 0.0,
	// "chain_rateideal14": 0.0,
	// "chain_rateideal15": 0.0,
	// "chain_rateideal16": 0.0,

	// "chain_opencore_6": "1",
	// "chain_opencore_7": "1",
	// "chain_opencore_8": "1"

	// "chain_offside_6": "0",
	// "chain_offside_7": "0",
	// "chain_offside_8": "0",

	// "chain_xtime6": "{}",
	// "chain_xtime7": "{}",
	// "chain_xtime8": "{}",
}

type Summary struct {
	Accepted               int64
	BestShare              int64   `json:"Best Share"`
	DeviceHardwarePercent  float64 `json:"Device Hardware%"`
	DeviceRejectedPercent  float64 `json:"Device Rejected%"`
	DifficultyAccepted     float64 `json:"Difficulty Accepted"`
	DifficultyRejected     float64 `json:"Difficulty Rejected"`
	DifficultyStale        float64 `json:"Difficulty Stale"`
	Discarded              int64
	Elapsed                int64
	FoundBlocks            int64 `json:"Found Blocks"`
	GetFailures            int64 `json:"Get Failures"`
	Getworks               int64
	HardwareErrors         int64   `json:"Hardware Errors"`
	LocalWork              int64   `json:"Local Work"`
	MHS5s                  float64 `json:"MHS 5s"`
	MHSav                  float64 `json:"MHS av"`
	NetworkBlocks          int64   `json:"Network Blocks"`
	PoolRejectedPercentage float64 `json:"Pool Rejected%"`
	PoolStalePercentage    float64 `json:"Pool Stale%"`
	Rejected               int64
	RemoteFailures         int64 `json:"Remote Failures"`
	Stale                  int64
	TotalMH                float64 `json:"Total MH"`
	Utilty                 float64
	WorkUtility            float64 `json:"Work Utility"`
}

type Devs struct {
	GPU                 int64
	Enabled             string
	Status              string
	Temperature         float64
	FanSpeed            int     `json:"Fan Speed"`
	FanPercent          int64   `json:"Fan Percent"`
	GPUClock            int64   `json:"GPU Clock"`
	MemoryClock         int64   `json:"Memory Clock"`
	GPUVoltage          float64 `json:"GPU Voltage"`
	Powertune           int64
	MHSav               float64 `json:"MHS av"`
	MHS5s               float64 `json:"MHS 5s"`
	Accepted            int64
	Rejected            int64
	HardwareErrors      int64 `json:"Hardware Errors"`
	Utility             float64
	Intensity           string
	LastSharePool       int64   `json:"Last Share Pool"`
	LashShareTime       int64   `json:"Lash Share Time"`
	TotalMH             float64 `json:"TotalMH"`
	Diff1Work           int64   `json:"Diff1 Work"`
	DifficultyAccepted  float64 `json:"Difficulty Accepted"`
	DifficultyRejected  float64 `json:"Difficulty Rejected"`
	LastShareDifficulty float64 `json:"Last Share Difficulty"`
	LastValidWork       int64   `json:"Last Valid Work"`
	DeviceHardware      float64 `json:"Device Hardware%"`
	DeviceRejected      float64 `json:"Device Rejected%"`
	DeviceElapsed       int64   `json:"Device Elapsed"`
}

/*
       "Stratum Active": false,
       "Difficulty Accepted": 0.0,
       "Pool Rejected%": 100.0,
       "Difficulty Rejected": 4096.0,
       "Diff1 Shares": 0,
       "Status": "Alive",
       "Proxy Type": "",
       "Best Share": 0,
       "Pool Stale%": 0.0,
       "Quota": 1,
       "Rejected": 2,
       "Stratum URL": "",
       "Proxy": "",
       "Long Poll": "N",
       "Accepted": 0,
       "User": "cryptotrain.133",
       "Get Failures": 0,
       "Difficulty Stale": 0.0,
       "URL": "stratum+tcp://s1.theblocksfactory.com:9001",
       "Discarded": 0,
       "Has Stratum": true,
       "Last Share Time": "0",
       "Stale": 0,
       "POOL": 2,
       "Priority": 2,
       "Getworks": 1,
       "Has GBT": false,
       "Last Share Difficulty": 0.0,
       "Diff": "2.05K",
       "Remote Failures": 0
   }
*/
type Pool struct {
	Accepted               int64
	BestShare              int64   `json:"Best Share"`
	Diff1Shares            int64   `json:"Diff1 Shares"`
	DifficultyAccepted     float64 `json:"Difficulty Accepted"`
	DifficultyRejected     float64 `json:"Difficulty Rejected"`
	DifficultyStale        float64 `json:"Difficulty Stale"`
	Discarded              int64
	GetFailures            int64 `json:"Get Failures"`
	Getworks               int64
	HasGBT                 bool    `json:"Has GBT"`
	HasStratum             bool    `json:"Has Stratum"`
	LastShareDifficulty    float64 `json:"Last Share Difficulty"`
	LastShareTime          string  `json:"Last Share Time"`
	LongPoll               string  `json:"Long Poll"`
	Pool                   int64   `json:"POOL"`
	PoolRejectedPercentage float64 `json:"Pool Rejected%"`
	PoolStalePercentage    float64 `json:"Pool Stale%"`
	Priority               int64
	ProxyType              string `json:"Proxy Type"`
	Proxy                  string
	Quota                  int64
	Rejected               int64
	RemoteFailures         int64 `json:"Remote Failures"`
	Stale                  int64
	Status                 string
	StratumActive          bool   `json:"Stratum Active"`
	StratumURL             string `json:"Stratum URL"`
	URL                    string `json:"URL"`
	User                   string `json:"User"`
	Works                  int64
}

type statsResponse struct {
	Status []status `json:"STATUS"`
	Stats  []Stats  `json:"STATS"`
	Id     int64    `json:"id"`
}

type summaryResponse struct {
	Status  []status  `json:"STATUS"`
	Summary []Summary `json:"SUMMARY"`
	Id      int64     `json:"id"`
}

type devsResponse struct {
	Status []status `json:"STATUS"`
	Devs   []Devs   `json:"DEVS"`
	Id     int64    `json:"id"`
}

type poolsResponse struct {
	Status []status `json:"STATUS"`
	Pools  []Pool   `json:"POOLS"`
	Id     int64    `json:"id"`
}

type addPoolResponse struct {
	Status []status `json:"STATUS"`
	Id     int64    `json:"id"`
}

// New returns a CGMiner pointer, which is used to communicate with a running
// CGMiner instance. Note that New does not attempt to connect to the miner.
func New(hostname string, port int64, timeout int) *CGMiner {
	miner := new(CGMiner)
	miner.server = fmt.Sprintf("%s:%d", hostname, port)
	miner.timeout = time.Second * time.Duration(timeout)
	return miner
}

func (miner *CGMiner) runCommand(command, argument string) (string, error) {
	conn, err := net.DialTimeout("tcp", miner.server, miner.timeout)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	type commandRequest struct {
		Command   string `json:"command"`
		Parameter string `json:"parameter,omitempty"`
	}

	request := &commandRequest{
		Command: command,
	}

	if argument != "" {
		request.Parameter = argument
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	fmt.Fprintf(conn, "%s", requestBody)
	result, err := bufio.NewReader(conn).ReadString('\x00')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(result, "\x00"), nil
}

// Devs returns basic information on the miner. See the Devs struct.
func (miner *CGMiner) Devs() (*[]Devs, error) {
	result, err := miner.runCommand("devs", "")
	if err != nil {
		return nil, err
	}

	var devsResponse devsResponse
	err = json.Unmarshal([]byte(result), &devsResponse)
	if err != nil {
		return nil, err
	}

	var devs = devsResponse.Devs
	return &devs, err
}

// Summary returns basic information on the miner. See the Summary struct.
func (miner *CGMiner) Summary() (*Summary, error) {
	result, err := miner.runCommand("summary", "")
	if err != nil {
		return nil, err
	}

	var summaryResponse summaryResponse
	err = json.Unmarshal([]byte(result), &summaryResponse)
	if err != nil {
		return nil, err
	}

	if len(summaryResponse.Summary) != 1 {
		return nil, errors.New("Received multiple Summary objects")
	}

	var summary = summaryResponse.Summary[0]
	return &summary, err
}

// Stats returns basic information on the miner. See the Stats struct.
func (miner *CGMiner) Stats() (*Stats, error) {
	result, err := miner.runCommand("stats", "")
	if err != nil {
		return nil, err
	}

	var statsResponse statsResponse
	// fix incorrect json response from miner "}{"
	fixResponse := bytes.Replace([]byte(result), []byte("}{"), []byte(","), 1)
	err = json.Unmarshal(fixResponse, &statsResponse)
	if err != nil {
		return nil, err
	}
	var stats = statsResponse.Stats[0]
	return &stats, nil
}

// Pools returns a slice of Pool structs, one per pool.
func (miner *CGMiner) Pools() ([]Pool, error) {
	result, err := miner.runCommand("pools", "")
	if err != nil {
		return nil, err
	}

	var poolsResponse poolsResponse
	err = json.Unmarshal([]byte(result), &poolsResponse)
	if err != nil {
		return nil, err
	}

	var pools = poolsResponse.Pools
	return pools, nil
}

// AddPool adds the given URL/username/password combination to the miner's
// pool list.
func (miner *CGMiner) AddPool(url, username, password string) error {
	// TODO: Don't allow adding a pool that's already in the pool list
	// TODO: Escape commas in the URL, username, and password
	parameter := fmt.Sprintf("%s,%s,%s", url, username, password)
	result, err := miner.runCommand("addpool", parameter)
	if err != nil {
		return err
	}

	var addPoolResponse addPoolResponse
	err = json.Unmarshal([]byte(result), &addPoolResponse)
	if err != nil {
		// If there an error here, it's possible that the pool was actually added
		return err
	}

	status := addPoolResponse.Status[0]

	if status.Status != "S" {
		return errors.New(fmt.Sprintf("%d: %s", status.Code, status.Description))
	}

	return nil
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
