package internal

import (
	"fmt"

	tester_utils "github.com/codecrafters-io/tester-utils"
)

func testReceiveQuestionInResponse(stageHarness *tester_utils.StageHarness) error {
	cancels, err := startDNSServers(stageHarness)
	for _, cancel := range cancels {
		defer cancel()
	}
	if err != nil {
		return err
	}

	logger := stageHarness.Logger

	response, err := sendDNSQueryWithId(logger, uint16(DEFAULT_PKT_ID), DEFAULT_DOMAIN)
	if err != nil {
		return fmt.Errorf("Error sending DNS query: %s\n", err)
	}

	if response.Id != DEFAULT_PKT_ID {
		return fmt.Errorf("Expected ID to be 1234, got %d", response.Id)
	}
	if len(response.Question) == 0 {
		return fmt.Errorf("Expected question section to have one entry got %d", len(response.Question))
	}
	if response.Question[0].Name != DEFAULT_DOMAIN {
		return fmt.Errorf("Expected question domain name to be `%v` got `%v`", DEFAULT_DOMAIN, response.Question[0].Name)
	}

	return nil
}
