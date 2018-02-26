/*
SBRParser: SMS Backup & Restore Android app parser

Copyright (c) 2018 Dan O'Day <d@4n68r.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
 */

// Package main for command-line SMS Backup & Restore parser.
// This tool parses SMS Backup & Restore Android app XML output.
package main

import (
	"fmt"
	"flag"
	"os"
	"io/ioutil"
	"encoding/xml"
	"github.com/danzek/sms-backup-and-restore-parser/smsbackuprestore"
	"time"
	"path/filepath"
	"strings"
)

// SMSOutput calls GenerateSMSOutput() and prints status/errors.
func SMSOutput(m *smsbackuprestore.Messages, outputDir string) {
	// generate sms
	fmt.Println("\nCreating SMS output...")
	err := smsbackuprestore.GenerateSMSOutput(m, outputDir)
	if err != nil {
		fmt.Printf("Error encountered:\n%q\n", err)
	} else {
		fmt.Println("Finished generating SMS output")
		fmt.Println("sms.tsv file contains tab-separated values (TSV), i.e. use tab character as the delimiter")
	}
}

// MMSOutput calls DecodeImages() and GenerateMMSOutput() and prints status/errors.
func MMSOutput(m *smsbackuprestore.Messages, outputDir string) {
	// decode and output mms images
	fmt.Println("\nCreating images output...")
	numImagesIdentified, numImagesSuccessfullyWritten, imgOutputErrors := smsbackuprestore.DecodeImages(m, outputDir)
	if imgOutputErrors != nil && len(imgOutputErrors) > 0 {
		for e := range imgOutputErrors {
			fmt.Printf("\t%q\n", e)
		}
	}
	fmt.Println("Finished decoding images")
	fmt.Printf("%d images were identified and %d were successfully written to file\n", numImagesIdentified, numImagesSuccessfullyWritten)
	fmt.Println("Image file names are in format: <original file name (if known)>_<mms index>-<sms index>.<file extension>")

	// generate mms output
	fmt.Println("\nCreating MMS output...")
	mmsOutputErr := smsbackuprestore.GenerateMMSOutput(m, outputDir)
	if mmsOutputErr != nil {
		fmt.Printf("Error encountered:\n%q\n", mmsOutputErr)
	} else {
		fmt.Println("Finished generating MMS output")
		fmt.Println("mms.tsv file contains tab-separated values (TSV), i.e. use tab character as the delimiter")
	}
}

// CallsOutput calls GenerateCallOutput() and prints status/errors.
func CallsOutput (c *smsbackuprestore.Calls, outputDir string) {
	// generate calls
	fmt.Println("\nCreating calls output...")
	err := smsbackuprestore.GenerateCallOutput(c, outputDir)
	if err != nil {
		fmt.Printf("Error encountered:\n%q\n", err)
	} else {
		fmt.Println("Finished generating calls output")
		fmt.Println("calls.tsv file contains tab-separated values (TSV), i.e. use tab character as the delimiter")
	}
}

// GetExecutablePath returns the absolute path to the location where this executable is being ran from
func GetExecutablePath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return ".", fmt.Errorf("Error: Try running this application from another location: %q\n", err)
	}

	exePath, err := filepath.Abs(filepath.Dir(exe))
	if err != nil {
		return ".", fmt.Errorf("Error: Try running this application from another location: %q\n", err)
	}

	return exePath, nil
}

// main function for command-line SMS Backup & Restore app XML output parser.
func main() {
	// time program execution
	start := time.Now()

	// get executable path
	exePath, err := GetExecutablePath()
	if err != nil {
		panic(err)
	}

	// parse command-line args/flags
	pOutputDirectory := flag.String("d", exePath, "Directory path for parsed output (current executable directory is default)")
	flag.Parse()

	// validate output directory
	if outputDirInfo, err := os.Stat(*pOutputDirectory); os.IsNotExist(err) || !outputDirInfo.IsDir() {
		fmt.Fprintf(os.Stderr, "Invalid output directory path: %s", *pOutputDirectory)
		return
	}
	fmt.Printf("Output directory set to %s\n", *pOutputDirectory)

	if len(flag.Args()) > 0 {
		for _, xmlFilePath := range flag.Args() {
			// ensure file is valid (file path to xml file with sms backup and restore output)
			fileInfo, err := os.Stat(xmlFilePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error with path to XML file: %q\n", err)
				return
			} else if fileInfo.IsDir() {
				fmt.Fprint(os.Stderr, "XML path must point to specific XML filename, not to a directory.\n")
				return
			}

			// get just file name and perform verification checks (assumes default lowercase naming convention)
			fileName := filepath.Base(xmlFilePath)
			if (strings.HasPrefix(fileName, "calls-") || strings.HasPrefix(fileName, "sms-")) && filepath.Ext(fileName) == ".xml" {
				// open xml file
				f, err := os.Open(xmlFilePath)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error opening XML file: %s\n", xmlFilePath)
					break
				}
				defer f.Close()

				// status message
				fmt.Printf("\nLoading %s into memory and parsing (this may take a little while) ...\n", xmlFilePath)

				// read entire file into data variable
				data, fileReadErr := ioutil.ReadFile(xmlFilePath)
				if fileReadErr != nil {
					panic(fileReadErr)
				}

				// determine file type
				if strings.HasPrefix(fileName, "sms-") {
					// sms backup
					// instantiate messages object
					m := new(smsbackuprestore.Messages)
					if err := xml.Unmarshal(data, m); err != nil {
						panic(err)
					}

					// print validation / qc / stats to stdout
					m.PrintMessageCountQC()

					// generate sms
					SMSOutput(m, *pOutputDirectory)

					// generate mms
					MMSOutput(m, *pOutputDirectory)
				} else {
					// calls backup
					// instantiate calls object
					c := new(smsbackuprestore.Calls)
					if err := xml.Unmarshal(data, c); err != nil {
						panic(err)
					}

					// print validation / qc / stats to stdout
					c.PrintCallCountQC()

					// generate calls output
					CallsOutput(c, *pOutputDirectory)
				}
			} else {
				fmt.Fprintf(os.Stderr, "Unexpected file name: %s\n", fileName)
			}
		}
	} else {
		fmt.Fprint(os.Stderr, "Missing required argument: Specify path to xml backup file(s).\n" +
			"Example: sbrparser.exe C:\\Users\\4n68r\\Documents\\sms-20180213135542.xml\n")  // todo -- use name of executable
		return
	}

	// print completion messages
	fmt.Printf("\nCompleted in %.2f seconds.\n", time.Since(start).Seconds())
	fmt.Printf("Output saved to %s\n", *pOutputDirectory)
}
