package exporters

import (
	"encoding/json"
	"fmt"

	cli "github.com/hedibertosilva/pgdump-mapper/internal/cli"
	models "github.com/hedibertosilva/pgdump-mapper/models"
)

func JSON(schema string, tables []models.Table, isJsonPretty bool) {
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

	var output []byte
	var err error

	if len(tablesToExport) > 0 {
		if isJsonPretty {
			output, err = json.MarshalIndent(tablesToExport, "", "  ")
		} else {
			output, err = json.Marshal(tablesToExport)
		}

		if err != nil {
			cli.ReturnError(err)
		}

		fmt.Println(string(output))
	}
}
