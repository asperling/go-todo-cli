package cli

import (
	"fmt"
	"os"
	"strconv"
)

func ValidateArgs(minLength int, intIndices []int, usage string) ([]int, error) {
	if len(os.Args) < minLength {
		return nil, fmt.Errorf("%s", usage)
	}
	var intValues []int
	for _, index := range intIndices {
		if index >= len(os.Args) || index < 0 {
			return nil, fmt.Errorf("Argument at index %d is out of range", index)
		}
		intValue, err := strconv.Atoi(os.Args[index])
		if err != nil {
			return nil, fmt.Errorf("Argument at index %d must be a valid number", index)
		}
		intValues = append(intValues, intValue)
	}
	return intValues, nil
}
