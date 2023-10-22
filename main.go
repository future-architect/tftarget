/*
Copyright © 2023 Takafumi Miyanaga miya.org.0309@gmai.com
*/
package main

import (
	"time"

	"github.com/future-architect/tftarget/cmd"
)


func main() {
	cmd.SetVersionInfo(cmd.Version, time.Now().String())
	cmd.Execute()
}
