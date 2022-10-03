package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
)

type TestHeapMap struct {
	testUpsert func(string)
	testPop    func() (string, int, error)
}

func (testHeapMap TestHeapMap) Upsert(sequence string) {
	testHeapMap.testUpsert(sequence)
}

func (testHeapMap TestHeapMap) Pop() (string, int, error) {
	return testHeapMap.testPop()
}

func createTestFile(data string) (*os.File, error) {
	filename := "/tmp/test_file.txt"
	if wErr := os.WriteFile(filename, []byte(data), 0666); wErr != nil {
		err := fmt.Errorf("failed to create write file in main_test; [error_msg: %v]", wErr)
		return nil, err
	}

	fileReader, oErr := os.Open(filename)
	if oErr != nil {
		err := fmt.Errorf("failed to open created file %v; [error_msg: %v]", filename, oErr)
		return nil, err
	}

	return fileReader, nil
}

func Test_parseFile(t *testing.T) {
	testFileData := "This is the first line of the test file.\n\n\t\tThis is the second line of the test file full of punctuations: [!@#$%^&*()-=_+|;':\",.<>?']."
	testFileReader, fileErr := createTestFile(testFileData)
	defer os.Remove(testFileReader.Name())
	if fileErr != nil {
		t.Error(fileErr)
	}

	testCh := make(chan []string)
	var testwg sync.WaitGroup
	testwg.Add(1)
	go parseFile(&testwg, testFileReader, testCh)

	lineArray1 := <-testCh
	lineArray2 := <-testCh
	line1 := strings.Join(lineArray1, " ")
	line2 := strings.Join(lineArray2, " ")
	expectedLine1 := "this is the first line of the test file"
	//line 2 should be wrapped with last words from 1st
	expectedLine2 := "test file this is the second line of the test file full of punctuations"

	if line1 != expectedLine1 {
		t.Errorf("error, unexpected line1 value read from channel; [expected: %v] [received: %v]", expectedLine1, line1)
	}

	if line2 != expectedLine2 {
		t.Errorf("error, unexpected line2 value read from channel; [expected: %v] [received: %v]", expectedLine2, line2)
	}
}

func Test_parseFile_OneWord(t *testing.T) {
	testFileData := "\t\tThis"
	testFileReader, fileErr := createTestFile(testFileData)
	defer os.Remove(testFileReader.Name())
	if fileErr != nil {
		t.Error(fileErr)
	}

	testCh := make(chan []string)
	var testwg sync.WaitGroup
	testwg.Add(1)
	go parseFile(&testwg, testFileReader, testCh)

	lineArray1 := <-testCh
	lineArray2 := <-testCh
	line1 := strings.Join(lineArray1, " ")
	line2 := strings.Join(lineArray2, " ")
	expectedLine1 := "this"
	expectedLine2 := ""

	if line1 != expectedLine1 {
		t.Errorf("error, unexpected line1 value read from channel; [expected: %v] [received: %v]", expectedLine1, line1)
	}

	if line2 != "" {
		t.Errorf("error, unexpected line2 value read from channel; [expected: %v] [received: %v]", expectedLine2, line2)
	}
}

func Test_storeSequences(t *testing.T) {
	var testHeapMap TestHeapMap
	testHeapMap.testUpsert = func(sequence string) {
		if sequence != "this is a test" {
			t.Errorf("error, unexpected sequence '%v' received from upsert", sequence)
		}
	}

	var testwg sync.WaitGroup
	testwg.Add(1)
	testCh := make(chan []string)
	testLineArray := []string{"this", "is", "a", "test"}
	testSize := 4
	go storeSequences(&testwg, testCh, testHeapMap, testSize)
	testCh <- testLineArray
	close(testCh)
}

func Test_getTopKSequences(t *testing.T) {
	var testHeapMap TestHeapMap
	testHeapMap.testPop = func() (string, int, error) {
		sequence := "this is a test"
		return sequence, 5, nil
	}

	defer func() {
		if r := recover(); r != nil {
			t.Error("error, unexpected panic from getTopKSequences")
		}
	}()

	getTopKSequences(testHeapMap, 2, "testfilename")
}

func Test_getTopKSequences_PopErr(t *testing.T) {
	var testHeapMap TestHeapMap
	testHeapMap.testPop = func() (string, int, error) {
		popErr := fmt.Errorf("test popErr")
		return "", -1, popErr
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("error, expected panic from getTopKSequences")
		}
	}()

	getTopKSequences(testHeapMap, 2, "testfilename")
}
