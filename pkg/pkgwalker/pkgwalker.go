package pkgwalker

type Next int

const (
	Continue Next = iota
	StopPkg
	StopAll
)
