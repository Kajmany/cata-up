package main

import (
	"fmt"
	"os"
)

func main() {
	var _ *Config
	_, err := GetConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
