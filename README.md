## Encrypted Content Management

[![Build Status](https://github.com/vkuznet/ecm/actions/workflows/go.yml/badge.svg)](https://github.com/vkuznet/ecm/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/vkuznet/ecm)](https://goreportcard.com/report/github.com/vkuznet/ecm)
[![GoDoc](https://godoc.org/github.com/vkuznet/ecm?status.svg)](https://godoc.org/github.com/vkuznet/ecm)

[Encrypted Content Management (ECM)](https://vkuznet.github.io/ecm/) is a generic toolkit for encrypting any kind of digital content (files, passwords, notes, etc.). It can be viewed as an alternative to password managers similar to 1Password, LastPass, ButterCup, etc. but it is not restricted to secure only the meta-data. Any content can be encrypted with ECM.

ECM fills the gap among different solutions and provides CLI, Web server, browser extension and UI based on the GoLang crypto library. It depends on [rclone](https://rclone.org/) library to provide sync capabilities to various cloud providers (such as Dropbox, PCloud, etc.) as well as syncing to specific hosts (via sftp protocol). All components of ECM toolkit are written in Go, and therefore are portable across different architectures and hardwares. The CLI interface allows it to work with ECM from a terminal, and is therefore portable across hosts. The web server (and wasm extension) provides a fully featured web server with 2fa authentication. Finally, thanks to [fyne.io](https://fyne.io/) library the UI provides a consistent interface among platform and mobile devices.

### Motivation
We started this project after many years of experience with the 1Password solution.
Originally, there were few missing pieces with 1Password such as:
- command line interface (provided in version 8)
- web server (there is no publicly available server so far)
- support of multiple cloud providers and usage of on/off-site premises

But more importantly, the 1Password has changed their license based approach to subscription model starting version 8, see their [reply](https://1password.community/discussion/133705/does-1password8-support-non-subscription-mode/p1?new=1). Even though it is profitable for that company we considered that over the time it is not a sustainable solution, e.g. the current pricing model of $5/month leads to $60/year and since such a manager is mandatory in our digital life it can lead to substantial expenses over the long run. Despite the monthly fee, 1Password does not provide the ability to use its own infrastructure, there is lack of support for different cloud storage providers and 1Password is closed source. All of these factors lead to seeking alternative solutions and the idea of implementing password manager without aforementioned limitation. The Go-language seems to be an excellent choice for such implementation because:
- it is open-source
- it has solid crypto library as part of its Standard Library
- there are numerous native GUI solutions, e.g.
(via [fyne.io](https://fyne.io/))
- Go-based code can be ported mobile platform,
e.g.  [gomobile](https://pkg.go.dev/golang.org/x/mobile/cmd/gomobile).
- it supports various architectures, such as AMD, ARM, Power8, and [WebAssembly](https://www.wikiwand.com/en/WebAssembly) which allows it to implement desired functionality in a single language and port it across multiple hardware platforms.

Therefore, after a few iterations the ECM toolkit was released and it is released here as an open-source solution.

### Architectore
Below you can find current architecture of ECM toolkit:
![ECM architecture](doc/images/ecm-architecture.png)

The ECM toolkit has the following list of requirements:
- OS and architecture agnostic
- support different environments, work as a CLI tool, provide HTTP service, etc.
- support multiple cipher's implementation (currently supports AES and NaCl)
- support flexible data formats, e.g. pre-defined Login/Password records,
  or notes, or event entire files
- support multiple vaults
- support multiple storage systems, e.g. local FS, various cloud platforms,
remote storage, storage on remote nodes, etc.
- easy customization including vaults location, transfer, synchronisations,
and record attributes, etc.

This work is in progress and can only be viewed as an alpha release.

### Current functionality
So far, the ECM toolkit works as a CLI and web service. It supports different ciphers (AES and NaCl are implemented). It allows you to add login records, as well as user-based meta-data, it also allows you to add any kind of file to the vault. It provides basic search capabilities based on regex matching, record editing, etc. Since the vault resides in a specific directory, and records stored in individual encrypted files, the sync procedure with any destination is very simple and can be organized via `rsync` tool.

### Implementations
- [crypt](crypt/README.md) library used by ECM
- [vault](vault/README.md) library used by ECM
- [cli](cli/README.md) interface for ECM
- [server](server/README.md) implementation of ECM
- [term](term/README.md) based implementation of ECM
- [ui](ui/README.md) implementation of ECM
- [wasm](wasm/README.md) implementation of ECM
