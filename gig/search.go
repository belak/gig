package main

import "gopkg.in/alecthomas/kingpin.v1"

var (
	searchCommand = kingpin.Command("search", "search for packages").Dispatch(search)
)

func search(c *kingpin.ParseContext) error {
	return nil
}
