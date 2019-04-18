package utils

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
)

const bufferCapcityInBytes = 128 * 1024 // 128 MB

// ProcessLinePairs process lines from filename in pairs
func ProcessLinePairs(filename string, allowLeftover bool, parser func(string, string) error) (err error) {
	// load the file
	file, err := os.Open(filename)
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
			if err != nil {
				return err
			}

			firstLine = ""
			readFirstLine = false
		}
	}

	if readFirstLine && !allowLeftover {
		return errors.New("File did not contain an even number of lines")
	}

	// if something broke, throw an error
	return scanner.Err()
}

//IterateFiles iterates over files in a directory with a given extension
func IterateFiles(dir string, extension string, callback func(string) error) (err error) {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == extension {
			return callback(path)
		}
		return nil
	})
}
