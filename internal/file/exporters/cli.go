package exporters

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/hedibertosilva/pgdump-mapper/internal/cli"
	"github.com/hedibertosilva/pgdump-mapper/models"
)

func Contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func CLI(schema string, AllTables []models.Table) {
	for _, table := range AllTables {
		tableName := table.Name

		// Filter table
		if cli.Filters.TableName != "" &&
			(table.Name != cli.Filters.TableName ||
				table.Schema != schema) {
			continue
		}

		columns := table.Columns
		var data []map[string]string
		if table.Data != nil {
			data = table.Data
		} else {
			data = []map[string]string{}
		}

		// Filter columns
		var selectedColumns []string
		if len(cli.Filters.Columns) > 0 {
			for _, col := range cli.Filters.Columns {
				if Contains(columns, col) {
					selectedColumns = append(selectedColumns, col)
				}
			}
		} else {
			selectedColumns = columns
		}

		fmt.Printf("\nTable: %s\n\n", schema+"."+tableName)
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		fmt.Fprintln(w, strings.Join(selectedColumns, "\t"))

		for _, row := range data {
			values := []string{}
			for _, col := range selectedColumns {
				values = append(values, row[col])
			}
			fmt.Fprintln(w, strings.Join(values, "\t"))
		}

		w.Flush()
	}
}
