package internal

import (
	"fmt"

	tester_utils "github.com/codecrafters-io/tester-utils"
	"github.com/miekg/dns"
)

func testReceiveAnswerInResponse(stageHarness *tester_utils.StageHarness) error {
	b := NewDnsServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}
	logger := stageHarness.Logger
	if err := retryDialUntilSuccess(logger); err != nil {
		return err
	}

	queryDomain := "codecrafters.io."
	packetIdentifier := 1234

	dnsMsg, err := sendDNSQueryWithId(logger, uint16(packetIdentifier), queryDomain)
	if err != nil {
		return fmt.Errorf("Error sending DNS query: %s\n", err)
	}

	// id is 1234
	if dnsMsg.Id != 1234 {
		return fmt.Errorf("Expected ID to be 1234, got %d", dnsMsg.Id)
	}
	if len(dnsMsg.Question) != 1 {
		return fmt.Errorf("Expected question section to have one entry got %d", len(dnsMsg.Question))
	}
	if dnsMsg.Question[0].Name != queryDomain {
		return fmt.Errorf("Expected question domain name to be `%v` got `%v`", queryDomain, dnsMsg.Question[0].Name)
	}
	if len(dnsMsg.Answer) == 0 {
		return fmt.Errorf("Expected answer section to have at least one entry got %d", len(dnsMsg.Answer))
	}
	record := dnsMsg.Answer[0]

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
