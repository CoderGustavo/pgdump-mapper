package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	// "regexp"
	// "strings"
	// "text/scanner"

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
//		"foreign_key": map[string]string{},
// }
//

func parseCopy(line string, tbl *map[string]interface{}) {
	reData := regexp.MustCompile(`([^\t]+)`)
	var columns []string = (*tbl)["columns"].([]string)
	var numColumns int = len(columns)
	var values []string = reData.FindAllString(line, -1)
	var numValues int = len(values)

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
	// Get Primary Key
	rePKey := regexp.MustCompile(`ADD CONSTRAINT (\w+) PRIMARY KEY \((\w+)\);`)
	matchPKey := rePKey.FindStringSubmatch(line)

	if len(matchPKey) == 3 {
		return matchPKey[2]
	}

	return ""
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

	var currentTable map[string]interface{}
	var allTables []map[string]interface{}
	var tmpAlterTable map[string]string
	found := false
	state := "IDLE"
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
				targetTable := map[string]string{
					"name":   metadata[2],
					"schema": metadata[1],
				}
				if objTable, exist := findTable(allTables, targetTable); exist {
					found = true
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
				if !found {
					allTables = append(allTables, currentTable)
				} else {
					found = false
				}
				currentTable = make(map[string]interface{})
				state = "IDLE"
			}
		}

		if state == "ALTER-TABLE" {
			reAlterTable := regexp.MustCompile(`ALTER TABLE ONLY (\w+)\.(\w+)`)
			matchAlterTable := reAlterTable.FindStringSubmatch(line)
			if len(matchAlterTable) == 3 {
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
		}

	}
	j, _ := json.Marshal(allTables)
	fmt.Println(string(j))

}

func Export() {
	fmt.Println("exporting", *Input)
}
