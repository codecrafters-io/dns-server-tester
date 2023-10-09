package internal

import (
	"fmt"

	tester_utils "github.com/codecrafters-io/tester-utils"
	"github.com/miekg/dns"
)

func testReceiveAnswerInResponse(stageHarness *tester_utils.StageHarness) error {
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
		return fmt.Errorf("Error sending DNS query: %s\n", err)
	}

	if len(response.Question) != 1 {
		return fmt.Errorf("Expected question section to have one entry got %d", len(response.Question))
	}
	if response.Question[0].Name != DEFAULT_DOMAIN {
		return fmt.Errorf("Expected question domain name to be `%v` got `%v`", DEFAULT_DOMAIN, response.Question[0].Name)
	}
	if len(response.Answer) == 0 {
		return fmt.Errorf("Expected answer section to have at least one entry got %d", len(response.Answer))
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
