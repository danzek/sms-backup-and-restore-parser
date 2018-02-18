# SMS Backup & Restore Parser

This tool parses the XML output from the [SMS Backup & Restore Android app](https://play.google.com/store/apps/details?id=com.riteshsahu.SMSBackupRestore).

## Existing Parsers
The SMS Backup & Restore Android app is currently maintained by [SyncTech](http://synctech.com.au/), and they offer both [paid and free versions](http://synctech.com.au/sms-backup-restore/) of the app as well as [an online parser](http://synctech.com.au/view-or-edit-sms-call-log-files-on-computer/). They also have [some documentation for the XML format used by the app on their website](http://synctech.com.au/fields-in-xml-backup-files/). In addition, [they documented various tools and methods for parsing the data.](http://synctech.com.au/view-or-edit-backup-files-on-computer/)

[Matt Joseph maintains an online tool for parsing these backups](https://mattj.io/sms-backup-reader/) ([GitHub project](https://github.com/devadvance/sms-backup-reader-2)). He also wrote [a legacy Java application for parsing the backups](https://mattj.io/sms-backup-reader/) ([GitHub project](https://github.com/devadvance/smsbackupreader)).

## With all of these existing parsers, why did you write your own?

**Stability:** I encountered backup XML files that were several gigabytes in size. These caused existing solutions to fail. The backup base64-encodes images from MMS messages within the file itself. I found that lines of data in these files exceeded the maximum buffer size for line-reader-based solutions to reading the file contents and either caused truncation or complete failure in several solutions (including Excel). Also, it didn't make sense to me why parsing should take so long (my solution processed a 1 GB XML SMS/MMS backup in less then 15 seconds&mdash;which includes exporting out all of the images).

**Privacy:** Both SyncTech and Matt Joseph's recommended solutions involve uploading the backup file(s) to a website. While they both claim that the data remains local and is processed locally on your computer (and I have no reason to doubt these claims), I simply don't trust web-based solutions in this regard as a future maintainer could silently change this policy/behavior.

**Flexibility:** I'm not a fan of having to view the data in a web browser. I want to read the data in a standard delimited format so that I can manipulate/filter it on the command line or in a spreadsheet application such as Microsoft Excel.

## License & Disclaimer

This project has been released open source under the **MIT License**.

Copyright &copy; 2018 Dan O'Day (d@4n68r.com)

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
