package pkgwalker

import (
	"fmt"
	"go/build"
)

type BreadthFirstWalker struct {
	OnPackage OnPackageFunc
	OnErr     func(err error) Next

	seen map[string]bool
	q    []pair
}

type OnPackageFunc func(pkg, fromPkg string) Next

type pair struct {
	pkg     string
	fromPkg string
}

func (w *BreadthFirstWalker) Walk(pkg string) {
	if w.seen == nil {
		w.seen = map[string]bool{}
	}
	w.q = []pair{{pkg: pkg}}
	w.walk()
}

func (w *BreadthFirstWalker) walk() {
	for {
		if len(w.q) == 0 {
			return
		}
		_pair := w.q[0]
		w.q = w.q[1:]
		pkg, fromPkg := _pair.pkg, _pair.fromPkg
		if pkg == "C" {
			continue
		}
		switch next := w.OnPackage(pkg, fromPkg); next {
		case Continue:
			if w.seen[pkg] {
				continue
			}
			w.seen[pkg] = true
			p, err := build.Import(pkg, "", 0)
			if err != nil {
				if w.OnErr != nil {
					if next := w.OnErr(err); next == StopAll {
						return
					}
				}
			} else {
				for _, toPkg := range p.Imports {
					w.q = append(w.q, pair{pkg: toPkg, fromPkg: pkg})
				}
			}
		case StopPkg:
		case StopAll:
			return
		default:
			if w.OnErr != nil {
				w.OnErr(fmt.Errorf("%d not a valid Next", next))
			}
			return
		}
	}
}
