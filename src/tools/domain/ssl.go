package domain

import (
	"crypto/tls"
	"fmt"
	"time"
)

func getSSLDetails(domain string) (string, error) {
	conn, err := tls.Dial("tcp", domain+":443", &tls.Config{ServerName: domain})
	if err != nil {
		return "", err
	}
	defer conn.Close()

	var result string
	certs := conn.ConnectionState().PeerCertificates
	for _, c := range certs {
		if c.IsCA {
			continue
		}
		result += fmt.Sprintln(c.NotBefore.Format(time.DateOnly), "->", c.NotAfter.Format(time.DateOnly), "/", int(time.Until(c.NotAfter).Hours()/24), "days left")
		result += fmt.Sprintln("Subject:", c.Subject)
		result += fmt.Sprintln("Issuer:", c.Issuer)
		result += fmt.Sprintln("DNS Names:", c.DNSNames)
	}
	return result, nil
}
