package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"./config"
	"./parser"
)

var conf *config.Config

type CoreConfig struct {
	Prefixdir string
}

var confVals *CoreConfig

func main() {
	// TODO: make right
	if len(os.Args) < 2 {
		fmt.Println("Usage: gig <file.tune>\n")
		os.Exit(1)
	}

	conf, err := config.NewConfig("../configs/gig.toml")
	if err != nil {
		fmt.Printf("Error loading config file, %s\n", err)
		os.Exit(1)
	}

	confVals = &CoreConfig{}
	err = conf.Load("core", confVals)
	if err != nil {
		fmt.Printf("Error loading config values, %s\n", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "fetch":
		if len(os.Args) < 3 {
			fmt.Println("Usage: gig fetch <package>\n")
			os.Exit(1)
		}
		fetchSource(os.Args[2])
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
	env, err := parser.NewEnv()
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
		if file.Name() == name+".tune" {
			fmt.Printf("Package %s found\n", name)
			return
		}
	}

	fmt.Printf("Package %s not found\n", name)
}

// Downloads and extracts package source
func fetchSource(name string) {
	filename := "tunes/" + name + ".tune"

	env, err := parseTunefile(filename)
	if err != nil {
		fmt.Printf("Error loading tunefile, %s\n", err)
		os.Exit(1)
	}

	url, err := env.GetString("pkg-url")
	if err != nil {
		fmt.Printf("Error retrieving package URL, %s\n", err)
		os.Exit(1)
	}

	checksum, err := env.GetString("pkg-sha1")
	if err != nil {
		fmt.Printf("Error retrieving package checksum, %s\n", err)
		os.Exit(1)
	}

	downloadSource(url, checksum)
}

func downloadSource(url, checksum string) {
	// Check if prefixdir exists and create if it doesn't
	src, err := os.Stat(confVals.Prefixdir)
	if err != nil {
		// Create prefix directory
		err = os.MkdirAll(confVals.Prefixdir, 0755)
		if err != nil {
			fmt.Println("Error creating prefix directory, %s\n", err)
			os.Exit(1)
		}
	} else {
		if !src.IsDir() {
			fmt.Println("Prefix directory exists and is not a directory")
			os.Exit(1)
		}
	}

	base := path.Base(url)
	fmt.Printf("Downloading %s...\n", url)

	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error establishing connection, %s\n", err)
		os.Exit(1)
	}

	data, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("Error downloading source, %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Downloaded %d bytes\n", len(data))

	fmt.Printf("%s\n", checksum)

	// Compare checksums
	if checksum != calcChecksum(data) {
		fmt.Printf("Checksums do not match, exiting\n")
		os.Exit(1)
	}

	// Ungzip, untar, and write
	var reader io.Reader
	reader = bytes.NewReader(data)

	if strings.HasSuffix(base, ".gz") || strings.HasSuffix(base, ".tgz") {
		reader = gunzip(reader)
	}

	tarReader := untar(reader)

	idx := strings.LastIndex(base, ".")
	if idx > 0 {
		base = base[0:idx]
	}

	dirname, err := ioutil.TempDir(confVals.Prefixdir, base+"_")

	if err != nil {
		fmt.Printf("Error creating temporary directory, %s\n", err)
		os.Exit(1)
	}

	err = os.Chmod(dirname, 0755)
	if err != nil {
		fmt.Printf("Error chmodding temporary directory, %s\n", err)
		os.Exit(1)
	}

	for {
		header, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			os.Exit(1)
		}

		// get the individual filename and extract to the current directory
		filename := header.Name

		switch header.Typeflag {
		case tar.TypeDir:
			// handle directory
			fmt.Println("Creating directory :", filename)
			err = os.MkdirAll(dirname+"/"+filename, 0755) // os.FileMode(header.Mode)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		case tar.TypeReg:
			fallthrough
		case tar.TypeRegA:
			// handle normal file
			fmt.Println("Untarring :", filename)
			writer, err := os.Create(dirname + "/" + filename)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			io.Copy(writer, tarReader)

			err = os.Chmod(dirname+"/"+filename, os.FileMode(header.Mode))

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			writer.Close()
		default:
			fmt.Printf("Unable to untar type : %c in file %s\n", header.Typeflag, dirname+"/"+filename)
		}
	}
}

func calcChecksum(data []byte) string {
	hash := sha1.New()

	io.WriteString(hash, string(data))

	fmt.Printf("%x\n", hash.Sum(nil))

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func gunzip(data io.Reader) io.Reader {
	fmt.Println("Unzipping...")

	reader, err := gzip.NewReader(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer reader.Close()

	return reader
}

func untar(data io.Reader) *tar.Reader {
	fmt.Println("Untarring...")

	reader := tar.NewReader(data)

	return reader
}
