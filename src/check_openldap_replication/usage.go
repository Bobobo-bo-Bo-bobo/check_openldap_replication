package main

import (
	"fmt"
)

func showUsage() {
	showVersion()
	fmt.Printf(helpText, name, defaultCriticalDifference, defaultWarnDifference)
}
