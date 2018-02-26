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

package smsbackuprestore

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"strconv"
)

// GenerateCallOutput outputs a tab-delimited file named "calls.tsv" containing parsed calls from the backup file.
func GenerateCallOutput(c *Calls, outputDir string) error {
	callOutput, err := os.Create(filepath.Join(outputDir, "calls.tsv"))
	if err != nil {
		return fmt.Errorf("Unable to create file: calls.tsv\n%q", err)
	}
	defer callOutput.Close()

	// print header row
	headers := []string{
		"Call Index #",
		"Number",
		"Duration (Seconds)",
		"Date",
		"Type",
		"Readable Date",
		"Contact Name",
	}
	fmt.Fprintf(callOutput, "%s\n", strings.Join(headers, "\t"))

	// iterate over calls
	for i, call := range c.Calls {
		row := []string{
			strconv.Itoa(i),
			call.Number.String(),
			strconv.Itoa(call.Duration),
			call.Date.String(),
			call.Type.String(),
			call.ReadableDate,
			RemoveCommasBeforeSuffixes(call.ContactName),
		}
		fmt.Fprintf(callOutput, "%s\n", strings.Join(row, "\t"))
	}

	return nil
}
