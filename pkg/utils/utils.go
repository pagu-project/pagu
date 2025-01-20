package utils

import (
	"fmt"
	"os"

	"golang.org/x/exp/constraints"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return err == nil
}

// SetFlag applies mask to the flags.
func SetFlag[T constraints.Integer](flags, mask T) T {
	return flags | mask
}

// UnsetFlag removes mask from the flags.
func UnsetFlag[T constraints.Integer](flags, mask T) T {
	return flags & ^mask
}

// IsFlagSet checks if the mask is set for the given flags.
func IsFlagSet[T constraints.Integer](flags, mask T) bool {
	return flags&mask == mask
}

// MarshalEnum serializes an enum value into its string representation using the provided `toString` map.
// Returns an error if the value does not have a corresponding string.
func MarshalEnum[T comparable](value T, toString map[T]string) (string, error) {
	str, ok := toString[value]
	if !ok {
		return "", fmt.Errorf("unknown enum value: %v", value)
	}

	return str, nil
}

// UnmarshalEnum deserializes a string into an enum value using the provided `toString` map.
// Returns an error if the string does not match any known enum value.
func UnmarshalEnum[T comparable](str string, toString map[T]string) (T, error) {
	for key, val := range toString {
		if val == str {
			return key, nil
		}
	}

	var zero T

	return zero, fmt.Errorf("unknown enum type: %s", str)
}
