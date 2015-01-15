package tunes

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"

	"../config"
	"../parser"
)

type CoreConfig struct {
	Prefixdir string
}

const (
	BUFSIZE = 8192 // Yes, this is completely arbitrary. Sue me.
)

func Download(tune *parser.Env, conf *config.Config) error {
	confVals := &CoreConfig{}
	err := conf.Load("core", confVals)
	if err != nil {
		return err
	}

	url, err := tune.GetString("pkg-url")
	if err != nil {
		return err
	}

	checksum, err := tune.GetString("pkg-sha1")
	if err != nil {
		return err
	}

	// Check if prefixdir exists and create if it doesn't
	src, err := os.Stat(confVals.Prefixdir + "src")
	if err != nil {
		// Create prefix directory
		err = os.MkdirAll(confVals.Prefixdir+"src", 0755)
		if err != nil {
			return err
		}
	} else {
		if !src.IsDir() {
			fmt.Errorf("Prefix directory exists and is not a directory")
		}
	}

	base := path.Base(url)
	fmt.Printf("Downloading %s...\n\n", url)

	res, err := http.Get(url)
	if err != nil {
		return err
	}

	outFile, err := os.OpenFile(confVals.Prefixdir+"src/"+base,
		os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer outFile.Close()

	totalBytes, err := strconv.ParseInt(res.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return err
	}

	hash := sha1.New()
	downloadedBytes := 0
	buffer := make([]byte, BUFSIZE)

	fmt.Printf("0/%d bytes (0%%)", totalBytes)
	for {
		bytes, err := res.Body.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		hash.Write(buffer[:bytes])
		outFile.Write(buffer[:bytes])
		downloadedBytes += bytes
		fmt.Printf("\r%d/%d bytes (%d%%)", downloadedBytes, totalBytes, int(float64(downloadedBytes)/float64(totalBytes)*100.0))

		if err == io.EOF {
			break
		}
	}

	if checksum != fmt.Sprintf("%x", hash.Sum(nil)) {
		return fmt.Errorf("Checksums not equal")
	}

	fmt.Printf("\nDownloaded %d bytes\n", downloadedBytes)

	return nil
}
