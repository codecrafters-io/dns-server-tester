package internal

import (
	"fmt"
	"net"

	tester_utils "github.com/codecrafters-io/tester-utils"
	logger "github.com/codecrafters-io/tester-utils/logger"
	"github.com/miekg/dns"
)

func testMoreRecords(stageHarness *tester_utils.StageHarness) error {
	cancels, err := startDNSServers(stageHarness)
	for _, cancel := range cancels {
		defer cancel()
	}
	if err != nil {
		return err
	}

	logger := stageHarness.Logger

	queryDomain := randomDomainWithType(dns.TypeAAAA)
	if err := testAAAARecrod(logger, queryDomain); err != nil {
		return err
	}
	queryDomain = randomDomainWithType(dns.TypeMX)
	if err := testMXRecord(logger, queryDomain); err != nil {
		return err
	}
	queryDomain = randomDomainWithType(dns.TypeNS)
	if err := testNSRecord(logger, queryDomain); err != nil {
		return err
	}

	return nil
}

func testAAAARecrod(logger *logger.Logger, queryDomain string) error {
	request := new(dns.Msg)
	request.SetQuestion(dns.Fqdn(queryDomain), dns.TypeAAAA)
	response, err := sendQuery(logger, request)
	if err != nil {
		return err
	}
	expectedIP := net.ParseIP(dnsRecords[queryDomain][dns.TypeAAAA])

	if len(response.Answer) == 0 {
		return fmt.Errorf("Expected some answer record to be present. Got none")
	}

	for _, record := range response.Answer {

		if record.Header().Name != queryDomain {
			return fmt.Errorf("Expected answer domain name to be `%v` got `%v`", queryDomain, record.Header().Name)
		}
		if record.Header().Rrtype != dns.TypeAAAA {
			return fmt.Errorf("Expected answer type to be 1 got %d", record.Header().Rrtype)
		}

		if aaaaRecord, ok := record.(*dns.AAAA); ok {
			if !aaaaRecord.AAAA.Equal(expectedIP) {
				return fmt.Errorf("Expected IPv4 address to be %v, got %v", expectedIP, aaaaRecord.AAAA)
			} else {
				return nil
			}
		} else {
			return fmt.Errorf("Expected answer record to be of type A (IPv4) got %T", record)
		}

	}
	return nil
}

func testMXRecord(logger *logger.Logger, queryDomain string) error {
	request := new(dns.Msg)
	request.SetQuestion(dns.Fqdn(queryDomain), dns.TypeMX)
	response, err := sendQuery(logger, request)
	if err != nil {
		return err
	}

	if len(response.Answer) == 0 {
		return fmt.Errorf("Expected some answer record to be present. Got none")
	}

	for _, record := range response.Answer {

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
	request := new(dns.Msg)
	request.SetQuestion(dns.Fqdn(queryDomain), dns.TypeNS)
	response, err := sendQuery(logger, request)
	if err != nil {
		return err
	}

	if len(response.Answer) == 0 {
		return fmt.Errorf("Expected some answer record to be present. Got none")
	}

	for _, record := range response.Answer {

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
