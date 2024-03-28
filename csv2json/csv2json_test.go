package main

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

func Test_getFileData(t *testing.T) {
	// Declare test slice  
	tests := []struct {
		name 	string
		want 	inputFile 
		wantErr bool
		osArgs	[]string
	} {
		{"Default parameters", inputFile{"test.csv", "comma", false}, false, []string{"cmd", "test.csv"}},
		{"No parameters", inputFile{}, true, []string{"cmd"}},
		{"Semicolon enabled", inputFile{"test.csv", "semicolon", false}, false, []string{"cmd", "--separator=semicolon", "test.csv"}},
		{"Pretty enabled", inputFile{"test.csv", "comma", true}, false, []string{"cmd", "--pretty", "test.csv"}},
		{"Pretty and semicolon enabled", inputFile{"test.csv", "semicolon", true}, false, []string{"cmd", "--pretty", "--separator=semicolon", "test.csv"}},
		{"Separator not identified", inputFile{}, true, []string{"cmd", "--separator=pipe", "test.csv"}},
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
			if (err !=nil) != test.wantErr {
				t.Errorf("getFileData() error = %v, wantErr %v", err, test.wantErr)
			}

			// assert whether or not we get correct wanted value
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("getFileData() = %v, want %v", got, test.want)
			}

		})
	}
}