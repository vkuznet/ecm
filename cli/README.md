## Encrypted Content Manager

[![Build Status](https://github.com/vkuznet/ecm/actions/workflows/go.yml/badge.svg)](https://github.com/vkuznet/ecm/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/vkuznet/ecm)](https://goreportcard.com/report/github.com/vkuznet/ecm)
[![GoDoc](https://godoc.org/github.com/vkuznet/ecm?status.svg)](https://godoc.org/github.com/vkuznet/ecm)

Encrypted Content Manager (ECM) is a generic tool for encrypting any kind of
digital content (files). It can be viewed as alternative to password managers similar to
1Password, LastPass, ButterCup, etc. but it is not restricted to secure
meta-data. Any content can be encrypted with ECM.

Since majority of password management tools are designed to work with your
browsers and lack of CLI and web server support we decided to create a more
flexible version of password manager tool which will be suitable in different
environments, i.e. work in a terminal (CLI version), work as a HTTP service,
support multiple architectures, etc.

The current list of requirements is the following:
- OS and architecture agnostic
- support different environments, work as CLI tool, provide HTTP service, etc.
- support multiple cipher's implementation (currently supports AES and NaCl)
- support flexible data formats, e.g. pre-defined Login/Password records,
  or notes, or event entire files
- support multiple vaults
- support multiple storage systems, e.g. local FS, various cloud platforms,
remote storage, storage on remote nodes, etc.
- easy customization including themes, record attributes, etc.

This work is in progress and can only be viewed as pre-alpha release.

### Current functionality
This area contains CLI implementation of ECM toolkit.

Here is few examples:
```
# get help
./ecm -help
Usage of ./ecm:
  -cipher string
    	cipher to use (aes, nacl)
  -decrypt string
    	decrypt given file to stdout
  -edit string
    	edit record with given ID
  -encrypt string
    	encrypt given file and place it into vault
  -examples
    	show examples
  -export string
    	export vault records to given file (ECM JSON native format)
  -import string
    	import records from a given file. Support: CSV, JSON, or ecm.json (native format)
  -info
    	show vault info
  -lock int
    	lock interval in seconds (default 60)
  -pat string
    	search pattern in vault records
  -pcopy string
    	extract given attribute from the record and copy to clipboard
  -recreate
    	recreate vault and its records with new password/cipher
  -rid string
    	show record with given ID and copy its password to clipboard
  -vault string
    	vault name
  -verbose int
    	verbose level
  -version
    	show version

# show examples
./ecm -examples

# list all vault records
./ecm
Enter vault secret:

------------
ID:		    6b346bbd-a8a5-4af8-b9c6-c308c42bcd86
Name:		Record
Login:		test
Password:	************
URL:
Tags:
Note:

# show individual record
./ecm -rid fb26fd73-ea17-49f5-b38b-cf17575f1264

# edit individual record
./ecm -edit fb26fd73-ea17-49f5-b38b-cf17575f1264

# recreate (re-encrypt) vault
./ecm -recreate

# import 1Password records and export them to records.json (ECM JSON data-format)
# at this point you can edit records.json in your favorite editor
./ecm -import 1password.csv -export ./records.json

# import ECM JSON to the vault (ecm.json must be used and it
# should contain ECM JSON data-format)
./ecm -import ecm.json

# encrypt given file and store it into the vault
./ecm -encrypt myfile.txt

# show vault info
./ecm -info
Enter vault secret:
vault /Users/vk/.ecm/Primary
Last modified: 2022-05-22 10:23:15.381822738 -0400 EDT
Size 288 (288.0B), mode drwxr-xr-x
6 records, encrypted with aes cipher

# decrypt given vault record
./ecm -decrypt ~/.ecm/Primary/2dface67-e5a8-44f7-ad58-adfa0f54b954.aes
```
Here is a typical structure of `ecm` vault(s):
```
tree ~/.ecm
/Users/users/.ecm
├── Primary
│   ├── 6b346bbd-a8a5-4af8-b9c6-c308c42bcd86.aes
│   ├── 71488552-1023-4480-9aa4-a909b23726ee.aes
│   ├── 9636520f-63ad-478c-92f7-3ed3b4eb579f.aes
│   ├── acb8a9f7-6140-42d2-bb32-f730f7ab572f.aes
│   ├── backups
│   │   ├── 6b346bbd-a8a5-4af8-b9c6-c308c42bcd86.aes-2021-08-30T18:02:56-04:00
│   │   ├── 6b346bbd-a8a5-4af8-b9c6-c308c42bcd86.aes-2021-08-30T18:03:10-04:00
│   └── fb26fd73-ea17-49f5-b38b-cf17575f1264.aes
....
```
