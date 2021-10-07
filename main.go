package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	PlainText string = "text"
	JSON string = "json"
)

func main() {
	outputType := flag.String("t", PlainText, "output type (json or text)")
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("please provide the log file")
		return
	}

	logFilePath := flag.Arg(0)

	if !(*outputType == PlainText || *outputType == JSON) {
		fmt.Println("invalid output type")
		return
	}

	_, err := os.Stat(logFilePath)
	if err != nil && os.IsNotExist(err) {
		fmt.Println("the log file does not exist")
		return
	}
}