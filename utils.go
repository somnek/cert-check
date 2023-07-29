package main

import (
	"crypto/tls"
)

func getInfo(domain string) (ssl, error) {
	conn, err := tls.Dial("tcp", domain+":443", nil)
	if err != nil {
		return ssl{}, err
	}
	defer conn.Close()

	state := conn.ConnectionState()
	leafCert := state.PeerCertificates[0]
	issuedOn := leafCert.NotBefore.String()
	expiresOn := leafCert.NotAfter.String()
	issuer := leafCert.Issuer.CommonName
	commonName := leafCert.Subject.CommonName

	return ssl{
		domain:     domain,
		issuedOn:   issuedOn,
		expiresOn:  expiresOn,
		issuer:     issuer,
		commonName: commonName,
	}, nil

}
