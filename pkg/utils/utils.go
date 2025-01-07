package utils

import (
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
