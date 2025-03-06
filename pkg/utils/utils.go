package utils

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/ed25519"
	"github.com/pactus-project/pactus/util/bech32m"
	"github.com/pagu-project/pagu/internal/entity"
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

// IsDefinedOnBotID checks if any of the given bot IDs is defined on the target bot ID.
func IsDefinedOnBotID(botIDs []entity.BotID, target entity.BotID) bool {
	for _, botID := range botIDs {
		if botID == target {
			return true
		}
	}

	return false
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

// TestnetPrivateKeyFromString parses a testnet private key and returns the Private Key object.
// Note that for Testnet Private Keys, the HRP (Human-Readable Part) is set to `TSECRET1`,
// which differs from the Pactus Mainnet where the HRP is set to `SECRET1`.
// This function is a workaround to parse testnet Private Keys alongside mainnet Private Keys.
func TestnetPrivateKeyFromString(text string) (*ed25519.PrivateKey, error) {
	// Decode the bech32m encoded private key.
	hrp, typ, data, err := bech32m.DecodeToBase256WithTypeNoLimit(text)
	if err != nil {
		return nil, err
	}

	// Check if hrp is valid
	if hrp != "tsecret" {
		return nil, crypto.InvalidHRPError(hrp)
	}

	if typ != crypto.SignatureTypeEd25519 {
		return nil, crypto.InvalidSignatureTypeError(typ)
	}

	return ed25519.PrivateKeyFromBytes(data)
}

// TestnetAddressFromString parses a testnet address and returns the Address object.
// Note that for Testnet addresses, the HRP (Human-Readable Part) is set to `tpc1`,
// which differs from the Pactus Mainnet where the HRP is set to `pc1`.
// This function is a workaround to parse testnet addresses alongside mainnet addresses.
func TestnetAddressFromString(text string) (crypto.Address, error) {
	// Decode the bech32m encoded address.
	hrp, typ, data, err := bech32m.DecodeToBase256WithTypeNoLimit(text)
	if err != nil {
		return crypto.Address{}, err
	}

	// Check if hrp is valid
	if hrp != "tpc" {
		return crypto.Address{}, crypto.InvalidHRPError(hrp)
	}

	// check type is valid
	validTypes := []crypto.AddressType{
		crypto.AddressTypeValidator,
		crypto.AddressTypeBLSAccount,
		crypto.AddressTypeEd25519Account,
	}
	if !slices.Contains(validTypes, crypto.AddressType(typ)) {
		return crypto.Address{}, crypto.InvalidAddressTypeError(typ)
	}

	// check length is valid
	if len(data) != 20 {
		return crypto.Address{}, crypto.InvalidLengthError(len(data) + 1)
	}

	var addr crypto.Address
	addr[0] = typ
	copy(addr[1:], data)

	return addr, nil
}

func TestnetAddressToString(addr crypto.Address) string {
	str, _ := bech32m.EncodeFromBase256WithType(
		"tpc",
		addr[0],
		addr[1:])

	return str
}

// FormatDuration formats the duration into a human-readable string.
func FormatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	if hours > 0 {
		return fmt.Sprintf("%d hours and %d minutes", hours, minutes)
	}

	return fmt.Sprintf("%d minutes", minutes)
}

// IsToggleEnabled converts a toggle-like string to a boolean value.
// It returns true for "true", "yes", "on" and "1",
// otherwise tt returns false.
func IsToggleEnabled(toggleStr string) bool {
	toggleStr = strings.ToLower(toggleStr)
	return toggleStr == "true" || toggleStr == "yes" || toggleStr == "on" || toggleStr == "1"
}
