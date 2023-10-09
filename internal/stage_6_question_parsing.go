package internal

import (
	"fmt"

	tester_utils "github.com/codecrafters-io/tester-utils"
	"github.com/miekg/dns"
)

func testBasicQuestionParsing(stageHarness *tester_utils.StageHarness) error {
	b := NewDnsServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}
	logger := stageHarness.Logger
	if err := retryDialUntilSuccess(logger); err != nil {
		return err
	}

	queryDomain := "mail.google.com."
	response, err := sendQuery(logger, queryDomain, dns.TypeA)
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
