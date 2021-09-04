## (Generic | Global) Password Manager

[![Build Status](https://github.com/vkuznet/gpm/actions/workflows/go.yml/badge.svg)](https://github.com/vkuznet/gpm/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/vkuznet/gpm)](https://goreportcard.com/report/github.com/vkuznet/gpm)
[![GoDoc](https://godoc.org/github.com/vkuznet/gpm?status.svg)](https://godoc.org/github.com/vkuznet/gpm)

(Generic | Global) Password Manager (GPM) is a password manager similar
to 1Password, LastPass, ButterCup, and similar password management tools.
Since majority of these tools are designed to work with your browsers we
wanted to create a more flexible version of password manager tool which
will be suitable in different environments, i.e. work in a terminal (CLI
version), work as a HTTP service, etc.

Here is a list of requirements:
- OS and architecture agnostic
- work in different environment, as a CLI tool, as a service, etc.
- support multiple cipher's implementation (currently supports AES and NaCl)
- support flexible data formats, e.g. pre-defined Login/Password records,
  or notes, or entire files
- support multiple vaults
- easy vault sync management, e.g. on local FS, on multiple cloud platforms
- easy customization (themes, record attributes, etc.)

This work is in progress and can only be viewed as alpha release.

Here is a preview of current functionality:
![Preview](doc/images/gpm.gif)

### Current functionality
So far, the tool is capable of creation vaults with different ciphers (AES and
NaCl are supported), the different records as well as entire files can be
added to the vault. It provides search capabilities, record editing, etc.
Since vault resides in specific directory and records resides in
individual encrypted files, the sync procedure with any destination is very
simple and can be organized via `rsync` tool.

Here is few examples:
```
# get help
./gpm -help
Usage of ./gpm:
  -cipher string
    	cipher to use (aes, nacl)
  -decrypt string
    	decrypt given file name
  -encrypt string
    	encrypt given file and place it into vault
  -vault string
    	vault name
  -verbose int
    	verbose level
  -version
    	Show version

# start vault, by default it will be created in $HOME/.gpm/Primary
# you may define your own location via GPM_HOME
./gpm

# encrypt given file and store it into the vault
./gpm -encrypt myfile.txt

# vault records
ls ~/.gpm/Primary

# decrypt given vault record
./gpm -decrypt ~/.gpm/Primary/2dface67-e5a8-44f7-ad58-adfa0f54b954.aes
```
Here is a typical structure of `gpm` vault(s):
```
tree ~/.gpm
/Users/vk/.gpm
├── MyNacl
│   ├── a589a35c-4d65-46e6-b18b-590a0623ded4.naci
│   └── backups
│       └── a589a35c-4d65-46e6-b18b-590a0623ded4.naci-2021-09-01T07:41:05-04:00
├── Primary
│   ├── 6b346bbd-a8a5-4af8-b9c6-c308c42bcd86.aes
│   ├── 71488552-1023-4480-9aa4-a909b23726ee.aes
│   ├── 9636520f-63ad-478c-92f7-3ed3b4eb579f.aes
│   ├── acb8a9f7-6140-42d2-bb32-f730f7ab572f.aes
│   ├── backups
│   │   ├── 6b346bbd-a8a5-4af8-b9c6-c308c42bcd86.aes-2021-08-30T18:02:56-04:00
│   │   ├── 6b346bbd-a8a5-4af8-b9c6-c308c42bcd86.aes-2021-08-30T18:03:10-04:00
│   └── fb26fd73-ea17-49f5-b38b-cf17575f1264.aes
├── config.json
├── gpm.log-20210831
....
```
It consists of two vaults (MyNacl) and (Primary) which holds different
encrypted records and their possible backups, the config.json, the logs.

So far the following keys are assigned to manage the vault records:
- `Ctrl-N` next widget
- `Ctrl-P` previous widget
- `Ctrl-F` switch to Search/Find input
- `Ctrl-L` switch to Records widget
- `Ctrl-E` record edit mode widget
- `Ctrl-G` generate password
- `Ctrl-P` copy password to clipboard
- `Ctrl-Q` Exit the app

### GPM Server
We add ability to start GPM server. It can be done as simple as following
```

# creeate server config file, e.g. server_config.json
# adjust your vault area to where your actual vault is
{
    "port": 8888,
    "vault_area": "/path/.gpm",
    "verbose": 1
}

# start the server
./gpm -server server_config.json
```

The GPM server support the following list of APIs
- GET URL/Vault provides list of records
- GET URL/Vault/recordID provides encrypted data record
- DELETE URL/Vault/recordID delete data record
- POST URL/Vault -d payload, upload record to the server
- GET URL/Vault/token provides token to use in API requests

For example,
```
# to get records from the vault named Primary
curl http;//localhost:8888/vault/Primary

# to get specific record content from vault Primary:
curl http;//localhost:8888/vault/Primary/fb26fd73-ea17-49f5-b38b-cf17575f1264.aes

# to post record to the vault Primary:
curl -X POST -d@your_record.json http;//localhost:8888/vault/Primary

# to delete record from the vault Primary
curl -X DELETE http;//localhost:8888/vault/Primary/fb26fd73-ea17-49f5-b38b-cf17575f1264.aes 

```
