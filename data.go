package main

type VaultData struct {
	Name  string
	Value string
}

type VaultRecord struct {
	Name    string
	Aliases []string
	Data    []VaultData
	Note    string
}
