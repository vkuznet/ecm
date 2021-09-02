## (Generic|Global|GNU) Password Manager

[![Build Status](https://github.com/vkuznet/pwm/actions/workflows/go.yml/badge.svg)](https://github.com/vkuznet/pwm/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/vkuznet/pwm)](https://goreportcard.com/report/github.com/vkuznet/pwm)
[![GoDoc](https://godoc.org/github.com/vkuznet/pwm?status.svg)](https://godoc.org/github.com/vkuznet/pwm)

(Generic | Global | GNU) Password Manager (GPM) is a password manager similar
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

Here is a preview of current functionality:
![Preview](doc/images/gpm.gif)
