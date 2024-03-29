package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"reflect"
	"testing"
)

func Test_getFileData(t *testing.T) {
	// Define test slice
	tests := []struct {
		name    string
		want    inputFile
		wantErr bool
		osArgs  []string
	}{
		{"Default parameters", inputFile{"test.csv", "comma", false}, false, []string{"cmd", "test.csv"}},
		{"No parameters", inputFile{}, true, []string{"cmd"}},
		{"Semicolon enabled", inputFile{"test.csv", "semicolon", false}, false, []string{"cmd", "--seperator=semicolon", "test.csv"}},
		{"Pretty enabled", inputFile{"test.csv", "comma", true}, false, []string{"cmd", "--pretty", "test.csv"}},
		{"Pretty and semicolon enabled", inputFile{"test.csv", "semicolon", true}, false, []string{"cmd", "--pretty", "--seperator=semicolon", "test.csv"}},
		{"Separator not identified", inputFile{}, true, []string{"cmd", "--seperator=pipe", "test.csv"}},
	}

	// Iterate over test slices
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualOsArgs := os.Args

			// defer func to handle cleanup
			defer func() {
				// restore original os.Args ref
				os.Args = actualOsArgs
				// reset the flag command line to parse it again
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			}()

			// set specific command for the test
			os.Args = test.osArgs

			// run the test
			got, err := getFileData()

			// assert whether or now we get the correct error value
			if (err != nil) != test.wantErr {
				t.Errorf("getFileData() error = %v, wantErr %v", err, test.wantErr)
			}

			// assert whether or not we get correct wanted value
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("getFileData() = %v, want %v", got, test.want)
			}

		})
	}
}

func Test_checkIfValidFile(t *testing.T) {
	// Create temp file and emptyu csv file
	tmpfile, err := ioutil.TempFile("", "test*.csv")
	if err != nil {
		panic(err)
	}

	// delete temp file
	defer os.Remove((tmpfile.Name()))

	// Define test struct
	tests := []struct {
		name     string
		filename string
		want     bool
		wantErr  bool
	}{
		// Define test cases
		{"file exist", tmpfile.Name(), true, false},
		{"file not exist", "nowhere/test.csv", false, true},
		{"file not csv", "test.csv", false, true},
	}

	// iterate over test
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := checkIfValidFile(test.filename)
			if (err != nil) != test.wantErr {
				t.Errorf("checkIfValidFIle() error = %v, wanterr = %v", err, test.wantErr)
				return
			}

			// Check the returning val
			if got != test.want {
				t.Errorf("checkIfValidFIle() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_processCsvFile(t *testing.T) {
	// Define map that is to be expected from out functino
	wantMapSLice := []map[string]string{
		{"COL1": "1", "COL2": "2", "COL3": "3"},
		{"COL1": "4", "COL2": "5", "COL3": "6"},
	}

	// define test cases
	tests := []struct {
		name      string // Name of test
		csvString string // the content of our csv file
		seperator string // the seperator used for test case
	}{
		{"Comma separator", "COL1,COL2,COL3\n1,2,3\n4,5,6\n", "comma"},
		{"Semicolon separator", "COL1;COL2;COL3\n1;2;3\n4;5;6\n", "semicolon"},
	}

	// iterate over test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create temp csv file
			tempfile, err := ioutil.TempFile("", "test*.csv")
			check(err)

			// remove test csv file
			defer os.Remove(tempfile.Name())

			// write content of csv file
			_, err = tempfile.WriteString(test.csvString)
			// persisting data on disk
			tempfile.Sync()

			// define inputfile struct
			testFileData := inputFile{
				filepath:  tempfile.Name(),
				pretty:    false,
				seperator: test.seperator,
			}

			// define writer channel
			writerCh := make(chan map[string]string)

			// call targeted fiole as goroutinme
			go processCsvFile(testFileData, writerCh)

			// iterate over slice containing the expected map values
			for _, wantMap := range wantMapSLice {
				record := <-writerCh
				if !reflect.DeepEqual(record, wantMap) {
					t.Errorf("processCsvFile() = %v, wnat %v", record, wantMap)
				}
			}
		})
	}
}

func Test_writeJSONFile(t *testing.T) {
	// Define map of data to be processed to json
	dataMap := []map[string]string{
		{"COL1": "1", "COL2": "2", "COL3": "3"},
		{"COL1": "4", "COL2": "5", "COL3": "6"},
	}

	// define test case
	tests := []struct {
		csvPath  string
		jsonPath string
		pretty   bool
		name     string
	}{
		{"compact.csv", "compact.json", false, "Compact JSON"},
		{"pretty.csv", "pretty.json", true, "Pretty JSON"},
	}

	// Iterate over the test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// creatae mock channel
			writerCh := make(chan map[string]string)
			done := make(chan bool)

			// Create go routines
			go func() {
				for _, record := range dataMap {
					writerCh <- record
				}
				close(writerCh)
			}()

			// Running targetted functions
			go writeJSONFile(test.csvPath, writerCh, done, test.pretty)

			// waiting for past function to end
			<-done

			// Getting text from JSON file created from previous function
			testOutputput, err := ioutil.ReadFile(test.jsonPath)

			if err != nil {
				t.Errorf("writeJSONFile(), output file got error: %v", err)
			}

			// Cleanup
			defer os.Remove(test.jsonPath)

			// getting the text from the json file with expected data
			wantOutput, err := ioutil.ReadFile(filepath.Join("testJsonFiles", test.jsonPath))

			check(err)

			// Making the assertion between our generated JSON file content and the expected JSON file content
			if (string(testOutputput)) != (string(wantOutput)) {
				t.Errorf("writeJsonFile() = %v, want %v", string(testOutputput), string(wantOutput))
			}
		})
	}
}
