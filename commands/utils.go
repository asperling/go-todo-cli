package commands

import (
	"fmt"
	"strconv"
)

func ValidateArgs(args []string, minLength int, intIndices []int, usage string) ([]int, error) {
	if len(args) < minLength {
		return nil, fmt.Errorf("%s", usage)
	}
	var intValues []int
	for _, index := range intIndices {
		if index >= len(args) || index < 0 {
			return nil, fmt.Errorf("argument at index %d is out of range", index)
		}
		intValue, err := strconv.Atoi(args[index])
		if err != nil {
			return nil, fmt.Errorf("argument at index %d must be a valid number", index)
		}
		intValues = append(intValues, intValue)
	}
	return intValues, nil
}
