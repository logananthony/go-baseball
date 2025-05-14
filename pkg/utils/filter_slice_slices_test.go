package utils

import (
	"reflect"
	"testing"
)

func TestFilterSliceSlices(t *testing.T) {
	input := [][]int{
		{10, 2024},
		{20, 2023},
		{30, 2022},
	}

	expected := [][]int{
		{10, 2024},
		{30, 2022},
	}

	result := FilterSliceSlices(input, 20)
  t.Log("Result: ", result)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

