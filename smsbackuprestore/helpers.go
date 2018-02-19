package smsbackuprestore

import (
	"strings"
	"unicode"
)

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

func RemoveCommasBeforeSuffixes(contacts string) string {
	// recursively strip commas before suffixes such as MD (doctors) to keep as single contact name
	hit := false
	suffixMap := map[string][]string{
		"MD": {", MD", ",MD", ", M.D", ",M.D"},
		"DO": {", DO", ",DO", ", D.O", ",D.O"},
		"NP": {", NP", ",NP", ", N.P", ",N.P"},
		"RN": {", RN", ",RN", ", R.N", ",R.N"},
		"JR": {", JR", ",JR", ", J.R", ",J.R"},
		"SR": {", SR", ",SR", ", S.R", ",S.R"},
		"II": {", II", ",II"},
		"III": {", III", ",III"},
		"INC": {", INC", ",INC"},
		"LLP": {", LLP", ",LLP", ", L.L.P", ",L.L.P"},
		"LLC": {", LLC", ",LLC", ", L.L.C", ",L.L.C"},
		"ACSW": {", ACSW", ",ACSW", ", A.C.S.W", ",A.C.S.W"},
		"LCSW": {", LCSW", ",LCSW", ", L.C.S.W", ",L.C.S.W"},
		"PHD": {", PHD", ",PHD", ", PH.D", ",PH.D", "P.H.D", ",P.H.D"},
	}

	for s, combos := range suffixMap {
		if strings.Contains(strings.Replace(strings.ToUpper(contacts), ".", "", -1), s) {
			for _, suffix := range combos {
				searchHitIndex := strings.Index(strings.ToUpper(contacts), suffix)
				if searchHitIndex >= 0 {
					hit = true

					// special logic for DO false positives
					if s == "DO" {
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

func GetFileExtensionFromContentType(contentType string) string {
	// content type is like "image/png", so this extracts "png" in this case
	ext := contentType
	si := strings.Index(contentType, "/")
	if si >= 0 {
		ext = contentType[si+1:]
	}
	return ext
}

func CleanupMessageBody(body string) string {
	// strip unwanted characters from SMS body
	body = strings.Replace(body, "\n", " ", -1)
	body = strings.Replace(body, "\r", " ", -1)
	body = strings.Replace(body, "\t", " ", -1)
	return body
}
