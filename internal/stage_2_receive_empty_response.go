package internal

import (
	"fmt"

	tester_utils "github.com/codecrafters-io/tester-utils"
	"github.com/miekg/dns"
)

func testReceiveEmptyResponse(stageHarness *tester_utils.StageHarness) error {
	// b := NewDnsServerBinary(stageHarness)
	// if err := b.Run(); err != nil {
	// 	return err
	// }

	queryDomain := "codecrafters.io."
	packetIdentifier := 1234

	dnsMsg, err := sendDNSQuery(uint16(packetIdentifier), queryDomain, SERVER_ADDR)
	if err != nil {
		return fmt.Errorf("Error sending DNS query: %s\n", err)
	}
	fmt.Println(dnsMsg)

	if dnsMsg.Id != uint16(packetIdentifier) {
		return fmt.Errorf("Expected ID to be 1234, got %d", dnsMsg.Id)
	}
	// // question section is empty
	// if len(dnsMsg.Question) != 0 {
	// 	return fmt.Errorf("Expected question section to be empty, got %d", len(dnsMsg.Question))
	// }
	// // answer section is empty
	// if len(dnsMsg.Answer) != 0 {
	// 	return fmt.Errorf("Expected answer section to be empty, got %d", len(dnsMsg.Answer))
	// }

	return nil
}

func sendDNSQuery(id uint16, queryDomain, dnsServerAddr string) (*dns.Msg, error) {
	c := new(dns.Client)

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(queryDomain), dns.TypeA)
	msg.Id = id

	response, _, err := c.Exchange(msg, dnsServerAddr)
	if err != nil {
		return nil, fmt.Errorf("DNS query failed: %s", err)
	}

	return response, nil
}
