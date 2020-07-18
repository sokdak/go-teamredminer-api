package cgminer

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"path"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/phayes/freeport"
)

const (
	proto        = "tcp"
	ip           = "127.0.0.1"
	minerTimeout = 5 * time.Second
)

type accepted struct {
	conn net.Conn
	err  error
}

func wait(amout int) {
	time.Sleep(time.Duration(amout) * 10 * time.Millisecond)
}

func mockTCPServer(ctx context.Context, ip string, port int, payload []byte) {
	addr := fmt.Sprintf("%s:%d", ip, port)
	listener, err := net.Listen(proto, addr)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
	c := make(chan accepted, 1)
	go func() {
		conn, err := listener.Accept()
		c <- accepted{conn, err}
	}()
	for {
		select {
		case <-ctx.Done():
			listener.Close()
			return
		case a := <-c:
			if a.err != nil {
				log.Fatal(err)
			}
			go handleConn(a.conn, payload)
		default: //nolint
		}
	}
}

func handleConn(conn net.Conn, payload []byte) {
	// nolint
	_, _ = conn.Write(payload)
	_ = conn.Close()
}

func getPort() int {
	port, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}
	return port
}

func getFixture(filename string) []byte {
	payload, err := ioutil.ReadFile(path.Join("testdata", filename))
	if err != nil {
		log.Fatal(err)
	}
	payload = append(payload, byte(0x00))
	return payload
}
func TestConnFailed(t *testing.T) {
	port := getPort()
	miner := NewCGMiner(ip, port, minerTimeout)
	_, err := miner.Version()
	if err == nil {
		t.Fatal("must throw err: connection refused")
	}

	if _, ok := err.(ConnectError); !ok {
		t.Fatalf("expected ConnectError type, got %T", err)
	}
}

func TestVersion(t *testing.T) {
	testCaseValue := getFixture("TestVersion.json")
	expected := &Version{
		BMMiner:     "2.0.0",
		API:         "3.1",
		Miner:       "16.8.1.3",
		CompileTime: "Fri Nov 17 17:57:49 CST 2017",
		Type:        "Antminer S9",
	}
	ctx, finish := context.WithCancel(context.Background())
	port := getPort()
	go mockTCPServer(ctx, ip, port, testCaseValue)
	wait(1)
	miner := NewCGMiner(ip, port, minerTimeout)
	version, err := miner.Version()
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(version, expected); diff != nil {
		t.Error(diff)
	}
	finish()
	wait(1)
}

func TestVersionEmpty(t *testing.T) {
	testCaseValue := getFixture("TestVersionEmpty.json")
	ctx, finish := context.WithCancel(context.Background())
	port := getPort()
	go mockTCPServer(ctx, ip, port, testCaseValue)
	wait(1)
	miner := NewCGMiner(ip, port, minerTimeout)
	_, err := miner.Version()
	if err == nil {
		t.FailNow()
	}
	finish()
	wait(1)
}
func TestVersionMany(t *testing.T) {
	testCaseValue := getFixture("TestVersionMany.json")
	ctx, finish := context.WithCancel(context.Background())
	port := getPort()
	go mockTCPServer(ctx, ip, port, testCaseValue)
	wait(1)
	miner := NewCGMiner(ip, port, minerTimeout)
	_, err := miner.Version()
	if err == nil {
		t.FailNow()
	}
	finish()
	wait(1)
}

func TestVersionBadJSON(t *testing.T) {
	testCaseValue := getFixture("TestVersionBadJSON.json")
	ctx, finish := context.WithCancel(context.Background())
	port := getPort()
	go mockTCPServer(ctx, ip, port, testCaseValue)
	wait(1)
	miner := NewCGMiner(ip, port, minerTimeout)
	_, err := miner.Version()
	if err == nil {
		t.FailNow()
	}
	finish()
	wait(1)
}
func TestVersionBadJSON2(t *testing.T) {
	testCaseValue := []byte(`{"d,k`)
	ctx, finish := context.WithCancel(context.Background())
	port := getPort()
	go mockTCPServer(ctx, ip, port, testCaseValue)
	wait(1)
	miner := NewCGMiner(ip, port, minerTimeout)
	_, err := miner.Version()
	if err == nil {
		t.FailNow()
	}
	finish()
	wait(1)
}

func TestVersionStatusError(t *testing.T) {
	testCaseValue := getFixture("TestVersionStatusError.json")
	ctx, finish := context.WithCancel(context.Background())
	port := getPort()
	go mockTCPServer(ctx, ip, port, testCaseValue)
	wait(1)
	miner := NewCGMiner(ip, port, minerTimeout)
	_, err := miner.Version()
	if err == nil {
		t.FailNow()
	}
	finish()
	wait(1)
}
func TestVersionStatusFatal(t *testing.T) {
	testCaseValue := getFixture("TestVersionStatusFatal.json")
	ctx, finish := context.WithCancel(context.Background())
	port := getPort()
	go mockTCPServer(ctx, ip, port, testCaseValue)
	wait(1)
	miner := NewCGMiner(ip, port, minerTimeout)
	_, err := miner.Version()
	if err == nil {
		t.FailNow()
	}
	finish()
	wait(1)
}

