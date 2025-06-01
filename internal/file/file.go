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

var Input *string
var Options models.Options

// currentTable Template:
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
//

func parseCopy(line string, tbl *map[string]interface{}) {
	reData := regexp.MustCompile(`([^\t]+)`)
	columns := (*tbl)["columns"].([]string)
	numColumns := len(columns)
	values := reData.FindAllString(line, -1)
	numValues := len(values)

	if numValues == numColumns {
		data := make(map[string]string)
		for i := 0; i < numColumns; i++ {
			data[columns[i]] = values[i]
		}
		if _, exist := (*tbl)["data"]; !exist {
			(*tbl)["data"] = []map[string]string{}
		}
		(*tbl)["data"] = append((*tbl)["data"].([]map[string]string), data)
		if _, exist := (*tbl)["values"]; !exist {
			(*tbl)["values"] = [][]string{}
		}
		(*tbl)["values"] = append((*tbl)["values"].([][]string), values)
	}

}

func parsePKey(line string) string {
	re := regexp.MustCompile(`ADD CONSTRAINT (\w+) PRIMARY KEY \((\w+)\);`)
	match := re.FindStringSubmatch(line)
	if len(match) == 3 {
		// It's expected 3 elements:
		// [0] Original line
		// [1] PKey name
		// [2] Pkey column
		return match[2]
	}
	return ""
}

func parseFKey(line string) map[string]string {
	re := regexp.MustCompile(`ADD CONSTRAINT (\w+) FOREIGN KEY \((\w+)\) REFERENCES (\w+).(\w+)\((\w+)\)`)
	match := re.FindStringSubmatch(line)
	if len(match) == 6 {
		// It's expected 6 elements:
		// [0] Original line
		// [1] FKey name
		// [2] From column
		// [3] Target schema
		// [4] Target table
		// [5] Target column
		return map[string]string{
			"from":   match[2],
			"target": match[3] + "." + match[4] + "." + match[5],
		}
	}
	return nil
}

func findTable(allTables []map[string]interface{}, tmpAlterTable map[string]string) (*map[string]interface{}, bool) {
	for _, table := range allTables {
		if table["name"] == tmpAlterTable["name"] && table["schema"] == tmpAlterTable["schema"] {
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
		currentTable  map[string]interface{}
		allTables     []map[string]interface{}
		tmpAlterTable map[string]string
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
				tmpAlterTable = map[string]string{
					"schema": matchAlterTable[1],
					"name":   matchAlterTable[2],
				}
			}
			if pkey := parsePKey(line); pkey != "" {
				if objTable, exist := findTable(allTables, tmpAlterTable); exist {
					(*objTable)["primary_key"] = pkey
				} else {
					currentTable = map[string]interface{}{
						"name":        tmpAlterTable["name"],
						"schema":      tmpAlterTable["schema"],
						"primary_key": pkey,
					}
					allTables = append(allTables, currentTable)
				}
				state = "IDLE"
			}
			if fkey := parseFKey(line); fkey != nil {
				fkeys := []map[string]string{}
				if objTable, exist := findTable(allTables, tmpAlterTable); exist {
					if objFkeys, exist := (*objTable)["foreign_key"]; exist {
						(*objTable)["foreign_key"] = append(objFkeys.([]map[string]string), fkey)
					} else {
						(*objTable)["foreign_key"] = append(fkeys, fkey)
					}
				} else {
					currentTable = map[string]interface{}{
						"name":        tmpAlterTable["name"],
						"schema":      tmpAlterTable["schema"],
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
