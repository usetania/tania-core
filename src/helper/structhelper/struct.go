package structhelper

import (
	"reflect"
	"strings"
)

func GetName(val interface{}) string {
	name := reflect.TypeOf(val).String()
	splitted := strings.Split(name, ".")

	return splitted[len(splitted)-1]
}
