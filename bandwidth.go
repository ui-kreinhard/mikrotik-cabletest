package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type BandWidthResult struct {
	Status            string
	Duration          string
	TxCurrent         float64
	Tx10SecondAverage float64
	TxTotalAverage    float64
	RxCurrent         float64
	Rx10SecondAverage float64
	RxTotalAverage    float64
	LostPackets       int
	RandomData        string
	Direction         string
	TxSize            int
	RxSize            int
	ConnectionCount   int
	LocalCPULoad      string
	RemoteCPULoad     string
}

func (b *BandWidthResult) CheckPackageLoss() bool {
	return b.LostPackets < 2500
}

func (b *BandWidthResult) CheckTxAverage() bool {
	return b.Tx10SecondAverage > 150
}

func (b *BandWidthResult) CheckRXAverage() bool {
	return b.Rx10SecondAverage > 150
}

func (b *BandWidthResult) Report() string {
	ret := fmt.Sprintln("Package loss", toOKNOK(b.CheckPackageLoss()), b.LostPackets)
	ret += fmt.Sprintln("TX-Average", toOKNOK(b.CheckTxAverage()), b.Tx10SecondAverage)
	ret += fmt.Sprintln("RX-Average", toOKNOK(b.CheckRXAverage()), b.Rx10SecondAverage)
	return ret
}

func toOKNOK(in bool) string {
	if in {
		return "OK"
	}
	return "NOK"
}

func parseBandwidthResult(input string) (BandWidthResult, error) {

	var result BandWidthResult

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "status:"):
			result.Status = strings.TrimSpace(strings.TrimPrefix(line, "status:"))
		case strings.HasPrefix(line, "duration:"):
			result.Duration = strings.TrimSpace(strings.TrimPrefix(line, "duration:"))
		case strings.HasPrefix(line, "tx-current:"):
			value, err := extractFloat(line)
			if err != nil {
				return result, err
			}
			result.TxCurrent = value
		case strings.HasPrefix(line, "tx-10-second-average:"):
			value, err := extractFloat(line)
			if err != nil {
				return result, err
			}
			result.Tx10SecondAverage = value
		case strings.HasPrefix(line, "tx-total-average:"):
			value, err := extractFloat(line)
			if err != nil {
				return result, err
			}
			result.TxTotalAverage = value
		case strings.HasPrefix(line, "rx-current:"):
			value, err := extractFloat(line)
			if err != nil {
				return result, err
			}
			result.RxCurrent = value
		case strings.HasPrefix(line, "rx-10-second-average:"):
			value, err := extractFloat(line)
			if err != nil {
				return result, err
			}
			result.Rx10SecondAverage = value
		case strings.HasPrefix(line, "rx-total-average:"):
			value, err := extractFloat(line)
			if err != nil {
				return result, err
			}
			result.RxTotalAverage = value
		case strings.HasPrefix(line, "lost-packets:"):
			value, err := extractInt(line)
			if err != nil {
				return result, err
			}
			result.LostPackets = value
		case strings.HasPrefix(line, "random-data:"):
			result.RandomData = strings.TrimSpace(strings.TrimPrefix(line, "random-data:"))
		case strings.HasPrefix(line, "direction:"):
			result.Direction = strings.TrimSpace(strings.TrimPrefix(line, "direction:"))
		case strings.HasPrefix(line, "tx-size:"):
			value, err := extractInt(line)
			if err != nil {
				return result, err
			}
			result.TxSize = value
		case strings.HasPrefix(line, "rx-size:"):
			value, err := extractInt(line)
			if err != nil {
				return result, err
			}
			result.RxSize = value
		case strings.HasPrefix(line, "connection-count:"):
			value, err := extractInt(line)
			if err != nil {
				return result, err
			}
			result.ConnectionCount = value
		case strings.HasPrefix(line, "local-cpu-load:"):
			result.LocalCPULoad = strings.TrimSpace(strings.TrimPrefix(line, "local-cpu-load:"))
		case strings.HasPrefix(line, "remote-cpu-load:"):
			result.RemoteCPULoad = strings.TrimSpace(strings.TrimPrefix(line, "remote-cpu-load:"))
		}
	}

	return result, nil
}

func extractFloat(line string) (float64, error) {
	segments := strings.Split(line, ":")
	if len(segments) < 2 {
		return 0, errors.New("invalid segment length")
	}
	re := regexp.MustCompile(`[+-]?([0-9]*[.])?[0-9]+`)
	match := re.FindString(segments[1])
	value, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func extractInt(line string) (int, error) {
	re := regexp.MustCompile(`[0-9]+`)
	match := re.FindString(line)
	value, err := strconv.Atoi(match)
	if err != nil {
		return 0, err
	}
	return value, nil
}
