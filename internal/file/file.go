package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/tabwriter"

	yaml "gopkg.in/yaml.v3"

	cli "github.com/hedibertosilva/pgdump-mapper/internal/cli"
	parser "github.com/hedibertosilva/pgdump-mapper/internal/file/parser"
	models "github.com/hedibertosilva/pgdump-mapper/models"
)

type Table struct {
	Name       string              `json:"name"`
	Schema     string              `json:"schema"`
	Data       []map[string]string `json:"data"`
	Columns    []string            `json:"columns"`
	Values     [][]string          `json:"values"`
	PrimaryKey string              `json:"primary_key"`
	ForeignKey []map[string]string `json:"foreign_key"`
}

var (
	Input         *string
	Options       models.Options
	AllTables     []Table
	DBFile        *os.File
	TmpSqliteFile string = "pgdump-mapper.sqlite.txt"
	TmpCacheDir   string = "/tmp/pgdump-mapper"
	TmpCacheFile  string = ""
	cwd, _               = os.Getwd()
)

func findTable(allTables []Table, targetTable Table) (*Table, bool) {
	for _, table := range allTables {
		if table.Name == targetTable.Name &&
			table.Schema == targetTable.Schema {
			return &table, true
		}
	}
	return nil, false
}

func contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func Read() {

	if Options.Cache {
		basename := filepath.Base(*Input)
		TmpCacheFile = filepath.Join(TmpCacheDir, basename)

		fileBytes, err := os.ReadFile(TmpCacheFile)
		if err != nil {
			// No cache found. Process and save one later.
		}

		err = json.Unmarshal(fileBytes, &AllTables)
		if err == nil {
			// Cache loaded.
			return
		}
	}

	// No cache found or requested

	file, err := os.Open(*Input)
	if err != nil {
		cli.ReturnError(err)
	}
	defer file.Close()

	var (
		currentTable    Table
		cacheAlterTable Table
	)

	if Options.Sqlite {
		DBFile, err = os.Create(TmpSqliteFile)
		if err != nil {
			cli.ReturnError(err)
		}
	}

	state := "IDLE"
	foundTable := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "COPY") {
			state = "COPY"
		}

		if strings.HasPrefix(line, "ALTER TABLE") {
			state = "ALTER-TABLE"
		}

		if Options.Sqlite && strings.HasPrefix(line, "CREATE TABLE") {
			state = "CREATE-TABLE"
		}

		if state == "CREATE-TABLE" {
			tmpLine := strings.ReplaceAll(line, "public.", "") + "\n"
			// Workaround for last column with comma
			if strings.HasPrefix(line, "    CONSTRAINT") {
				tmpLine = "    CONSTRAINT temp"
			}
			_, err := DBFile.WriteString(tmpLine)
			if err != nil {
				cli.ReturnError(err)
			}

			err = DBFile.Sync()
			if err != nil {
				cli.ReturnError(err)
			}

			if line == ");" {
				state = "IDLE"
			}
		}

		if state == "COPY" {
			// Get Metadata
			reMetadata := regexp.MustCompile(`COPY (\w+)\.(\w+) \((.+)\) FROM stdin;`)
			metadata := reMetadata.FindStringSubmatch(line)
			if len(metadata) == 4 {
				// It's expected 4 elements:
				// [0] Original line
				// [1] Schema
				// [2] Table
				// [3] Columns
				targetTable := Table{
					Name:   metadata[2],
					Schema: metadata[1],
				}
				if objTable, exist := findTable(AllTables, targetTable); exist {
					foundTable = true
					currentTable = *objTable
					currentTable.Columns = strings.Split(metadata[3], ", ")
				} else {
					currentTable = Table{
						Name:    targetTable.Name,
						Schema:  targetTable.Schema,
						Columns: strings.Split(metadata[3], ", "),
					}
				}
			} else {
				// Convert currentTable to map[string]interface{} for parser.Copy
				tableMap := map[string]interface{}{
					"name":        currentTable.Name,
					"schema":      currentTable.Schema,
					"data":        currentTable.Data,
					"columns":     currentTable.Columns,
					"values":      currentTable.Values,
					"primary_key": currentTable.PrimaryKey,
					"foreign_key": currentTable.ForeignKey,
				}
				parser.Copy(line, &tableMap)
				// Update currentTable from tableMap if needed
				if v, ok := tableMap["data"].([]map[string]string); ok {
					currentTable.Data = v
				}
				if v, ok := tableMap["values"].([][]string); ok {
					currentTable.Values = v
				}
				if strings.HasPrefix(line, "\\.") {
					if !foundTable {
						AllTables = append(AllTables, currentTable)
					} else {
						foundTable = false
					}
					currentTable = Table{}
					state = "IDLE"
				}
			}
		}

		if state == "ALTER-TABLE" {
			reAlterTable := regexp.MustCompile(`ALTER TABLE ONLY (\w+)\.(\w+)`)
			matchAlterTable := reAlterTable.FindStringSubmatch(line)
			if len(matchAlterTable) == 3 {
				// It's expected 3 elements:
				// [0] Original line
				// [1] Schema.
				// [2] Table
				cacheAlterTable = Table{
					Name:   matchAlterTable[2],
					Schema: matchAlterTable[1],
				}
			}
			if pkey := parser.PKey(line); pkey != "" {
				if objTable, exist := findTable(AllTables, cacheAlterTable); exist {
					(*objTable).PrimaryKey = pkey
				} else {
					currentTable = Table{
						Name:       cacheAlterTable.Name,
						Schema:     cacheAlterTable.Schema,
						PrimaryKey: pkey,
					}
					AllTables = append(AllTables, currentTable)
				}
				state = "IDLE"
			}
			if fkey := parser.FKey(line); fkey != nil {
				fkeys := []map[string]string{}
				if objTable, exist := findTable(AllTables, cacheAlterTable); exist {
					fromName := (*objTable).Name
					fromSchema := (*objTable).Schema
					fkey["from"] = fromSchema + "." + fromName + "." + fkey["from"]
					if objFkeys := (*objTable).ForeignKey; objFkeys != nil {
						(*objTable).ForeignKey = append(objFkeys, fkey)
					} else {
						(*objTable).ForeignKey = append(fkeys, fkey)
					}
				} else {
					currentTable = Table{
						Name:       cacheAlterTable.Name,
						Schema:     cacheAlterTable.Schema,
						ForeignKey: append(fkeys, fkey),
					}
					AllTables = append(AllTables, currentTable)
				}
				state = "IDLE"
			}
		}

	}

	// Save Cache

	if Options.Cache {
		err := os.MkdirAll(TmpCacheDir, os.ModePerm)
		if err != nil {
			cli.ReturnError(err)
		}
		file, err := os.Create(TmpCacheFile)
		if err != nil {
			cli.ReturnError(err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		if err := encoder.Encode(AllTables); err != nil {
			cli.ReturnError(err)
		}
	}
}

func Export() {

	var schema string
	if cli.Filters.Schema != "" {
		schema = cli.Filters.Schema
	} else {
		schema = "public"
	}

	if Options.Json || Options.JsonPretty {
		var tablesToExport []Table
		if cli.Filters.TableName != "" {
			for _, table := range AllTables {
				if table.Name == cli.Filters.TableName &&
					table.Schema == schema {
					tablesToExport = append(tablesToExport, table)
				}
			}
		} else {
			tablesToExport = AllTables
		}

		var output []byte
		var err error

		if len(tablesToExport) > 0 {
			if Options.JsonPretty {
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

	if Options.Yaml {
		var tablesToExport []Table
		if cli.Filters.TableName != "" {
			for _, table := range AllTables {
				if table.Name == cli.Filters.TableName &&
					table.Schema == schema {
					tablesToExport = append(tablesToExport, table)
				}
			}
		} else {
			tablesToExport = AllTables
		}

		out, err := yaml.Marshal(tablesToExport)
		if err != nil {
			cli.ReturnError(err)
		}
		fmt.Println(string(out))
	}

	if Options.Html {
		// Parse template
		tmpl, err := template.New("index").Parse(htmlTemplate)
		if err != nil {
			cli.ReturnError(err)
		}

		// Create output file
		outputFile, err := os.Create("index.html")
		if err != nil {
			cli.ReturnError(err)
		}
		defer outputFile.Close()

		// Execute template with data
		err = tmpl.Execute(outputFile, AllTables)
		if err != nil {
			cli.ReturnError(err)
		}

		fmt.Printf("%s/index.html created!\n", cwd)
	}

	if Options.Sqlite {
		for _, table := range AllTables {
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

				_, err := DBFile.WriteString(insertStmt)
				if err != nil {
					cli.ReturnError(err)
				}

			}

		}

		fmt.Printf("%s/%s created!\n", cwd, TmpSqliteFile)
		fmt.Printf("Import it using: sqlite3 <db-name>.sqlite3 < <%s path>\n", TmpSqliteFile)

		if DBFile != nil {
			DBFile.Close()
		}
	}

	if Options.Cli {
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
					if contains(columns, col) {
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
}
