package internal

import (
	"context"
	"fmt"

	logger "github.com/codecrafters-io/tester-utils/logger"
	"github.com/miekg/dns"
)

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)

	for _, question := range r.Question {
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

func handleCompressingDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.Compress = true
	m.SetReply(r)
	fmt.Println(m.String())

	for _, question := range r.Question {
		fmt.Printf("Received question: %s\n", question.Name)
		data := "172.0.0.1"
		rr, err := dns.NewRR(fmt.Sprintf("%s IN %s %s", question.Name, recordTypeToString(dns.TypeA), data))
		if err != nil {
			fmt.Println("Error creating DNS record:", err)
			continue
		}
		m.Answer = append(m.Answer, rr)
	}

	if err := w.WriteMsg(m); err != nil {
		fmt.Println("Error writing DNS response:", err)
	}
}

func startDNSServer(ctx context.Context, logger *logger.Logger, addr string) {
	server := &dns.Server{Addr: addr, Net: "udp"}
	for _, domain := range longAssDomainNames() {
		dns.HandleFunc(domain, handleCompressingDNSRequest)
	}
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
