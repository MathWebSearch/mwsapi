package args

import (
	"flag"
	"fmt"
	"os"
)

// Args represents command-line arguments
type Args struct {
	ElasticHost string
	ElasticPort int

	IndexDir string
}

// ElasticURL returns the url to elasticsearch
func (args *Args) ElasticURL() string {
	return fmt.Sprintf("http://%s:%d", args.ElasticHost, args.ElasticPort)
}

// ParseArgs parses arguments from a list of strings
func ParseArgs(args []string) *Args {
	var flags Args

	// create a new flagset
	// that prints it's usage on --help
	flagSet := flag.NewFlagSet(args[0], flag.ExitOnError)
	flagSet.Usage = func() {
		fmt.Fprintf(flagSet.Output(), "Usage of %s:\n", args[0])
		flagSet.PrintDefaults()
	}

	flagSet.StringVar(&flags.ElasticHost, "elastic-host", "0.0.0.0", "Host to use for elasticsearch")
	flagSet.IntVar(&flags.ElasticPort, "elastic-port", 9200, "Port to use for elasticsearch")

	defaultIndexDir := "/index/"
	flagSet.StringVar(&flags.IndexDir, "index-dir", defaultIndexDir, "Directory to use for Indexes")

	// parse and exit
	flagSet.Parse(args[1:])
	return &flags
}

// Validate validates the command-line arguments or panics
func (args *Args) Validate() bool {

	fmt.Printf("elastic-host: %q\n", args.ElasticHost)

	if args.ElasticPort <= 0 || args.ElasticPort > 65535 {
		fmt.Printf("elastic-port: %d is not a valid port number", args.ElasticPort)
		return false
	}
	fmt.Printf("elastic-port: %d\n", args.ElasticPort)

	if !ensureDirectory(args.IndexDir) {
		fmt.Printf("index-dir: %q is not a directory\n", args.IndexDir)
		return false
	}
	fmt.Printf("index-dir: %q\n", args.IndexDir)

	fmt.Println("------------------------------------------")
	return true
}

// ensureDirectoryOrPanic ensures that caniddate is a directory or otherwise panics with message
func ensureDirectory(candidate string) bool {
	fi, err := os.Stat(candidate)
	if err != nil {
		return false
	}

	mode := fi.Mode()
	if !mode.IsDir() {
		return false
	}

	return true
}
