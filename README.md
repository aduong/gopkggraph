# gopkggraph

Utility that prints the package graph of Go packages.

## Usage

    go run github.com/aduong/gopkggraph -help

Print the graph with defaults (as pairs, no filtering):

    $ go run github.com/aduong/gopkggraph github.com/aduong/gopkggraph
    github.com/aduong/gopkggraph -> errors
    github.com/aduong/gopkggraph -> flag
    github.com/aduong/gopkggraph -> fmt
    github.com/aduong/gopkggraph -> github.com/aduong/gopkggraph/pkg/pkgwalker
    github.com/aduong/gopkggraph -> io
    github.com/aduong/gopkggraph -> os
    github.com/aduong/gopkggraph -> regexp
    errors -> internal/reflectlite
    flag -> reflect
    flag -> sort
    flag -> strconv
    flag -> strings
    flag -> time
    ...

Print the graph in a tree-like format:

    $ go run github.com/aduong/gopkggraph -format treelike github.com/aduong/gopkggraph
    github.com/aduong/gopkggraph
    ├─ errors
    ├─ flag
    ├─ fmt
    ├─ github.com/aduong/gopkggraph/pkg/pkgwalker
    ├─ io
    ├─ os
    └─ regexp
    errors
    └─ internal/reflectlite
    ...

Print the graph as pairs filtering only for your package:

    $ go run github.com/aduong/gopkggraph -match github.com/aduong/gopkggraph github.com/aduong/gopkggraph
    github.com/aduong/gopkggraph -> github.com/aduong/gopkggraph/pkg/pkgwalker

## Acknowledgements

Inspired by Dave Cheney's [graphpkg](https://github.com/davecheney/graphpkg) which produces an SVG of package graphs.
