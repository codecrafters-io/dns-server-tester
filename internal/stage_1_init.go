package internal

import (
	"context"
	"fmt"
	"net"
	"time"

	logger "github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
	"github.com/miekg/dns"
)

const (
	SERVER_ADDR    = "127.0.0.1:2053"
	RESOLVER_ADDR  = "127.0.0.1:5354"
	DEFAULT_DOMAIN = "codecrafters.io."
	DEFAULT_PKT_ID = 1234
)

func testInit(stageHarness *test_case_harness.TestCaseHarness) error {
	cancels, err := startDNSServers(stageHarness)
	for _, cancel := range cancels {
		defer cancel()
	}
	if err != nil {
		return err
	}

	return nil
}

func startDNSServers(stageHarness *test_case_harness.TestCaseHarness) ([]context.CancelFunc, error) {
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

	if err := retryUDPCommunicationUntilSuccess(logger, b); err != nil {
		return cancels, err
	}
	return cancels, nil
}

func retryUDPCommunicationUntilSuccess(logger *logger.Logger, b *DnsServerBinary) error {
	var e error
	retries := 0
	logger.Infof("Connecting to %s using UDP", SERVER_ADDR)
	for retries < 5 {
		if retries > 2 {
			if b.HasExited() {
				return fmt.Errorf("Looks like your program has terminated. A DNS server is expected to be a long-running process.")
			}
			logger.Infof("Did not receive response from DNS server, retrying")
			logger.Errorf("%s", e)

		}

		err := doUDPCommunication()
		if err == nil {
			return nil
		}
		e = err
		retries += 1
		time.Sleep(1 * time.Second)
	}

	return e
}

func doUDPCommunication() error {
	conn, err := net.Dial("udp", SERVER_ADDR)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = sendDnsPacket(conn)
	if err != nil {
		return err
	}
	err = readAnyResponse(conn)
	if err != nil {
		return err
	}
	return nil
}

func sendDnsPacket(conn net.Conn) error {
	request, err := getDnsMsgBytes()
	if err != nil {
		return err
	}
	_, err = conn.Write(request)
	return err
}

func readAnyResponse(conn net.Conn) error {
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	responseBuffer := make([]byte, 1024)
	_, err := conn.Read(responseBuffer)
	return err
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
