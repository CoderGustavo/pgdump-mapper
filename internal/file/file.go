package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	cli "github.com/hedibertosilva/pgdump-mapper/internal/cli"
	models "github.com/hedibertosilva/pgdump-mapper/models"
)

// Table Template:
//
//  map[string]interface{}{
//		"name":        "",
//		"schema":      "",
//		"data":        []map[string]string{},
//		"columns":     map[string]string{},
//		"values":      [][]string{},
//		"primary_key": "",
//		"foreign_key": []map[string]string{},
// }

var (
	Input   *string
	Options models.Options
)

func findTable(allTables []map[string]interface{},
	cacheAlterTable map[string]string) (*map[string]interface{}, bool) {
	for _, table := range allTables {
		if table["name"] == cacheAlterTable["name"] &&
			table["schema"] == cacheAlterTable["schema"] {
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
		allTables       []map[string]interface{}
		cacheAlterTable map[string]string
	)

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
				if objTable, exist := findTable(allTables, targetTable); exist {
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
			}
			parseCopy(line, &currentTable)
			if strings.HasPrefix(line, "\\.") {
				if !foundTable {
					allTables = append(allTables, currentTable)
				} else {
					foundTable = false
				}
				currentTable = make(map[string]interface{})
				state = "IDLE"
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
				if objTable, exist := findTable(allTables, cacheAlterTable); exist {
					(*objTable)["primary_key"] = pkey
				} else {
					currentTable = map[string]interface{}{
						"name":        cacheAlterTable["name"],
						"schema":      cacheAlterTable["schema"],
						"primary_key": pkey,
					}
					allTables = append(allTables, currentTable)
				}
				state = "IDLE"
			}
			if fkey := parseFKey(line); fkey != nil {
				fkeys := []map[string]string{}
				if objTable, exist := findTable(allTables, cacheAlterTable); exist {
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
					allTables = append(allTables, currentTable)
				}
				state = "IDLE"
			}
		}

	}
	j, _ := json.Marshal(allTables)
	fmt.Println(string(j))

}

func Export() {
	fmt.Println("Exporting", *Input)
}
