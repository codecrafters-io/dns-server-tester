package internal

import (
	"fmt"
	"math/rand"

	tester_utils "github.com/codecrafters-io/tester-utils"
	"github.com/miekg/dns"
)

func testCompressedQuestionParsing(stageHarness *tester_utils.StageHarness) error {
	b := NewDnsServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}
	logger := stageHarness.Logger
	if err := retryDialUntilSuccess(logger); err != nil {
		return err
	}

	queryDomain1 := "abc.codecrafters.io."
	queryDomain2 := "def.codecrafters.io."
	packetIdentifier := uint16(rand.Uint32())

	c := new(dns.Client)
	retryDialUntilSuccess(logger)

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(queryDomain1), dns.TypeA)
	msg.Question = append(msg.Question, dns.Question{Name: dns.Fqdn(queryDomain2), Qtype: dns.TypeA, Qclass: dns.ClassINET})
	msg.Id = packetIdentifier
	msg.RecursionDesired = true
	msg.Compress = true
	logger.Infof("Querying `A` record for [%s, %s]", queryDomain1, queryDomain2)
	logger.Debugf("Sending Request: (Messages with >>> prefix are part of this log)")
	logDnsMsg(logger, msg)

	dnsMsg, _, err := c.Exchange(msg, SERVER_ADDR)
	if err != nil {
		return fmt.Errorf("DNS query failed: %s.\nIf you are seeing this after a while then it is likely that your server is not responding with appropriate id", err)
	}
	logger.Debugf("Received Response: (Messages with >>> prefix are part of this log)")
	logDnsMsg(logger, dnsMsg)

	for i, record := range dnsMsg.Answer {

		if record.Header().Name != msg.Question[i].Name {
			return fmt.Errorf("Expected answer domain name to be `%v` got `%v`", msg.Question[i].Name, record.Header().Name)
		}
		if record.Header().Rrtype != dns.TypeA {
			return fmt.Errorf("Expected answer type to be 1 got %d", record.Header().Rrtype)
		}

		if _, ok := record.(*dns.A); ok {
			continue
		} else {
			return fmt.Errorf("Expected answer record to be of type A (IPv4) got %T", record)
		}

	}
	return nil
}
