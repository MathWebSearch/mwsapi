package args

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

// Args represents command-line arguments
type Args struct {
	ElasticHost string
	ElasticPort int

	DocumentPhaseOnly bool
	Count             bool

	Text string

	idStrings string
	IDs       []int64

	From int64
	Size int64
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

	flagSet.BoolVar(&flags.DocumentPhaseOnly, "document-phase-only", false, "When set, stop after the document query phase")
	flagSet.BoolVar(&flags.Count, "count", false, "When set, only count number of results instead of actually running the query")

	flagSet.StringVar(&flags.Text, "text", "", "Text to query for")
	flagSet.StringVar(&flags.idStrings, "ids", "", "Comma-seperated MathWebSearch IDs to query for")

	flagSet.Int64Var(&flags.From, "from", 0, "Slice to start results at")
	flagSet.Int64Var(&flags.Size, "size", 10, "Maximum number of results to return")

	// parse and exit
	flagSet.Parse(args[1:])
	return &flags
}

// Validate validates the command-line arguments or panics
func (args *Args) Validate() bool {

	//fmt.Printf("elastic-host: %q\n", args.ElasticHost)

	if args.ElasticPort <= 0 || args.ElasticPort > 65535 {
		fmt.Printf("elastic-port: %d is not a valid port number", args.ElasticPort)
		return false
	}
	//fmt.Printf("elastic-port: %d\n", args.ElasticPort)

	//fmt.Printf("document-phase-only: %t\n", args.DocumentPhaseOnly)
	// fmt.Printf("count: %t\n", args.Count)

	//fmt.Printf("text: %q\n", args.Text)

	// parse the ids
	if args.idStrings != "" {
		for _, e := range strings.Split(args.idStrings, ",") {
			n, err := strconv.ParseInt(e, 10, 64)
			if err != nil {
				fmt.Printf("ids: %q is not a valid number: %s\n", e, err.Error())
				return false
			}
			args.IDs = append(args.IDs, n)
		}
	} else {
		args.IDs = []int64{}
	}
	//fmt.Printf("ids: %#v\n", args.IDs)

	if args.From < 0 {
		fmt.Println("from must be at least 0")
		return false
	}
	//fmt.Printf("from: %d\n", args.From)

	if args.Size < 1 {
		fmt.Println("size must be at least 1")
		return false
	}
	//fmt.Printf("size: %d\n", args.Size)

	//fmt.Println("------------------------------------------")
	return true
}
