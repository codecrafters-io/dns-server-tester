package internal

import (
	"fmt"

	"github.com/codecrafters-io/tester-utils/test_case_harness"
	"github.com/miekg/dns"
)

func testReceiveQuestionInResponse(stageHarness *test_case_harness.TestCaseHarness) error {
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

	if len(response.Question) != 1 {
		return friendlyQuestionErr(response)
	}

	actualQuestion := response.Question[0]

	if err := validateQuestion(DEFAULT_DOMAIN, &actualQuestion); err != nil {
		return err
	}

	return nil
}

func validateQuestion(domain string, question *dns.Question) error {
	if question.Name != DEFAULT_DOMAIN {
		return fmt.Errorf("Expected question domain name to be `%v` got `%v`", domain, question.Name)
	}

	if question.Qtype != dns.TypeA {
		more_info := ""
		if question.Qtype == 256 {
			more_info = "\nIt is likely you are missing a null terminator in your domain name."
		}
		return fmt.Errorf("Expected question type to be 1 got %d%s", question.Qtype, more_info)
	}

	if question.Qclass != dns.ClassINET {
		return fmt.Errorf("Expected question class to be 1 got %d", question.Qclass)
	}

	return nil
}
