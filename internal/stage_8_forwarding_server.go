package internal

import (
	"fmt"
	"net"

	tester_utils "github.com/codecrafters-io/tester-utils"
	"github.com/miekg/dns"
)

func testForwarding(stageHarness *tester_utils.StageHarness) error {
	// b := NewDnsServerBinary(stageHarness)
	// if err := b.Run(); err != nil {
	// 	return err
	// }

	// Generate
	queryDomain := "codecrafters.io."

	if err := testARecord(stageHarness, queryDomain, net.IPv4(76, 76, 21, 21)); err != nil {
		return err
	}

	if err := testARecord(stageHarness, "google.com.", net.IPv4(142, 250, 183, 14)); err != nil {
		return err
	}

	return nil
}

func testARecord(stageHarness *tester_utils.StageHarness, queryDomain string, expectedIP net.IP) error {
	c := new(dns.Client)
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(queryDomain), dns.TypeA)
	msg.RecursionDesired = true

	dnsMsg, _, err := c.Exchange(msg, SERVER_ADDR)
	if err != nil {
		return fmt.Errorf("DNS query failed: %s.\nIf you are seeing this after a while then it is likely that your server is not responding with appropriate id", err)
	}

	fmt.Println(dnsMsg)

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
