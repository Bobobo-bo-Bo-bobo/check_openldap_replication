package main

import (
	"time"
)

func buildConfiguration(m string, s string, b string, i bool, c string, t uint) configuration {
	return configuration{
		MasterURI:   m,
		SlaveURI:    s,
		InsecureSSL: i,
		Base:        b,
		CACert:      c,
		Timeout:     time.Duration(t) * time.Duration(time.Second),
	}
}
