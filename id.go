package platform

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// IDLength is the exact length a byte slice must have in order to be decoded into an ID
const IDLength = 16

// ID is a unique identifier.
//
// Its zero value is not a valid ID.
type ID uint64

// IDGenerator represents a generator for IDs.
type IDGenerator interface {
	// ID creates unique byte slice ID.
	ID() ID
}

// IDFromString creates an ID from a given string.
//
// It errors if the input string does not match a valid ID.
func IDFromString(str string) (*ID, error) {
	var id ID
	err := id.DecodeFromString(str)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

// Decode parses b as a hex-encoded byte-slice-string.
//
// It errors if the input byte slice does not have the correct length
// or if it contains all zeros.
func (i *ID) Decode(b []byte) error {
	if len(b) != IDLength {
		return fmt.Errorf("input must be an array of %d bytes", IDLength)
	}
	if bytes.Equal(b, make([]byte, IDLength)) {
		return fmt.Errorf("all 0s is not a valid ID")
	}

	dst := make([]byte, hex.DecodedLen(IDLength))
	_, err := hex.Decode(dst, b)
	if err != nil {
		return err
	}
	*i = ID(binary.BigEndian.Uint64(dst))
	return nil
}

// DecodeFromString parses s as a hex-encoded string.
func (i *ID) DecodeFromString(s string) error {
	return i.Decode([]byte(s))
}

// Encode converts ID to a hex-encoded byte-slice-string.
//
// It errors if the receiving ID holds its zero value.
func (i ID) Encode() ([]byte, error) {
	if i == 0 {
		return nil, fmt.Errorf("all 0s is not a valid ID")
	}

	b := make([]byte, hex.DecodedLen(IDLength))
	binary.BigEndian.PutUint64(b, uint64(i))

	dst := make([]byte, hex.EncodedLen(len(b)))
	hex.Encode(dst, b)
	return dst, nil
}

// Valid checks whether the receiving ID is a valid one or not.
func (i ID) Valid() bool {
	if _, err := i.Encode(); err != nil {
		return false
	}
	return true
}

// String returns the ID as a hex encoded string.
//
// Returns an empty string in the case the ID is invalid.
func (i ID) String() string {
	enc, _ := i.Encode()
	return string(enc)
}

// UnmarshalJSON implements JSON unmarshaller for IDs.
func (i *ID) UnmarshalJSON(b []byte) error {
	b = b[1 : len(b)-1]
	return i.Decode(b)
}

// MarshalJSON implements JSON marshaller for IDs.
func (i ID) MarshalJSON() ([]byte, error) {
	enc, err := i.Encode()
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(enc))
}
