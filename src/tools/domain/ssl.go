package domain

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"
)

func GetCertDetails(host string) ([]*x509.Certificate, error) {
	conn, err := tls.Dial("tcp", host+":443", &tls.Config{ServerName: host})
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}
	defer conn.Close()
	return conn.ConnectionState().PeerCertificates, nil
}

func printSSLDetails(domain string) {
	infos, err := GetCertDetails(domain)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i, c := range infos {
		fmt.Printf("=== Certificate #%d ===\n", i+1)
		fmt.Printf("Subject:      %s\n", c.Subject)
		fmt.Printf("Issuer:       %s\n", c.Issuer)
		fmt.Printf("DNS Names:    %v\n", c.DNSNames)
		fmt.Printf("Valid From:   %s\n", c.NotBefore.Format("2006-01-02"))
		fmt.Printf("Valid Until:  %s\n", c.NotAfter.Format("2006-01-02"))
		fmt.Printf("Days Left:    %d\n", int(time.Until(c.NotAfter).Hours()/24))
		fmt.Printf("Is CA:        %v\n", c.IsCA)
		fmt.Printf("Serial No:    %s\n", c.SerialNumber)
	}
}
