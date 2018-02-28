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
	"strings"
	"unicode"
	"regexp"
)

// ReplaceAllBytesSubmatchFunc replaces all bytes in the byte slice that match the specified pattern.
//
// This is being done in an attempt to render emoji's properly due to SMS Backup & Restore app rendering of emoji's as
// HTML entitites in decimal (SLOW).
//
// Function is based on http://elliot.land/post/go-replace-string-with-regular-expression-callback
func ReplaceAllBytesSubmatchFunc(re *regexp.Regexp, b []byte, repl func([][]byte) []byte) []byte {
	var result []byte
	lastIndex := 0

	for _, v := range re.FindAllSubmatchIndex(b, -1) {
		var groups [][]byte
		for i := 0; i < len(v); i += 2 {
			groups = append(groups, b[v[i]:v[i+1]])
		}

		result = append(result, b[lastIndex:v[0]]...)
		result = append(result, repl(groups)...)
		lastIndex = v[1]
	}
	result = append(result, b[lastIndex:]...)
	return result
}

// NormalizePhoneNumber attempts to normalize phone numbers in the format 13125551212, ignoring input with multiple
// numbers delimited by a tilde ('~') character.
func NormalizePhoneNumber(number string) string {
	if strings.Contains(number, "~") {
		// don't parse when multiple numbers are provided (fail-safe)
		return number
	}

	// remove unwanted characters (for now not stripping all non-numeric chars because of odd cases)
	number = strings.TrimSpace(number)
	number = strings.Replace(number, "-", "", -1)
	number = strings.Replace(number, "(", "", -1)
	number = strings.Replace(number, ")", "", -1)
	number = strings.Replace(number, "+", "", -1)
	number = strings.Replace(number, " ", "", -1)

	// try to ensure all numbers have format 13125551212
	if len(number) == 11 {
		if strings.HasPrefix(number, "1") {
			return number
		}
	} else {
		if len(number) == 10 && !strings.HasPrefix(number, "1") {
			number = "1" + number
			return number
		}
	}
	return number
}

// RemoveCommasBeforeSuffixes recursively strips commas before suffixes such as M.D. to prevent contact names from
// being split by a comma in the middle of a name and suffix.
func RemoveCommasBeforeSuffixes(contacts string) string {
	// recursively strip commas before suffixes such as MD (doctors) to keep as single contact name
	// obviously this list is insufficient -- must be tailored to individual data sets
	// (this function could be its own project)
	hit := false
	suffixMap := map[string][]string{
		"MD": 		{", MD", ",MD", ", M.D", ",M.D"},
		"DO": 		{", DO", ",DO", ", D.O", ",D.O"},
		"NP": 		{", NP", ",NP", ", N.P", ",N.P"},
		"RN": 		{", RN", ",RN", ", R.N", ",R.N"},
		"JR": 		{", JR", ",JR", ", J.R", ",J.R"},
		"SR": 		{", SR", ",SR", ", S.R", ",S.R"},
		"II": 		{", II", ",II"},
		"III": 		{", III", ",III"},
		"INC": 		{", INC", ",INC"},
		"LLP": 		{", LLP", ",LLP", ", L.L.P", ",L.L.P"},
		"LLC": 		{", LLC", ",LLC", ", L.L.C", ",L.L.C"},
		"LPN": 		{", LPN", ",LPN", ", L.P.N", ",L.P.N"},
		"ACSW": 	{", ACSW", ",ACSW", ", A.C.S.W", ",A.C.S.W"},
		"LCSW": 	{", LCSW", ",LCSW", ", L.C.S.W", ",L.C.S.W"},
		"MA":		{", MA", ",MA", ", M.A", ",M.A"},
		"PHD": 		{", PHD", ",PHD", ", PH.D", ",PH.D", ", P.H.D", ",P.H.D"},
	}

	for s, combos := range suffixMap {
		if strings.Contains(strings.Replace(strings.ToUpper(contacts), ".", "", -1), s) {
			for _, suffix := range combos {
				searchHitIndex := strings.Index(strings.ToUpper(contacts), suffix)
				if searchHitIndex >= 0 {
					hit = true

					// special logic for DO/MA false positives
					if s == "DO" || s == "MA" {
						if len(contacts) > searchHitIndex+len(suffix) {  // not at end of string
							if unicode.IsLetter(rune(contacts[searchHitIndex+len(suffix)])) {
								hit = false  // false positive
							}
						}
					}

					// remove comma before suffix
					if hit {
						contacts = contacts[:searchHitIndex] + contacts[searchHitIndex+1:]
					}
				}
			}
		}
	}

	if hit {
		// call recursively to find any additional hits
		contacts = RemoveCommasBeforeSuffixes(contacts)
	}
	return contacts
}

// GetFileExtensionFromContentType determines the file extension of the base64-encoded file based on the content type.
func GetFileExtensionFromContentType(contentType string) string {
	// content type is like "image/png", so this extracts "png" in this case
	ext := contentType
	si := strings.Index(contentType, "/")
	if si >= 0 {
		ext = contentType[si+1:]
	}
	return ext
}

// CleanupMessageBody removes newlines and tabs from strings.
func CleanupMessageBody(body string) string {
	// strip unwanted characters from SMS body
	body = strings.Replace(body, "\n", " ", -1)
	body = strings.Replace(body, "\r", " ", -1)
	body = strings.Replace(body, "\t", " ", -1)
	return body
}
