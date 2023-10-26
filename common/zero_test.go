package common_test

import (
	"github.com/CharLemAznable/resilience4go/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZero(t *testing.T) {
	// Test case 1: Test with int type
	expectedInt := 0
	resultInt := common.Zero[int]()
	assert.Equal(t, expectedInt, resultInt)

	// Test case 2: Test with string type
	expectedString := ""
	resultString := common.Zero[string]()
	assert.Equal(t, expectedString, resultString)

	// Test case 3: Test with custom struct type
	type CustomStruct struct {
		Name string
		Age  int
	}
	expectedCustomStruct := CustomStruct{}
	resultCustomStruct := common.Zero[CustomStruct]()
	assert.Equal(t, expectedCustomStruct, resultCustomStruct)
}
