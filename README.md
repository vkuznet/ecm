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
So far, the ECM toolkit works as a CLI and web service. It supports different ciphers (AES and
NaCl are implemented). It allows to add login records, as well as user-based
meta-data, it also allow to add any kind of file to the vault.
It provides basic search capabilities based on regpex matching, record editing, etc.
Since vault resides in specific directory, and records stored in
individual encrypted files, the sync procedure with any destination is very
simple and can be organized via `rsync` tool.

### Implementations
- [crypt](crypt/README.md) library used by ECM
- [vault](vault/README.md) library used by ECM
- [cli](cli/README.md) interface for ECM
- [server](server/README.md) implementation of ECM
- [term](term/README.md) based implementation of ECM
- [ui](ui/README.md) implementation of ECM
