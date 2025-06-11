package models

type Options struct {
	Help   bool
	Json   bool
	Yaml   bool
	Html   bool
	Sqlite bool
}

var CatalogOptions = map[string]bool{
	"-h":       false,
	"--help":   false,
	"--json":   false,
	"--yaml":   false,
	"--html":   false,
	"--sqlite": false,
}
