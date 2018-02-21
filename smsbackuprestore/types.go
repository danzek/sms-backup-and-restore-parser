package smsbackuprestore

import (
	"encoding/xml"
	"strconv"
	"time"
	"fmt"
	"encoding/base64"
	"io/ioutil"
	"os"
)

type PhoneNumber 		string
type SMSMessageType 	int
type SMSStatus			int
type AndroidTS 			string
type BoolValue			int
type ReadStatus			int

type Messages struct {
	XMLName 			xml.Name 		`xml:"smses"`
	Count 				string 			`xml:"count,attr"`
	BackupSet			string			`xml:"backup_set,attr"`
	BackupDateString	string			`xml:"backup_date,attr"`
	SMS 				[]SMS			`xml:"sms"`
	MMS 				[]MMS			`xml:"mms"`
}

type SMS struct {
	XMLName 			xml.Name 		`xml:"sms"`
	Protocol			string			`xml:"protocol,attr"`
	Address				PhoneNumber		`xml:"address,string,attr"`
	Type				SMSMessageType	`xml:"type,string,attr"`
	Subject				string			`xml:"subject,attr"`
	Body				string			`xml:"body,attr"`
	ServiceCenter		PhoneNumber		`xml:"service_center,string,attr"`
	Status				SMSStatus		`xml:"status,string,attr"`
	Read				ReadStatus		`xml:"read,string,attr"`
	Date				AndroidTS 		`xml:"date,string,attr"`  // consider reading in as int
	Locked				BoolValue		`xml:"locked,string,attr"`
	DateSent			AndroidTS		`xml:"date_sent,string,attr"`
	ReadableDate		string			`xml:"readable_date,attr"`
	ContactName			string			`xml:"contact_name,attr"`
}

type MMS struct {
	XMLName 			xml.Name 		`xml:"mms"`
	TextOnly			int				`xml:"text_only,string,attr"`
	Read				int				`xml:"read,string,attr"`
	Date				AndroidTS		`xml:"date,string,attr"`  // consider reading in as int
	Locked				BoolValue		`xml:"locked,string,attr"`
	DateSent			AndroidTS		`xml:"date_sent,string,attr"`
	ReadableDate		string			`xml:"readable_date,attr"`
	ContactName			string			`xml:"contact_name,attr"`
	Seen				BoolValue		`xml:"seen,string,attr"`
	FromAddress			PhoneNumber		`xml:"from_address,string,attr"`
	Address				PhoneNumber		`xml:"address,string,attr"`
	MessageClassifier	string			`xml:"m_cls,attr"`
	MessageSize			string			`xml:"m_size,attr"`
	Parts				[]Part			`xml:"parts>part"`
	Addresses			[]Address		`xml:"addrs>addr"`
}

type Part struct {
	XMLName 			xml.Name 		`xml:"part"`
	ContentType			string			`xml:"ct,attr"`
	Name				string			`xml:"name,attr"`
	FileName			string			`xml:"fn,attr"`
	ContentDisplay		string			`xml:"cd,attr"`
	Text				string			`xml:"text,attr"`
	Base64Data			string			`xml:"data,attr"`
}

type Address struct {
	XMLName 			xml.Name 		`xml:"addr"`
	Address				PhoneNumber		`xml:"address,string,attr"`
}

func (smt SMSMessageType) String() string {
	// see http://synctech.com.au/fields-in-xml-backup-files/
	// type – 1 = Received, 2 = Sent, 3 = Draft, 4 = Outbox, 5 = Failed, 6 = Queued
	smsMsgType := []string{"Received", "Sent", "Draft", "Outbox", "Failed", "Queued"}
	if smt > 0 && smt < 7 {
		return smsMsgType[smt-1]
	}
	return strconv.Itoa(int(smt))  // ignoring error
}

func (ss SMSStatus) String() string {
	// see http://synctech.com.au/fields-in-xml-backup-files/
	// status – None = -1, Complete = 0, Pending = 32, Failed = 64
	switch ss {
	case -1:
		return "None"
	case 0:
		return "Complete"
	case 32:
		return "Pending"
	case 64:
		return "Failed"
	default:
		return ""
	}
}

func (rs ReadStatus) String() string {
	// see http://synctech.com.au/fields-in-xml-backup-files/
	// read – Read Message = 1, Unread Message = 0
	if rs == 0 {
		return "Unread"
	} else if rs == 1 {
		return "Read"
	}
	return ""
}

func (timestamp AndroidTS) String() string {
	i, err := strconv.ParseInt(string(timestamp), 10, 64)
	if err != nil {
		return string(timestamp)
	}
	t := time.Unix(i/1000, 0).UTC()
	return t.String()
}

func (bv BoolValue) String() string {
	if bv == 0 {
		return "False"
	} else if bv == 1 {
		return "True"
	}
	return strconv.Itoa(int(bv))
}

func (p PhoneNumber) String() string {
	return NormalizePhoneNumber(string(p))
}

func (p Part) ImageFileName(mmsIndex int, partIndex int) string {
	ext := GetFileExtensionFromContentType(p.ContentType)
	if ext == "jpeg" {
		ext = "jpg"
	}
	fileName := p.Name
	if fileName == "null" {
		fileName = p.FileName
	}
	return fmt.Sprintf("%s_%d-%d.%s", fileName, mmsIndex, partIndex, ext)
}

func (p Part) DecodeAndWriteImage(outputPath string) error {
	// decode base64 image string as byte slice
	decodedImageBytes, err := base64.StdEncoding.DecodeString(p.Base64Data)
	if err != nil {
		return fmt.Errorf("Error decoding base64 image %s: %q\n", outputPath, err)
	}

	// write decoded byte slice to file
	fileWriteErr := ioutil.WriteFile(outputPath, decodedImageBytes, os.ModePerm)
	if fileWriteErr != nil {
		return fmt.Errorf("Error writing image %s to file: %q\n", outputPath, fileWriteErr)
	}

	return nil
}

func (m *Messages) PrintMessageCountQC() {
	lengthSMS := len(m.SMS)
	lengthMMS := len(m.MMS)

	fmt.Println("\nXML File Validation / QC")
	fmt.Println("===============================================================")
	fmt.Printf("Message count reported by SMS Backup and Restore app: %s\n", m.Count)

	// convert reportedCount to int for later comparison/validation
	count, err := strconv.Atoi(m.Count)
	if err != nil {
		fmt.Printf("Error converting reported count to integer: %s", m.Count)
		count = 0
	}

	fmt.Printf("Actual # SMS messages identified: %d\n", lengthSMS)
	fmt.Printf("Actual # MMS messages identified: %d\n", lengthMMS)
	fmt.Printf("Total actual messages identified: %d ... ", lengthSMS + lengthMMS)
	if lengthSMS + lengthMMS == count {
		fmt.Print("OK\n")
	} else {
		fmt.Print("DISCREPANCY DETECTED\n")
	}
}
