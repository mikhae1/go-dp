package config

import (
	"fmt"

	"github.com/minkolazer/gp/lib"
)

// Logger provides basic logging support
type Logger struct {
}

func (writer Logger) Write(bytes []byte) (int, error) {
	return fmt.Print(lib.ColorInfo("[DEBUG]") + " " + string(bytes))
}
