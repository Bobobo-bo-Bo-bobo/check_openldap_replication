package main

import (
	"fmt"
	"net/url"
)

func buildConfiguration(m string, s string, b string, i bool, c string) (configuration, error) {
	var err error
	var result = configuration{
		MasterURI:   m,
		SlaveURI:    s,
		InsecureSSL: i,
		Base:        b,
		CACert:      c,
	}

	result.masterAddr, result.masterSSL, err = parseURI(m)
	if err != nil {
		return result, err
	}

	result.slaveAddr, result.slaveSSL, err = parseURI(s)
	if err != nil {
		return result, err
	}

	return result, nil
}

func parseURI(s string) (string, bool, error) {
	var ssl bool

	u, err := url.Parse(s)
	if err != nil {
		return "", ssl, err
	}
	if u.Host == "" {
		return "", ssl, fmt.Errorf("Host of URI is empty")
	}

	switch u.Scheme {
	case "ldap":
		ssl = false
	case "ldaps":
		ssl = true
	default:
		return "", ssl, fmt.Errorf("Unsupported URI scheme")
	}

	return u.Host, ssl, nil
}
