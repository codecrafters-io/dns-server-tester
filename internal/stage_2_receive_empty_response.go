package internal

import (
	"fmt"
	"strings"

	tester_utils "github.com/codecrafters-io/tester-utils"
	logger "github.com/codecrafters-io/tester-utils/logger"
	"github.com/miekg/dns"
)

func testReceiveEmptyResponse(stageHarness *tester_utils.StageHarness) error {
	b := NewDnsServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}
	logger := stageHarness.Logger
	if err := retryDialUntilSuccess(logger); err != nil {
		return err
	}

	queryDomain := "codecrafters.io."
	packetIdentifier := 1234

	_, err := sendDNSQueryWithId(logger, uint16(packetIdentifier), queryDomain)
	if err != nil {
		return fmt.Errorf("Error sending DNS query: %s\n", err)
	}

	return nil
}

func sendDNSQueryWithId(logger *logger.Logger, id uint16, queryDomain string) (*dns.Msg, error) {
	c := new(dns.Client)

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(queryDomain), dns.TypeA)
	msg.Id = id
	logger.Infof("Querying `A` record for %s", queryDomain)
	logger.Debugf("Sending Request: (Messages with >>> prefix are part of this log)")
	logDnsMsg(logger, msg)

	response, _, err := c.Exchange(msg, SERVER_ADDR)
	if err != nil {
		return nil, fmt.Errorf("DNS query failed: %s.\nIf you are seeing this after a while then it is likely that your server is not responding with appropriate id", err)
	}

	logger.Debugf("Received Response: (Messages with >>> prefix are part of this log)")
	logDnsMsg(logger, response)

	return response, nil
}

func logDnsMsg(logger *logger.Logger, msg *dns.Msg) {
	for _, line := range strings.Split(msg.String(), "\n") {
		logger.Debugf(">>> %s", line)
	}
}
