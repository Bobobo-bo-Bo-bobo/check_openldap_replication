package main

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
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

	port, err := getPortFromHost(u.Host)
	if err != nil {
		return "", ssl, err
	}

	switch u.Scheme {
	case "ldap":
		ssl = false
		if port == 0 {
			u.Host += ":389"
		}

	case "ldaps":
		ssl = true
		if port == 0 {
			u.Host += ":636"
		}
	default:
		return "", ssl, fmt.Errorf("Unsupported URI scheme")
	}

	return u.Host, ssl, nil
}

func getPortFromHost(h string) (uint, error) {
	// ldap.Dial and ldap.DialTLS require a port in the address
	// If missing, use default ports (389/tcp for ldap and 636/tcp for ldaps)
	spl := strings.Split(strings.TrimSpace(h), ":")

	// Neither port nor IPv6 address
	if len(spl) == 1 {
		return 0, nil
	}

	// Names are in format assemble.the.minions or assemble.the.minions:3389
	// IPv4 addresses are in format 10.10.10.10 or 10.10.10.10:6636
	// IPv6 addresses are in format [fe80::c2d6:d200:2751:7dd9] or [fe80::c2d6:d200:2751:7dd9]:3389
	p := spl[len(spl)-1]
	if p[len(p)-1] == ']' {
		return 0, nil
	}

	port, err := strconv.ParseInt(p, 10, 0)
	if err != nil {
		return 0, err
	}

	if port <= 0 || port > 65535 {
		return 0, fmt.Errorf("Invalid port")
	}

	return uint(port), nil
}
