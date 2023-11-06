package internal

import (
	"fmt"
	"net"
	"strings"

	"github.com/miekg/dns"
)

func friendlyErr(err error) error {
	if err == dns.ErrShortRead {
		return fmt.Errorf("%s\nLooks like your server is responding with partial DNS packet", err)
	}
	if err == dns.ErrRdata {
		return fmt.Errorf("%s\nBad rdata error can happen because of a lot of reasons.\nIf you are seeing this in an early stage that probably indicates that the packet you sent was not properly parsed.", err)
	}
	if strings.Contains(err.Error(), "overflowing header size") {
		return fmt.Errorf("%s\nThis means that the length you have specified in the header of a Record is too long for what followed in the result.", err)
	}
	netErr, ok := err.(net.Error)
	if ok && netErr.Timeout() {
		return fmt.Errorf("%s\nThe most common cause of i/o timeout in this case is that the server might not be responding with appropriate id. Please make sure that the id is set to 1234.", err)
	}
	return err
}

func friendlyQuestionErr(response *dns.Msg) error {
	helpMsg := "If you have question set in question section and are still seeing this error, make sure that you are setting the QCOUNT to 1 in the header section too."
	return fmt.Errorf("Expected question section to have one entry got %d\n%s", len(response.Question), helpMsg)
}

func friendlyAnswerErr(response *dns.Msg) error {
	helpMsg := "If you have answer set in answer section and are still seeing this error, make sure that you are setting the ANCOUNT to 1 in the header section too."
	return fmt.Errorf("Expected answer section to have one entry got %d\n%s", len(response.Answer), helpMsg)
}
