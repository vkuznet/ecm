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

Here is a preview of current functionality:
![Preview](doc/images/gpm.gif)

### Current functionality
So far, the ECM tool works as a CLI and web service. It supports different ciphers (AES and
NaCl are implemented). It allows to add login records, as well as user-based
meta-data, it also allow to add any kind of file to the vault.
It provides basic search capabilities based on regpex matching, record editing, etc.
Since vault resides in specific directory, and records stored in
individual encrypted files, the sync procedure with any destination is very
simple and can be organized via `rsync` tool.

Here is few examples:
```
# get help
./ecm -help
Usage of ./ecm:
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

# start vault, by default it will be created in $HOME/.ecm/Primary
# you may define your own location via ECM_HOME
./ecm

# encrypt given file and store it into the vault
./ecm -encrypt myfile.txt

# vault records
ls ~/.ecm/Primary

# decrypt given vault record
./ecm -decrypt ~/.ecm/Primary/2dface67-e5a8-44f7-ad58-adfa0f54b954.aes
```
Here is a typical structure of `ecm` vault(s):
```
tree ~/.ecm
/Users/users/.ecm
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
├── ecm.log-20210831
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

### ECM Server
We add ability to start ECM server. It can be done as simple as following
```

# creeate server config file, e.g. server_config.json
# adjust your vault area to where your actual vault is
{
    "port": 8888,
    "vault_area": "/path/.ecm",
    "verbose": 1
}

# start the server
./ecm -server server_config.json
```

The ECM server support the following list of APIs
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
