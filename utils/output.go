package utils

import (
	"encoding/json"
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
	bytes, _ := json.MarshalIndent(res, "", "  ")
	if err != nil {
		fmt.Printf("%#v\n", err)
		os.Exit(1)
	}
	os.Stdout.Write(bytes)
	fmt.Println("")
}
