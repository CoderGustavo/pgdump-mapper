package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"

	yaml "gopkg.in/yaml.v3"

	cli "github.com/hedibertosilva/pgdump-mapper/internal/cli"
	parser "github.com/hedibertosilva/pgdump-mapper/internal/file/parser"
	models "github.com/hedibertosilva/pgdump-mapper/models"
)

// Table Template:
//
//  map[string]interface{}{
// 		"name":        "",
// 		"schema":      "",
// 		"data":        []map[string]string{},
// 		"columns":     map[string]string{},
// 		"values":      [][]string{},
// 		"primary_key": "",
// 		"foreign_key": []map[string]string{},
// }

var (
	Input         *string
	Options       models.Options
	AllTables     []map[string]interface{}
	DBFile        *os.File
	TmpSqliteFile string = "pgdump-mapper.sqlite.txt"
	cwd, _               = os.Getwd()
)

func findTable(allTables []map[string]interface{},
	targetTable map[string]string) (*map[string]interface{}, bool) {
	for _, table := range allTables {
		if table["name"] == targetTable["name"] &&
			table["schema"] == targetTable["schema"] {
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
	file, err := os.Open(*Input)
	if err != nil {
		cli.ReturnError(err)
	}
	defer file.Close()

	var (
		currentTable    map[string]interface{}
		cacheAlterTable map[string]string
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
				targetTable := map[string]string{
					"name":   metadata[2],
					"schema": metadata[1],
				}
				if objTable, exist := findTable(AllTables, targetTable); exist {
					foundTable = true
					currentTable = *objTable
					currentTable["columns"] = strings.Split(metadata[3], ", ")
				} else {
					currentTable = map[string]interface{}{
						"name":    targetTable["name"],
						"schema":  targetTable["schema"],
						"columns": strings.Split(metadata[3], ", "),
					}
				}
			} else {
				parser.Copy(line, &currentTable)
				if strings.HasPrefix(line, "\\.") {
					if !foundTable {
						AllTables = append(AllTables, currentTable)
					} else {
						foundTable = false
					}
					currentTable = make(map[string]interface{})
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
				cacheAlterTable = map[string]string{
					"schema": matchAlterTable[1],
					"name":   matchAlterTable[2],
				}
			}
			if pkey := parser.PKey(line); pkey != "" {
				if objTable, exist := findTable(AllTables, cacheAlterTable); exist {
					(*objTable)["primary_key"] = pkey
				} else {
					currentTable = map[string]interface{}{
						"name":        cacheAlterTable["name"],
						"schema":      cacheAlterTable["schema"],
						"primary_key": pkey,
					}
					AllTables = append(AllTables, currentTable)
				}
				state = "IDLE"
			}
			if fkey := parser.FKey(line); fkey != nil {
				fkeys := []map[string]string{}
				if objTable, exist := findTable(AllTables, cacheAlterTable); exist {
					fromName := (*objTable)["name"].(string)
					fromSchema := (*objTable)["schema"].(string)
					fkey["from"] = fromSchema + "." + fromName + "." + fkey["from"]
					if objFkeys, exist := (*objTable)["foreign_key"]; exist {
						(*objTable)["foreign_key"] = append(objFkeys.([]map[string]string), fkey)
					} else {
						(*objTable)["foreign_key"] = append(fkeys, fkey)
					}
				} else {
					currentTable = map[string]interface{}{
						"name":        cacheAlterTable["name"],
						"schema":      cacheAlterTable["schema"],
						"foreign_key": append(fkeys, fkey),
					}
					AllTables = append(AllTables, currentTable)
				}
				state = "IDLE"
			}
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
		var tablesToExport []map[string]interface{}
		if cli.Filters.TableName != "" {
			for _, table := range AllTables {
				if table["name"].(string) == cli.Filters.TableName &&
					table["schema"].(string) == schema {
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
		var tablesToExport []map[string]interface{}
		if cli.Filters.TableName != "" {
			for _, table := range AllTables {
				if table["name"].(string) == cli.Filters.TableName &&
					table["schema"].(string) == schema {
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
			tableName := table["name"]
			columns := table["columns"].([]string)
			var data []map[string]string
			if table["data"] != nil {
				data = table["data"].([]map[string]string)
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
			tableName := table["name"].(string)

			// Filter table
			if cli.Filters.TableName != "" && (table["name"].(string) == cli.Filters.TableName &&
				table["schema"].(string) == schema) {
				continue
			}

			columns := table["columns"].([]string)
			var data []map[string]string
			if table["data"] != nil {
				data = table["data"].([]map[string]string)
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

			fmt.Printf("\nTable: %s\n", tableName)
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
