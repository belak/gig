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
	defer res.Body.Close()

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
	var downloadedBytes int64 = 0

	w := io.MultiWriter(outFile, hash)

	fmt.Printf("0/%d bytes (0%%)", totalBytes)
	for {
		// TODO: calculate checksum here and compare at end
		bytes, err := io.CopyN(w, res.Body, BUFSIZE)
		downloadedBytes += bytes

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		fmt.Printf("\r%d/%d bytes (%d%%)", downloadedBytes, totalBytes, int(float64(downloadedBytes)/float64(totalBytes)*100.0))
	}

	fmt.Printf("\r%d/%d bytes (%d%%)", downloadedBytes, totalBytes, int(float64(downloadedBytes)/float64(totalBytes)*100.0))
	fmt.Printf("\nDownloaded %d bytes\n", downloadedBytes)
	fmt.Printf("expected: %s\ncalcul'd: %x\n", checksum, hash.Sum(nil))

	if checksum != fmt.Sprintf("%x", hash.Sum(nil)) {
		return fmt.Errorf("Checksums not equal")
	}

	return nil
}
