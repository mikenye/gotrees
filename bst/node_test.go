package bst

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNode_String(t *testing.T) {
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
