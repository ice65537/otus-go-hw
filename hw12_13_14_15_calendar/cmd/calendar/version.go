package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	release   = "12.1"
	buildDate = "2023-06-12"
	gitHash   = "6a8c53d1692a79da3bb2f23ecc6bbb91e695996f"
)

func printVersion() {
	if err := json.NewEncoder(os.Stdout).Encode(struct {
		Release   string
		BuildDate string
		GitHash   string
	}{
		Release:   release,
		BuildDate: buildDate,
		GitHash:   gitHash,
	}); err != nil {
		fmt.Printf("error while decode version info: %v\n", err)
	}
}
