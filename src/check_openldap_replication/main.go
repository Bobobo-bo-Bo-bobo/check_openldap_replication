package main

import (
	"flag"
	"fmt"
	"math"
	"os"
)

func main() {
	var _showVersion = flag.Bool("version", false, "Show version")
	var _showHelp = flag.Bool("help", false, "Show help text")
	var warn = flag.Uint("warning", defaultWarnDifference, "Warn level")
	var crit = flag.Uint("critical", defaultCriticalDifference, "Critical level")
	var master = flag.String("master", "", "LDAP Master")
	var slave = flag.String("slave", "", "LDAP slave")
	var base = flag.String("base", "", "LDAP base")
	var insecure = flag.Bool("insecure", false, "Skip SSL verification")
	var cacert = flag.String("ca-cert", "", "CA certificate for SSL")

	flag.Usage = showUsage
	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Fprintln(os.Stderr, "Error: Too many arguments")
		os.Exit(UNKNOWN)
	}

	if *_showHelp {
		showUsage()
		os.Exit(OK)
	}

	if *_showVersion {
		showVersion()
		os.Exit(OK)
	}

	if *base == "" {
		fmt.Fprintln(os.Stderr, "Error: Missing mandatory parameter for LDAP search base")
		os.Exit(UNKNOWN)
	}
	if *master == "" {
		fmt.Fprintln(os.Stderr, "Error: Missing mandatory parameter for LDAP master URI")
		os.Exit(UNKNOWN)
	}
	if *slave == "" {
		fmt.Fprintln(os.Stderr, "Error: Missing mandatory parameter for LDAP slave URI")
		os.Exit(UNKNOWN)
	}

	// sanity checks
	if *warn == 0 {
		fmt.Fprintln(os.Stderr, "Error: Warn limit must be greater than zero")
		os.Exit(UNKNOWN)
	}
	if *crit == 0 {
		fmt.Fprintln(os.Stderr, "Error: Critical limit must be greater than zero")
		os.Exit(UNKNOWN)
	}
	if *warn > *crit {
		fmt.Fprintln(os.Stderr, "Error: Warn limit must be less or equal to critical limit")
		os.Exit(UNKNOWN)
	}

	cfg, err := buildConfiguration(*master, *slave, *base, *insecure, *cacert)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: "+err.Error())
		os.Exit(UNKNOWN)
	}

	// get contextCSN from master
	mcon, err := connect(cfg.masterAddr, cfg.masterSSL, cfg.InsecureSSL, cfg.CACert)
	if err != nil {
		fmt.Println("CRITICAL - " + err.Error())
		os.Exit(CRITICAL)
	}
	defer mcon.Close()

	mcsn, err := getContextCSN(mcon, cfg.Base)
	if err != nil {
		fmt.Println("CRITICAL - " + err.Error())
		os.Exit(CRITICAL)
	}

	mt, err := parseCSN(mcsn)
	if err != nil {
		fmt.Println("CRITICAL - " + err.Error())
		os.Exit(CRITICAL)
	}

	// get contextCSN from slave
	scon, err := connect(cfg.slaveAddr, cfg.slaveSSL, cfg.InsecureSSL, cfg.CACert)
	if err != nil {
		fmt.Println("CRITICAL - " + err.Error())
		os.Exit(CRITICAL)
	}
	defer scon.Close()

	scsn, err := getContextCSN(scon, cfg.Base)
	if err != nil {
		fmt.Println("CRITICAL - " + err.Error())
		os.Exit(CRITICAL)
	}

	st, err := parseCSN(scsn)
	if err != nil {
		fmt.Println("CRITICAL - " + err.Error())
		os.Exit(CRITICAL)
	}

	delta := mt.Sub(st).Seconds()
	perfdata := buildNagiosPerfData(delta, *warn, *crit)
	delta = math.Abs(delta)

	if delta >= float64(*crit) {
		fmt.Printf("CRITICAL - LDAP directories are out of sync by %.3f seconds | %s\n", delta, perfdata)
		os.Exit(CRITICAL)
	} else if delta >= float64(*warn) {
		fmt.Printf("WARNING - LDAP directories are out of sync by %.3f seconds | %s\n", delta, perfdata)
		os.Exit(WARNING)
	} else {
		fmt.Printf("OK - LDAP directories are in sync (time difference is %.3f seconds) | %s\n", delta, perfdata)
		os.Exit(OK)
	}
	os.Exit(UNKNOWN)
}
