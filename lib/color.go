package lib

import (
	"hash/fnv"

	fcolor "github.com/fatih/color"
)

var defaultColors = []fcolor.Attribute{
	fcolor.FgCyan,
	fcolor.FgYellow,
	fcolor.FgMagenta,
	fcolor.FgBlue,

	fcolor.FgHiGreen,
	fcolor.FgHiYellow,
	fcolor.FgHiBlue,
	fcolor.FgHiMagenta,
	fcolor.FgHiCyan,
}

// Color colors string
func Color(str string) (coloredStr string) {
	hash := fnv.New32a()

	hash.Write([]byte(str))

	colorAtrr := defaultColors[hash.Sum32()%uint32(len(defaultColors))]

	coloredStr = fcolor.New(colorAtrr).SprintFunc()(str)

	return
}

// ColorErr prints red string
func ColorErr(s string) string {
	return fcolor.New(fcolor.FgRed).SprintFunc()(s)
}

// ColorOK pronts green string
func ColorOK(s string) string {
	return fcolor.New(fcolor.FgGreen).SprintFunc()(s)
}

// ColorInfo prints color string
func ColorInfo(s string) string {
	return fcolor.New(fcolor.FgCyan).SprintFunc()(s)
}

// ColorStrong prints bold string
func ColorStrong(s string) string {
	return fcolor.New(fcolor.Bold).SprintFunc()(s)
}
