// Package searchpackage provides functionalities for searching and processing text in files.
//
// The package includes functions to search for a specified word within text files (both .txt and .csv formats),
// process each line or record containing the word, and optionally colorize the word for better visibility.
// It offers the flexibility to either print the matching lines to the standard output or write them to a destination file.
//
// Usage:
// 1. Initialize the search criteria using InitSearch function, providing the word to search, input file path, and optional destination file path.
// 2. Execute the search using PerformSearch function, which processes the input file based on its extension (.txt or .csv).
//
// Example:
//	// Initialize search with a word, input file, and optional destination file
//	searchpackage.InitSearch("example", "input.txt", "output.txt")
//
//	// Perform the search
//	searchpackage.PerformSearch()
//
// In this example, the package searches for the word "example" in the "input.txt" file, colorizes occurrences of the word,
// and writes the matching lines to the "output.txt" file.
//
// Note: This package utilizes goroutines to concurrently process lines or records in the input file,
// enhancing performance for large files.

package searchpackage

// searchpackage is a simple package for searching and processing text in files.

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fatih/color"
)

var (
	word, inputFile, distPath string
	outputFile                *os.File
	wg                        sync.WaitGroup
	mutex                     sync.Mutex
)

func InitSearch(searchWord, file, destPath string) {
	word = searchWord
	inputFile = file
	distPath = destPath

	if len(word) < 1 {
		log.Fatal("no word is provided to be searched")
		os.Exit(1)
	}
	if len(inputFile) < 1 {
		log.Fatal("no input file is provided")
		os.Exit(1)
	}
}

func PerformSearch() {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	extension := fileExtension(inputFile)

	if distPath != "" {
		outputFile, err = os.Create(distPath)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer outputFile.Close()

	if extension == ".txt" {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			wg.Add(1)
			go processLine(line)
		}
	}
	if extension == ".csv" {
		csvReader := csv.NewReader(file)
		for {
			line, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal("error in reading csv file----> ", err)
			}
			wg.Add(1)
			go processCSV(line)
		}
	}

	wg.Wait()
}

func fileExtension(file string) string {
	return filepath.Ext(file)
}

func processCSV(lines []string) {
	defer wg.Done()
	for _, line := range lines {

		if containWord(line) {
			if distPath == "" {
				line = colorizeWord(line)
				fmt.Println(line)
			}
		}
	}
}

func processLine(line string) {
	defer wg.Done()
	if containWord(line) {
		if distPath == "" {
			line = colorizeWord(line)
			fmt.Println(line)
		}
		if distPath != "" {
			writeToFile(outputFile, line)
		}

	}
}

func containWord(textLine string) bool {
	return strings.Contains(strings.ToLower(textLine), strings.ToLower(word))
}

func writeToFile(outputFile *os.File, line string) {
	mutex.Lock()
	defer mutex.Unlock()
	fmt.Fprintln(outputFile, line)
}

func colorizeWord(line string) string {
	return strings.ReplaceAll(line, word, color.RedString(word))
}
