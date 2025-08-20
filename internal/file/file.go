package file

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	cli "github.com/hedibertosilva/pgdump-mapper/internal/cli"
	exporters "github.com/hedibertosilva/pgdump-mapper/internal/file/exporters"
	parsers "github.com/hedibertosilva/pgdump-mapper/internal/file/parsers"
	models "github.com/hedibertosilva/pgdump-mapper/models"
)

var (
	Input         *string
	Options       models.Options
	tables        []models.Table
	dbFile        *os.File
	tmpSqliteFile string = "pgdump-mapper.sqlite.txt"
	tmpCacheDir   string = "/tmp/pgdump-mapper"
	tmpCacheFile  string = ""
	rootPath, _          = os.Getwd()
)

func FindTable(tables []models.Table, targetTable models.Table) (*models.Table, bool) {
	for _, table := range tables {
		if table.Name == targetTable.Name &&
			table.Schema == targetTable.Schema {
			return &table, true
		}
	}
	return nil, false
}

func Read() {

	if Options.Cache {
		basename := filepath.Base(*Input)
		tmpCacheFile = filepath.Join(tmpCacheDir, basename)

		fileBytes, err := os.ReadFile(tmpCacheFile)
		if err != nil {
			// No cache found. Process and save one later.
		}

		err = json.Unmarshal(fileBytes, &tables)
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
		currentTable    models.Table
		cacheAlterTable models.Table
	)

	if Options.Sqlite {
		dbFile, err = os.Create(tmpSqliteFile)
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
			_, err := dbFile.WriteString(tmpLine)
			if err != nil {
				cli.ReturnError(err)
			}

			err = dbFile.Sync()
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
				targetTable := models.Table{
					Name:   metadata[2],
					Schema: metadata[1],
				}
				if objTable, exist := FindTable(tables, targetTable); exist {
					foundTable = true
					currentTable = *objTable
					currentTable.Columns = strings.Split(metadata[3], ", ")
				} else {
					currentTable = models.Table{
						Name:    targetTable.Name,
						Schema:  targetTable.Schema,
						Columns: strings.Split(metadata[3], ", "),
					}
				}
			} else {
				// Convert currentTable to map[string]interface{} for parsers.Copy
				tableMap := map[string]interface{}{
					"name":        currentTable.Name,
					"schema":      currentTable.Schema,
					"data":        currentTable.Data,
					"columns":     currentTable.Columns,
					"values":      currentTable.Values,
					"primary_key": currentTable.PrimaryKey,
					"foreign_key": currentTable.ForeignKey,
				}
				parsers.Copy(line, &tableMap)
				// Update currentTable from tableMap if needed
				if v, ok := tableMap["data"].([]map[string]string); ok {
					currentTable.Data = v
				}
				if v, ok := tableMap["values"].([][]string); ok {
					currentTable.Values = v
				}
				if strings.HasPrefix(line, "\\.") {
					if !foundTable {
						tables = append(tables, currentTable)
					} else {
						foundTable = false
					}
					currentTable = models.Table{}
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
				cacheAlterTable = models.Table{
					Name:   matchAlterTable[2],
					Schema: matchAlterTable[1],
				}
			}
			if pkey := parsers.PKey(line); pkey != "" {
				if objTable, exist := FindTable(tables, cacheAlterTable); exist {
					(*objTable).PrimaryKey = pkey
				} else {
					currentTable = models.Table{
						Name:       cacheAlterTable.Name,
						Schema:     cacheAlterTable.Schema,
						PrimaryKey: pkey,
					}
					tables = append(tables, currentTable)
				}
				state = "IDLE"
			}
			if fkey := parsers.FKey(line); fkey != nil {
				fkeys := []map[string]string{}
				if objTable, exist := FindTable(tables, cacheAlterTable); exist {
					fromName := (*objTable).Name
					fromSchema := (*objTable).Schema
					fkey["from"] = fromSchema + "." + fromName + "." + fkey["from"]
					if objFkeys := (*objTable).ForeignKey; objFkeys != nil {
						(*objTable).ForeignKey = append(objFkeys, fkey)
					} else {
						(*objTable).ForeignKey = append(fkeys, fkey)
					}
				} else {
					currentTable = models.Table{
						Name:       cacheAlterTable.Name,
						Schema:     cacheAlterTable.Schema,
						ForeignKey: append(fkeys, fkey),
					}
					tables = append(tables, currentTable)
				}
				state = "IDLE"
			}
		}

	}

	// Save Cache

	if Options.Cache {
		err := os.MkdirAll(tmpCacheDir, os.ModePerm)
		if err != nil {
			cli.ReturnError(err)
		}
		file, err := os.Create(tmpCacheFile)
		if err != nil {
			cli.ReturnError(err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		if err := encoder.Encode(tables); err != nil {
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
		exporters.JSON(schema, tables, Options.JsonPretty)
	}

	if Options.Yaml {
		exporters.YAML(schema, tables)
	}

	if Options.Html {
		exporters.HTML(tables, rootPath)
	}

	if Options.Sqlite {
		exporters.SQLite(schema, tables, dbFile, rootPath, tmpSqliteFile)
	}

	if Options.Cli {
		exporters.CLI(schema, tables)
	}
}
