package internal

import (
	"fmt"

	tester_utils "github.com/codecrafters-io/tester-utils"
	logger "github.com/codecrafters-io/tester-utils/logger"
	"github.com/miekg/dns"
)

func testHeaderParsing(stageHarness *tester_utils.StageHarness) error {
	b := NewDnsServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	queryDomain := "codecrafters.io."

	if _, err := sendQuery(stageHarness.Logger, queryDomain, dns.TypeA); err != nil {
		return err
	}

	return nil
}

func sendQuery(logger *logger.Logger, queryDomain string, recordType uint16) (*dns.Msg, error) {
	c := new(dns.Client)

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(queryDomain), recordType)
	logger.Infof("Querying `%s` record for %s", recordTypeToString(recordType), queryDomain)
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
