package internal

import (
	"fmt"
	"net"

	tester_utils "github.com/codecrafters-io/tester-utils"
	logger "github.com/codecrafters-io/tester-utils/logger"
	"github.com/miekg/dns"
)

func testForwarding(stageHarness *tester_utils.StageHarness) error {
	cancels, err := startDNSServers(stageHarness)
	for _, cancel := range cancels {
		defer cancel()
	}
	if err != nil {
		return err
	}

	logger := stageHarness.Logger

	queryDomain := randomDomainWithType(dns.TypeA)

	if err := testARecord(logger, queryDomain); err != nil {
		return err
	}

	return nil
}

func testARecord(logger *logger.Logger, queryDomain string) error {
	request := new(dns.Msg)
	request.SetQuestion(dns.Fqdn(queryDomain), dns.TypeA)

	response, err := sendQuery(logger, request)
	if err != nil {
		return err
	}
	expectedIP := net.ParseIP(dnsRecords[queryDomain][dns.TypeA])

	for _, record := range response.Answer {
		if record.Header().Name != queryDomain {
			return fmt.Errorf("Expected answer domain name to be `%v` got `%v`", queryDomain, record.Header().Name)
		}
		if record.Header().Rrtype != dns.TypeA {
			return fmt.Errorf("Expected answer type to be 1 got %d", record.Header().Rrtype)
		}

		if aRecord, ok := record.(*dns.A); ok {
			if !aRecord.A.Equal(expectedIP) {
				return fmt.Errorf("Expected IPv4 address to be %v, got %v", expectedIP, aRecord.A)
			} else {
				return nil
			}
		} else {
			return fmt.Errorf("Expected answer record to be of type A (IPv4) got %T", record)
		}

	}

	return fmt.Errorf("Expected IPv4 address to be %v. No such match found.", expectedIP)
}
