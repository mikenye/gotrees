package bst_test

import (
	"fmt"
	"github.com/mikenye/gotrees/bst"
	"github.com/mikenye/gotrees/rbtree"
)

func ExampleTree_Delete() {

	// create the tree with integer keys and string values
	tree := bst.New[int, string, struct{}](func(a, b int) bool {
		return a < b
	})

	// insert some nodes in the tree
	node3, _ := tree.Insert(3, "three")
	node1, _ := tree.Insert(1, "one")
	node5, _ := tree.Insert(5, "five")
	tree.Insert(0, "zero")
	tree.Insert(2, "two")
	tree.Insert(4, "four")
	node7, _ := tree.Insert(7, "seven")
	tree.Insert(6, "six")
	node9, _ := tree.Insert(9, "nine")
	tree.Insert(8, "eight")
	tree.Insert(10, "ten")

	// delete the odd nodes
	tree.Delete(node1)
	tree.Delete(node3)
	tree.Delete(node5)
	tree.Delete(node7)
	tree.Delete(node9)

	// show the tree
	fmt.Printf("Tree:\n%s", tree)

	// Output:
	// Tree:
	//       ╭── 0: zero [{}]
	//  ╭── 2: two [{}]
	// 4: four [{}]
	//  │    ╭── 6: six [{}]
	//  ╰── 8: eight [{}]
	//       ╰── 10: ten [{}]
}

func ExampleTree_Insert() {

	// create the tree with integer keys and string values
	tree := bst.New[int, string, struct{}](func(a, b int) bool {
		return a < b
	})

	// insert some nodes in the tree
	tree.Insert(3, "three")
	tree.Insert(1, "one")
	tree.Insert(5, "five")
	tree.Insert(0, "zero")
	tree.Insert(2, "two")
	tree.Insert(4, "four")
	tree.Insert(7, "seven")
	tree.Insert(6, "six")
	tree.Insert(9, "nine")
	tree.Insert(8, "eight")
	tree.Insert(10, "ten")

	// show the tree
	fmt.Printf("Tree after insert:\n%s", tree)

	// Output:
	// Tree after insert:
	//       ╭── 0: zero [{}]
	//  ╭── 1: one [{}]
	//  │    ╰── 2: two [{}]
	// 3: three [{}]
	//  │    ╭── 4: four [{}]
	//  ╰── 5: five [{}]
	//       │    ╭── 6: six [{}]
	//       ╰── 7: seven [{}]
	//            │    ╭── 8: eight [{}]
	//            ╰── 9: nine [{}]
	//                 ╰── 10: ten [{}]
}

func ExampleTree_Successor_traversal() {

	// create the tree with integer keys and string values
	tree := rbtree.New[int, string](func(a, b int) bool {
		return a < b
	})

	// insert some nodes in the tree
	tree.Insert(0, "zero")
	tree.Insert(1, "one")
	tree.Insert(2, "two")
	tree.Insert(3, "three")
	tree.Insert(4, "four")
	tree.Insert(5, "five")
	tree.Insert(6, "six")
	tree.Insert(7, "seven")
	tree.Insert(8, "eight")
	tree.Insert(9, "nine")
	tree.Insert(10, "ten")

	fmt.Println("Traversing the tree in ascending order:")

	// traverse the tree in ascending order.
	// for loop init statement:
	// `node := tree.Min(tree.Root())` sets `node` to the minimum in the tree (smallest key)
	//
	// for loop condition expression:
	// `!tree.IsNil(node)` loops while `none` is not nil.
	//
	// for loop post statement:
	// `node = tree.Successor(node)` will set the node to the in-order successor,
	// and will return the sentinel nil after the maximum in the tree
	for node := tree.Min(tree.Root()); !tree.IsNil(node); node = tree.Successor(node) {
		fmt.Printf(
			"Node with key %d has value %s (and color: %s)\n",
			tree.Key(node),
			tree.Value(node),
			tree.Metadata(node),
		)
	}

	// Output:
	// Traversing the tree in ascending order:
	// Node with key 0 has value zero (and color: ⬛)
	// Node with key 1 has value one (and color: ⬛)
	// Node with key 2 has value two (and color: ⬛)
	// Node with key 3 has value three (and color: ⬛)
	// Node with key 4 has value four (and color: ⬛)
	// Node with key 5 has value five (and color: ⬛)
	// Node with key 6 has value six (and color: ⬛)
	// Node with key 7 has value seven (and color: 🟥)
	// Node with key 8 has value eight (and color: 🟥)
	// Node with key 9 has value nine (and color: ⬛)
	// Node with key 10 has value ten (and color: 🟥)
}

