package internal

import (
	"fmt"
	"strings"

	tester_utils "github.com/codecrafters-io/tester-utils"
	logger "github.com/codecrafters-io/tester-utils/logger"
	"github.com/miekg/dns"
)

func testReceiveEmptyResponse(stageHarness *tester_utils.StageHarness) error {
	cancels, err := startDNSServers(stageHarness)
	for _, cancel := range cancels {
		defer cancel()
	}
	if err != nil {
		return err
	}
	logger := stageHarness.Logger

	_, err = sendDNSQueryWithId(logger, uint16(DEFAULT_PKT_ID), DEFAULT_DOMAIN)
	if err != nil {
		return fmt.Errorf("Error sending DNS query: %s\n", err)
	}

	return nil
}

func sendDNSQueryWithId(logger *logger.Logger, id uint16, queryDomain string) (*dns.Msg, error) {
	c := new(dns.Client)

	request := new(dns.Msg)
	request.SetQuestion(dns.Fqdn(queryDomain), dns.TypeA)
	request.Id = id
	logger.Infof("Querying `A` record for %s", queryDomain)
	logger.Debugf("Sending Request: (Messages with >>> prefix are part of this log)")
	logDNSPacket(logger, request)

	response, _, err := c.Exchange(request, SERVER_ADDR)
	if err != nil {
		return nil, fmt.Errorf("DNS query failed: %s.\nIf you are seeing this after a while then it is likely that your server is not responding with appropriate id", err)
	}

	logger.Debugf("Received Response: (Messages with >>> prefix are part of this log)")
	logDNSPacket(logger, response)

	return response, nil
}

func logDNSPacket(logger *logger.Logger, packet *dns.Msg) {
	for _, line := range strings.Split(packet.String(), "\n") {
		logger.Debugf(">>> %s", line)
	}
}
