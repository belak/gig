package main

import (
	"fmt"
	"os"
	"strings"

	"./parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: gig <file.tune>\n")
		os.Exit(1)
	}

	env := parser.NewEnv()
	node, err := env.ParseFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error parsing file, %s\n", err)
		os.Exit(1)
	}

	env.Eval(node)

	env.Invoke("pkg-install", []interface{}{})

	desc, err := env.GetString("pkg-description")
	if err != nil {
		fmt.Printf("Error fetching description, %s\n", err)
		os.Exit(1)
	}
	fmt.Println(desc)

	deps, err := env.GetStringArray("pkg-dependencies")
	if err != nil {
		fmt.Printf("Error fetching dependencies, %s\n", err)
		os.Exit(1)
	}
	fmt.Println(strings.Join(deps, ", "))
}
