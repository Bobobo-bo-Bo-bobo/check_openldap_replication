package main

import (
	"fmt"
)

func buildNagiosPerfData(d float64, w uint, c uint) string {
	// Format of the performance data is defined in https://nagios-plugins.org/doc/guidelines.html#AEN200
	return fmt.Sprintf("delta=%.3fs;%d;%d", d, w, c)
}
