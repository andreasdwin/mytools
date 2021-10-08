package main

import (
	"flag"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestUserInputLogFileNotProvided(t *testing.T) {
	os.Args = []string{"mytools"}
	defer func() {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}()

	inp, err := getUserInput()
	if err == nil {
		t.Errorf("expected error, got %v", inp)
	}
}

func TestUserInputLogDefaultFlag(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "dummy_*.log")
	if err != nil {
		panic(err)
	}

	defer os.Remove(tmpfile.Name())

	os.Args = []string{"mytools", tmpfile.Name()}
	defer func() {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}()

	inp, err := getUserInput()

	expectedInput := userInput{
		logFilePath: tmpfile.Name(),
		outputType:  PlainText,
	}

	if err != nil {
		t.Errorf("expected not error, got %v", err)
		return
	}
	if !reflect.DeepEqual(inp, expectedInput) {
		t.Errorf("expected %v, got %v", expectedInput, inp)
	}
}

func TestUserInputLogJSONType(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "dummy_*.log")
	if err != nil {
		panic(err)
	}

	defer os.Remove(tmpfile.Name())

	os.Args = []string{"mytools", "-t=json", tmpfile.Name()}
	defer func() {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}()

	inp, err := getUserInput()

	expectedInput := userInput{
		logFilePath: tmpfile.Name(),
		outputType:  JSON,
	}

	if err != nil {
		t.Errorf("expected not error, got %v", err)
		return
	}
	if !reflect.DeepEqual(inp, expectedInput) {
		t.Errorf("expected %v, got %v", expectedInput, inp)
	}
}

func TestUserInputLogInvalidType(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "dummy_*.log")
	if err != nil {
		panic(err)
	}

	defer os.Remove(tmpfile.Name())

	os.Args = []string{"mytools", "-t=invalid", tmpfile.Name()}
	defer func() {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}()

	inp, err := getUserInput()

	if err == nil {
		t.Errorf("expected error, got %v", inp)
	}
}

func TestUserInputLogDefaultTypeWithOutputFile(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "dummy_*.log")
	if err != nil {
		panic(err)
	}

	defer os.Remove(tmpfile.Name())

	os.Args = []string{"mytools", "-o=dummy_log.txt", tmpfile.Name()}
	defer func() {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}()

	inp, err := getUserInput()

	expectedInput := userInput{
		logFilePath:    tmpfile.Name(),
		outputType:     PlainText,
		outputFilePath: "dummy_log.txt",
	}

	if err != nil {
		t.Errorf("expected not error, got %v", err)
		return
	}
	if !reflect.DeepEqual(inp, expectedInput) {
		t.Errorf("expected %v, got %v", expectedInput, inp)
	}
}

func TestUserInputLogJSONWithOutputFile(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "dummy_*.log")
	if err != nil {
		panic(err)
	}

	defer os.Remove(tmpfile.Name())

	os.Args = []string{"mytools", "-t=json", "-o=dummy_log.json", tmpfile.Name()}
	defer func() {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}()

	inp, err := getUserInput()

	expectedInput := userInput{
		logFilePath:    tmpfile.Name(),
		outputType:     JSON,
		outputFilePath: "dummy_log.json",
	}

	if err != nil {
		t.Errorf("expected not error, got %v", err)
		return
	}
	if !reflect.DeepEqual(inp, expectedInput) {
		t.Errorf("expected %v, got %v", expectedInput, inp)
	}
}

func TestReadLogFile(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "dummy_*.log")
	if err != nil {
		panic(err)
	}

	defer os.Remove(tmpfile.Name())

	logFile := []string{"08 Oct 2021 10:40:36.508 # Server initialized", "08 Oct 2021 10:40:36.513 # Ready to accept connections"}
	for _, log := range logFile {
		tmpfile.WriteString(log + "\n")
	}

	outputChannel := make(chan string)

	go func() {
		defer close(outputChannel)
		readLogFile(tmpfile.Name(), outputChannel)
	}()

	for _, log := range logFile {
		txt := <-outputChannel
		if log != txt {
			t.Errorf("expected %v, got %v", log, txt)
			return
		}
	}
}

