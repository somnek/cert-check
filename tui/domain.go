package tui

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	dialTimeout = 5 * time.Second
)

// GetConfigPath returns the path to the config file
func GetConfigPath(folder, file string) string {
	return filepath.Join(os.Getenv("HOME"), folder, file)
}

// FileExists checks if a file exists at the given path
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DeleteDomain remove domain from config file
func DeleteFromConfig(domain, path string) error {
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
		return nil
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

// SaveDomain saves a new domain to the config file.
func SaveDomain(domain, path string) error {
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

func GetSavedDomains(path string) ([]string, error) {
	domains := []string{}

	f, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var c userConfig
	if err := yaml.Unmarshal(f, &c); err != nil {
		return nil, fmt.Errorf("error unmarshalling file: %v", err)
	}

	domains = append(domains, c.Domains...)
	return domains, err
}

// dial domains and return ssl info
func DialDomains(ssls *[]ssl, domains []string) error {
	var wg sync.WaitGroup
	ch := make(chan chDailRes, len(domains))

	for _, domain := range domains {

		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			fmt.Printf("...Dailing %s\n", d)
			GetInfo(d, ch)
		}(domain)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for res := range ch {
		if res.err != nil {
			return res.err
		}
		*ssls = append(*ssls, res.ssl)
	}
	return nil
}

// GetInfo returns the ssl info for a given domain
func GetInfo(domain string, ch chan chDailRes) {
	// create a custom Dialer with a timeout value
	Dialer := &net.Dialer{
		Timeout: dialTimeout,
	}

	conn, err := tls.DialWithDialer(Dialer, "tcp", domain+":443", nil)
	if err != nil {
		// check whether it is a timeout error
		if e, ok := err.(net.Error); ok && e.Timeout() {
			ch <- chDailRes{ssl{}, fmt.Errorf("timeout error: %v", err)}
		}
		ch <- chDailRes{ssl{}, fmt.Errorf("error dialing: %v", err)}
	}
	defer conn.Close()

	state := conn.ConnectionState()
	leafCert := state.PeerCertificates[0]
	issuedOn := leafCert.NotBefore.String()
	expiresOn := leafCert.NotAfter.String()
	issuer := leafCert.Issuer.CommonName
	commonName := leafCert.Subject.CommonName
	daysLeft := CalculateDaysLeft(leafCert.NotAfter)

	ch <- chDailRes{
		ssl{
			domain:     domain,
			issuedOn:   issuedOn,
			expiresOn:  expiresOn,
			issuer:     issuer,
			commonName: commonName,
			daysLeft:   daysLeft,
		},
		nil,
	}
}

// CreateConfig creates a config file with dummy content
func CreateConfig(configFolder, configFile string) error {
	configFolderPath := GetConfigPath(configFolder, "")
	if err := os.MkdirAll(configFolderPath, 0755); err != nil {
		return fmt.Errorf("error creating folder: %v", err)
	}

	configFilePath := GetConfigPath(configFolder, configFile)
	f, err := os.Create(configFilePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)

	}
	defer f.Close()

	dummyContent := `domains:
- google.com
- x.com`

	if _, err = f.WriteString(dummyContent); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}