func ExampleTree_Predecessor_traversal() {

	// create the tree with integer keys and string values
	tree := rbtree.New[int, string](func(a, b int) bool {
		return a < b
	})

	// insert some nodes in the tree
	tree.Insert(0, "zero")
	tree.Insert(1, "one")
	tree.Insert(2, "two")
	tree.Insert(3, "three")
	tree.Insert(4, "four")
	tree.Insert(5, "five")
	tree.Insert(6, "six")
	tree.Insert(7, "seven")
	tree.Insert(8, "eight")
	tree.Insert(9, "nine")
	tree.Insert(10, "ten")

	fmt.Println("Traversing the tree in descending order:")

	// traverse the tree in ascending order.
	// for loop init statement:
	// `node := tree.Max(tree.Root())` sets `node` to the maximum in the tree (largest key)
	//
	// for loop condition expression:
	// `!tree.IsNil(node)` loops while `none` is not nil.
	//
	// for loop post statement:
	// `node = tree.Predecessor(node)` will set the node to the in-order predecessor,
	// and will return the sentinel nil after the minimum in the tree
	for node := tree.Max(tree.Root()); !tree.IsNil(node); node = tree.Predecessor(node) {
		fmt.Printf(
			"Node with key %d has value %s (and color: %s)\n",
			tree.Key(node),
			tree.Value(node),
			tree.Metadata(node),
		)
	}

	// Output:
	// Traversing the tree in descending order:
	// Node with key 10 has value ten (and color: 🟥)
	// Node with key 9 has value nine (and color: ⬛)
	// Node with key 8 has value eight (and color: 🟥)
	// Node with key 7 has value seven (and color: 🟥)
	// Node with key 6 has value six (and color: ⬛)
	// Node with key 5 has value five (and color: ⬛)
	// Node with key 4 has value four (and color: ⬛)
	// Node with key 3 has value three (and color: ⬛)
	// Node with key 2 has value two (and color: ⬛)
	// Node with key 1 has value one (and color: ⬛)
	// Node with key 0 has value zero (and color: ⬛)
}

func ExampleTree_Floor() {
	// Create a tree with even numbers
	tree := bst.New[int, string, struct{}](func(a, b int) bool {
		return a < b
	})

	tree.Insert(2, "two")
	tree.Insert(4, "four")
	tree.Insert(6, "six")
	tree.Insert(8, "eight")
	tree.Insert(10, "ten")

	// Find floor value for various inputs

	// Exact match
	if node, found := tree.Floor(6); found {
		fmt.Printf("Floor(6) = %d: %s\n", tree.Key(node), tree.Value(node))
	} else {
		fmt.Println("Floor(6) not found")
	}

	// Floor value for an odd number (finds the next lower value)
	if node, found := tree.Floor(7); found {
		fmt.Printf("Floor(7) = %d: %s\n", tree.Key(node), tree.Value(node))
	} else {
		fmt.Println("Floor(7) not found")
	}

	// Floor value for a value smaller than the minimum
	if node, found := tree.Floor(1); found {
		fmt.Printf("Floor(1) = %d: %s\n", tree.Key(node), tree.Value(node))
	} else {
		fmt.Println("Floor(1) not found")
	}

	// Floor value for a value larger than the maximum
	if node, found := tree.Floor(12); found {
		fmt.Printf("Floor(12) = %d: %s\n", tree.Key(node), tree.Value(node))
	} else {
		fmt.Println("Floor(12) not found")
	}

	// Output:
	// Floor(6) = 6: six
	// Floor(7) = 6: six
	// Floor(1) not found
	// Floor(12) = 10: ten
}

func ExampleTree_Ceiling() {
	// Create a tree with even numbers
	tree := bst.New[int, string, struct{}](func(a, b int) bool {
		return a < b
	})

	tree.Insert(2, "two")
	tree.Insert(4, "four")
	tree.Insert(6, "six")
	tree.Insert(8, "eight")
	tree.Insert(10, "ten")

	// Find ceiling value for various inputs

	// Exact match
	if node, found := tree.Ceiling(6); found {
		fmt.Printf("Ceiling(6) = %d: %s\n", tree.Key(node), tree.Value(node))
	} else {
		fmt.Println("Ceiling(6) not found")
	}

	// Ceiling value for an odd number (finds the next higher value)
	if node, found := tree.Ceiling(5); found {
		fmt.Printf("Ceiling(5) = %d: %s\n", tree.Key(node), tree.Value(node))
	} else {
		fmt.Println("Ceiling(5) not found")
	}

	// Ceiling value for a value smaller than the minimum
	if node, found := tree.Ceiling(1); found {
		fmt.Printf("Ceiling(1) = %d: %s\n", tree.Key(node), tree.Value(node))
	} else {
		fmt.Println("Ceiling(1) not found")
	}

	// Ceiling value for a value larger than the maximum
	if node, found := tree.Ceiling(12); found {
		fmt.Printf("Ceiling(12) = %d: %s\n", tree.Key(node), tree.Value(node))
	} else {
		fmt.Println("Ceiling(12) not found")
	}

	// Output:
	// Ceiling(6) = 6: six
	// Ceiling(5) = 6: six
	// Ceiling(1) = 2: two
	// Ceiling(12) not found
}
