package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Domains []string `yaml:"domains"`
}

func check(domain string) error {
	conn, err := tls.Dial("tcp", domain+":443", nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	state := conn.ConnectionState()
	leafCert := state.PeerCertificates[0]
	issuedOn := leafCert.NotBefore
	expiresOn := leafCert.NotAfter
	org := leafCert.Issuer.Organization[0]
	commonName := leafCert.Subject.CommonName

	fmt.Println("---")
	fmt.Println("Domain:", domain)
	fmt.Println("Common Name:", commonName)
	fmt.Println("Organization:", org)
	fmt.Println("Issued On:", issuedOn)
	fmt.Println("Expires On:", expiresOn)

	return nil
}

func main() {
	file, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatal("Failed to read config file", err)
	}

	var config Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Fatal("Failed to parse config file", err)
	}

	for _, domain := range config.Domains {
		if err := check(domain); err != nil {
			fmt.Println(err)
		}
	}
}
