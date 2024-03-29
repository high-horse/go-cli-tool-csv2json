package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type inputFile struct {
	filepath  string
	seperator string
	pretty    bool
}

func main() {

}

func check(err error) {
	if err != nil {
		exitGracefully(err)
	}
}

func exitGracefully(err error) {
	fmt.Fprintf(os.Stderr, "error %v\n", err)
	os.Exit(1)
}

func getFileData() (inputFile, error) {
	// validate correct number of args.
	if len(os.Args) < 2 {
		return inputFile{}, errors.New("a filepath argument is required")
	}

	// Defining option flags.
	// 3 args : flag name, default value, and short desc (desplayed with --help flag)
	seperator := flag.String("seperator", "comma", "Column seperator")
	pretty := flag.Bool("pretty", false, "Generate pretty JSON")

	// Parse the command line arguments
	flag.Parse()

	// File location
	fileLocation := flag.Arg(0)

	// Check whether the seperator is defined or not
	if !(*seperator == "comma" || *seperator == "semicolon") {
		return inputFile{}, errors.New("Only comma or semicolon are allowed")
	}

	// return struct with required data
	return inputFile{fileLocation, *seperator, *pretty}, nil
}

func checkIfValidFile(filename string) (bool, error) {
	if fileExtension := filepath.Ext(filename); fileExtension != ".csv" {
		return false, fmt.Errorf("file %s is not  CSV", filename)
	}

	// check filepath belong to existing file
	if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
		return false, fmt.Errorf("file %s does not exist", filename)
	}

	return true, nil
}

// Process csv args filedata struct, and writer channel
func processCsvFile(fileData inputFile, writerCh chan<- map[string]string) {
	// OPen file to read
	file, err := os.Open(fileData.filepath)
	check(err)

	// close channel
	defer file.Close()

	// Define headers and line slice
	var headers, line []string

	// Initialize csv headers
	reader := csv.NewReader(file)

	if fileData.seperator == "semicolon" {
		reader.Comma = ';'
	}

	// read first line i.e headers
	headers, err = reader.Read()
	check(err)

	// iterate over each line in csv file
	for {
		line, err = reader.Read()

		// if reach end of file, break ch and loop
		if err == io.EOF {
			close(writerCh)
			break
		} else if err != nil {
			exitGracefully(err)
		}

		// Process csv line
		record, err := processLine(headers, line)

		if err != nil {
			fmt.Printf("Line: %v, Error %s\n", line, err)
			continue
		}

		writerCh <- record
	}
}

func processLine(headers []string, dataList []string) (map[string]string, error) {
	if len(dataList) != len(headers) {
		return nil, errors.New("line doesnot match headers format, skipping")
	}

	recordMap := make(map[string]string)

	for i, name := range headers {
		recordMap[name] = dataList[i]
	}

	return recordMap, nil
}
