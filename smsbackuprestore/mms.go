package smsbackuprestore

import (
	"path/filepath"
	"os"
	"strings"
	"fmt"
	"strconv"
)

func DecodeImages(m *Messages) (numImagesIdentified, numImagesSuccessfullyWritten int, errors []error) {
	numImagesIdentified = 0
	numImagesSuccessfullyWritten = 0

	// create output directory for images
	outputDir := filepath.Join(".", "images")
	os.MkdirAll(outputDir, os.ModePerm)

	for mmsIndex, mms := range m.MMS {
		for partIndex, part := range mms.Parts {
			if strings.Contains(part.ContentType, "image/") {
				numImagesIdentified++
				outputImgFilename := part.ImageFileName(mmsIndex, partIndex)

				// decode base64 image string as byte slice and write decoded byte slice to file
				outputPath := filepath.Join(outputDir, outputImgFilename)
				err := part.DecodeAndWriteImage(outputPath)
				if err != nil {
					errors = append(errors, err)
				} else {
					numImagesSuccessfullyWritten++
				}
			}
		}
	}
	return numImagesIdentified, numImagesSuccessfullyWritten, errors
}

func GenerateMMSOutput(m *Messages) error {
	mmsOutput, err := os.Create("mms.tsv")
	if err != nil {
		return fmt.Errorf("Unable to create file: mms.tsv\n%q", err)
	}
	defer mmsOutput.Close()

	// print header row
	headers := []string{
		"MMS Index #",
		"MMS Part Index #",
		"Text Only",
		"Read",
		"Date",
		"Locked",
		"Date Sent",
		"Readable Date",
		"Contact Name",
		"Seen",
		"From Address",
		"Address",
		"Addresses",
		"Message Classifier",
		"Message Size",
		"Part Content Type",
		"Part Name",
		"Part File Name",
		"Part Text",
		"Part Content Display",
		"Part Output Image Name",
	}
	fmt.Fprintf(mmsOutput, "%s\n", strings.Join(headers, "\t"))

	// iterate over mms
	for mmsIndex, mms := range m.MMS {
		var names []string
		var numbers []string
		var addresses []string
		var contactNameList string
		var addressList string
		addressesList := ""

		groupMessage := false
		if strings.Contains(mms.ContactName, ",") || strings.Contains(mms.Address.String(), "~") {
			groupMessage = true

			// get names
			for _, name := range strings.Split(RemoveCommasBeforeSuffixes(mms.ContactName), ",") {
				names = append(names, strings.TrimSpace(name))
			}

			// get/normalize phone numbers
			for _, number := range strings.Split(mms.Address.String(), "~") {
				numbers = append(numbers, PhoneNumber(number).String())
			}
		}

		// create semicolon-delimited output for group messages
		if groupMessage {
			// semicolon-delimited fields
			contactNameList = strings.Join(names, ";")
			addressList = strings.Join(numbers, ";")
		} else {
			contactNameList = RemoveCommasBeforeSuffixes(mms.ContactName)
			addressList = mms.Address.String()
		}

		// get any addresses for group message
		for _, addr := range mms.Addresses {
			addresses = append(addresses, addr.Address.String())
		}
		if len(addresses) > 0 {
			addressesList = strings.Join(addresses, ";")
		}

		for partIndex, part := range mms.Parts {
			imageFile := "N/A"
			if strings.Contains(part.ContentType, "image/") {
				imageFile = part.ImageFileName(mmsIndex, partIndex)
			}

			row := []string{
				strconv.Itoa(mmsIndex),
				strconv.Itoa(partIndex),
				mms.TextOnly.String(),
				mms.Read.String(),
				mms.Date.String(),
				mms.Locked.String(),
				mms.DateSent.String(),
				mms.ReadableDate,
				contactNameList,
				mms.Seen.String(),
				mms.FromAddress.String(),
				addressList,
				addressesList,
				mms.MessageClassifier,
				mms.MessageSize,
				part.ContentType,
				part.Name,
				part.FileName,
				CleanupMessageBody(part.Text),
				part.ContentDisplay,
				imageFile,
			}
			fmt.Fprintf(mmsOutput, "%s\n", strings.Join(row, "\t"))
		}
	}
	return nil
}
