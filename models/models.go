package models

type Options struct {
	Help bool
	Raw  bool
}

var ValidOptions = map[string]bool{
	"--help": true,
	"-h":     true,
	"--raw":  true,
	"-r":     true,
}
