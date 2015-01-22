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
	url, err := tune.GetString("pkg-url")
	if err != nil {
		return err
	}

	checksum, err := tune.GetString("pkg-sha1")
	if err != nil {
		return err
	}

	defaultDirInterface, err := conf.Get("prefixdir")
	if err != nil {
		return err
	}

	var defaultDir string
	var ok bool
	if defaultDir, ok = defaultDirInterface.(string); !ok {
		return fmt.Errorf("Error converting defaultDir to string")
	}

	fmt.Println(defaultDir)

	// Check if prefixdir exists and create if it doesn't
	src, err := os.Stat(defaultDir)
	if err != nil {
		// Create prefix directory
		err = os.MkdirAll(defaultDir, 0755)
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

	outFile, err := os.OpenFile(defaultDir+base,
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
		bytes, err := io.CopyN(w, res.Body, BUFSIZE)
		downloadedBytes += bytes
		fmt.Printf("\r%d/%d bytes (%d%%)", downloadedBytes, totalBytes, int(float64(downloadedBytes)/float64(totalBytes)*100.0))

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
	}

	if checksum != fmt.Sprintf("%x", hash.Sum(nil)) {
		return fmt.Errorf("Checksums not equal")
	}

	return nil
}
