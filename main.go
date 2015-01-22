package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"./config"
	"./parser"
	"./tunes"
)

func main() {
	// TODO: make right
	if len(os.Args) < 2 {
		fmt.Println("Usage: gig <file.tune>\n")
		os.Exit(1)
	}

	var conf *config.Config

	// TODO: take config arg
	conf, err := config.NewDefaultConfig()
	if err != nil {
		fmt.Printf("Error loading config file, %s\n", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "fetch":
		if len(os.Args) < 3 {
			fmt.Println("Usage: gig fetch <package>\n")
			os.Exit(1)
		}
		fetchSource(conf, os.Args[2])
	case "search":
		if len(os.Args) < 3 {
			fmt.Println("Usage: gig search <package>\n")
			os.Exit(1)
		}
		search(os.Args[2])
	default:
		parseTunefile(os.Args[1])
	}
}

func parseTunefile(filename string) (*parser.Env, error) {
	env, err := parser.NewBootstrappedEnv()
	if err != nil {
		return nil, fmt.Errorf("Error creating new environment, %s\n", err)
	}

	node, err := env.LoadTune(filename)
	if err != nil {
		return nil, fmt.Errorf("Error parsing file, %s\n", err)
	}

	_, err = env.Eval(node)
	if err != nil {
		return nil, fmt.Errorf("Error running tunefile, %s\n", err)
	}

	return env, nil
}

// Searches tunes for package
func search(name string) {
	files, err := ioutil.ReadDir("tunes")
	if err != nil {
		fmt.Printf("Error searching tunes, %s\n", err)
		os.Exit(1)
	}

	for _, file := range files {
		if file.Name()[0] == '.' {
			continue
		}

		if file.Name() == name+".tune" {
			fmt.Printf("Package %s found\n", name)
			return
		}
	}

	fmt.Printf("Package %s not found\n", name)
}

// Downloads and extracts package source
func fetchSource(conf *config.Config, name string) {
	filename := "tunefiles/" + name + ".tune"

	tune, err := parseTunefile(filename)
	if err != nil {
		fmt.Printf("Error loading tunefile, %s\n", err)
		os.Exit(1)
	}

	err = tunes.Download(tune, conf)
	if err != nil {
		fmt.Printf("Error downloading tune, %s\n", err)
		os.Exit(1)
	}
}