func TestPrintOutputFilePlainText(t *testing.T) {
	logFile := []string{"08 Oct 2021 10:40:36.508 # Server initialized", "08 Oct 2021 10:40:36.513 # Ready to accept connections"}

	outputChannel := make(chan string)
	done := make(chan bool)
	go func() {
		for _, log := range logFile {
			outputChannel <- log
		}
		close(outputChannel)
	}()

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	go printLogFile(PlainText, outputChannel, done)

	<-done

	w.Close()
	output, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	expectedOutput := "08 Oct 2021 10:40:36.508 # Server initialized\n08 Oct 2021 10:40:36.513 # Ready to accept connections\n"

	if string(output) != expectedOutput {
		t.Errorf("expected %v, got %v", expectedOutput, string(output))
	}
}

func TestWriteOutputFilePlainText(t *testing.T) {
	outputFileName := "dummy_log.txt"

	logFile := []string{"08 Oct 2021 10:40:36.508 # Server initialized", "08 Oct 2021 10:40:36.513 # Ready to accept connections"}

	outputChannel := make(chan string)
	done := make(chan bool)
	go func() {
		for _, log := range logFile {
			outputChannel <- log
		}
		close(outputChannel)
	}()

	go writeOutputFile(outputFileName, PlainText, outputChannel, done)

	<-done

	outputFile, err := ioutil.ReadFile(outputFileName)
	if err != nil {
		t.Errorf("error opening output file, got %v", err)
		return
	}
	defer os.Remove(outputFileName)

	expectedFile := "08 Oct 2021 10:40:36.508 # Server initialized\n08 Oct 2021 10:40:36.513 # Ready to accept connections\n"

	if string(outputFile) != expectedFile {
		t.Errorf("expected %v, got %v", expectedFile, string(outputFile))
	}
}

func TestPrintOutputFileJSON(t *testing.T) {
	logFile := []string{"08 Oct 2021 10:40:36.508 # Server initialized", "08 Oct 2021 10:40:36.513 # Ready to accept connections"}

	outputChannel := make(chan string)
	done := make(chan bool)
	go func() {
		for _, log := range logFile {
			outputChannel <- log
		}
		close(outputChannel)
	}()

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	go printLogFile(JSON, outputChannel, done)

	<-done

	w.Close()
	output, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	expectedOutput := "[\n    \"08 Oct 2021 10:40:36.508 # Server initialized\",\n    \"08 Oct 2021 10:40:36.513 # Ready to accept connections\"\n]\n"

	if string(output) != expectedOutput {
		t.Errorf("expected %v, got %v", expectedOutput, string(output))
	}
}

func TestWriteOutputFileJSON(t *testing.T) {
	outputFileName := "dummy_log.json"

	logFile := []string{"08 Oct 2021 10:40:36.508 # Server initialized", "08 Oct 2021 10:40:36.513 # Ready to accept connections"}

	outputChannel := make(chan string)
	done := make(chan bool)
	go func() {
		for _, log := range logFile {
			outputChannel <- log
		}
		close(outputChannel)
	}()

	go writeOutputFile(outputFileName, JSON, outputChannel, done)

	<-done

	outputFile, err := ioutil.ReadFile(outputFileName)
	if err != nil {
		t.Errorf("error opening output file, got %v", err)
		return
	}
	defer os.Remove(outputFileName)

	expectedFile := "[\n    \"08 Oct 2021 10:40:36.508 # Server initialized\",\n    \"08 Oct 2021 10:40:36.513 # Ready to accept connections\"\n]\n"

	if string(outputFile) != expectedFile {
		t.Errorf("expected %v, got %v", expectedFile, string(outputFile))
	}
}
