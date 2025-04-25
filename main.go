package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tfmodtest/azopaform/pkg"
	"github.com/tfmodtest/azopaform/pkg/shared"
)

func main() {
	singlePath := flag.String("path", "", "The path of policy definition file")
	dir := flag.String("dir", "", "The dir which contains policy definitions")
	packageName := flag.String("package", "main", "The package name for the generated Rego files")
	utilRegoFileName := flag.String("util-file-name", "util.rego", "The name of the util Rego file (cannot be set together with util-library-package-name)")
	utilLibraryPackageName := flag.String("util-library-package-name", "", "The name of the util library package (if set, util file won't be generated; cannot be set together with util-file-name)")
	flag.Parse()

	if *utilLibraryPackageName != "" && *utilRegoFileName != "util.rego" {
		_, _ = fmt.Fprintln(os.Stderr, "Cannot set both `util-file-name` and `util-library-package-name` flags simultaneously.")
		os.Exit(1)
	}

	options := shared.Options{
		PackageName:            *packageName,
		UtilRegoFileName:       *utilRegoFileName,
		UtilLibraryPackageName: *utilLibraryPackageName,
	}

	ctx := shared.NewContextWithOptions(options)

	if err := pkg.AzurePolicyToRego(*singlePath, *dir, ctx); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
