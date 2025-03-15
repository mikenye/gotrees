package bst

import (
	"fmt"
	"strings"
)

// Node represents a single element within the binary search tree (BST).
//
// Each node stores a key-value pair and maintains references to its parent
// and child nodes, allowing for hierarchical structuring within the tree.
//
// The BST maintains its structure based on the ordering function defined
// when the tree is created, ensuring efficient search, insertion, and deletion operations.
type Node[K, V, M any] struct {
	key                 K
	value               V
	parent, left, right *Node[K, V, M]
	metadata            M
}

// String returns a string representation of the node.
//
// The output format is "key: value [metadata]", where both key and value
// are converted to strings. If the key or value implements fmt.Stringer,
// its String() method is used; otherwise, fmt.Sprintf is used.
// Metadata is only included if the metadata type implements fmt.Stringer.
//
// Returns:
//   - A string representation of the node in "key: value [metadata]" format.
func (n *Node[K, V, M]) String() string {
	builder := new(strings.Builder)

	// write node key
	if s, ok := any(n.key).(fmt.Stringer); ok {
		builder.WriteString(s.String())
	} else {
		builder.WriteString(fmt.Sprintf("%v", n.key))
	}

	// separator between node & value
	builder.WriteString(": ")

	// write node value
	if any(n.value) != nil {
		if s, ok := any(n.value).(fmt.Stringer); ok {
			builder.WriteString(s.String())
		} else {
			builder.WriteString(fmt.Sprintf("%v", n.value))
		}
	} else {
		builder.WriteString("<nil>")
	}

	// write node metadata
	builder.WriteString(" [")
	if s, ok := any(n.metadata).(fmt.Stringer); ok {
		builder.WriteString(s.String())
	} else {
		builder.WriteString(fmt.Sprintf("%v", n.metadata))
	}
	builder.WriteString("]")

	return builder.String()
}
