package internal

import (
	"fmt"

	tester_utils "github.com/codecrafters-io/tester-utils"
)

func testReceiveQuestionInResponse(stageHarness *tester_utils.StageHarness) error {
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

	return nil
}
