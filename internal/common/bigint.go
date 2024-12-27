package common

import (
	"fmt"
	"math/big"
	"strings"
)

type BigInt big.Int

// MarshalJSON converts *big.Int to a JSON string
func (b *BigInt) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", (*big.Int)(b).String())), nil
}

// UnmarshalJSON parses a JSON string into *big.Int
func (b *BigInt) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), "\"")
	bigInt, ok := new(big.Int).SetString(str, 10)
	if !ok {
		return fmt.Errorf("invalid big.Int format: %s", str)
	}
	*b = BigInt(*bigInt)
	return nil
}
