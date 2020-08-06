# Preface
In larger OpenLDAP setups the data is replicated to one or more OpenLDAP slave servers (e.g. per datacenter and or location).
To ensure data integrity and consistency the state of the replication should be checked at regular intervals.

This check allows for checking the state of replication of an OpenLDAP slave server and can be integrated
in a Nagios based monitoring solution like [Icinga2](https://icinga.com/products/)

# Build requirements
This check is implemented in Go so, obviously, a Go compiler is required.

Building this programm requires the [go-ldap.v3](https://github.com/go-ldap/ldap/) library.

# Command line parameters
Bbecause replication meta data is stored in the operational data of the BaseDN and a single server can replicate different BaseDNs from different servers, the BaseDN is mandatory. To compare the replication information the URI of the LDAP master and LDAP slave is required too.

| *Paraemter* | *Description* | *Default* | *Comment* |
|:------------|:--------------|:---------:|:----------|
| `--base=<base>` | LDAP search base | - | **manatory** |
| `--ca-cert=<file>` | CA certificate for validation of the server SSL certificate | - |
| `--critical=<sec>` | Report critical state if difference is <sec> or higher | 20 | - |
| `--help` | Show help text | - | - |
| `--insecure` | Skip SSL verification of server certificate | - | Do not use in a production environment |
| `--master=<uri>` | URI of LDAP master | - | **mandatory** |
| `--slave=<uri>` | URI of LDAP slave | - | **mandatory** |
| `--timeout=<sec>` | LDAP connection and search timeout in seconds | 15 | - |
| `--version` | Show version information | - | - |
| `--warning=<sec>` | Report warning state if difference is <sec> or higher | 20 | - |

# Licenses
## check_openldap_replication
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.

## go-ldap.v3 (https://github.com/go-ldap/ldap/)
The MIT License (MIT)

Copyright (c) 2011-2015 Michael Mitton (mmitton@gmail.com)
Portions copyright (c) 2015-2016 go-ldap Authors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

