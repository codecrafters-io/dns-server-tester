package internal

import (
	"fmt"

	tester_utils "github.com/codecrafters-io/tester-utils"
	logger "github.com/codecrafters-io/tester-utils/logger"
	"github.com/miekg/dns"
)

func testMoreRecords(stageHarness *tester_utils.StageHarness) error {
	b := NewDnsServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}
	logger := stageHarness.Logger
	if err := retryDialUntilSuccess(logger); err != nil {
		return err
	}

	queryDomain := "codecrafters.io."

	if err := testAAAARecrod(logger, "google.com."); err != nil {
		return err
	}
	if err := testMXRecord(logger, queryDomain); err != nil {
		return err
	}
	if err := testNSRecord(logger, queryDomain); err != nil {
		return err
	}

	return nil
}

func testAAAARecrod(logger *logger.Logger, queryDomain string) error {
	dnsMsg, err := sendQuery(logger, queryDomain, dns.TypeAAAA)
	if err != nil {
		return err
	}

	if len(dnsMsg.Answer) == 0 {
		return fmt.Errorf("Expected some answer record to be present. Got none")
	}

	for _, record := range dnsMsg.Answer {

		if record.Header().Name != queryDomain {
			return fmt.Errorf("Expected answer domain name to be `%v` got `%v`", queryDomain, record.Header().Name)
		}
		if record.Header().Rrtype != dns.TypeAAAA {
			return fmt.Errorf("Expected answer type to be 1 got %d", record.Header().Rrtype)
		}

		if _, ok := record.(*dns.AAAA); ok {
			continue
		} else {
			return fmt.Errorf("Expected answer record to be of type A (IPv4) got %T", record)
		}

	}
	return nil
}

func testMXRecord(logger *logger.Logger, queryDomain string) error {
	dnsMsg, err := sendQuery(logger, queryDomain, dns.TypeMX)
	if err != nil {
		return err
	}

	if len(dnsMsg.Answer) == 0 {
		return fmt.Errorf("Expected some answer record to be present. Got none")
	}

	for _, record := range dnsMsg.Answer {

		if record.Header().Name != queryDomain {
			return fmt.Errorf("Expected answer domain name to be `%v` got `%v`", queryDomain, record.Header().Name)
		}
		if record.Header().Rrtype != dns.TypeMX {
			return fmt.Errorf("Expected answer type to be 1 got %d", record.Header().Rrtype)
		}

		if _, ok := record.(*dns.MX); ok {
			continue
		} else {
			return fmt.Errorf("Expected answer record to be of type A (IPv4) got %T", record)
		}

	}
	return nil
}

func testNSRecord(logger *logger.Logger, queryDomain string) error {
	dnsMsg, err := sendQuery(logger, queryDomain, dns.TypeNS)
	if err != nil {
		return err
	}

	if len(dnsMsg.Answer) == 0 {
		return fmt.Errorf("Expected some answer record to be present. Got none")
	}

	for _, record := range dnsMsg.Answer {

		if record.Header().Name != queryDomain {
			return fmt.Errorf("Expected answer domain name to be `%v` got `%v`", queryDomain, record.Header().Name)
		}
		if record.Header().Rrtype != dns.TypeNS {
			return fmt.Errorf("Expected answer type to be 1 got %d", record.Header().Rrtype)
		}

		if _, ok := record.(*dns.NS); ok {
			continue
		} else {
			return fmt.Errorf("Expected answer record to be of type A (IPv4) got %T", record)
		}

	}
	return nil
}
