package lib

import "net"

func IsNotFoundDnsErr(err error) bool {
	if err == nil {
		return false
	}
	dnsErr, ok := err.(*net.DNSError)
	return ok && dnsErr.IsNotFound
}
