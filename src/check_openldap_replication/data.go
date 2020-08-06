package main

import (
	"time"
)

type configuration struct {
	Base        string
	MasterURI   string
	SlaveURI    string
	InsecureSSL bool
	CACert      string
	Timeout     time.Duration
}
