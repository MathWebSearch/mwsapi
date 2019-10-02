package utils

import (
	"github.com/json-iterator/go"
	"fmt"
	"os"
)

// OutputJSONOrErr outputs res as json if no error message is provided
func OutputJSONOrErr(res interface{}, err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// stdout the json
	bytes, _ := jsoniter.MarshalIndent(res, "", "  ")
	os.Stdout.Write(bytes)
	fmt.Println("")
}
