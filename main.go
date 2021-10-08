package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
)

const (
	PlainText string = "text"
	JSON string = "json"
)

type userInput struct {
	logFilePath string
	outputType string
	outputFilePath string
}

func getUserInput() (userInput, error) {
	outputType := flag.String("t", PlainText, "output type (json or text)")
	outputFilePath := flag.String("o", "", "output file path")
	flag.Parse()

	if len(os.Args) < 2 {
		return userInput{}, errors.New("log file required")
	}

	logFilePath := flag.Arg(0)

	if !(*outputType == PlainText || *outputType == JSON) {
		return userInput{}, errors.New("invalid output type")
	}

	_, err := os.Stat(logFilePath)
	if err != nil && os.IsNotExist(err) {
		return userInput{}, errors.New("log file not exist")
	}

	inp := userInput{
		logFilePath: logFilePath,
		outputType: *outputType,
		outputFilePath: *outputFilePath,
	}

	return inp, nil
}

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

func printLogFile(outputType string, outputChannel <-chan string, done chan<- bool) {
	if outputType == PlainText {
		for txt := range outputChannel {
			fmt.Println(txt)
		}
	} else {
		fmt.Print("[")

		first := true
		for str := range outputChannel {

			if first {
				first = false
			} else {
				fmt.Print(",")
			}
			jbyte, _ := json.Marshal(str)
			fmt.Print("\n    " + string(jbyte))
		}
		if first {
			fmt.Println("]")
		} else {
			fmt.Println("\n]")
		}
	}

	done <- true
}

func writeOutputFile(filePath string, outputType string, outputChannel <-chan string, done chan<- bool) {
	of, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer of.Close()

	if outputType == PlainText {
		for txt := range outputChannel {
			of.WriteString(txt + "\n")
		}
	} else {
		of.WriteString("[")

		first := true
		for str := range outputChannel {

			if first {
				first = false
			} else {
				of.WriteString(",")
			}
			jbyte, _ := json.Marshal(str)
			of.WriteString("\n    " + string(jbyte))
		}
		if first {
			of.WriteString("]")
		} else {
			of.WriteString("\n]\n")
		}
	}

	done <- true
}

func main() {
	input, err := getUserInput()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	outputChannel := make(chan string)
	done := make(chan bool)

	go func() {
		defer close(outputChannel)
		readLogFile(input.logFilePath, outputChannel)
	}()

	if input.outputFilePath != "" {
  		go writeOutputFile(input.outputFilePath, input.outputType, outputChannel, done)
 	} else {
 		go printLogFile(input.outputType, outputChannel, done)
 	}

	<-done
	close(done)

}
