package internal

import (
	"fmt"

	"github.com/codecrafters-io/tester-utils/test_case_harness"
	"github.com/miekg/dns"
)

func testReceiveAnswerInResponse(stageHarness *test_case_harness.TestCaseHarness) error {
	cancels, err := startDNSServers(stageHarness)
	for _, cancel := range cancels {
		defer cancel()
	}
	if err != nil {
		return err
	}

	logger := stageHarness.Logger

	response, err := sendDNSQueryWithId(logger, uint16(DEFAULT_PKT_ID), DEFAULT_DOMAIN)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	if err := validateQuestion(DEFAULT_DOMAIN, &response.Question[0]); err != nil {
		return err
	}

	if len(response.Answer) == 0 {
		return friendlyAnswerErr(response)
	}
	record := response.Answer[0]

	if record.Header().Name != DEFAULT_DOMAIN {
		return fmt.Errorf("Expected answer domain name to be `%v` got `%v`", DEFAULT_DOMAIN, record.Header().Name)
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
