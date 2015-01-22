package main

import "gopkg.in/alecthomas/kingpin.v1"

var (
	updateCommand = kingpin.Command("update", "update the package database").Dispatch(update)
)

func update(c *kingpin.ParseContext) error {
	return nil
}
