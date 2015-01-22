package main

import "gopkg.in/alecthomas/kingpin.v1"

var (
	upgradeCommand = kingpin.Command("upgrade", "upgrade packages").Dispatch(upgrade)
)

func upgrade(c *kingpin.ParseContext) error {
	return nil
}
