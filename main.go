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
		printUsage(os.Stderr)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func mainE() error {
	match := flag.String("match", ".*", "regex to accept packages")
	format := flag.String("format", "pairs", "format to print graph: list, pairs (default), treelike")
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
		return err
	}
	if flag.NArg() == 0 {
		return errors.New("no packages specified")
	}
	var printer func(pkg, fromPkg string)
	switch *format {
	case "pairs":
		printer = pairsPrinter(" -> ")
	case "treelike":
		printer = treeLikePrinter()
	case "list":
		printer = listPrinter()
	default:
		return fmt.Errorf("unsupported format: %s", *format)
	}
	w := &pkgwalker.BreadthFirstWalker{
		OnPackage: compose(withMatcher(matchRE), withPrinter(printer))(nop),
		OnErr: func(err error) pkgwalker.Next {
			fmt.Fprintln(os.Stderr, err)
			return pkgwalker.Continue
		},
	}
	for _, pkg := range flag.Args() {
		w.Walk(pkg)

	}
	return nil
}

func listPrinter() func(pkg, fromPkg string) {
	return func(pkg, _ string) {
		fmt.Println(pkg)
	}
}

// treeLikePrinter prints in a format similiar to the tool `tree`.
// NB: relies on the fact that the walk is breadth-first.
func treeLikePrinter() func(pkg, fromPkg string) {
	lastFromPkg := ""
	var pkgs []string
	return func(pkg, fromPkg string) {
		if fromPkg != lastFromPkg {
			if lastFromPkg != "" {
				fmt.Println(lastFromPkg)
			}
			for i, pkg := range pkgs {
				if i != len(pkgs)-1 {
					fmt.Printf("├─ %s\n", pkg)
				} else {
					fmt.Printf("└─ %s\n", pkg)
				}
			}
			lastFromPkg = fromPkg
			pkgs = pkgs[:0]
		}
		pkgs = append(pkgs, pkg)
	}
}

func pairsPrinter(sep string) func(pkg, fromPkg string) {
	return func(pkg, fromPkg string) {
		fmt.Printf("%s%s%s\n", fromPkg, sep, pkg)
	}
}

type OnPackageTransform func(pkgwalker.OnPackageFunc) pkgwalker.OnPackageFunc

func withMatcher(matchRE *regexp.Regexp) OnPackageTransform {
	return func(f pkgwalker.OnPackageFunc) pkgwalker.OnPackageFunc {
		return func(pkg, fromPkg string) pkgwalker.Next {
			if !matchRE.MatchString(pkg) {
				return pkgwalker.StopPkg
			}
			return f(pkg, fromPkg)
		}
	}
}

func withPrinter(print func(pkg, fromPkg string)) OnPackageTransform {
	return func(f pkgwalker.OnPackageFunc) pkgwalker.OnPackageFunc {
		return func(pkg, fromPkg string) pkgwalker.Next {
			if fromPkg == "" {
				return pkgwalker.Continue
			}
			print(pkg, fromPkg)
			return f(pkg, fromPkg)
		}
	}
}

// nop is the do nothing OnPackage function
func nop(_, _ string) pkgwalker.Next {
	return pkgwalker.Continue
}

func compose(ts ...OnPackageTransform) OnPackageTransform {
	return func(f pkgwalker.OnPackageFunc) pkgwalker.OnPackageFunc {
		for i := len(ts) - 1; i >= 0; i-- {
			f = ts[i](f)
		}
		return f
	}
}

func printUsage(writer io.Writer) {
	fmt.Fprint(writer, `Usage: gopkggraph PKG...

PKG is a package path like github.com/aduong/gopkggraph/pkg/pkgwalker
`)
}
