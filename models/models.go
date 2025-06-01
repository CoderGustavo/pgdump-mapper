package models

type Options struct {
	Help bool
	Raw  bool
}

var DefaultOptions = map[string]bool{
	"--help": false,
	"-h":     false,
	"--raw":  true,
	"-r":     true,
}
