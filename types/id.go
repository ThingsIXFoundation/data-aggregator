package types

import (
	"encoding/hex"
	"fmt"
)

type ID [32]byte

func (id ID) String() string {
	return fmt.Sprintf("0x%x", id[:])
}

func (id ID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

func (id *ID) UnmarshalText(input []byte) error {
	if len(input) != 64 && len(input) != 66 {
		return fmt.Errorf("id has invalid length")
	}

	if len(input) == 66 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X') {
		input = input[2:]
	} else if len(input) != 64 {
		return fmt.Errorf("id doesn't start with 0x")
	}

	decoded, err := hex.DecodeString(string(input))
	if err != nil {
		return fmt.Errorf("invalid id")
	}
	copy(id[:], decoded)
	return nil
}
