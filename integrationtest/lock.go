package integrationtest

import "path"
import filemutex "github.com/alexflint/go-filemutex"

var mutex *filemutex.FileMutex

// LockTests prevents any other processes from running tests at the same time
func LockTests() {
	if err := mutex.Lock(); err != nil {
		panic(err)
	}
}

// UnLockTests reverses a previous Locktests call
func UnLockTests() {
	if err := mutex.Unlock(); err != nil {
		panic(err)
	}
}

func init() {
	var err error
	mutex, err = filemutex.New(path.Join(TestDataPath, "lockfile"))
	if err != nil {
		panic(err)
	}
}
