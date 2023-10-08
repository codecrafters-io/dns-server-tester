package internal

import (
	"fmt"

	tester_utils "github.com/codecrafters-io/tester-utils"
	"github.com/miekg/dns"
)

func testMoreRecords(stageHarness *tester_utils.StageHarness) error {
	// b := NewDnsServerBinary(stageHarness)
	// if err := b.Run(); err != nil {
	// 	return err
	// }

	// Generate
	queryDomain := "codecrafters.io."

	if err := testAAAARecrod(stageHarness, "google.com."); err != nil {
		return err
	}
	if err := testMXRecord(stageHarness, queryDomain); err != nil {
		return err
	}
	if err := testNSRecord(stageHarness, queryDomain); err != nil {
		return err
	}

	return nil
}

func testAAAARecrod(stageHarness *tester_utils.StageHarness, queryDomain string) error {
	c := new(dns.Client)
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(queryDomain), dns.TypeAAAA)
	msg.RecursionDesired = true

	dnsMsg, _, err := c.Exchange(msg, SERVER_ADDR)
	if err != nil {
		return fmt.Errorf("DNS query failed: %s.\nIf you are seeing this after a while then it is likely that your server is not responding with appropriate id", err)
	}

	fmt.Println(dnsMsg)

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

func testMXRecord(stageHarness *tester_utils.StageHarness, queryDomain string) error {
	c := new(dns.Client)
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(queryDomain), dns.TypeMX)
	msg.RecursionDesired = true

	dnsMsg, _, err := c.Exchange(msg, SERVER_ADDR)
	if err != nil {
		return fmt.Errorf("DNS query failed: %s.\nIf you are seeing this after a while then it is likely that your server is not responding with appropriate id", err)
	}

	fmt.Println(dnsMsg)

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

func testNSRecord(stageHarness *tester_utils.StageHarness, queryDomain string) error {
	c := new(dns.Client)
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(queryDomain), dns.TypeNS)
	msg.RecursionDesired = true

	dnsMsg, _, err := c.Exchange(msg, SERVER_ADDR)
	if err != nil {
		return fmt.Errorf("DNS query failed: %s.\nIf you are seeing this after a while then it is likely that your server is not responding with appropriate id", err)
	}

	fmt.Println(dnsMsg)

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
