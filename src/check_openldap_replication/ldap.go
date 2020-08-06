package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"

	ldap "gopkg.in/ldap.v3"
)

func connect(u string, i bool, ca string, t time.Duration) (*ldap.Conn, error) {
	var c *ldap.Conn
	var err error

	dc := &net.Dialer{Timeout: t}
	tlscfg := &tls.Config{}

	if i {
		tlscfg.InsecureSkipVerify = true
	}

	if ca != "" {
		tlscfg.RootCAs = x509.NewCertPool()
		cadata, err := ioutil.ReadFile(ca)
		if err != nil {
			return c, err
		}
		tlsok := tlscfg.RootCAs.AppendCertsFromPEM(cadata)
		if !tlsok {
			return c, fmt.Errorf("Internal error while adding CA data to CA pool")
		}
	}

	c, err = ldap.DialURL(u, ldap.DialWithDialer(dc), ldap.DialWithTLSConfig(tlscfg))
	return c, err
}

func getContextCSN(c *ldap.Conn, b string, t int) (string, error) {
	var search *ldap.SearchRequest

	// assuming access to operational attributes of the base entry can be done anonymously
	// Note: We limit the number of results to one, becasue entryCSN can be present only once or none at all
	search = ldap.NewSearchRequest(b, ldap.ScopeBaseObject, ldap.NeverDerefAliases, 1, t, false, "(objectClass=*)", ldapSearchEntryCSN, nil)

	result, err := c.Search(search)
	if err != nil {
		return "", err
	}

	if len(result.Entries) == 0 {
		return "", fmt.Errorf("No contextCSN attribute found")
	}

	// XXX: There should be always one and only one result (if any)
	return result.Entries[0].GetAttributeValue("contextCSN"), nil
}

func parseCSN(csn string) (time.Time, error) {
	var result time.Time
	var err error

	/*
		    The CSN (Change Sequence Number) is defined in https://www.openldap.org/faq/data/cache/1145.html as:

			" CSN Representation
			Right now, a CSN is represented as:

			    GT '#' COUNT '#' SID '#' MOD

			    GT: Generalized Time with microseconds resolution,
			        without timezone/daylight saving:

			        YYYYmmddHHMMSS.uuuuuuZ

			        YYYY:   4-digit year (0001-9999)
			        mm:     2-digit month (01-12)
			        dd:     2-digit day (01-31)
			        HH:     2-digit hours (00-23)
			        MM:     2-digit minutes (00-59)
			        SS:     2-digit seconds (00-59; 00-60 for leap?)
			        .:      literal dot ('.')
			        uuuuuu: 6-digit microseconds (000000-999999)
			        Z:      literal capital zee ('Z')

			    COUNT: 6-hex change counter (000000-ffffff); used to distinguish multiple
			         changes occurring within the same time quantum.

			    SID: 3-hex Server ID (000-fff)

			    MOD: 6-hex (000000-ffffff); used for ordering the modifications within
			         an LDAP Modify operation (right now, in OpenLDAP it's always 000000)

			NOTE: in previous implementations, the Generalized Time string did not have the microseconds portion. This is tolerated by OpenLDAP 2.4.
			NOTE: in OpenLDAP 2.2-2.3, the SID was 2 digits only. This is tolerated by OpenLDAP 2.4.
			NOTE: in <draft-ietf-ldup-model>'s suggested format, the SID field was not required to be a number (actually, it was required to be the distinguished value of the naming attribute of the Replica Subentry representing the replica). The COUNT and MOD fields were both 4-digit hexadecimal numbers, and were prefixed by "0x".
			NOTE: the CSN format used by OpenLDAP 2.1 is not supported by any later release. It was based on <draft-ietf-ldup-model-03>, but used a decimal number (of unspecified length) for the SID field. "
	*/

	splitted := strings.Split(csn, "#")
	if len(splitted) != 4 {
		return result, fmt.Errorf("Invalid CSN")
	}

	// We don't care about COUNT, SID and MOD
	result, err = time.Parse(csnTimeFormat, splitted[0])
	if err != nil {
		return result, err
	}

	return result, nil
}
