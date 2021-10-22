package pages

import (
	"fmt"
	"forum/utils"
)

func templateNotFound(name string) {
	err := fmt.Errorf("could not find the template for %v.html", name)
	utils.FatalErr(err)
}
