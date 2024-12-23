package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type CustomType struct {
}

type Custom struct {
}
type Custom1 struct {
}

func Test_parseTypeName_case(t *testing.T) {
	assert.Equal(t, ParseTypeName[CustomType](), "custom-type")
}
func Test_parseTypeName_case2(t *testing.T) {
	assert.Equal(t, ParseTypeName[Custom](), "custom")
}
func Test_parseTypeName_case3(t *testing.T) {
	assert.Equal(t, ParseTypeName[Custom1](), "custom1")
}
