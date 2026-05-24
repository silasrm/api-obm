package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildFilterParams_DefaultLimit(t *testing.T) {
	fp := buildFilterParams(0, "")
	assert.Equal(t, 20, fp.Limit)
	assert.Equal(t, "", fp.Cursor)
}

func TestBuildFilterParams_NegativeLimit(t *testing.T) {
	fp := buildFilterParams(-5, "")
	assert.Equal(t, 20, fp.Limit)
}

func TestBuildFilterParams_ValidLimit(t *testing.T) {
	fp := buildFilterParams(50, "cursor123")
	assert.Equal(t, 50, fp.Limit)
	assert.Equal(t, "cursor123", fp.Cursor)
}
