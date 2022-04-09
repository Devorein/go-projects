package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func checkDomain(domain string) {
	var hasMx, hasDMARC, hasSPF bool

	var spfRecord, dmarcRecord string

	mxRecord, err := net.LookupMX(domain)

	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	if len(mxRecord) > 0 {
		hasMx = true
	}

	txtRecords, err := net.LookupTXT(domain)

	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	for _, txtRecord := range txtRecords {
		if strings.HasPrefix(txtRecord, "v=spf1") {
			hasSPF = true
			spfRecord = txtRecord
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}

	fmt.Printf("%v %v %v %v %v %v", domain, hasMx, hasSPF, spfRecord, hasDMARC, dmarcRecord)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("domain, hasMx, hasSPF, spfRecord, hasDMARC, dmarcRecord\n")

	for scanner.Scan() {
		checkDomain(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error: Could not read from input: %v\n", err)
	}
}
