package lib

import (
	"fmt"
	"os"
)

// Exception prints error and exit
func Exception(err error) {
	if err == nil {
		return
	}

	fmt.Println(ColorErr(err.Error()))
	os.Exit(1)
}
