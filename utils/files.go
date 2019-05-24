package utils

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const bufferCapcityInBytes = 128 * 1024 // 128 MB

// ProcessLinePairs process lines from filename in pairs
func ProcessLinePairs(filename string, allowLeftover bool, parser func(string, string) error) (err error) {
	// load the file
	file, err := os.Open(filename)
	err = errors.Wrap(err, "os.Open failed")
	if err != nil {
		return err
	}
	defer file.Close()

	// read a file line by line
	scanner := bufio.NewScanner(file)
	//adjust the capacity to your need (max characters in line)
	buf := make([]byte, bufferCapcityInBytes*1024)
	scanner.Buffer(buf, bufferCapcityInBytes*1024)

	readFirstLine := false
	var firstLine string
	for scanner.Scan() {
		// we have to read the first line first
		if !readFirstLine {
			firstLine = scanner.Text()
			readFirstLine = true

			// we read the first one already, so read the second one
		} else {
			err := parser(firstLine, scanner.Text())
			err = errors.Wrap(err, "parser failed")
			if err != nil {
				return err
			}

			firstLine = ""
			readFirstLine = false
		}
	}

	if readFirstLine && !allowLeftover {
		return errors.New("[ProcessLinePairs] File did not contain an even number of lines")
	}

	// if something broke, throw an error
	err = scanner.Err()
	err = errors.Wrap(err, "scanner.Err failed")
	return
}

//IterateFiles iterates over files in a directory with a given extension
func IterateFiles(dir string, extension string, callback func(string) error) (err error) {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == extension {
			err = callback(path)
			err = errors.Wrap(err, "callback failed")
			return err
		}
		return nil
	})
}

// HashFile computes the hash of a segment
func HashFile(filename string) (hash string, err error) {
	// the hasher implementation
	hasher := sha256.New()

	// open the segment
	f, err := os.Open(filename)
	err = errors.Wrap(err, "os.open failed")
	if err != nil {
		return
	}
	defer f.Close()

	// start hashing the file
	if _, err := io.Copy(hasher, f); err != nil {
		err = errors.Wrap(err, "io.Copy failed")
		return "", err
	}

	// turn it into a string
	hash = hex.EncodeToString(hasher.Sum(nil))
	return
}
