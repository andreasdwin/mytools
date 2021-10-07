package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

const (
	PlainText string = "text"
	JSON string = "json"
)

func readLogFile(logFilePath string, outputChannel chan<- string) {
	lf, err := os.Open(logFilePath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer lf.Close()

	scanner := bufio.NewScanner(lf)
	for scanner.Scan() {
		outputChannel <- scanner.Text()
	}
}

func printLogFile(outputChannel <-chan string, done chan<- bool) {
	for txt := range outputChannel {
		fmt.Println(txt)
	}

	done <- true
}

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

	outputChannel := make(chan string)
	done := make(chan bool)

	go func() {
		defer close(outputChannel)
		readLogFile(logFilePath, outputChannel)
	}()

	go printLogFile(outputChannel, done)

	<-done
	close(done)

}
