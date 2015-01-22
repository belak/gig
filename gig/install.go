package main

import "gopkg.in/alecthomas/kingpin.v1"

var (
	installCommand = kingpin.Command("install", "install packages with dependencies").Dispatch(install)
)

func install(c *kingpin.ParseContext) error {
	return nil
}
