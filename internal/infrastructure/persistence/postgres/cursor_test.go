package postgres

import (
	"encoding/base64"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeCursor_Empty(t *testing.T) {
	assert.Equal(t, int64(0), decodeCursor(""))
}

func TestDecodeCursor_InvalidBase64(t *testing.T) {
	assert.Equal(t, int64(0), decodeCursor("!!!invalid!!!"))
}

func TestDecodeCursor_Valid(t *testing.T) {
	encoded := base64.StdEncoding.EncodeToString([]byte("100"))
	assert.Equal(t, int64(100), decodeCursor(encoded))
}

func TestDecodeCursor_InvalidNumber(t *testing.T) {
	encoded := base64.StdEncoding.EncodeToString([]byte("not-a-number"))
	assert.Equal(t, int64(0), decodeCursor(encoded))
}

func TestEncodeCursor(t *testing.T) {
	encoded := encodeCursor(50)
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	assert.NoError(t, err)
	num, err := strconv.ParseInt(string(decoded), 10, 64)
	assert.NoError(t, err)
	assert.Equal(t, int64(50), num)
}

func TestEncodeDecodeCursor_Roundtrip(t *testing.T) {
	for _, offset := range []int64{0, 20, 100, 999999} {
		encoded := encodeCursor(offset)
		decoded := decodeCursor(encoded)
		assert.Equal(t, offset, decoded)
	}
}
