package internal

import (
	"math/rand"

	"github.com/miekg/dns"
)

var dnsRecords = map[string]map[uint16]string{
	"example.com.": {
		dns.TypeA:     "127.0.0.1",
		dns.TypeAAAA:  "2001:db8::1",
		dns.TypeCNAME: "example-cname.com.",
		dns.TypeMX:    "10 mail.example.com.",
	},
	"google.com.": {
		dns.TypeA:    "172.217.3.110",
		dns.TypeAAAA: "2607:f8b0:4006:801::200e",
		dns.TypeMX:   "5 alt1.aspmx.l.google.com.",
	},
	"facebook.com.": {
		dns.TypeA:     "31.13.65.36",
		dns.TypeAAAA:  "2a03:2880:f12f:83:face:b00c::25de",
		dns.TypeCNAME: "star-mini.c10r.facebook.com.",
	},
	"github.com.": {
		dns.TypeA:    "140.82.113.3",
		dns.TypeAAAA: "2606:50c0:8000::64",
	},
	"stackoverflow.com.": {
		dns.TypeA:     "151.101.129.69",
		dns.TypeAAAA:  "151.101.1.69",
		dns.TypeCNAME: "stackoverflow.design",
		dns.TypeMX:    "10 aspmx.l.google.com.",
	},
	"reddit.com.": {
		dns.TypeA:    "151.101.65.140",
		dns.TypeAAAA: "151.101.1.140",
		dns.TypeMX:   "10 mail.reddit.com.",
	},
	"youtube.com.": {
		dns.TypeA:    "172.217.3.110",
		dns.TypeAAAA: "2607:f8b0:4006:801::200e",
	},
	"example-cname.com.": {
		dns.TypeA: "192.168.0.1",
	},
	"mail.example.com.": {
		dns.TypeA: "192.168.0.2",
	},
	"mail.google.com.": {
		dns.TypeA: "173.194.221.27",
	},
	"alt1.aspmx.l.google.com.": {
		dns.TypeA: "173.194.221.27",
	},
	"star-mini.c10r.facebook.com.": {
		dns.TypeA: "31.13.65.36",
	},
	"stackoverflow.design.": {
		dns.TypeA: "151.101.1.69",
	},
	"aspmx.l.google.com.": {
		dns.TypeA: "108.177.119.26",
	},
	"mail.reddit.com.": {
		dns.TypeA: "151.101.65.140",
	},
	"codecrafters.io.": {
		dns.TypeA:  "76.76.21.21",
		dns.TypeMX: "10 mail.codecrafters.io.",
	},
	"abc.codecrafters.io.": {
		dns.TypeA:  "76.76.21.21",
		dns.TypeMX: "10 mail.codecrafters.io.",
	},
	"def.codecrafters.io.": {
		dns.TypeA:  "76.76.21.21",
		dns.TypeMX: "10 mail.codecrafters.io.",
	},
}

func domainsWithType(recordType uint16) []string {
	var result []string
	for domain, records := range dnsRecords {
		if _, exists := records[recordType]; exists {
			result = append(result, domain)
		}
	}
	return result
}

func randomDomainWithType(recordType uint16) string {
	domains := domainsWithType(recordType)
	randomInt := rand.Intn(len(domains))
	return domains[randomInt]
}
