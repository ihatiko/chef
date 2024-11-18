package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// ParseTypeName Example: Custom -> custom, CustomType -> custom-type
func ParseTypeName[Z any]() string {
	tp := getType[Z]()
	re := regexp.MustCompile(`[A-Z][^A-Z]*`)
	splitType := re.FindAllString(tp, -1)
	for index, value := range splitType {
		splitType[index] = strings.ToLower(value)
	}
	n := strings.ToLower(strings.Join(splitType, "-"))
	return n
}

func getType[T any]() string {
	s := fmt.Sprintf("%T", new(T))
	return strings.Split(s, ".")[1]
}
