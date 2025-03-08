package bst_test

import (
	"fmt"
	"github.com/mikenye/gotrees/bst"
)

func ExampleNew() {
	// Create a BST with:
	//   - Key (K) type of int
	//   - Value (V) type of string
	//   - No metadata (empty struct)
	tree := bst.New[int, string, struct{}](func(a, b int) bool { return a < b })
	fmt.Println(tree)

	// Output:
	// Empty Tree
}
