package internal

import (
	"fmt"
	"net"
	"time"

	tester_utils "github.com/codecrafters-io/tester-utils"
	logger "github.com/codecrafters-io/tester-utils/logger"
	"github.com/miekg/dns"
)

const (
	SERVER_ADDR = "127.0.0.1:2053"
)

// Example from the grep course
func testInit(stageHarness *tester_utils.StageHarness) error {

	b := NewDnsServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	err := retryDialUntilSuccess(logger)
	if err != nil {
		logger.Infof("All retries failed.")
		return err
	}

	logger.Infof("Success.")
	return nil
}

func retryDialUntilSuccess(logger *logger.Logger) error {
	var err error
	retries := 0
	logger.Infof("Connecting to %s using UDP", SERVER_ADDR)
	for retries < 5 {
		if retries > 1 {
			logger.Infof("Failed to connect to port 2053, retrying in 1s")
		}

		conn, err := net.Dial("udp", SERVER_ADDR)
		if err != nil {
			continue
		}
		defer conn.Close()

		msg, err := getDnsMsgBytes()
		if err != nil {
			continue
		}
		_, err = conn.Write(msg)
		if err != nil {
			continue
		}

		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		buffer := make([]byte, 1024)
		_, err = conn.Read(buffer)

		if err != nil {
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() {
				logger.Debugf("No ICMP response, port is likely open.")
				break
			}
		} else {
			logger.Debugf("Got a response")
			break
		}

		retries += 1
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		return err
	}

	return nil
}

func getDnsMsgBytes() ([]byte, error) {
	msg := new(dns.Msg)
	msg.Id = dns.Id()
	msg.SetQuestion("codecrafters.io.", dns.TypeA)

	buf, err := msg.Pack()
	if err != nil {
		return nil, fmt.Errorf("Error encoding DNS message: %s\n", err)
	}

	return buf, nil
}
