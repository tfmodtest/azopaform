package main

import (
	"flag"
	"fmt"
	"json-rule-finder/pkg"
	"json-rule-finder/pkg/shared"
	"os"
)

func main() {
	singlePath := flag.String("path", "", "The path of policy definition file")
	dir := flag.String("dir", "", "The dir which contains policy definitions")
	packageName := flag.String("package", "main", "The package name for the generated Rego files")
	utilRegoFileName := flag.String("util-file-name", "util.rego", "The name of the util Rego file")
	flag.Parse()

	options := pkg.Options{
		PackageName:      *packageName,
		UtilRegoFileName: *utilRegoFileName,
	}

	if err := pkg.AzurePolicyToRego(*singlePath, *dir, options, shared.NewContext()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		//os.Exit(1)
	}
}
