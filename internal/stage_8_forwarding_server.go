package internal

import (
	"fmt"
	"net"

	tester_utils "github.com/codecrafters-io/tester-utils"
	logger "github.com/codecrafters-io/tester-utils/logger"
	"github.com/miekg/dns"
)

func testForwarding(stageHarness *tester_utils.StageHarness) error {
	b := NewDnsServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}
	logger := stageHarness.Logger
	if err := retryDialUntilSuccess(logger); err != nil {
		return err
	}

	queryDomain := "codecrafters.io."

	if err := testARecord(stageHarness.Logger, queryDomain, net.IPv4(76, 76, 21, 21)); err != nil {
		return err
	}

	if err := testARecord(stageHarness.Logger, "google.com.", net.IPv4(142, 250, 183, 14)); err != nil {
		return err
	}

	return nil
}

func testARecord(logger *logger.Logger, queryDomain string, expectedIP net.IP) error {
	dnsMsg, err := sendQuery(logger, queryDomain, dns.TypeA)
	if err != nil {
		return err
	}

	for _, record := range dnsMsg.Answer {
		if record.Header().Name != queryDomain {
			return fmt.Errorf("Expected answer domain name to be `%v` got `%v`", queryDomain, record.Header().Name)
		}
		if record.Header().Rrtype != dns.TypeA {
			return fmt.Errorf("Expected answer type to be 1 got %d", record.Header().Rrtype)
		}

		if _, ok := record.(*dns.A); ok {
			// TODO: Actually test this once we have our home grown DNS server which the program under test will call
			return nil
		} else {
			return fmt.Errorf("Expected answer record to be of type A (IPv4) got %T", record)
		}

	}

	return fmt.Errorf("Expected IPv4 address to be %v. No such match found.", expectedIP)
}
