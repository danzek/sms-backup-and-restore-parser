package smsbackuprestore

import (
	"fmt"
	"path/filepath"
	"os"
	"strings"
	"encoding/base64"
	"io/ioutil"
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
