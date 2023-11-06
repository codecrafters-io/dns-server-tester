package internal

import (
	"fmt"
	"net"

	"github.com/miekg/dns"
)

func friendlyErr(err error) error {
	if err == dns.ErrShortRead {
		return fmt.Errorf("%s\nLooks like your server is responding with partial DNS packet", err)
	}
	netErr, ok := err.(net.Error)
	if ok && netErr.Timeout() {
		return fmt.Errorf("%s\nThe most common cause of i/o timeout in this case is that the server might not be responding with appropriate id. Please make sure that the id is set to 1234.", err)
	}
	return err
}
