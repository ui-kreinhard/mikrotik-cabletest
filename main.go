package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/helloyi/go-sshclient"
	"github.com/ui-kreinhard/mikrotik-cabletest/config"
)

func main() {
	config := config.LoadConfig()
	fmt.Println("Using config values:")
	fmt.Println(config)
	fmt.Println("Connecting to", config.SwitchIp, "port", config.SshPort)
	client, err := sshclient.DialWithPasswd(config.SwitchIp+":"+strconv.Itoa(config.SshPort), config.SwitchUsername, config.SwitchPassword)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
	fmt.Println("connected")
	instruct("Insert looped back casble in port " + config.PortToTest)
	fmt.Println("Starting wiring test")
	time.Sleep(10 * time.Second)
	resultCableTest := cableTestPort(client, config.PortToTest)
	if !resultCableTest.IsNormal() {
		fmt.Println("Cable test failed - check for wiring errors")
		fmt.Printf("%+v\n", resultCableTest)
		return
	}
	fmt.Println("Wiring of cable OK")
	instruct("Remove Cable loop and add other switch for speed test")
	fmt.Println("Starting speed test")
	time.Sleep(10 * time.Second)
	resultCableTestWithLink := cableTestPort(client, config.PortToTest)
	if !resultCableTestWithLink.LinkState {
		fmt.Println("no link established", resultCableTestWithLink)
	}
	fmt.Println("link established between switches")
	bandWidthResult := bandWidthTest(client)
	fmt.Println(bandWidthResult.Report())
}

func instruct(message string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(message + " - Enter to continue, ctrl+c to cancel")
	reader.ReadString('\n')
}

func bandWidthTest(client *sshclient.Client) BandWidthResult {
	out, _ := client.Cmd("/tool bandwidth-test 192.168.88.1 password=admin direction=both duration=30s").Output()
	result, err := parseBandwidthResult(string(out))
	if err != nil {
		log.Fatalln(err)
	}
	return result
}

func cableTestPort(client *sshclient.Client, port string) CableTestReport {
	out, err := client.Cmd("/interface ethernet cable-test " + port + " once").Output()
	if err != nil {
		log.Fatalln(err)
	}
	crt := CableTestReport{}
	err = crt.parse(string(out))
	if err != nil {
		log.Fatalln(err)
	}
	return crt
}

type CableTestReport struct {
	LinkState bool
	Pair1     CablePair
	Pair2     CablePair
	Pair3     CablePair
	Pair4     CablePair
}

func (c *CableTestReport) IsNormal() bool {
	return c.Pair1.normal() && c.Pair2.normal() && c.Pair3.normal() && c.Pair4.normal()
}

func (c *CableTestReport) parse(raw string) error {
	lines := strings.Split(raw, "\n")
	if len(lines) < 3 {
		return errors.New("Invalid line count " + raw)
	}
	linkState, err := extractLinkState(lines[1])
	if err != nil {
		return err
	}
	c.LinkState = linkState
	if linkState {
		return nil
	}
	pairs, err := extractCablePairs(lines[2])
	if err != nil {
		return err
	}
	c.Pair1 = pairs[0]
	c.Pair2 = pairs[1]
	c.Pair3 = pairs[2]
	c.Pair4 = pairs[3]
	return nil
}

func extractCablePairs(rawLine string) (pairs []CablePair, err error) {
	outerSegments := strings.Split(rawLine, ": ")
	if len(outerSegments) < 2 {
		err = errors.New("invalid segment count " + rawLine)
		return
	}
	innerSegments := strings.Split(outerSegments[1], ",")
	for _, innerSegment := range innerSegments {
		pair, err1 := extractCablePair(innerSegment)
		if err1 != nil {
			err = err1
			return
		}
		pairs = append(pairs, pair)
	}
	return
}

func extractCablePair(rawLine string) (pair CablePair, err error) {
	segment := strings.Split(rawLine, ":")
	if len(segment) < 2 {
		err = errors.New("invalid segment count " + rawLine)
		return
	}
	pair.Type = segment[0]
	length, err := strconv.ParseInt(strings.TrimSpace(segment[1]), 10, 32)
	if err != nil {
		return
	}
	pair.Length = int(length)

	return
}

func extractLinkState(rawLine string) (bool, error) {
	segments := strings.Split(rawLine, ": ")
	if len(segments) < 2 {
		return false, errors.New("invalid segment count " + rawLine)
	}
	if strings.TrimSpace(segments[1]) == "no-link" {
		return false, nil
	}
	return true, nil
}

type CablePair struct {
	Type   string
	Length int
}

func (c *CablePair) normal() bool {
	return c.Type == "normal"
}
