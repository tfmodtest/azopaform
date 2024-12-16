package main

import (
	"flag"
	"fmt"
	"json-rule-finder/pkg"
	"os"
)

var noResourceTypeFound int

func main() {
	singlePath := flag.String("path", "", "The path of policy definition file")
	dir := flag.String("dir", "", "The dir which contains policy definitions")
	flag.Parse()
	pkg.NoResourceTypeFound = noResourceTypeFound
	if err := pkg.AzurePolicyToRego(*singlePath, *dir, pkg.NewContext()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		//os.Exit(1)
	}
}
