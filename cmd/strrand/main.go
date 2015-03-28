package main

import (
	"fmt"
	"os"

	"github.com/Songmu/strrand"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	ptn := os.Args[1]

	str, err := strrand.RandomString(ptn)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(str)
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `Usage:
    $ strrand '[1-3]{2}random[!?]'
    18random!`)
}
