package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v3"

	cli "github.com/hedibertosilva/pgdump-mapper/internal/cli"
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
	TmpSqliteFile string = "/tmp/pgdump-mapper.db.sqlite.txt"
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
			if strings.HasPrefix(line, "    CONSTRAINT") {
				tmpLine = "    CONSTRAINT tmp"
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
				parseCopy(line, &currentTable)
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
			if pkey := parsePKey(line); pkey != "" {
				if Options.Sqlite {
					DBFile.WriteString(fmt.Sprintf("ALTER TABLE %s\n", cacheAlterTable["name"]))
					DBFile.WriteString(fmt.Sprintf("%s\n", line))

					err = DBFile.Sync()
					if err != nil {
						cli.ReturnError(err)
					}
				}
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
			if fkey := parseFKey(line); fkey != nil {
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
	if Options.Json {
		j, _ := json.Marshal(AllTables)
		fmt.Println(string(j))
	}

	if Options.Yaml {
		out, err := yaml.Marshal(AllTables)
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

		cwd, err := os.Getwd()
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

		if DBFile != nil {
			DBFile.Close()
		}
	}
}
