## (Generic|Global|GNU) Password Manager

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