func TestSummary(t *testing.T) {
	testCaseValue := getFixture("TestSummary.json")
	expected := &Summary{
		Elapsed:               40206,
		GHS5s:                 13581.93,
		GHSav:                 13580.83,
		FoundBlocks:           0,
		Getworks:              1375,
		Accepted:              7629,
		Rejected:              2,
		HardwareErrors:        354,
		Utility:               11.38,
		Discarded:             21530,
		Stale:                 0,
		GetFailures:           0,
		LocalWork:             2010877,
		RemoteFailures:        0,
		NetworkBlocks:         66,
		TotalMH:               546017369885.0,
		WorkUtility:           186557.18,
		DifficultyAccepted:    124993536.0,
		DifficultyRejected:    18432.0,
		DifficultyStale:       0.0,
		BestShare:             99242396,
		DeviceHardwarePercent: 0.0003,
		DeviceRejectedPercent: 0.0147,
		PoolRejectedPercent:   0.0147,
		PoolStalePercent:      0.0,
		LastGetWork:           1521044524,
	}
	ctx, finish := context.WithCancel(context.Background())
	port := getPort()
	go mockTCPServer(ctx, ip, port, testCaseValue)
	wait(1)
	miner := NewCGMiner(ip, port, minerTimeout)
	summary, err := miner.Summary()
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(summary, expected); diff != nil {
		t.Error(diff)
	}
	finish()
	wait(1)
}

func TestSummaryStatusError(t *testing.T) {
	testCaseValue := getFixture("TestSummaryStatusError.json")
	ctx, finish := context.WithCancel(context.Background())
	port := getPort()
	go mockTCPServer(ctx, ip, port, testCaseValue)
	miner := NewCGMiner(ip, port, minerTimeout)
	_, err := miner.Summary()
	if err == nil {
		t.FailNow()
	}
	finish()
	wait(1)
}

func TestSummaryStatusFatal(t *testing.T) {
	testCaseValue := getFixture("TestSummaryStatusFatal.json")
	ctx, finish := context.WithCancel(context.Background())
	port := getPort()
	go mockTCPServer(ctx, ip, port, testCaseValue)
	miner := NewCGMiner(ip, port, minerTimeout)
	_, err := miner.Summary()
	if err == nil {
		t.FailNow()
	}
	finish()
	wait(1)
}

func TestSummaryEmpty(t *testing.T) {
	testCaseValue := getFixture("TestSummaryEmpty.json")
	ctx, finish := context.WithCancel(context.Background())
	port := getPort()
	go mockTCPServer(ctx, ip, port, testCaseValue)
	miner := NewCGMiner(ip, port, minerTimeout)
	_, err := miner.Summary()
	if err == nil {
		t.FailNow()
	}
	finish()
	wait(1)
}

func TestSummaryMany(t *testing.T) {
	testCaseValue := getFixture("TestSummaryMany.json")
	ctx, finish := context.WithCancel(context.Background())
	port := getPort()
	go mockTCPServer(ctx, ip, port, testCaseValue)
	miner := NewCGMiner(ip, port, minerTimeout)
	_, err := miner.Summary()
	if err == nil {
		t.FailNow()
	}
	finish()
	wait(1)
}

func TestSummaryBadJSON(t *testing.T) {
	testCaseValue := getFixture("TestSummaryBadJSON.json")
	ctx, finish := context.WithCancel(context.Background())
	port := getPort()
	go mockTCPServer(ctx, ip, port, testCaseValue)
	miner := NewCGMiner(ip, port, minerTimeout)
	_, err := miner.Summary()
	if err == nil {
		t.FailNow()
	}
	t.Log(err)
	finish()
	wait(1)
}

func TestSummaryBadJSON2(t *testing.T) {
	testCaseValue := []byte(`{;s$`)
	ctx, finish := context.WithCancel(context.Background())
	port := getPort()
	go mockTCPServer(ctx, ip, port, testCaseValue)
	miner := NewCGMiner(ip, port, minerTimeout)
	_, err := miner.Summary()
	if err == nil {
		t.FailNow()
	}
	t.Log(err)
	finish()
	wait(1)
}

func TestSummaryConnRefused(t *testing.T) {
	port := getPort()
	miner := NewCGMiner(ip, port, minerTimeout)
	_, err := miner.Summary()
	if err == nil {
		t.FailNow()
	}
}

