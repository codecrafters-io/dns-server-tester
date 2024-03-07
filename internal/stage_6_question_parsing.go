package internal

import (
	"fmt"

	"github.com/codecrafters-io/tester-utils/test_case_harness"
	"github.com/miekg/dns"
)

func testBasicQuestionParsing(stageHarness *test_case_harness.TestCaseHarness) error {
	cancels, err := startDNSServers(stageHarness)
	for _, cancel := range cancels {
		defer cancel()
	}
	if err != nil {
		return err
	}

	logger := stageHarness.Logger

	request := new(dns.Msg)
	queryDomain := randomDomainWithType(dns.TypeA)
	request.SetQuestion(dns.Fqdn(queryDomain), dns.TypeA)
	response, err := sendQuery(logger, request)
	if err != nil {
		return err
	}

	if len(response.Answer) == 0 {
		return fmt.Errorf("Expected answer section to have at least one entry got %d", len(response.Answer))
	}

	for _, record := range response.Answer {

		if record.Header().Name != queryDomain {
			return fmt.Errorf("Expected answer domain name to be `%v` got `%v`", queryDomain, record.Header().Name)
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
