package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/aduong/gopkggraph/pkg/pkgwalker"
)

func main() {
	if err := mainE(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func mainE() error {
	match := flag.String("match", ".*", "regex to accept packages")
	help := flag.Bool("help", false, "help")
	flag.Parse()
	if *help {
		printUsage(os.Stdout)
		flag.CommandLine.SetOutput(os.Stdout)
		flag.PrintDefaults()
		return nil
	}
	matchRE, err := regexp.Compile(*match)
	if err != nil {
		printUsage(os.Stderr)
		flag.PrintDefaults()
		return err
	}
	pkg := flag.Arg(0)
	if pkg == "" {
		printUsage(os.Stderr)
		return errors.New("not enough args")
	}
	w := &pkgwalker.BreadthFirstWalker{
		OnPackage: func(pkg, fromPkg string) pkgwalker.Next {
			if !matchRE.MatchString(pkg) {
				return pkgwalker.StopPkg
			}
			if fromPkg != "" {
				fmt.Printf("%s -> %s\n", fromPkg, pkg)
			}
			return pkgwalker.Continue
		},
		OnErr: func(err error) pkgwalker.Next {
			fmt.Fprint(os.Stderr, err)
			return pkgwalker.Continue
		},
	}
	w.Walk(pkg)
	return nil
}

func printUsage(writer io.Writer) {
	fmt.Fprint(writer, `Usage: gopkggraph PKG

PKG is a package path like github.com/aduong/gopkggraph/pkg/pkgwalker
`)
}
