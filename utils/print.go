package utils

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type printData struct {
	Name  string
	Value interface{}
}

// PrintStruct print data by name
func PrintStruct(t interface{}, sorted bool) {
	getType := reflect.TypeOf(t).Elem()
	getValue := reflect.ValueOf(t).Elem()
	pList := make([]printData, 0, getValue.NumField())
	maxLen := 0
	for i := 0; i < getValue.NumField(); i++ {
		field := getType.Field(i)
		value := getValue.FieldByName(field.Name)
		pList = append(pList, printData{field.Name, value})
		if len(field.Name) > maxLen {
			maxLen = len(field.Name)
		}
	}
	if sorted {
		sort.Slice(pList, func(i, j int) bool {
			return strings.ToUpper(pList[i].Name) < strings.ToUpper(pList[j].Name)
		})
	}
	fmt.Println("*******************************************************************")
	for _, item := range pList {
		space := strings.Repeat(" ", maxLen+5-len(item.Name))
		fmt.Printf("%s:%s%+v\n", item.Name, space, item.Value)
	}
	fmt.Printf("###################################################################\n\n")
}
