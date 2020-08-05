package main

const name = "check_openldap_replication"
const version = "1.0.0"

const versionText = `%s version %s
Copyright (C) 2020 by Andreas Maus <maus@ypbind.de>
This program comes with ABSOLUTELY NO WARRANTY.

%s is distributed under the Terms of the GNU General
Public License Version 3. (http://www.gnu.org/copyleft/gpl.html)

Build with go version: %s

`

const defaultWarnDifference = 10
const defaultCriticalDifference = 20

const helpText = `Usage: %s --base=<base> [--ca-cert=<file>] [--critical=<sec>] [--help] --insecure --master=<uri> --slave=<uri> [--version] [--warning=<sec>]

  --base=<base>     LDAP search base. This option is mandatory.

  --ca-cert=<file>  Use <file> as CA certificate for SSL verification.

  --critical=<sec>  Report critical state if difference is <sec> or higher.
                    Default: %d

  --help            This text.

  --insecure        Skip SSL verification of server certificate.

  --master=<uri>    URI of LDAP master. This option is mandatory.

  --slave=<uri>     URI of LDAP slave. This option is mandatory.

  --version         Show version information.

  --warning=<sec>   Report warning state if difference is <sec> or higher.
                    Default: %d

`

const (
	// OK - Nagios exit code
	OK int = iota
	// WARNING - Nagios exit code
	WARNING
	// CRITICAL - Nagios exit code
	CRITICAL
	// UNKNOWN - Nagios exit code
	UNKNOWN
)

const csnTimeFormat = "20060102150405.000000Z"

var ldapSearchEntryCSN = []string{"contextCSN"}
