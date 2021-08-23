package main

type VaultData struct {
	Name  string
	Value string
}

type VaultRecord struct {
	ID      int64
	Name    string
	Aliases []string
	Data    []VaultData
	Note    string
}
