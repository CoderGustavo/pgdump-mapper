package parser

import "regexp"

func FKey(line string) map[string]string {
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
