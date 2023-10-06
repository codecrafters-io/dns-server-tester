package internal

import (
	"fmt"
	"math/rand"

	tester_utils "github.com/codecrafters-io/tester-utils"
	"github.com/miekg/dns"
)

func testHeaderParsing(stageHarness *tester_utils.StageHarness) error {
	// b := NewDnsServerBinary(stageHarness)
	// if err := b.Run(); err != nil {
	// 	return err
	// }

	queryDomain := "codecrafters.io."
	packetIdentifier := uint16(rand.Uint32())

	c := new(dns.Client)
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(queryDomain), dns.TypeA)
	msg.Id = packetIdentifier
	msg.RecursionDesired = rand.Intn(2) == 0

	dnsMsg, _, err := c.Exchange(msg, SERVER_ADDR)
	if err != nil {
		return fmt.Errorf("DNS query failed: %s.\nIf you are seeing this after a while then it is likely that your server is not responding with appropriate id", err)
	}

	fmt.Println(dnsMsg)

	// this check is not required as dns gomodule does it
	if dnsMsg.Id != packetIdentifier {
		return fmt.Errorf("Expected ID to be 1234, got %d", dnsMsg.Id)
	}

	if dnsMsg.RecursionDesired != msg.RecursionDesired {
		return fmt.Errorf("Expected RecursionAvailable to be %t, got %t", msg.RecursionDesired, dnsMsg.RecursionAvailable)
	}

	return nil
}
