package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// getConfigPath returns the path to the config file
func getConfigPath(folder, file string) string {
	return filepath.Join(os.Getenv("HOME"), folder, file)
}

// fileExists checks if a file exists at the given path
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// deleteDomain remove domain from config file
func deleteDomain(domain, path string) error {
	f, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// unmarshal
	var c userConfig
	if err := yaml.Unmarshal(f, &c); err != nil {
		return fmt.Errorf("error unmarshalling file: %v", err)
	}

	// remove domain from slice
	idx := Find(c.Domains, domain)
	if idx == -1 {
		return fmt.Errorf("domain not found")
	}
	c.Domains = Delete(c.Domains, idx)

	// marshal
	newF, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("error marshalling file: %v", err)
	}

	// write
	if err := os.WriteFile(path, newF, 0644); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

// saveDomain saves a new domain to the config file.
func saveDomain(domain, path string) error {
	f, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	// unmarshal
	var c userConfig
	if err := yaml.Unmarshal(f, &c); err != nil {
		return fmt.Errorf("error unmarshalling file: %v", err)
	}

	// prepend domain to slice
	c.Domains = append([]string{domain}, c.Domains...)

	// marshal
	newF, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("error marshalling file: %v", err)
	}
	if err = os.WriteFile(path, newF, 0644); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	return nil
}

// getSavedDomains reads the config file and returns a slice of ssl structs.
func getSavedDomains(ssls *[]ssl, path string) error {
	f, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	var c userConfig
	if err := yaml.Unmarshal(f, &c); err != nil {
		return fmt.Errorf("error unmarshalling file: %v", err)
	}

	for _, domain := range c.Domains {
		var info ssl

		info, err = GetInfo(domain)
		if err != nil {
			return err
		}
		*ssls = append(*ssls, info)
	}
	return nil
}

// createConfig creates a config file with dummy content
func createConfig(configFolder, configFile string) error {
	configFolderPath := getConfigPath(configFolder, "")
	if err := os.MkdirAll(configFolderPath, 0755); err != nil {
		return fmt.Errorf("error creating folder: %v", err)
	}

	configFilePath := getConfigPath(configFolder, configFile)
	f, err := os.Create(configFilePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)

	}
	defer f.Close()

	dummyContent := `domains:
- google.com
- github.com
- x.com`

	if _, err = f.WriteString(dummyContent); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

// GetInfo returns the ssl info for a given domain
func GetInfo(domain string) (ssl, error) {
	// create a custom Dialer with a timeout value
	Dialer := &net.Dialer{
		Timeout: 1 * time.Second,
	}

	conn, err := tls.DialWithDialer(Dialer, "tcp", domain+":443", nil)
	if err != nil {
		// check whether it is a timeout error
		if e, ok := err.(net.Error); ok && e.Timeout() {
			return ssl{}, fmt.Errorf("timeout error: %v", err)
		}
		return ssl{}, fmt.Errorf("error dialing: %v", err)
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
