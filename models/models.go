package models

type Options struct {
	Help       bool
	JSON       bool
	JSONPretty bool
	YAML       bool
	HTML       bool
	SQLite     bool
	CLI        bool
	Cache      bool
}

type FilterOptions struct {
	Schema    string
	TableName string
	Columns   []string
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
