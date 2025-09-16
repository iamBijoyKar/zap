package out

import (
	"log"

	"github.com/fatih/color"
)

func LogError(err string) {
	log.Print(color.RedString(err))
}
