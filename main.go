/*
Copyright © 2023 Takafumi Miyanaga miya.org.0309@gmai.com
*/
package main

import (
	"time"

	"github.com/future-architect/tftarget/cmd"
)

var (
	version = "0.0.2"
)

func main() {
	cmd.SetVersionInfo(version, time.Now().String())
	cmd.Execute()
}
