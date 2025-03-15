package bst

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestNode_String_stringer tests that where key, value and metadata types implement fmt.Stringer,
// their string representation is correctly printed.
func TestNode_String_stringer(t *testing.T) {
	d := time.Date(2006, 01, 02, 03, 04, 05, 00, time.UTC)
	n := &Node[time.Time, time.Time, time.Time]{
		key:      d,
		value:    d,
		metadata: d,
	}
	assert.Equal(t,
		"2006-01-02 03:04:05 +0000 UTC: 2006-01-02 03:04:05 +0000 UTC [2006-01-02 03:04:05 +0000 UTC]",
		n.String())
}

func TestNode_String_nil(t *testing.T) {
	n := &Node[int, *time.Time, struct{}]{
		key:      1,
		value:    nil,
		metadata: struct{}{},
	}
	assert.Equal(t,
		"1: <nil> [{}]",
		n.String())
}
