package internal

import (
	"fmt"

	tester_utils "github.com/codecrafters-io/tester-utils"
	"github.com/miekg/dns"
)

func testCompressedQuestionParsing(stageHarness *tester_utils.StageHarness) error {
	cancels, err := startDNSServers(stageHarness)
	for _, cancel := range cancels {
		defer cancel()
	}
	if err != nil {
		return err
	}

	logger := stageHarness.Logger

	queryDomain1 := "example.com."
	queryDomain2 := "google.com"

	c := new(dns.Client)

	request := new(dns.Msg)
	request.SetQuestion(dns.Fqdn(queryDomain1), dns.TypeA)
	request.Question = append(request.Question, dns.Question{Name: dns.Fqdn(queryDomain2), Qtype: dns.TypeA, Qclass: dns.ClassINET})
	request.RecursionDesired = true
	// request.Compress = true
	logger.Infof("Querying `A` record for [%s, %s]", queryDomain1, queryDomain2)
	logger.Debugf("Sending Request: (Messages with >>> prefix are part of this log)")
	logDNSPacket(logger, request)

	response, _, err := c.Exchange(request, SERVER_ADDR)
	if err != nil {
		return fmt.Errorf("DNS query failed: %s.\nIf you are seeing this after a while then it is likely that your server is not responding with appropriate id", err)
	}
	logger.Debugf("Received Response: (Messages with >>> prefix are part of this log)")
	logDNSPacket(logger, response)

	if len(response.Answer) < 2 {
		return fmt.Errorf("Expected answer section to have at least 2 entries got %d", len(response.Answer))
	}

	for i, record := range response.Answer {

		if record.Header().Name != request.Question[i].Name {
			return fmt.Errorf("Expected answer domain name to be `%v` got `%v`", request.Question[i].Name, record.Header().Name)
		}
		if record.Header().Rrtype != dns.TypeA {
			return fmt.Errorf("Expected answer type to be 1 got %d", record.Header().Rrtype)
		}

		if _, ok := record.(*dns.A); ok {
			continue
		} else {
			return fmt.Errorf("Expected answer record to be of type A (IPv4) got %T", record)
		}

	}
	return nil
}
