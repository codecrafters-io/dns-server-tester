package internal

import (
	"context"
	"fmt"

	logger "github.com/codecrafters-io/tester-utils/logger"
	"github.com/miekg/dns"
)

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	fmt.Printf("Received DNS request, %v", m)
	m.SetReply(r)

	for _, question := range r.Question {
		fmt.Println("Question:", question)
		records, exists := dnsRecords[question.Name]
		if !exists {
			continue
		}

		for recordType, data := range records {
			if recordType == question.Qtype {
				rr, err := dns.NewRR(fmt.Sprintf("%s IN %s %s", question.Name, recordTypeToString(recordType), data))
				if err != nil {
					fmt.Println("Error creating DNS record:", err)
					continue
				}
				m.Answer = append(m.Answer, rr)
			}
		}
	}

	if err := w.WriteMsg(m); err != nil {
		fmt.Println("Error writing DNS response:", err)
	}
}

func startDNSServer(ctx context.Context, logger *logger.Logger, addr string) {
	server := &dns.Server{Addr: addr, Net: "udp"}
	dns.HandleFunc(".", handleDNSRequest)

	logger.Debugf("DNS resolver listening on %s", addr)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Errorf("Error starting DNS server: %s", err)
		}
	}()

	select {
	case <-ctx.Done():
		logger.Debugf("Shutting down DNS resolver server...")
		server.Shutdown()
	}
}
