# SMS Backup & Restore Parser

This tool parses the XML output from the [SMS Backup & Restore Android app](https://play.google.com/store/apps/details?id=com.riteshsahu.SMSBackupRestore).

## Usage

This tool assumes the default file naming convention of the [SMS Backup & Restore Android app](https://play.google.com/store/apps/details?id=com.riteshsahu.SMSBackupRestore), e.g.,

    calls-20180101000000.xml
    sms-20180101000000.xml

Simply pass the file name(s) of the XML backup file(s) you wish to parse and the tool will correctly identify the type of backup based on the file name. The parser can be ran with one or both files as parameters and will output data to the directory where the tool is located by default. Below are examples of running the compiled application on *nix and Windows systems, respectively:

    ./sbrparser calls-20180101000000.xml
    sbrparser.exe calls-20180101000000.xml sms-20180101000000.xml

To change the output directory (the default is the location of the program itself), use the `-d` parameter *before* passing the XML backup filename(s). For example, to output to the current working directory (`pwd`) on *nix:

    ./sbrparser -d . calls-20180101000000.xml sms-20180101000000.xml

And an example on Windows, directing output to the Desktop:

    sbrparser.exe -d C:\Users\4n68r\Desktop calls-20180101000000.xml sms-20180101000000.xml

## Expected Outputs

For the **calls backup file**, expected output is:

 - `calls.tsv` &mdash; tab-separated parsed calls data.


For the **SMS backup file**, expected outputs are:

 - `sms.tsv` &mdash; tab-separated parsed SMS data.
 - `mms.tsv` &mdash; tab-separated parsed MMS data.
 - `images/` &mdash; directory containing decoded images from MMS messages, saved with original file name plus MMS and Part indices to ensure a unique file name. File name format:

       <original file name>_<MMS Message Index>-<MMS Message Part Index>.<File Extension>

   A column named "`Part Output Image Name`" in the MMS output contains the precise file name of the outputted image.

## Existing Parsers
The SMS Backup & Restore Android app is currently maintained by [SyncTech](http://synctech.com.au/), and they offer both [paid and free versions](http://synctech.com.au/sms-backup-restore/) of the app as well as [an online parser](http://synctech.com.au/view-or-edit-sms-call-log-files-on-computer/). They also have [some documentation for the XML format used by the app on their website](http://synctech.com.au/fields-in-xml-backup-files/). In addition, [they documented various tools and methods for parsing the data.](http://synctech.com.au/view-or-edit-backup-files-on-computer/)

[Matt Joseph maintains an online tool for parsing these backups](https://mattj.io/sms-backup-reader/) ([GitHub project](https://github.com/devadvance/sms-backup-reader-2)). He also wrote [a legacy Java application for parsing the backups](https://mattj.io/sms-backup-reader/) ([GitHub project](https://github.com/devadvance/smsbackupreader)).

## With all of these existing parsers, why did you write your own?

**Stability:** I encountered backup XML files that were several gigabytes in size. These large backup files caused some existing solutions to fail. The backup format base64-encodes images from MMS messages within the XML file itself. I found that lines of data in these files exceeded the maximum buffer size for line-reader-based solutions and either caused truncation or complete failure in several solutions (including Excel).

**Performance:** Some of the existing solutions hang for a long time without responding and/or take a really long time to finish (if they finished at all). I wanted something simpler, faster, and more reliable.

**Privacy:** Both SyncTech and Matt Joseph's recommended solutions involve uploading the backup file(s) to a website. While they both claim that the data remain local and are processed locally on your computer (and I have no reason to doubt these claims), I simply don't trust web-based solutions in this regard as a future maintainer could silently change this policy/behavior.

**Flexibility:** I'm not a fan of having to view the data in a web browser (I want to be able to work with the data *anywhere* I want to). I also want to limit any dependencies (such as configuring a web server and/or installing a runtime environment). I value making the data available in a standard delimited format so that I can manipulate/filter records on the command line or in a spreadsheet application (such as Microsoft Excel).

## License & Disclaimer

This project has been released open source under the **MIT License**.

Copyright &copy; 2018 Dan O'Day (d@4n68r.com)

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

## Installation

1. [Download and install Go](https://go.dev/doc/install)
2. [Ensure your `GOPATH` is correctly configured](https://go.dev/doc/gopath_code)
3. Use `go install` to download, build, and install SMS Backup & Restore Parser (`sbrparser`) into [your `GOPATH`](https://go.dev/doc/gopath_code#GOPATH)



    go install 'github.com/danzek/sms-backup-and-restore-parser/cmd/sbrparser@latest'


This README is *not* a tutorial for how to use Go to install software. I've added these instructions as a courtesy for those unfamiliar with golang. If you need further assistance, [consult other resources for learning golang.](https://go.dev/doc/tutorial/getting-started)
