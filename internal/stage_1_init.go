package internal

import (
	"fmt"
	"net"
	"time"

	tester_utils "github.com/codecrafters-io/tester-utils"
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

	_, err = conn.Write([]byte("ping"))
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
