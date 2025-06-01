package file

import "regexp"

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
