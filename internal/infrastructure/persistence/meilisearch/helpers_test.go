package meilisearch

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonStr_Empty(t *testing.T) {
	assert.Equal(t, "", jsonStr(nil))
	assert.Equal(t, "", jsonStr(json.RawMessage{}))
}

func TestJsonStr_QuotedString(t *testing.T) {
	raw := json.RawMessage(`"hello world"`)
	assert.Equal(t, "hello world", jsonStr(raw))
}

func TestJsonStr_UnquotedString(t *testing.T) {
	raw := json.RawMessage(`plain`)
	assert.Equal(t, "plain", jsonStr(raw))
}

func TestJsonStr_EscapedQuotes(t *testing.T) {
	raw := json.RawMessage(`"say \"hi\""`)
	assert.Equal(t, `say "hi"`, jsonStr(raw))
}

func TestJsonInt_Empty(t *testing.T) {
	assert.Equal(t, int64(0), jsonInt(nil))
	assert.Equal(t, int64(0), jsonInt(json.RawMessage{}))
}

func TestJsonInt_PositiveNumber(t *testing.T) {
	raw := json.RawMessage(`42`)
	assert.Equal(t, int64(42), jsonInt(raw))
}

func TestJsonInt_LargeNumber(t *testing.T) {
	raw := json.RawMessage(`9999999`)
	assert.Equal(t, int64(9999999), jsonInt(raw))
}

func TestJsonInt_InvalidInput(t *testing.T) {
	raw := json.RawMessage(`abc`)
	assert.Equal(t, int64(0), jsonInt(raw))
}

func TestJsonInt_QuotedNumber(t *testing.T) {
	raw := json.RawMessage(`"123"`)
	assert.Equal(t, int64(0), jsonInt(raw))
}

func TestJsonInt_NegativeNumber(t *testing.T) {
	raw := json.RawMessage(`-5`)
	assert.Equal(t, int64(-5), jsonInt(raw))
}
