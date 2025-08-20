package exporters

import (
	"fmt"
	"os"
	"text/template"

	cli "github.com/hedibertosilva/pgdump-mapper/internal/cli"
	templates "github.com/hedibertosilva/pgdump-mapper/internal/data/templates"
	models "github.com/hedibertosilva/pgdump-mapper/models"
)

func HTML(tables []models.Table, rootPath string) {
	// Parse template
	tmpl, err := template.New("index").Parse(templates.HTML)
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
	err = tmpl.Execute(outputFile, tables)
	if err != nil {
		cli.ReturnError(err)
	}

	fmt.Printf("%s/index.html created!\n", rootPath)
}
