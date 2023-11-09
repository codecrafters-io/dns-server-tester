package internal

import (
	"fmt"

	tester_utils "github.com/codecrafters-io/tester-utils"
	"github.com/miekg/dns"
)

func testCompressedPacketParsing(stageHarness *tester_utils.StageHarness) error {
	cancels, err := startDNSServers(stageHarness)
	for _, cancel := range cancels {
		defer cancel()
	}
	if err != nil {
		return err
	}

	logger := stageHarness.Logger

	queryDomain1, queryDomain2 := domainPairForCompression()
	queryDomains := []string{queryDomain1, queryDomain2}

	request := new(dns.Msg)
	request.SetQuestion(dns.Fqdn(queryDomain1), dns.TypeA)
	request.Question = append(request.Question, dns.Question{Name: queryDomain2, Qtype: dns.TypeA, Qclass: dns.ClassINET})
	request.Compress = true

	response, err := sendMultiQuery(logger, request)
	if err != nil {
		return err
	}

	if len(response.Answer) != 2 {
		return fmt.Errorf("Expected answer section to have 2 entries got %d", len(response.Answer))
	}

	for i, record := range response.Answer {
		validateRecord(queryDomains[i], record)
	}

	return nil
}

func validateRecord(queryDomain string, record dns.RR) error {
	if record.Header().Name != queryDomain {
		return fmt.Errorf("Expected answer domain name to be `%v` got `%v`", queryDomain, record.Header().Name)
	}
	if record.Header().Rrtype != dns.TypeA {
		return fmt.Errorf("Expected answer type to be 1 got %d", record.Header().Rrtype)
	}

	if _, ok := record.(*dns.A); ok {
		return nil
	} else {
		return fmt.Errorf("Expected answer record to be of type A (IPv4) got %T", record)
	}
}
