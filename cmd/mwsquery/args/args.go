package args

import (
	"flag"
	"fmt"
)

// Args represents command-line arguments
type Args struct {
	MWSHost string
	MWSPort int

	Expressions []string

	MWSIdsOnly bool
	Count      bool

	From int64
	Size int64
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

	flagSet.StringVar(&flags.MWSHost, "mws-host", "0.0.0.0", "Host to use for mathwebsearch")
	flagSet.IntVar(&flags.MWSPort, "mws-port", 8080, "Port to use for mathwebsearch")

	flagSet.BoolVar(&flags.MWSIdsOnly, "ids", false, "When set, indicate for mathwebsearch to return only ids instead of full results")
	flagSet.BoolVar(&flags.Count, "count", false, "When set, only count number of results instead of actually running the query")

	flagSet.Int64Var(&flags.From, "from", 0, "Slice to start results at")
	flagSet.Int64Var(&flags.Size, "size", 10, "Maximum number of results to return")

	// parse and exit
	flagSet.Parse(args[1:])
	flags.Expressions = flagSet.Args()

	return &flags
}

// Validate validates the command-line arguments or panics
func (args *Args) Validate() bool {

	if len(args.Expressions) == 0 {
		fmt.Println("expressions: At least one argument is required")
		return false
	}

	if args.MWSPort <= 0 || args.MWSPort > 65535 {
		fmt.Printf("mws-port: %d is not a valid port number", args.MWSPort)
		return false
	}
	//fmt.Printf("mws-port: %d\n", args.MWSPort)

	// fmt.Printf("count: %t\n", args.Count)

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
