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
		outputType: PlainText,
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
		outputType: JSON,
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
		logFilePath: tmpfile.Name(),
		outputType: PlainText,
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
		logFilePath: tmpfile.Name(),
		outputType: JSON,
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