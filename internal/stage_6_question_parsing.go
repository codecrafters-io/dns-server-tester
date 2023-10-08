package internal

import (
	"fmt"
	"math/rand"

	tester_utils "github.com/codecrafters-io/tester-utils"
	"github.com/miekg/dns"
)

func testBasicQuestionParsing(stageHarness *tester_utils.StageHarness) error {
	// b := NewDnsServerBinary(stageHarness)
	// if err := b.Run(); err != nil {
	// 	return err
	// }

	// Generate
	queryDomain := "mail.google.com."
	packetIdentifier := uint16(rand.Uint32())

	c := new(dns.Client)
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(queryDomain), dns.TypeA)
	msg.Id = packetIdentifier
	msg.RecursionDesired = true

	dnsMsg, _, err := c.Exchange(msg, SERVER_ADDR)
	if err != nil {
		return fmt.Errorf("DNS query failed: %s.\nIf you are seeing this after a while then it is likely that your server is not responding with appropriate id", err)
	}

	fmt.Println(dnsMsg)

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
