package models

type Options struct {
	Help bool
	Raw  bool
	Json bool
	Yaml bool
	Html bool
}

var DefaultOptions = map[string]bool{
	"--help": false,
	"-h":     false,
	"--raw":  true,
	"-r":     true,
	"--json": false,
	"--yaml": false,
	"--html": false,
}
