package internal

import (
	"context"
	"fmt"
	"net"
	"time"

	tester_utils "github.com/codecrafters-io/tester-utils"
	logger "github.com/codecrafters-io/tester-utils/logger"
	"github.com/miekg/dns"
)

const (
	SERVER_ADDR    = "127.0.0.1:2053"
	RESOLVER_ADDR  = "127.0.0.1:5354"
	DEFAULT_DOMAIN = "codecrafters.io."
	DEFAULT_PKT_ID = 1234
)

func testInit(stageHarness *tester_utils.StageHarness) error {
	cancels, err := startDNSServers(stageHarness)
	for _, cancel := range cancels {
		defer cancel()
	}
	if err != nil {
		return err
	}

	return nil
}

func startDNSServers(stageHarness *tester_utils.StageHarness) ([]context.CancelFunc, error) {
	var cancels []context.CancelFunc
	logger := stageHarness.Logger
	ctx, cancel := context.WithCancel(context.Background())
	cancels = append(cancels, cancel)

	go startDNSServer(ctx, logger, RESOLVER_ADDR)

	logger.Infof("Starting DNS server on %s", SERVER_ADDR)
	b := NewDnsServerBinary(stageHarness)
	if err := b.Run("--resolver", RESOLVER_ADDR); err != nil {
		return cancels, err
	}

	if err := retryDialUntilSuccess(logger); err != nil {
		return cancels, err
	}
	return cancels, nil
}

func retryDialUntilSuccess(logger *logger.Logger) error {
	var e error
	retries := 0
	logger.Infof("Connecting to %s using UDP", SERVER_ADDR)
	for retries < 5 {
		if retries > 1 {
			logger.Infof("Failed to connect to port 2053, retrying in 1s")
		}

		conn, err := net.Dial("udp", SERVER_ADDR)
		if err != nil {
			e = err
			continue
		}
		defer conn.Close()

		request, err := getDnsMsgBytes()
		if err != nil {
			e = err
			continue
		}
		_, err = conn.Write(request)
		if err != nil {
			e = err
			continue
		}

		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		responseBuffer := make([]byte, 1024)
		_, err = conn.Read(responseBuffer)

		if err != nil {
			e = err
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() {
				e = nil
				logger.Debugf("No ICMP response, port is likely open.")
				break
			}
		} else {
			e = nil
			logger.Debugf("Got a response")
			break
		}

		retries += 1
		time.Sleep(1 * time.Second)
	}
	if e != nil {
		return e
	}

	return nil
}

func getDnsMsgBytes() ([]byte, error) {
	packet := new(dns.Msg)
	packet.Id = dns.Id()
	packet.SetQuestion(DEFAULT_DOMAIN, dns.TypeA)

	buf, err := packet.Pack()
	if err != nil {
		return nil, fmt.Errorf("Error encoding DNS message: %s\n", err)
	}

	return buf, nil
}
