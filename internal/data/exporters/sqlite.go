package exporters

import (
	"fmt"
	"os"
	"strings"

	cli "github.com/hedibertosilva/pgdump-mapper/internal/cli"
	models "github.com/hedibertosilva/pgdump-mapper/models"
)

func SQLite(schema string, allTables []models.Table, dbFile *os.File, rootPath string, tmpSQLiteFile string) {
	for _, table := range allTables {
		tableName := table.Name
		columns := table.Columns
		var data []map[string]string
		if table.Data != nil {
			data = table.Data
		} else {
			data = []map[string]string{}
		}

		for _, row := range data {
			values := []string{}
			for _, column := range columns {
				if value, ok := row[column]; ok {
					if value == "" {
						values = append(values, "NULL")
					} else if value == "\\N" {
						values = append(values, "NULL")
					} else {
						values = append(values, fmt.Sprintf("'%s'", strings.ReplaceAll(value, "'", "''")))
					}
				} else {
					values = append(values, "NULL")
				}
			}

			insertStmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);\n",
				tableName,
				strings.Join(columns, ", "),
				strings.Join(values, ", "))

			_, err := dbFile.WriteString(insertStmt)
			if err != nil {
				cli.ReturnError(err)
			}

		}

	}

	fmt.Printf("%s/%s created!\n", rootPath, tmpSQLiteFile)
	fmt.Printf("Import it using: sqlite3 <db-name>.sqlite3 < <%s path>\n", tmpSQLiteFile)

	if dbFile != nil {
		dbFile.Close()
	}
}
