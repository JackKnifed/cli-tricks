package cliTricks

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestBreakupStringArray(t *testing.T) {
	testData := []struct {
		input  string
		output []string
	}{
		{
			input:  "apple,banana,cherry",
			output: []string{"apple", "banana", "cherry"},
		}, {
			input:  "Dog, Eagle, Fox",
			output: []string{"Dog", "Eagle", "Fox"},
		}, {
			input:  "[Green Beans][Hot Tamales][Ice Cream]",
			output: []string{"Green Beans", "Hot Tamales", "Ice Cream"},
		}, {
			input:  "[JellyBean],[KitKat],[Marshmallow]",
			output: []string{"JellyBean", "KitKat", "Marshmallow"},
		}, {
			input:  "[\"Nutella\"],[\"Oatmeal\"],[\"Pie\"]",
			output: []string{"Nutella", "Oatmeal", "Pie"},
		},
	}

	for _, oneTest := range testData {
		assert.Equal(t, oneTest.output, BreakupStringArray(oneTest.input), "BreakupStringArray returned non-expected results")
	}
}

func ExampleGetItem() {
	testBytes := []byte(`{"Everything":"Awesome","Team":{"Everything":"Cool"}}`)
	var testData interface{}
	err := json.Unmarshal(testBytes, &testData)
	if err != nil {
		fmt.Printf("hit a snag unmarshalling the data - %v", err)
	}

	item, err := GetItem(testData, []string{"Team", "Everything"})
	if err != nil {
		fmt.Printf("hit a snag retrieving the item - %v", err)
	}
	fmt.Println(item)

	// Output:
	// Cool
}

// func TestGetInt(t *testing.T) {
// 	testData := []struct{
// 		input interface{}
// 		target []string
// 		output int
// 	}{
// 		{
// 			input: map[string]interface{}{
// 						"params": map[string]float64{
// 							"data": 63,
// 					},
// 			},
// 			target: []string{"params", "data",},
// 			output: 63,
// 		},
// 	}

// 	for _, oneTest := range testData {
// 		result, err := GetInt(oneTest.input, oneTest.target)
// 		assert.Equal(t, oneTest.output, result)
// 		assert.NoError(t, err)

// 		result2, err := GetItem(oneTest.input, []string{"params"})
// 		fmt.Println(oneTest.input)
// 		fmt.Println(oneTest.output)
// 		fmt.Println(result2)
// 	}
// }

func TestGetIntJSON(t *testing.T) {
	testData := []struct {
		input  []byte
		target []string
		output int
		status error
	}{
		{
			input:  []byte(`{"params":{"data":63}}`),
			target: []string{"params", "data"},
			output: 63,
			status: nil,
		}, {
			// we always round down
			input:  []byte(`{"params":{"data":63.9}}`),
			target: []string{"params", "data"},
			output: 63,
			status: nil,
		}, {
			input:  []byte(`{"params":{"data":"potato"}}`),
			target: []string{"params", "data"},
			output: -1,
			status: errors.New("got non-float item - potato"),
		}, {
			input:  []byte(`{"params":{"data":"potato"}}`),
			target: []string{"bad", "address"},
			output: -1,
			status: errors.New("bad item - bad address - [address]"),
		},
	}

	for _, oneTest := range testData {
		var testData interface{}
		err := json.Unmarshal(oneTest.input, &testData)
		assert.Nil(t, err, "Problems unmarshaling the input")

		result, err := GetInt(testData, oneTest.target)
		assert.Equal(t, oneTest.output, result)
		assert.Equal(t, oneTest.status, err)
	}
}

