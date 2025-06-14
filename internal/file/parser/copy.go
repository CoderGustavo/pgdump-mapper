package parser

import (
	"strings"
)

func Copy(line string, tbl *map[string]interface{}) {
	if line == "\\." {
		return
	}
	columns := (*tbl)["columns"].([]string)
	numColumns := len(columns)
	values := strings.Split(line, "\t")

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
