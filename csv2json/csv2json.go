package main

import (
	"errors"
	"flag"
	"os"
)

type inputFile struct {
	filepath  string
	seperator string
	pretty    bool
}

func main() {

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
	if !(*seperator == "comma" ||*seperator == "semicolon") {
		return inputFile{}, errors.New("Only comma or semicolon are allowed")
	}

	// return struct with required data
	return inputFile{fileLocation, *seperator, *pretty}, nil
}