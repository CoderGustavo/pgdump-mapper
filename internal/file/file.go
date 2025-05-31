package file

import (
	"fmt"

	models "github.com/hedibertosilva/pgdump-mapper/models"
)

var Input *string
var Options models.Options

func Read() {
	fmt.Println("reading", *Input)
}

func Export() {
	fmt.Println("exporting", *Input)
}
