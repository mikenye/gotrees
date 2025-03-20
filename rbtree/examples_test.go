package rbtree_test

import (
	"fmt"
	"github.com/mikenye/gotrees/rbtree"
)

func ExampleTree_Insert() {

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

	// show the tree
	fmt.Printf("Red-Black Tree after insert:\n%s", tree)

	// Output:
	// Red-Black Tree after insert:
	//       ╭── 0: zero [⬛]
	//  ╭── 1: one [⬛]
	//  │    ╰── 2: two [⬛]
	// 3: three [⬛]
	//  │    ╭── 4: four [⬛]
	//  ╰── 5: five [⬛]
	//       │    ╭── 6: six [⬛]
	//       ╰── 7: seven [🟥]
	//            │    ╭── 8: eight [🟥]
	//            ╰── 9: nine [⬛]
	//                 ╰── 10: ten [🟥]
}

func ExampleTree_Delete() {

	// create the tree with integer keys and string values
	tree := rbtree.New[int, string](func(a, b int) bool {
		return a < b
	})

	// insert some nodes in the tree
	tree.Insert(0, "zero")
	node1, _ := tree.Insert(1, "one")
	tree.Insert(2, "two")
	node3, _ := tree.Insert(3, "three")
	tree.Insert(4, "four")
	node5, _ := tree.Insert(5, "five")
	tree.Insert(6, "six")
	node7, _ := tree.Insert(7, "seven")
	tree.Insert(8, "eight")
	node9, _ := tree.Insert(9, "nine")
	tree.Insert(10, "ten")

	// delete the odd nodes
	tree.Delete(node1)
	tree.Delete(node3)
	tree.Delete(node5)
	tree.Delete(node7)
	tree.Delete(node9)

	// show the tree
	fmt.Printf("Red-Black Tree:\n%s", tree)

	// Output:
	// Red-Black Tree:
	//       ╭── 0: zero [⬛]
	//  ╭── 2: two [🟥]
	//  │    ╰── 4: four [⬛]
	// 6: six [⬛]
	//  │    ╭── 8: eight [🟥]
	//  ╰── 10: ten [⬛]
}

func ExampleTree_Floor_and_Ceiling() {
	// Create a red-black tree with even numbers
	tree := rbtree.New[int, string](func(a, b int) bool {
		return a < b
	})

	tree.Insert(2, "two")
	tree.Insert(4, "four")
	tree.Insert(6, "six")
	tree.Insert(8, "eight")
	tree.Insert(10, "ten")

	// Using inherited Floor and Ceiling methods from bst.Tree

	// Find the closest values to 5
	floorNode, floorFound := tree.Floor(5)
	ceilingNode, ceilingFound := tree.Ceiling(5)

	if floorFound {
		fmt.Printf("Floor(5) = %d: %s\n", tree.Key(floorNode), tree.Value(floorNode))
	} else {
		fmt.Println("Floor(5) not found")
	}

	if ceilingFound {
		fmt.Printf("Ceiling(5) = %d: %s\n", tree.Key(ceilingNode), tree.Value(ceilingNode))
	} else {
		fmt.Println("Ceiling(5) not found")
	}

	// Using Floor and Ceiling to implement a range query for keys between 3 and 7
	for node, found := tree.Ceiling(3); found && tree.Key(node) <= 7; node = tree.Successor(node) {
		fmt.Printf("Key in range [3,7]: %d\n", tree.Key(node))
	}

	// Output:
	// Floor(5) = 4: four
	// Ceiling(5) = 6: six
	// Key in range [3,7]: 4
	// Key in range [3,7]: 6
}
