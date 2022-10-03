package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/jhawk7/sequence-frequency/pkg/common"
	"github.com/jhawk7/sequence-frequency/pkg/sequence_heap_map"
)

func main() {
	if len(os.Args) > 1 {
		processFilePaths()
	} else {
		processStdin()
	}
}

func processStdin() {
	var wg sync.WaitGroup
	wg.Add(2)
	lineChan := make(chan []string)
	fileInfo, statErr := os.Stdin.Stat()
	sequenceSize := 3
	if statErr != nil {
		err := fmt.Errorf("failed to read from data from stdin; [err_msg: %v]", statErr)
		common.ErrorHandler(err, true)
	}

	sequenceHeapMap := sequence_heap_map.InitSequenceHeapMap()
	go parseFile(&wg, os.Stdin, lineChan) //pass stdin as io.Reader
	go storeSequences(&wg, lineChan, sequenceHeapMap, sequenceSize)
	wg.Wait()
	getTopKSequences(sequenceHeapMap, 100, fileInfo.Name())
}

func processFilePaths() {
	filenames := os.Args[1:]
	for i := 0; i < len(filenames); i++ {
		var wg sync.WaitGroup
		wg.Add(2)
		lineChan := make(chan []string)
		sequenceHeapMap := sequence_heap_map.InitSequenceHeapMap()
		sequenceSize := 3
		fileReader := common.OpenFile(filenames[i])
		defer fileReader.Close()
		go parseFile(&wg, fileReader, lineChan) //pass filereader as io.Reader
		go storeSequences(&wg, lineChan, sequenceHeapMap, sequenceSize)
		wg.Wait()
		getTopKSequences(sequenceHeapMap, 100, filenames[i])
	}
}

func parseFile(wg *sync.WaitGroup, fileReader io.Reader, lineChan chan<- []string) {
	scanner := bufio.NewScanner(fileReader)
	prevLastWords := ""

	//read file line by line
	for scanner.Scan() {
		lowercaseLine := strings.ToLower(scanner.Text())
		//remove punctuations and indentions
		replacer := strings.NewReplacer(".", "", "!", "", "\"", "", "'", "", "#", "", "$", "", "%", "", "&", "", "(", "", ")", "", "*", "", "+",
			"", ",", "", "-", "", "/", "", ":", "", ";", "", "?", "", "@", "", "[", "", "\\", "", "]", "", "^", "", "_", "", "`", "", "{", "", "|", "", "}", "", "~", "", "\t", "", "\n", "",
			"<", "", ">", "", "=", "")
		parsedLine := strings.TrimSpace(replacer.Replace(lowercaseLine))

		if len(parsedLine) == 0 {
			continue
		}

		//add last two words from prev line to beginning of new line
		var wrappedLine string
		if len(prevLastWords) > 0 {
			wrappedLine = prevLastWords + " " + parsedLine
		} else {
			wrappedLine = parsedLine
		}

		//split wrappedline into array and save last two words for next wrapped line
		lineArray := strings.Split(wrappedLine, " ")
		if len(lineArray)-2 >= 0 {
			prevLastWord1 := lineArray[len(lineArray)-2]
			prevLastWord2 := lineArray[len(lineArray)-1]
			prevLastWords = strings.TrimSpace(prevLastWord1 + " " + prevLastWord2)
		} else if len(lineArray)-1 >= 0 {
			prevLastWords = strings.TrimSpace(lineArray[len(lineArray)-1])
		} else {
			prevLastWords = ""
		}

		//feed lineArray into channel
		lineChan <- lineArray
	}
	close(lineChan)
	wg.Done()
}

func storeSequences(wg *sync.WaitGroup, lineChan <-chan []string, sequenceHeapMap sequence_heap_map.ISequenceHeapMap, size int) {
	//parse each lineArray from feed and store sequences/frequency
	for line := range lineChan {
		start := 0
		for end := 0; end < len(line); end++ {
			if (end - start + 1) < size {
				continue
			}
			sequence := strings.TrimSpace(strings.Join(line[start:end+1], " "))
			sequenceHeapMap.Upsert(sequence)
			start++
		}
	}
	wg.Done()
}

func getTopKSequences(sequenceHeapMap sequence_heap_map.ISequenceHeapMap, k int, filename string) {
	fmt.Printf("Top %v Seqences for File: %v\n", k, filename)
	for i := 1; i <= k; i++ {
		sequence, frequency, popErr := sequenceHeapMap.Pop()
		if popErr != nil {
			common.ErrorHandler(popErr, true)
			break
		}
		fmt.Printf("%v - %v\n", sequence, frequency)
	}
	fmt.Println("")
}
