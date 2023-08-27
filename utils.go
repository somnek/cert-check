package main

import (
	"crypto/tls"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func getConfigPath(folder, file string) string {
	return filepath.Join(os.Getenv("HOME"), folder, file)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func saveDomain(domain, path string) error {
	f, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	// unmarshal
	var c userConfig
	if err := yaml.Unmarshal(f, &c); err != nil {
		log.Fatalf("error unmarshalling file: %v", err)
	}
	c.Domains = append(c.Domains, domain)

	// marshal
	newF, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	if err = os.WriteFile(path, newF, 0644); err != nil {
		return err
	}
	return nil
}

// getSavedDomains reads a YAML file at the given path and unmarshals it into a config struct.
// It then appends each domain in the config to the slice of ssl structs.
func getSavedDomains(ssls *[]ssl, path string) error {
	f, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	var c userConfig
	if err := yaml.Unmarshal(f, &c); err != nil {
		log.Fatalf("error unmarshalling file: %v", err)
	}

	for _, domain := range c.Domains {
		var info ssl

		info, err = getInfo(domain)
		if err != nil {
			return err
		}
		*ssls = append(*ssls, info)
	}
	return nil
}

func createConfig(configFolder, configFile string) {
	configFolderPath := getConfigPath(configFolder, "")
	if err := os.MkdirAll(configFolderPath, 0755); err != nil {
		log.Fatalf("error creating folder: %v", err)
	}

	configFilePath := getConfigPath(configFolder, configFile)
	f, err := os.Create(configFilePath)
	if err != nil {
		log.Fatalf("error creating file: %v", err)
	}
	defer f.Close()

	dummyContent := `domains:
- google.com
- github.com
- x.com`

	if _, err = f.WriteString(dummyContent); err != nil {
		log.Fatalf("error writing to file: %v", err)
	}
}

func getInfo(domain string) (ssl, error) {
	conn, err := tls.Dial("tcp", domain+":443", nil)
	if err != nil {
		log.Fatalf("i cannot => %v", domain)
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
