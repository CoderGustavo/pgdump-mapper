package exporters

import (
	"fmt"

	"github.com/hedibertosilva/pgdump-mapper/internal/cli"
	"github.com/hedibertosilva/pgdump-mapper/models"
	"gopkg.in/yaml.v3"
)

func YAML(schema string, tables []models.Table) {
	var tablesToExport []models.Table
	if cli.Filters.TableName != "" {
		for _, table := range tables {
			if table.Name == cli.Filters.TableName &&
				table.Schema == schema {
				tablesToExport = append(tablesToExport, table)
			}
		}
	} else {
		tablesToExport = tables
	}

	out, err := yaml.Marshal(tablesToExport)
	if err != nil {
		cli.ReturnError(err)
	}
	fmt.Println(string(out))
}
