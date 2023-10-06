package internal

import (
	"fmt"
	"net"
	"time"

	tester_utils "github.com/codecrafters-io/tester-utils"
	"github.com/miekg/dns"
)

const (
	SERVER_ADDR = "127.0.0.1:2053"
)

// Example from the grep course
func testInit(stageHarness *tester_utils.StageHarness) error {
	// b := NewDnsServerBinary(stageHarness)
	// if err := b.Run(); err != nil {
	// 	return err
	// }

	conn, err := net.DialTimeout("udp", SERVER_ADDR, 2*time.Second)
	if err != nil {
		return fmt.Errorf("Error connecting to UDP server: %v\n", err)
	}
	defer conn.Close()

	msg, err := getDnsMsgBytes()
	if err != nil {
		return err
	}
	_, err = conn.Write(msg)
	if err != nil {
		return fmt.Errorf("Error sending UDP packet: %v\n", err)
	}

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)

	if err != nil {
		netErr, ok := err.(net.Error)
		if ok && netErr.Timeout() {
			fmt.Println("No ICMP response, port is likely open.")
			return nil
		}
		return fmt.Errorf("Error receiving UDP response: %v\n", err)
	}

	return nil
}

func getDnsMsgBytes() ([]byte, error) {
	msg := new(dns.Msg)
	msg.Id = dns.Id()
	msg.SetQuestion("codecrafters.io.", dns.TypeA)

	buf, err := msg.Pack()
	if err != nil {
		return nil, fmt.Errorf("Error encoding DNS message: %s\n", err)
	}

	return buf, nil
}
