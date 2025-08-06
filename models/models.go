package models

type Options struct {
	Help   bool
	Json   bool
	Yaml   bool
	Html   bool
	Sqlite bool
	Cli bool
}

type FilterOptions struct {
	TableName string
	Columns   []string
}

var CatalogOptions = map[string]bool{
	"-h":       false,
	"--help":   false,
	"--json":   false,
	"--yaml":   false,
	"--html":   false,
	"--sqlite": false,
	"--cli": false,
	"--table": false,
	"--columns": false,
}
