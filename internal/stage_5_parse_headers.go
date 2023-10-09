package internal

import (
	"fmt"

	tester_utils "github.com/codecrafters-io/tester-utils"
	logger "github.com/codecrafters-io/tester-utils/logger"
	"github.com/miekg/dns"
)

func testHeaderParsing(stageHarness *tester_utils.StageHarness) error {
	cancels, err := startDNSServers(stageHarness)
	for _, cancel := range cancels {
		defer cancel()
	}
	if err != nil {
		return err
	}

	logger := stageHarness.Logger

	if _, err := sendQuery(logger, DEFAULT_DOMAIN, dns.TypeA); err != nil {
		return err
	}

	return nil
}

func sendQuery(logger *logger.Logger, queryDomain string, recordType uint16) (*dns.Msg, error) {
	c := new(dns.Client)

	request := new(dns.Msg)
	request.SetQuestion(dns.Fqdn(queryDomain), recordType)
	logger.Infof("Querying `%s` record for %s", recordTypeToString(recordType), queryDomain)
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

func recordTypeToString(recordType uint16) string {
	switch recordType {
	case dns.TypeA:
		return "A"
	case dns.TypeAAAA:
		return "AAAA"
	case dns.TypeCNAME:
		return "CNAME"
	case dns.TypeMX:
		return "MX"
	case dns.TypeNS:
		return "NS"
	default:
		return fmt.Sprintf("%d", recordType)
	}
}
