package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/xml"
	"github.com/danzek/sms-backup-and-restore-parser/smsbackuprestore"
)

// main function for command-line SMS Backup & Restore app XML output parser
func main() {
	var xmlFilePath string

	// ensure required arg passed and file is valid (file path to xml file with sms backup and restore output)
	if len(os.Args) > 1 {
		xmlFilePath = os.Args[1]

		fileInfo, err := os.Stat(xmlFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error with path to XML file: %q\n", err)
			return
		} else if fileInfo.IsDir() {
			fmt.Fprint(os.Stderr, "XML path must point to specific XML filename, not to a directory.\n")
			return
		}
	} else {
		fmt.Fprint(os.Stderr, "Missing required argument: Specify path to xml backup file.\n" +
			"Example: sbrparser.exe C:\\Users\\4n68r\\Documents\\sms-20180213135542.xml\n")
		return
	}

	// open xml file
	f, err := os.Open(xmlFilePath)
	if err != nil {
		fmt.Fprint(os.Stderr, "Error opening XML file\n")
		panic(err)
	}
	defer f.Close()

	fmt.Printf("Parsing %s ...\n", xmlFilePath)

	// read entire file into data variable
	data, fileReadErr := ioutil.ReadFile(xmlFilePath)
	if fileReadErr != nil {
		panic(fileReadErr)
	}

	// instantiate messages object
	m := new(smsbackuprestore.Messages)
	if err := xml.Unmarshal(data, m); err != nil {
		panic(err)
	}

	// print validation / qc / stats to stdout
	m.PrintMessageCountQC()

	// generate sms
	fmt.Println("\nCreating SMS output...")
	err = smsbackuprestore.GenerateSMSOutput(m)
	if err != nil {
		fmt.Printf("Error encountered:\n%q\n", err)
	} else {
		fmt.Println("Finished generating SMS output")
		fmt.Println("sms.tsv file contains tab-separated values (TSV), i.e. use tab character as the delimiter")
	}

	// decode and output mms images
	fmt.Println("\nCreating images output...")
	numImagesIdentified, numImagesSuccessfullyWritten, imgOutputErrors := smsbackuprestore.DecodeImages(m)
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
	mmsOutputErr := smsbackuprestore.GenerateMMSOutput(m)
	if mmsOutputErr != nil {
		fmt.Printf("Error encountered:\n%q\n", mmsOutputErr)
	} else {
		fmt.Println("Finished generating MMS output")
		fmt.Println("mms.tsv file contains tab-separated values (TSV), i.e. use tab character as the delimiter")
	}
}
