package out

import (
	"fmt"

	"github.com/fatih/color"
)

func PrintError(err string) {
	fmt.Print(color.RedString(err))
}

func PrintInfo(text string) {
	fmt.Print(color.YellowString(text))
}

func PrintDefault(text string) {
	fmt.Print(color.WhiteString(text))
}