func TestStats(t *testing.T) {
	testCaseValue := getFixture("TestStatsS9.json")
	expected := &GenericStats{
		BMMiner:               "2.0.0",
		Miner:                 "16.8.1.3",
		CompileTime:           "Fri Nov 17 17:57:49 CST 2017",
		Type:                  "Antminer S9",
		Stats:                 0,
		ID:                    "BC50",
		Elapsed:               100986,
		Calls:                 0,
		Wait:                  0.0,
		Max:                   0.0,
		Min:                   99999999.0,
		Ghs5s:                 "13630.55",
		GhsAverage:            13569.28,
		MinerCount:            3,
		Frequency:             637,
		FanNum:                2,
		Fan1:                  0,
		Fan2:                  0,
		Fan3:                  4080,
		Fan4:                  0,
		Fan5:                  0,
		Fan6:                  4080,
		Fan7:                  0,
		Fan8:                  0,
		TempNum:               3,
		Temp1:                 0,
		Temp2:                 0,
		Temp3:                 0,
		Temp4:                 0,
		Temp5:                 0,
		Temp6:                 56,
		Temp7:                 52,
		Temp8:                 56,
		Temp9:                 0,
		Temp10:                0,
		Temp11:                0,
		Temp12:                0,
		Temp13:                0,
		Temp14:                0,
		Temp15:                0,
		Temp16:                0,
		Temp2_1:               0,
		Temp2_2:               0,
		Temp2_3:               0,
		Temp2_4:               0,
		Temp2_5:               0,
		Temp2_6:               71,
		Temp2_7:               67,
		Temp2_8:               71,
		Temp2_9:               0,
		Temp2_10:              0,
		Temp2_11:              0,
		Temp2_12:              0,
		Temp2_13:              0,
		Temp2_14:              0,
		Temp2_15:              0,
		Temp2_16:              0,
		Temp3_1:               0,
		Temp3_2:               0,
		Temp3_3:               0,
		Temp3_4:               0,
		Temp3_5:               0,
		Temp3_6:               0,
		Temp3_7:               0,
		Temp3_8:               0,
		Temp3_9:               0,
		Temp3_10:              0,
		Temp3_11:              0,
		Temp3_12:              0,
		Temp3_13:              0,
		Temp3_14:              0,
		Temp3_15:              0,
		Temp3_16:              0,
		FrequencyAvg1:         0.0,
		FrequencyAvg2:         0.0,
		FrequencyAvg3:         0.0,
		FrequencyAvg4:         0.0,
		FrequencyAvg5:         0.0,
		FrequencyAvg6:         627.57,
		FrequencyAvg7:         627.76,
		FrequencyAvg8:         627.09,
		FrequencyAvg9:         0.0,
		FrequencyAvg10:        0.0,
		FrequencyAvg11:        0.0,
		FrequencyAvg12:        0.0,
		FrequencyAvg13:        0.0,
		FrequencyAvg14:        0.0,
		FrequencyAvg15:        0.0,
		FrequencyAvg16:        0.0,
		TotalRateIdeal:        13501.38,
		TotalFrequencyAvg:     627.47,
		TotalAcn:              189,
		TotalRate:             13630.54,
		ChainRateIdeal1:       0.0,
		ChainRateIdeal2:       0.0,
		ChainRateIdeal3:       0.0,
		ChainRateIdeal4:       0.0,
		ChainRateIdeal5:       0.0,
		ChainRateIdeal6:       4500.33,
		ChainRateIdeal7:       4500.38,
		ChainRateIdeal8:       4500.66,
		ChainRateIdeal9:       0.0,
		ChainRateIdeal10:      0.0,
		ChainRateIdeal11:      0.0,
		ChainRateIdeal12:      0.0,
		ChainRateIdeal13:      0.0,
		ChainRateIdeal14:      0.0,
		ChainRateIdeal15:      0.0,
		ChainRateIdeal16:      0.0,
		TempMax:               56,
		DeviceHardwarePercent: 0.0004,
		NotMatchingWork:       1222,
		ChainAcn1:             0,
		ChainAcn2:             0,
		ChainAcn3:             0,
		ChainAcn4:             0,
		ChainAcn5:             0,
		ChainAcn6:             63,
		ChainAcn7:             63,
		ChainAcn8:             63,
		ChainAcn9:             0,
		ChainAcn10:            0,
		ChainAcn11:            0,
		ChainAcn12:            0,
		ChainAcn13:            0,
		ChainAcn14:            0,
		ChainAcn15:            0,
		ChainAcn16:            0,
		ChainAcs1:             "",
		ChainAcs2:             "",
		ChainAcs3:             "",
		ChainAcs4:             "",
		ChainAcs5:             "",
		ChainAcs6:             " oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo ooooooo",
		ChainAcs7:             " oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo ooooooo",
		ChainAcs8:             " oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo ooooooo",
		ChainAcs9:             "",
		ChainAcs10:            "",
		ChainAcs11:            "",
		ChainAcs12:            "",
		ChainAcs13:            "",
		ChainAcs14:            "",
		ChainAcs15:            "",
		ChainAcs16:            "",
		ChainHW1:              0,
		ChainHW2:              0,
		ChainHW3:              0,
		ChainHW4:              0,
		ChainHW5:              0,
		ChainHW6:              1184,
		ChainHW7:              22,
		ChainHW8:              15,
		ChainHW9:              0,
		ChainHW10:             0,
		ChainHW11:             0,
		ChainHW12:             0,
		ChainHW13:             0,
		ChainHW14:             0,
		ChainHW15:             0,
		ChainHW16:             0,
		ChainRate1:            0,
		ChainRate2:            0,
		ChainRate3:            0,
		ChainRate4:            0,
		ChainRate5:            0,
		ChainRate6:            4536.24,
		ChainRate7:            4545.53,
		ChainRate8:            4548.77,
		ChainRate9:            0,
		ChainRate10:           0,
		ChainRate11:           0,
		ChainRate12:           0,
		ChainRate13:           0,
		ChainRate14:           0,
		ChainRate15:           0,
		ChainRate16:           0,
		ChainXtime6:           "{X49=5}",
		ChainXtime7:           "{}",
		ChainXtime8:           "{}",
		ChainOffside6:         0,
		ChainOffside7:         0,
		ChainOffside8:         0,
		ChainOpenCore6:        0,
		ChainOpenCore7:        0,
		ChainOpenCore8:        1,
		MinerVersion:          "16.8.1.3",
		MinerID:               "80749dc610358854",
	}
	ctx, finish := context.WithCancel(context.Background())
	port := getPort()
	go mockTCPServer(ctx, ip, port, testCaseValue)
	wait(1)
	miner := NewCGMiner(ip, port, minerTimeout)
	result, err := miner.Stats()
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(result, expected); diff != nil {
		t.Error(diff)
	}
	finish()
	wait(1)
}
func TestStats2BoardsS9(t *testing.T) {
	testCaseValue := getFixture("TestStats2BoardsS9.json")
	expected := &GenericStats{
		BMMiner:               "2.0.0",
		Miner:                 "16.8.1.3",
		CompileTime:           "Tue Aug 15 11:37:49 CST 2017",
		Type:                  "Antminer S9",
		Stats:                 0,
		ID:                    "BC50",
		Elapsed:               176363,
		Calls:                 0,
		Wait:                  0.0,
		Max:                   0.0,
		Min:                   99999999.0,
		Ghs5s:                 "9074.024",
		GhsAverage:            13590.24,
		MinerCount:            3,
		Frequency:             637,
		FanNum:                2,
		Fan1:                  0,
		Fan2:                  0,
		Fan3:                  6120,
		Fan4:                  0,
		Fan5:                  0,
		Fan6:                  4080,
		Fan7:                  0,
		Fan8:                  0,
		TempNum:               3,
		Temp1:                 0,
		Temp2:                 0,
		Temp3:                 0,
		Temp4:                 0,
		Temp5:                 0,
		Temp6:                 42,
		Temp7:                 38,
		Temp8:                 51,
		Temp9:                 0,
		Temp10:                0,
		Temp11:                0,
		Temp12:                0,
		Temp13:                0,
		Temp14:                0,
		Temp15:                0,
		Temp16:                0,
		Temp2_1:               0,
		Temp2_2:               0,
		Temp2_3:               0,
		Temp2_4:               0,
		Temp2_5:               0,
		Temp2_6:               57,
		Temp2_7:               53,
		Temp2_8:               66,
		Temp2_9:               0,
		Temp2_10:              0,
		Temp2_11:              0,
		Temp2_12:              0,
		Temp2_13:              0,
		Temp2_14:              0,
		Temp2_15:              0,
		Temp2_16:              0,
		Temp3_1:               0,
		Temp3_2:               0,
		Temp3_3:               0,
		Temp3_4:               0,
		Temp3_5:               0,
		Temp3_6:               0,
		Temp3_7:               0,
		Temp3_8:               0,
		Temp3_9:               0,
		Temp3_10:              0,
		Temp3_11:              0,
		Temp3_12:              0,
		Temp3_13:              0,
		Temp3_14:              0,
		Temp3_15:              0,
		Temp3_16:              0,
		FrequencyAvg1:         0.0,
		FrequencyAvg2:         0.0,
		FrequencyAvg3:         0.0,
		FrequencyAvg4:         0.0,
		FrequencyAvg5:         0.0,
		FrequencyAvg6:         627,
		FrequencyAvg7:         628.1400146484375,
		FrequencyAvg8:         627,
		FrequencyAvg9:         0.0,
		FrequencyAvg10:        0.0,
		FrequencyAvg11:        0.0,
		FrequencyAvg12:        0.0,
		FrequencyAvg13:        0.0,
		FrequencyAvg14:        0.0,
		FrequencyAvg15:        0.0,
		FrequencyAvg16:        0.0,
		TotalRateIdeal:        13501.2998046875,
		TotalFrequencyAvg:     627.3800048828125,
		TotalAcn:              189,
		TotalRate:             9074.01953125,
		ChainRateIdeal1:       0.0,
		ChainRateIdeal2:       0.0,
		ChainRateIdeal3:       0.0,
		ChainRateIdeal4:       0.0,
		ChainRateIdeal5:       0.0,
		ChainRateIdeal6:       4499.97998046875,
		ChainRateIdeal7:       4500.68994140625,
		ChainRateIdeal8:       4500.60986328125,
		ChainRateIdeal9:       0.0,
		ChainRateIdeal10:      0.0,
		ChainRateIdeal11:      0.0,
		ChainRateIdeal12:      0.0,
		ChainRateIdeal13:      0.0,
		ChainRateIdeal14:      0.0,
		ChainRateIdeal15:      0.0,
		ChainRateIdeal16:      0.0,
		TempMax:               51,
		DeviceHardwarePercent: 0.0004,
		NotMatchingWork:       2382,
		ChainAcn1:             0,
		ChainAcn2:             0,
		ChainAcn3:             0,
		ChainAcn4:             0,
		ChainAcn5:             0,
		ChainAcn6:             63,
		ChainAcn7:             63,
		ChainAcn8:             63,
		ChainAcn9:             0,
		ChainAcn10:            0,
		ChainAcn11:            0,
		ChainAcn12:            0,
		ChainAcn13:            0,
		ChainAcn14:            0,
		ChainAcn15:            0,
		ChainAcn16:            0,
		ChainAcs1:             "",
		ChainAcs2:             "",
		ChainAcs3:             "",
		ChainAcs4:             "",
		ChainAcs5:             "",
		ChainAcs6:             " oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo ooooooo",
		ChainAcs7:             " oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo ooooooo",
		ChainAcs8:             " oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo ooooooo",
		ChainAcs9:             "",
		ChainAcs10:            "",
		ChainAcs11:            "",
		ChainAcs12:            "",
		ChainAcs13:            "",
		ChainAcs14:            "",
		ChainAcs15:            "",
		ChainAcs16:            "",
		ChainHW1:              0,
		ChainHW2:              0,
		ChainHW3:              0,
		ChainHW4:              0,
		ChainHW5:              0,
		ChainHW6:              2335,
		ChainHW7:              28,
		ChainHW8:              19,
		ChainHW9:              0,
		ChainHW10:             0,
		ChainHW11:             0,
		ChainHW12:             0,
		ChainHW13:             0,
		ChainHW14:             0,
		ChainHW15:             0,
		ChainHW16:             0,
		ChainRate1:            0,
		ChainRate2:            0,
		ChainRate3:            0,
		ChainRate4:            0,
		ChainRate5:            0,
		ChainRate6:            4524.63,
		ChainRate7:            4549.39,
		ChainRate8:            0,
		ChainRate9:            0,
		ChainRate10:           0,
		ChainRate11:           0,
		ChainRate12:           0,
		ChainRate13:           0,
		ChainRate14:           0,
		ChainRate15:           0,
		ChainRate16:           0,
		ChainXtime6:           "{X10=7}",
		ChainXtime7:           "{}",
		ChainXtime8:           "{X0=1,X1=1,X2=1,X3=1,X4=1,X5=1,X6=1,X7=1,X8=1,X9=1,X10=1,X11=1,X12=1,X13=1,X14=1,X15=1,X16=1,X17=1,X18=1,X19=1,X20=1,X21=1,X22=1,X23=1,X24=1,X25=1,X26=1,X27=1,X28=1,X29=1,X30=1,X31=1,X32=1,X33=1,X34=1,X35=1,X36=1,X37=1,X38=1,X39=1,X40=1,X41=1,X42=1,X43=1,X44=1,X45=1,X46=1,X47=1,X48=1,X49=1,X50=1,X51=1,X52=1,X53=1,X54=1,X55=1,X56=1,X57=1,X58=1,X59=1,X60=1,X61=1,X62=1}",
		ChainOffside6:         0,
		ChainOffside7:         0,
		ChainOffside8:         0,
		ChainOpenCore6:        0,
		ChainOpenCore7:        1,
		ChainOpenCore8:        1,
		MinerVersion:          "16.8.1.3",
		MinerID:               "80749dc610358854",
	}
	ctx, finish := context.WithCancel(context.Background())
	port := getPort()
	go mockTCPServer(ctx, ip, port, testCaseValue)
	wait(1)
	miner := NewCGMiner(ip, port, minerTimeout)
	result, err := miner.Stats()
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(result, expected); diff != nil {
		t.Error(diff)
	}
	finish()
	wait(1)
}
func TestStatsS9(t *testing.T) {
	testCaseValue := getFixture("TestStatsS9.json")
	expected := &StatsS9{
		BMMiner:               "2.0.0",
		Miner:                 "16.8.1.3",
		CompileTime:           "Fri Nov 17 17:57:49 CST 2017",
		Type:                  "Antminer S9",
		Stats:                 0,
		ID:                    "BC50",
		Elapsed:               100986,
		Calls:                 0,
		Wait:                  0.0,
		Max:                   0.0,
		Min:                   99999999.0,
		Ghs5s:                 13630.55,
		GhsAverage:            13569.28,
		MinerCount:            3,
		Frequency:             637,
		FanNum:                2,
		Fan3:                  4080,
		Fan6:                  4080,
		TempNum:               3,
		Temp6:                 56,
		Temp7:                 52,
		Temp8:                 56,
		Temp2_6:               71,
		Temp2_7:               67,
		Temp2_8:               71,
		FrequencyAvg6:         627.57,
		FrequencyAvg7:         627.76,
		FrequencyAvg8:         627.09,
		TotalRateIdeal:        13501.38,
		TotalFrequencyAvg:     627.47,
		TotalAcn:              189,
		TotalRate:             13630.54,
		ChainRateIdeal6:       4500.33,
		ChainRateIdeal7:       4500.38,
		ChainRateIdeal8:       4500.66,
		TempMax:               56,
		DeviceHardwarePercent: 0.0004,
		NotMatchingWork:       1222,
		ChainAcn6:             63,
		ChainAcn7:             63,
		ChainAcn8:             63,
		ChainAcs6:             " oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo ooooooo",
		ChainAcs7:             " oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo ooooooo",
		ChainAcs8:             " oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo ooooooo",
		ChainHW6:              1184,
		ChainHW7:              22,
		ChainHW8:              15,
		ChainRate6:            4536.24,
		ChainRate7:            4545.53,
		ChainRate8:            4548.77,
		ChainXtime6:           "{X49=5}",
		ChainXtime7:           "{}",
		ChainXtime8:           "{}",
		ChainOffside6:         0,
		ChainOffside7:         0,
		ChainOffside8:         0,
		ChainOpenCore6:        0,
		ChainOpenCore7:        0,
		ChainOpenCore8:        1,
		MinerVersion:          "16.8.1.3",
		MinerID:               "80749dc610358854",
	}
	ctx, finish := context.WithCancel(context.Background())
	port := getPort()
	go mockTCPServer(ctx, ip, port, testCaseValue)
	wait(1)
	miner := NewCGMiner(ip, port, minerTimeout)
	stats, err := miner.Stats()
	if err != nil {
		t.Fatal(err)
	}
	result, err := stats.S9()
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(result, expected); diff != nil {
		t.Error(diff)
	}
	finish()
	wait(1)
}
func TestStatsL3plus(t *testing.T) {
	testCaseValue := getFixture("TestStatsL3plus.json")
	expected := &StatsL3{
		CGMiner:               "4.9.0",
		Miner:                 "1.0.1.3",
		CompileTime:           "Fri Aug 25 17:28:57 CST 2017",
		Type:                  "Antminer L3+",
		Stats:                 0,
		ID:                    "L30",
		Elapsed:               204806,
		Calls:                 0,
		Wait:                  0.0,
		Max:                   0.0,
		Min:                   99999999.0,
		Ghs5s:                 580.455,
		GhsAverage:            577.49,
		MinerCount:            4,
		Frequency:             444,
		FanNum:                2,
		Fan1:                  5250,
		Fan2:                  4260,
		TempNum:               4,
		Temp1:                 40,
		Temp2:                 39,
		Temp3:                 38,
		Temp4:                 36,
		Temp2_1:               48,
		Temp2_2:               48,
		Temp2_3:               47,
		Temp2_4:               45,
		Temp3_1:               0,
		Temp3_2:               0,
		Temp3_3:               0,
		Temp3_4:               0,
		Temp4_1:               0,
		Temp4_2:               0,
		Temp4_3:               0,
		Temp4_4:               0,
		TempMax:               40,
		DeviceHardwarePercent: 0.0,
		NotMatchingWork:       6483,
		ChainAcn1:             72,
		ChainAcn2:             72,
		ChainAcn3:             72,
		ChainAcn4:             72,
		ChainAcs1:             " oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo",
		ChainAcs2:             " oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo",
		ChainAcs3:             " oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo",
		ChainAcs4:             " oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo",
		ChainHW1:              129,
		ChainHW2:              0,
		ChainHW3:              0,
		ChainHW4:              6354,
		ChainRate1:            145.57,
		ChainRate2:            144.69,
		ChainRate3:            145.10,
		ChainRate4:            145.09,
	}
	ctx, finish := context.WithCancel(context.Background())
	port := getPort()
	go mockTCPServer(ctx, ip, port, testCaseValue)
	wait(1)
	miner := NewCGMiner(ip, port, minerTimeout)
	stats, err := miner.Stats()
	if err != nil {
		t.Fatal(err)
	}
	result, err := stats.L3()
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(result, expected); diff != nil {
		t.Error(diff)
	}
	finish()
	wait(1)
}
func TestStatsD3(t *testing.T) {
	testCaseValue := getFixture("TestStatsD3.json")
	expected := &StatsD3{
		CGMiner:               "4.9.0",
		Miner:                 "1.0.0.6",
		CompileTime:           "Thu Aug 31 13:38:33 CST 2017",
		Type:                  "Antminer D3",
		Stats:                 0,
		ID:                    "D10",
		Elapsed:               122893,
		Calls:                 0,
		Wait:                  0.0,
		Max:                   0.0,
		Min:                   99999999.0,
		Ghs5s:                 17052.3,
		GhsAverage:            17541.89,
		MinerCount:            3,
		Frequency:             487,
		FanNum:                2,
		Fan1:                  5190,
		Fan2:                  5250,
		TempNum:               3,
		Temp1:                 60,
		Temp2:                 63,
		Temp3:                 64,
		Temp2_1:               75,
		Temp2_2:               79,
		Temp2_3:               80,
		TempMax:               64,
		DeviceHardwarePercent: 0.0001,
		NotMatchingWork:       4,
		ChainAcn1:             60,
		ChainAcn2:             60,
		ChainAcn3:             60,
		ChainAcs1:             " oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooo",
		ChainAcs2:             " oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooo",
		ChainAcs3:             " oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooo",
		ChainHW1:              0,
		ChainHW2:              3,
		ChainHW3:              1,
		ChainRate1:            5712.75,
		ChainRate2:            5712.88,
		ChainRate3:            5626.71,
	}
	ctx, finish := context.WithCancel(context.Background())
	port := getPort()
	go mockTCPServer(ctx, ip, port, testCaseValue)
	wait(1)
	miner := NewCGMiner(ip, port, minerTimeout)
	stats, err := miner.Stats()
	if err != nil {
		t.Fatal(err)
	}
	result, err := stats.D3()
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(result, expected); diff != nil {
		t.Error(diff)
	}
	finish()
	wait(1)
}
func TestStatsT9(t *testing.T) {
	testCaseValue := getFixture("TestStatsT9plus.json")
	expected := &StatsT9{
		BMMiner:               "2.0.0",
		Miner:                 "16.0.1.3",
		CompileTime:           "Fri Nov 24 23:19:16 EST 2017",
		Type:                  "Antminer T9+",
		Stats:                 0,
		ID:                    "BC50",
		Elapsed:               38439,
		Calls:                 0,
		Wait:                  0.0,
		Max:                   0.0,
		Min:                   99999999.0,
		Ghs5s:                 10585.82,
		GhsAverage:            10639.33,
		MinerCount:            9,
		Frequency:             575,
		FanNum:                2,
		Fan3:                  5880,
		Fan6:                  5400,
		TempNum:               9,
		Temp2:                 70,
		Temp3:                 73,
		Temp4:                 70,
		Temp9:                 70,
		Temp10:                70,
		Temp11:                73,
		Temp12:                73,
		Temp13:                70,
		Temp14:                70,
		Temp2_2:               85,
		Temp2_3:               88,
		Temp2_4:               85,
		Temp2_9:               85,
		Temp2_10:              85,
		Temp2_11:              88,
		Temp2_12:              88,
		Temp2_13:              85,
		Temp2_14:              85,
		FrequencyAvg2:         568.38,
		FrequencyAvg3:         568.38,
		FrequencyAvg4:         568.38,
		FrequencyAvg9:         568.38,
		FrequencyAvg10:        568.38,
		FrequencyAvg11:        568.38,
		FrequencyAvg12:        569.16,
		FrequencyAvg13:        569.16,
		FrequencyAvg14:        569.16,
		TotalRateIdeal:        10501.79,
		TotalFrequencyAvg:     568.64,
		TotalAcn:              162,
		TotalRate:             10585.82,
		ChainRateIdeal2:       1166.33,
		ChainRateIdeal3:       1166.33,
		ChainRateIdeal4:       1166.33,
		ChainRateIdeal9:       1166.33,
		ChainRateIdeal10:      1166.33,
		ChainRateIdeal11:      1166.33,
		ChainRateIdeal12:      1167.93,
		ChainRateIdeal13:      1167.93,
		ChainRateIdeal14:      1167.93,
		TempMax:               73,
		DeviceHardwarePercent: 0.0008,
		NotMatchingWork:       781,
		ChainAcn2:             18,
		ChainAcn3:             18,
		ChainAcn4:             18,
		ChainAcn9:             18,
		ChainAcn10:            18,
		ChainAcn11:            18,
		ChainAcn12:            18,
		ChainAcn13:            18,
		ChainAcn14:            18,
		ChainAcs2:             " oooooooo oooooooo oo",
		ChainAcs3:             " oooooooo oooooooo oo",
		ChainAcs4:             " oooooooo oooooooo oo",
		ChainAcs9:             " oooooooo oooooooo oo",
		ChainAcs10:            " oooooooo oooooooo oo",
		ChainAcs11:            " oooooooo oooooooo oo",
		ChainAcs12:            " oooooooo oooooooo oo",
		ChainAcs13:            " oooooooo oooooooo oo",
		ChainAcs14:            " oooooooo oooooooo oo",
		ChainHW2:              90,
		ChainHW3:              2,
		ChainHW4:              0,
		ChainHW9:              1,
		ChainHW10:             49,
		ChainHW11:             148,
		ChainHW12:             302,
		ChainHW13:             134,
		ChainHW14:             55,
		ChainRate2:            1163.84,
		ChainRate3:            1174.19,
		ChainRate4:            1173.73,
		ChainRate9:            1171.12,
		ChainRate10:           1176.45,
		ChainRate11:           1183.73,
		ChainRate12:           1182.49,
		ChainRate13:           1189.15,
		ChainRate14:           1171.12,
		ChainXtime2:           "{}",
		ChainXtime3:           "{}",
		ChainXtime4:           "{}",
		ChainXtime9:           "{}",
		ChainXtime10:          "{}",
		ChainXtime11:          "{}",
		ChainXtime12:          "{X7=1}",
		ChainXtime13:          "{}",
		ChainXtime14:          "{}",
		ChainOffside2:         0,
		ChainOffside3:         0,
		ChainOffside4:         0,
		ChainOffside9:         0,
		ChainOffside10:        0,
		ChainOffside11:        0,
		ChainOffside12:        0,
		ChainOffside13:        0,
		ChainOffside14:        0,
		ChainOpenCore2:        0,
		ChainOpenCore3:        0,
		ChainOpenCore4:        0,
		ChainOpenCore9:        0,
		ChainOpenCore10:       0,
		ChainOpenCore11:       0,
		ChainOpenCore12:       0,
		ChainOpenCore13:       0,
		ChainOpenCore14:       0,
		MinerVersion:          "16.0.1.3",
		MinerID:               "8134f54c6880881c",
	}
	ctx, finish := context.WithCancel(context.Background())
	port := getPort()
	go mockTCPServer(ctx, ip, port, testCaseValue)
	wait(1)
	miner := NewCGMiner(ip, port, minerTimeout)
	stats, err := miner.Stats()
	if err != nil {
		t.Fatal(err)
	}
	result, err := stats.T9()
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(result, expected); diff != nil {
		t.Error(diff)
	}
	finish()
	wait(1)
}
func TestStatsS7(t *testing.T) {
	testCaseValue := getFixture("TestStatsS7.json")
	expected := &StatsS7{
		CGMiner:     "4.8.0",
		Miner:       "3.5.3.0",
		CompileTime: "Mon May 23 14:53:27 CST 2016",
		Type:        "Antminer S7",
		Stats:       0,
		ID:          "BTM0",
		Elapsed:     3755,
		Calls:       0,
		Wait:        0.0,
		Max:         0.0,
		Min:         99999999.0,
		Ghs5s:       4697.23,
		GhsAverage:  4729.14,
		Baud:        115200,
		MinerCount:  3,
		AsicCount:   8,
		Timeout:     5,
		Frequency:   700,
		Voltage:     0.706,
		HWv1:        3,
		HWv2:        5,
		HWv3:        3,
		HWv4:        0,
		// WHY ??? 6?
		FanNum:                6,
		Fan1:                  3720,
		Fan3:                  3600,
		TempNum:               3,
		Temp1:                 58,
		Temp2:                 58,
		Temp3:                 54,
		TempAvg:               56,
		TempMax:               59,
		DeviceHardwarePercent: 0.0005,
		NotMatchingWork:       19,
		ChainAcn1:             45,
		ChainAcn2:             45,
		ChainAcn3:             45,
		ChainAcs1:             "oooooooo ooooo oooooooo oooooooo oooooooo oooooooo ",
		ChainAcs2:             "oooooooo ooooo oooooooo oooooooo oooooooo oooooooo ",
		ChainAcs3:             "oooooooo ooooo oooooooo oooooooo oooooooo oooooooo ",
		USBPipe:               0,
	}
	ctx, finish := context.WithCancel(context.Background())
	port := getPort()
	go mockTCPServer(ctx, ip, port, testCaseValue)
	wait(1)
	miner := NewCGMiner(ip, port, minerTimeout)
	stats, err := miner.Stats()
	if err != nil {
		t.Fatal(err)
	}
	result, err := stats.S7()
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(result, expected); diff != nil {
		t.Error(diff)
	}
	finish()
	wait(1)
}
