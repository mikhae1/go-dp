package cmd

import (
	"fmt"
	"os"
	"reflect"

	fcolor "github.com/fatih/color"
)

func colorErr(s string) string {
	return fcolor.New(fcolor.FgRed).SprintFunc()(s)
}

func colorOK(s string) string {
	return fcolor.New(fcolor.FgGreen).SprintFunc()(s)
}

func colorInfo(s string) string {
	return fcolor.New(fcolor.FgCyan).SprintFunc()(s)
}

func colorStrong(s string) string {
	return fcolor.New(fcolor.Bold).SprintFunc()(s)
}

func exception(err error) {
	if err == nil {
		return
	}

	fmt.Println(colorErr(err.Error()))
	os.Exit(1)
}

func arrayContains(array interface{}, val interface{}) (index int) {
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				return
			}
		}
	}

	return
}