func TestGetItemJSON(t *testing.T) {
	testData := []struct {
		input  []byte
		target []string
		output []byte
		status error
	}{
		{
			input:  []byte(`{"params":{"data":63}}`),
			target: []string{"params", "data"},
			output: []byte(`63`),
			status: nil,
		}, {
			// we always round down
			input:  []byte(`{"params":{"data":63.9}}`),
			target: []string{"params", "data"},
			output: []byte(`63.9`),
			status: nil,
		}, {
			input:  []byte(`{"params":{"data":"potato"}}`),
			target: []string{"params", "data"},
			output: []byte(`"potato"`),
			status: nil,
		}, {
			input:  []byte(`{"numbers":[4,8,15,16,23,42,63]}`),
			target: []string{"numbers", "3"},
			output: []byte(`16`),
			status: nil,
		}, {
			input:  []byte(`{"numbers":[4,8,15,16,23,42,63]}`),
			target: []string{"numbers", "potato"},
			output: []byte("null"),
			status: errors.New("got non-int address for []interface{}"),
		}, {
			input:  []byte(`[["apple","apricot","acorn"],"banana",["chestnut","cookie"]]`),
			target: []string{"0", "1"},
			output: []byte(`"apricot"`),
			status: nil,
		}, {
			input:  []byte(`{"params":{"data":"potato"}}`),
			target: []string{"bad", "address"},
			output: []byte("null"),
			status: errors.New("bad address - [address]"),
		},
	}

	for _, oneTest := range testData {
		var inputData, outputData interface{}
		err := json.Unmarshal(oneTest.input, &inputData)
		assert.Nil(t, err, "Problems unmarshaling the input - %q", oneTest.input)
		err = json.Unmarshal(oneTest.output, &outputData)
		assert.Nil(t, err, "Problems unmarshaling the output - input was %q and output was %q", oneTest.input, oneTest.output)

		result, err := GetItem(inputData, oneTest.target)
		assert.Equal(t, outputData, result)
		assert.Equal(t, oneTest.status, err)
	}
}

func TestSetItemJSON(t *testing.T) {
	testData := []struct {
		input  []byte
		target []string
		newVal []byte
		output []byte
		status error
	}{
		{
			input:  []byte(`{"params":{"data":63}}`),
			target: []string{"params", "data"},
			newVal: []byte("63"),
			output:  []byte(`{"params":{"data":63}}`),
			status: nil,
		}, {
			// we always round down
			input:  []byte(`{"params":{"data":42}}`),
			target: []string{"params", "data"},
			newVal: []byte(`63.9`),
			output:  []byte(`{"params":{"data":63}}`),
			status: nil,
		}, {
			input:  []byte(`{"params":{"data":"potato"}}`),
			target: []string{"params", "data"},
			newVal: []byte(`"banana"`),
			output:  []byte(`{"params":{"data":"banana"}}`),
			status: nil,
		}, {
			input:  []byte(`{"params":{"data":"potato"}}`),
			target: []string{"params", "magic"},
			newVal: []byte(`"banana"`),
			output:  []byte(`{"params":{"data":"banana","magic":"banana"}}`),
			status: nil,
		}, {
			input:  []byte(`{"numbers":[4,8,15,16,23,42]}`),
			target: []string{"numbers", "6"},
			newVal: []byte(`63`),
			output:  []byte(`{"numbers":[4,8,15,16,23,42,63]}`),
			status: nil,
		}, {
			input:  []byte(`[["apple","apricot"],"banana",["chestnut","cookie"]]`),
			target: []string{"0", "2"},
			newVal: []byte(`"acorn"`),
			output:  []byte(`[["apple","apricot","acorn"],"banana",["chestnut","cookie"]]`),
			status: nil,
		}, {
			input:  []byte(`[["apple","acorn"],"banana",["chestnut","cookie"]]`),
			target: []string{"0", "1"},
			newVal: []byte(`"apricot"`),
			output:  []byte(`[["apple","apricot","acorn"],"banana",["chestnut","cookie"]]`),
			status: nil,
		}, {
			input:  []byte(`{"params":{"data":"potato"}}`),
			target: []string{"bad", "address"},
			newVal: []byte(""),
			output:  []byte(`{"params":{"data":"potato"}}`),
			status: errors.New("bad address - [address]"),
		},
	}

	for _, oneTest := range testData {
		var inputData, newData, outputData interface{}
		err := json.Unmarshal(oneTest.input, &inputData)
		assert.Nil(t, err, "Problems unmarshaling the input")
		err = json.Unmarshal(oneTest.newVal, &newData)
		assert.Nil(t, err, "Problems unmarshaling the newData")
		err = json.Unmarshal(oneTest.output, &outputData)
		assert.Nil(t, err, "Problems unmarshaling the output")

		err = SetItem(&inputData, newData, oneTest.target)
		assert.Equal(t, outputData, inputData)
		assert.Equal(t, oneTest.status, err)
	}
}

func ExampleGetInt() {
	testBytes := []byte(`{"Everything":"Awesome","Team":{"Everything":"Cool", "Solution": 63}}`)
	var testData interface{}
	err := json.Unmarshal(testBytes, &testData)
	if err != nil {
		fmt.Printf("hit a snag unmarshalling the data - %v", err)
	}

	item, err := GetInt(testData, []string{"Team", "Solution"})
	if err != nil {
		fmt.Printf("hit a snag retrieving the item - %v", err)
		return
	}
	fmt.Println(item)
	fmt.Println(reflect.TypeOf(item))

	// Output:
	// 63
	// int
}