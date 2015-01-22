package main

import "gopkg.in/alecthomas/kingpin.v1"

var (
	infoCommand = kingpin.Command("info", "get info for a specific package").Dispatch(info)
)

func info(c *kingpin.ParseContext) error {
	return nil
}
