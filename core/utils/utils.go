package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// Example: Custom -> custom, CustomType -> custom-type
func ParseTypeName[Z any]() string {
	tp := getType[Z]()
	re := regexp.MustCompile(`[A-Z][^A-Z]*`)
	splittedType := re.FindAllString(tp, -1)
	for index, value := range splittedType {
		splittedType[index] = strings.ToLower(value)
	}
	n := strings.ToLower(strings.Join(splittedType, "-"))
	return n
}

func getType[T any]() string {
	s := fmt.Sprintf("%T", new(T))
	return strings.Split(s, ".")[1]
}
