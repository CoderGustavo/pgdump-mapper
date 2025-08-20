package models

type Options struct {
	Help       bool
	Json       bool
	JsonPretty bool
	Yaml       bool
	Html       bool
	Sqlite     bool
	Cli        bool
	Cache      bool
}

type FilterOptions struct {
	Schema    string
	TableName string
	Columns   []string
}

var CatalogOptions = map[string]bool{
	"-h":            false,
	"--help":        false,
	"--json":        false,
	"--json-pretty": false,
	"--yaml":        false,
	"--html":        false,
	"--sqlite":      false,
	"--cli":         false,
	"--cache":       false,
	"--schema":      false,
	"--table":       false,
	"--columns":     false,
}

type Table struct {
	Name       string              `json:"name"`
	Schema     string              `json:"schema"`
	Data       []map[string]string `json:"data"`
	Columns    []string            `json:"columns"`
	Values     [][]string          `json:"values"`
	PrimaryKey string              `json:"primary_key"`
	ForeignKey []map[string]string `json:"foreign_key"`
}
