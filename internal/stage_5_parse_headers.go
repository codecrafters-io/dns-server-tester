package internal

import (
	"fmt"

	logger "github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
	"github.com/miekg/dns"
)

func testHeaderParsing(stageHarness *test_case_harness.TestCaseHarness) error {
	cancels, err := startDNSServers(stageHarness)
	for _, cancel := range cancels {
		defer cancel()
	}
	if err != nil {
		return err
	}

	logger := stageHarness.Logger

	request := new(dns.Msg)
	request.SetQuestion(dns.Fqdn(DEFAULT_DOMAIN), dns.TypeA)

	for rd := 0; rd < 2; rd++ {
		request.MsgHdr.RecursionDesired = rd == 0
		for opcode := 0; opcode < 4; opcode++ {
			request.MsgHdr.Opcode = opcode

			response, err := sendQuery(logger, request)
			if err != nil {
				return err
			}
			if response.MsgHdr.RecursionDesired != request.MsgHdr.RecursionDesired {
				return fmt.Errorf("Expected RecursionDesired field to be %v got %v",
					request.MsgHdr.RecursionDesired, response.MsgHdr.RecursionDesired)
			}
			if response.MsgHdr.Opcode != request.MsgHdr.Opcode {
				return fmt.Errorf("Expected Opcode field to be %v got %v", request.MsgHdr.Opcode, response.MsgHdr.Opcode)
			}
			if request.MsgHdr.Opcode != 0 {
				if response.MsgHdr.Rcode != 4 {
					return fmt.Errorf("Expected Rcode field to be 4 got %v", response.MsgHdr.Rcode)
				}
			}
		}
	}

	return nil
}

func sendQuery(logger *logger.Logger, request *dns.Msg) (*dns.Msg, error) {
	c := new(dns.Client)
	question := request.Question[0]

	logger.Infof("Querying: %s", question.String())
	logger.Debugf("Sending Request: (Messages with >>> prefix are part of this log)")
	logDNSPacket(logger, request)

	response, _, err := c.Exchange(request, SERVER_ADDR)
	if err != nil {
		return nil, fmt.Errorf("DNS query failed: %s.\nIf you are seeing this after a while then it is likely that your server is not responding with appropriate id", err)
	}
	if !response.MsgHdr.Response {
		return nil, fmt.Errorf("Expected QR field to be set to 1. 1 indicates that it is a response. Got 0")
	}
	if response.MsgHdr.Authoritative {
		return nil, fmt.Errorf("Expected AA field to not be set. Got 1. 1 indicates that the response is authoritative which is not true for this server.")
	}

	logger.Debugf("Received Response: (Messages with >>> prefix are part of this log)")
	logDNSPacket(logger, response)

	return response, nil
}

func sendMultiQuery(logger *logger.Logger, request *dns.Msg) (*dns.Msg, error) {
	c := new(dns.Client)
	logger.Infof("Querying the following in the same request (Messages with >> prefix are part of this log)")
	for _, question := range request.Question {
		logger.Infof(">> %s", question.String())
	}

	logger.Debugf("Sending Request: (Messages with >>> prefix are part of this log)")
	logDNSPacket(logger, request)

	response, _, err := c.Exchange(request, SERVER_ADDR)
	if err != nil {
		return nil, fmt.Errorf("DNS query failed: %s.\nIf you are seeing this after a while then it is likely that your server is not responding with appropriate id", err)
	}
	if !response.MsgHdr.Response {
		return nil, fmt.Errorf("Expected QR field to be set to 1. 1 indicates that it is a response. Got 0")
	}
	if response.MsgHdr.Authoritative {
		return nil, fmt.Errorf("Expected AA field to not be set. Got 1. 1 indicates that the response is authoritative which is not true for this server.")
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
