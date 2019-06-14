package cmd

import (
	"flag"
	"fmt"

	"github.com/MathWebSearch/mwsapi/utils"
)

// Args represents command-line arguments
type Args struct {
	MWSHost string
	MWSPort int

	ElasticHost string
	ElasticPort int

	Host string
	Port int
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

	flagSet.IntVar(&flags.Port, "port", utils.GetenvInt("MWSAPID_PORT", 3000), "Port to listen on for queries")
	flagSet.StringVar(&flags.Host, "host", utils.Getenv("MWSAPID_HOST", "localhost"), "Host to listen on")

	flagSet.StringVar(&flags.MWSHost, "mws-host", utils.Getenv("MWSAPID_MWS_HOST", ""), "Host to use for mathwebsearch. If omitted, disable mathwebsearch support")
	flagSet.IntVar(&flags.MWSPort, "mws-port", utils.GetenvInt("MWSAPID_MWS_PORT", 8080), "Port to use for mathwebsearch")

	flagSet.StringVar(&flags.ElasticHost, "elastic-host", utils.Getenv("MWSAPID_ELASTIC_HOST", ""), "Host to use for elasticsearch")
	flagSet.IntVar(&flags.ElasticPort, "elastic-port", utils.GetenvInt("MWSAPID_ELASTIC_PORT", 9200), "Port to use for elasticsearch")

	// parse and exit
	flagSet.Parse(args[1:])

	return &flags
}

// Validate validates the command-line arguments or panics
func (args *Args) Validate() bool {

	if args.MWSPort <= 0 || args.MWSPort > 65535 {
		fmt.Printf("mws-port: %d is not a valid port number", args.MWSPort)
		return false
	}

	if args.ElasticPort <= 0 || args.ElasticPort > 65535 {
		fmt.Printf("elastic-port: %d is not a valid port number", args.MWSPort)
		return false
	}

	if args.Port <= 0 || args.Port > 65535 {
		fmt.Printf("port: %d is not a valid port number", args.Port)
		return false
	}

	return true
}
