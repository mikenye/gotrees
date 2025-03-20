// Package rbtree provides a generic, self-balancing Red-Black Binary Search Tree implementation.
//
// This package extends bst.Tree, adding automatic balancing by ensuring that:
//   - The tree remains approximately balanced, maintaining O(log n) insertions, deletions, and lookups.
//   - No two consecutive red nodes appear in a path.
//   - All paths from the root to leaves contain the same number of black nodes.
//
// # Key Features
//   - Self-Balancing: Uses Red-Black Tree rules to maintain efficiency.
//   - Generic Support: Works with any key (K) and value (V) types.
//
// # Usage Example
//
//	import "github.com/mikenye/gotrees/rbtree"
//
//	tree := rbtree.New[int, string](func(a, b int) bool { return a < b })
//	tree.Insert(10, "ten")
//	tree.Insert(20, "twenty")
//	node, found := tree.Search(10)
//
//	if found {
//		tree.Delete(node)
//	}
//
// # Safe Inherited Methods from bst.Tree
//
// The following methods are inherited from bst.Tree and can be used safely:
//   - [bst.Tree.Root]: Returns the root node.
//   - [bst.Tree.Search]: Finds a node by key.
//   - [bst.Tree.Successor]: Returns the next in-order node.
//   - [bst.Tree.Predecessor]: Returns the previous in-order node.
//   - [bst.Tree.TraverseInOrder]: In-order traversal.
//   - [bst.Tree.Min]: Returns the node with the smallest key.
//   - [bst.Tree.Max]: Returns the node with the largest key.
//   - [bst.Tree.Floor]: Returns the largest node with key ≤ given key.
//   - [bst.Tree.Ceiling]: Returns the smallest node with key ≥ given key.
//   - [bst.Tree.IsNil]: Checks if a node is the sentinel nil node.
//   - [bst.Tree.Parent]: Returns the parent of a node.
//
// # Unsafe Inherited Methods from bst.Tree
//
// The following methods from bst.Tree should not be used in rbtree, as they can violate Red-Black properties.
// They have been shadowed in rbtree, and modified to panic if used:
//
//   - [bst.Tree.MustSetMetadata]: ❌ Do not use
//   - [bst.Tree.SetKey]: ❌ Do not use
//   - [bst.Tree.SetLeft]: ❌ Do not use
//   - [bst.Tree.SetMetadata]: ❌ Do not use
//   - [bst.Tree.SetParent]: ❌ Do not use
//   - [bst.Tree.SetRight]: ❌ Do not use
//   - [bst.Tree.SetRoot]: ❌ Do not use
//   - [bst.Tree.Transplant]: ❌ Do not use
//
// ⚠️ Warning: Using any of these methods will likely break the Red-Black properties and cause undefined behavior.
//
// # Limitations
//
//   - Not Thread-Safe – Requires external synchronization for concurrent use.
//   - No Duplicate Keys – Keys must be unique.
package rbtree

import (
	"fmt"
	"github.com/mikenye/gotrees/bst"
)

// Color represents the color of a node in a Red-Black Tree.
//
// Nodes are either:
//   - Red (🟥), indicates a temporary imbalance during insertion/deletion.
//   - Black (⬛), maintains tree balancing properties.
//
// This enum is used to ensure Red-Black Tree properties are enforced correctly.
type Color bool

const (
	Red   Color = false // Red-colored node
	Black Color = true  // Black-colored node
)

// String returns a Unicode representation of the node color.
//
// Nodes are either:
//   - Red: function will return "🟥"
//   - Black: function will return "⬛"
func (c Color) String() string {
	if c == Black {
		return "⬛"
	} else {
		return "🟥"
	}
}

// Tree represents a Red-Black Tree, an extension of bst.Tree that maintains self-balancing properties.
//
// This tree ensures:
//   - O(log n) insertions, deletions, and lookups.
//   - Automatic re-balancing using the Red-Black Tree rules.
//   - Strict BST ordering with an additional node metadata Color for balancing.
//
// The tree embeds a generic Binary Search Tree bst.Tree, using Color as metadata
// to track whether a node is `Red` or `Black`. The `size` field keeps track of the total
// number of nodes.
type Tree[K, V any] struct {
	*bst.Tree[K, V, Color]     // Underlying BST structure
	size                   int // Total number of nodes
}

// isBlack returns true if the passed node is black or nil (nil leaves are considered black)
func (t *Tree[K, V]) isBlack(n *bst.Node[K, V, Color]) bool {
	if t.IsNil(n) || t.Metadata(n) != Red {
		return true
	}
	return false
}

// isRed returns true if the passed node is not nil and red
func (t *Tree[K, V]) isRed(n *bst.Node[K, V, Color]) bool {
	if !t.IsNil(n) && t.Metadata(n) == Red {
		return true
	}
	return false
}

// setColor sets the color of node n, if node n is not the sentinel nil node
func (t *Tree[K, V]) setColor(n *bst.Node[K, V, Color], c Color) {
	if !t.IsNil(n) {
		t.Tree.SetMetadata(n, c)
	}
}

// Delete removes the given node z from the Red-Black Tree while maintaining tree balance.
//
// Deleting a node modifies tree structure and may trigger rotation/recoloring
// to maintain Red-Black Tree properties.
func (t *Tree[K, V]) Delete(z *bst.Node[K, V, Color]) bool {
	// if nil input, don't delete anything and give nil output
	if t.IsNil(z) || z == nil {
		return false
	}

	var x, y *bst.Node[K, V, Color]

	// if node being deleted has one child
	if t.IsNil(t.Left(z)) || t.IsNil(t.Right(z)) {
		y = z // deletion case 1
	} else {
		y = t.Successor(z) // deletion case 2
	}

	if !t.IsNil(t.Left(y)) {
		// if node being deleted has left child, set x to left child
		x = t.Left(y)
	} else {
		// otherwise, set x to right child
		x = t.Right(y)
	}

	// update replacement node's parent
	t.Tree.SetParent(x, t.Parent(y))
	if t.IsNil(t.Parent(y)) {
		// if replacement has no parent, it becomes root
		t.SetRoot(x)
	} else {
		// update parent/child relationships
		if y == t.Left(t.Parent(y)) {
			// if y is a left child
			t.Tree.SetLeft(t.Parent(y), x)
		} else {
			// if y is a right child
			t.Tree.SetRight(t.Parent(y), x)
		}
	}
	if y != z {
		// copy y’s satellite data into z
		t.Tree.SetKey(z, t.Key(y))
		t.Tree.SetValue(z, t.Value(y))
	}

	// fixup
	if t.isBlack(y) {
		t.deleteFixup(x)
	}
	t.resetSentinel()
	t.size--
	return true
}

// deleteFixup restores Red-Black Tree properties after a node deletion.
//
// After deletion, the Red-Black Tree may violate one or more of the following properties:
// - Property 1: The root is always black.
// - Property 4: Red nodes cannot have red children.
// - Property 5: Every path from the root to a leaf must have the same number of black nodes.
//
// This function fixes violations by applying **four fixup cases**:
//
// Cases 1-4 (Fixing Double Black Issues)
// 1. Sibling is red: Perform rotation and recoloring.
// 2. Sibling and its children are black: Recolor sibling and move problem up the tree.
// 3. Sibling has one red child (far side is black): Rotate sibling and recolor.
// 4. Sibling has one red child (near side is red): Rotate parent, recolor, and fix final issues.
//
// The function proceeds iteratively, moving up the tree until balance is restored.
func (t *Tree[K, V]) deleteFixup(x *bst.Node[K, V, Color]) {
	for x != t.Root() && t.isBlack(x) {
		if x == t.Left(t.Parent(x)) { // is x a left child?
			w := t.Right(t.Parent(x))
			if t.isRed(w) {

				// case 1
				t.setColor(w, Black)
				t.setColor(t.Parent(x), Red)
				t.Tree.RotateLeft(t.Parent(x))
				w = t.Right(t.Parent(x))

			}
			if t.isBlack(t.Left(w)) && t.isBlack(t.Right(w)) {

				// case 2
				t.setColor(w, Red)
				x = t.Parent(x)
				//t.Tree.SetParent(x, t.Parent(t.Parent(z)))

			} else {

				if t.isBlack(t.Right(w)) {

					// case 3
					t.setColor(t.Left(w), Black)
					t.setColor(w, Red)
					t.Tree.RotateRight(w)
					w = t.Right(t.Parent(x))
				}

				// case 4
				t.setColor(w, t.Metadata(t.Parent(x)))
				t.setColor(t.Parent(x), Black)
				t.setColor(t.Right(w), Black)
				t.Tree.RotateLeft(t.Parent(x))
				x = t.Root()
			}
		} else {

			// same as above but with right and left exchanged

			w := t.Left(t.Parent(x))
			if t.isRed(w) {

				// case 1
				t.setColor(w, Black)
				t.setColor(t.Parent(x), Red)
				t.Tree.RotateRight(t.Parent(x))
				w = t.Left(t.Parent(x))

			}
			if t.isBlack(t.Right(w)) && t.isBlack(t.Left(w)) {

				// case 2
				t.setColor(w, Red)
				x = t.Parent(x)
				//t.Tree.SetParent(x, t.Parent(t.Parent(z)))

			} else {

				if t.isBlack(t.Left(w)) {

					// case 3
					t.setColor(t.Right(w), Black)
					t.setColor(w, Red)
					t.Tree.RotateLeft(w)
					w = t.Left(t.Parent(x))
				}

				// case 4
				t.setColor(w, t.Metadata(t.Parent(x)))
				t.setColor(t.Parent(x), Black)
				t.setColor(t.Left(w), Black)
				t.Tree.RotateRight(t.Parent(x))
				x = t.Root()
			}
		}
	}
	t.setColor(x, Black)
}

// Insert adds a new key-value pair to the Red-Black Tree while maintaining self-balancing properties.
//
//   - If the key already exists, its value is updated, and no fixup is needed.
//   - If the key is new, the node is inserted colored red, and the tree undergoes fixup rotations/recoloring
//     to maintain Red-Black Tree properties.
//
// Returns:
//   - The inserted or updated node.
//   - true if a new node was inserted, false if an existing node was updated.
func (t *Tree[K, V]) Insert(key K, value V) (*bst.Node[K, V, Color], bool) {
	n, updated := t.Tree.Insert(key, value)
	if !updated {
		return n, false
	}
	t.setColor(n, Red)

	// Fixup after insertion
	t.insertFixup(n)

	t.size++
	return n, true
}

// insertFixup performs recoloring/rotation of the red-black tree after an insertion takes place
//
// Red-Black Fixup Cases
// After inserting a red node, the tree may violate the Red-Black properties. The following cases
// are applied iteratively until balance is restored:
//
//  1. Parent and uncle are red: Recolor and move up the tree.
//  2. Parent is red, uncle is black, and inserted node is a right child: Rotate left.
//  3. Parent is red, uncle is black, and inserted node is a left child: Rotate right.
//
// The function also ensures that the root always remains black after insertion.
func (t *Tree[K, V]) insertFixup(z *bst.Node[K, V, Color]) {
	for t.isRed(t.Parent(z)) {
		if t.Parent(z) == t.Left(t.Parent(t.Parent(z))) { // If z's parent is a left child
			y := t.Right(t.Parent(t.Parent(z))) // y is z's uncle
			if t.isRed(y) {                     // Case 1: Parent & Uncle are Red
				t.setColor(t.Parent(z), Black)
				t.setColor(y, Black)
				t.setColor(t.Parent(t.Parent(z)), Red)
				z = t.Parent(t.Parent(z))
			} else {
				if z == t.Right(t.Parent(z)) { // Case 2: z is a right child
					z = t.Parent(z)
					t.Tree.RotateLeft(z)
				}
				// Case 3: z is a left child
				t.setColor(t.Parent(z), Black)
				t.setColor(t.Parent(t.Parent(z)), Red)
				t.Tree.RotateRight(t.Parent(t.Parent(z)))
			}
		} else {
			// Mirror the logic with left/right swapped
			y := t.Left(t.Parent(t.Parent(z)))
			if t.isRed(y) {
				t.setColor(t.Parent(z), Black)
				t.setColor(y, Black)
				t.setColor(t.Parent(t.Parent(z)), Red)
				z = t.Parent(t.Parent(z))
			} else {
				if z == t.Left(t.Parent(z)) {
					z = t.Parent(z)
					t.Tree.RotateRight(z)
				}
				t.setColor(t.Parent(z), Black)
				t.setColor(t.Parent(t.Parent(z)), Red)
				t.Tree.RotateLeft(t.Parent(t.Parent(z)))
			}
		}
	}
	t.setColor(t.Root(), Black)
}

// IsTreeValid verifies whether the Red-Black Tree maintains all BST and Red-Black properties.
//
// The function first validates the underlying BST structure, then applies the Red-Black Tree checks.
//
// This function checks the following five Red-Black Tree invariants:
//  1. Every node is either red or black: Enforced by using `Color` as metadata, which by its nature can only have Red (false) or Black (true).
//  2. The root is always black: If the root is red, the tree is invalid.
//  3. Every leaf (sentinel nil node) is black: Ensures correct tree termination.
//  4. Red nodes cannot have red children: Prevents consecutive red nodes (ensures balancing).
//  5. All paths from a node to its descendant leaves must have the same number of black nodes.
//
// Returns:
//   - nil if the tree is valid; or:
//   - An error describing the first detected violation if the tree is invalid.
func (t *Tree[K, V]) IsTreeValid() error {
	var err error

	// check underlying BST
	err = t.Tree.IsTreeValid()
	if err != nil {
		return fmt.Errorf("underlying BST is invalid: %v", err)
	}

	// check the red-black tree invariants
	// invariant 1: every node is either red or black.
	// this invariant is enforced due to t.Tree's M being type Color.

	// invariant 2: the root is black
	if !t.isBlack(t.Root()) {
		return fmt.Errorf("root node is not black")
	}

	// invariant 3: Every leaf (nil sentinel) is black.
	if t.Metadata(t.Parent(t.Root())) != Black {
		return fmt.Errorf("sentinel nil node is not black")
	}

	firstLeaf := true
	blackCount := 0

	t.TraverseInOrder(t.Root(), func(n *bst.Node[K, V, Color]) bool {

		// invariant 4: if a node is red, then both its children are black
		if t.isRed(n) && t.isRed(t.Left(n)) {
			err = fmt.Errorf("node %v is red and has red left child", t.Key(n))
			return false
		}
		if t.isRed(n) && t.isRed(t.Right(n)) {
			err = fmt.Errorf("node %v is red and has red right child", t.Key(n))
			return false
		}

		// invariant 5: For each node, all simple paths from the node to descendant
		// leaves contain the same number of black nodes.
		if !(t.IsLeaf(n) || t.IsUnary(n)) {
			return true // skip this check if not a leaf node
		}
		bc := 0
		for !t.IsNil(n) {
			if t.isBlack(n) {
				bc++
			}
			n = t.Parent(n)
		}
		if firstLeaf {
			blackCount = bc
			firstLeaf = false
			return true
		}
		if bc != blackCount {
			err = fmt.Errorf("node %v has black count mismatch", t.Key(n))
			return false
		}
		return true
	})
	if err != nil {
		return err
	}
	return nil
}

// Deprecated: Should not be called on an rbtree.Tree, doing so may corrupt the tree.
func (t *Tree[K, V]) MustSetMetadata() {
	panic(fmt.Errorf("MustSetMetadata should not be called on an rbtree.Tree, doing so may corrupt the tree"))
}

// resetSentinel re-initializes the sentinel nil node to maintain Red-Black Tree invariants.
//
// In a Red-Black Tree, the sentinel node serves as a placeholder for all nil references.
//
// This function ensures that the sentinel node:
//   - Has no left or right children (nil pointers).
//   - Has itself as its parent (ensuring a valid reference).
//   - Is always Black (as required by Red-Black Tree rules).
//
// This function should be called after deletions to maintain consistency.
func (t *Tree[K, V]) resetSentinel() {
	t.Tree.SetLeft(t.Sentinel(), nil)
	t.Tree.SetRight(t.Sentinel(), nil)
	t.Tree.SetParent(t.Sentinel(), t.Sentinel())
	t.setColor(t.Sentinel(), Black)
}

// Deprecated: Should not be called on an rbtree.Tree, doing so may corrupt the tree.
func (t *Tree[K, V]) RotateLeft() {
	panic(fmt.Errorf("RotateLeft should not be called on an rbtree.Tree, doing so may corrupt the tree"))
}

// Deprecated: Should not be called on an rbtree.Tree, doing so may corrupt the tree.
func (t *Tree[K, V]) RotateRight() {
	panic(fmt.Errorf("RotateRight should not be called on an rbtree.Tree, doing so may corrupt the tree"))
}

// Deprecated: Should not be called on an rbtree.Tree, doing so may corrupt the tree.
func (t *Tree[K, V]) SetLeft() {
	panic(fmt.Errorf("SetLeft should not be called on an rbtree.Tree, doing so may corrupt the tree"))
}

// Deprecated: Should not be called on an rbtree.Tree, doing so may corrupt the tree.
func (t *Tree[K, V]) SetMetadata() {
	panic(fmt.Errorf("SetMetadata should not be called on an rbtree.Tree, doing so may corrupt the tree"))
}

// Deprecated: Should not be called on an rbtree.Tree, doing so may corrupt the tree.
func (t *Tree[K, V]) SetParent() {
	panic(fmt.Errorf("SetParent should not be called on an rbtree.Tree, doing so may corrupt the tree"))
}

// Deprecated: Should not be called on an rbtree.Tree, doing so may corrupt the tree.
func (t *Tree[K, V]) SetRight() {
	panic(fmt.Errorf("SetRight should not be called on an rbtree.Tree, doing so may corrupt the tree"))
}

// Size returns the total number of nodes in the Red-Black Tree.
//
// This function provides an O(1) operation to retrieve the node count, which is
// maintained dynamically during insertions and deletions.
//
// Returns:
//   - The number of nodes currently stored in the tree.
func (t *Tree[K, V]) Size() int {
	return t.size
}

// Deprecated: Should not be called on an rbtree.Tree, doing so may corrupt the tree.
func (t *Tree[K, V]) Transplant() {
	panic(fmt.Errorf("Transplant should not be called on an rbtree.Tree, doing so may corrupt the tree"))
}

// New creates a new Red-Black Tree with the given key comparison function.
//
// This function initializes a self-balancing Red-Black Tree, which maintains
// O(log n) complexity for insertions, deletions, and lookups.
//
// Generics:
//   - K (Key type): Determines the ordering of nodes. Must be strictly ordered
//     via the provided bst.LessFunc function.
//   - V (Value type): The associated value stored in each node. If no value is needed,
//     struct{} can be used for zero memory overhead.
//
// Parameters:
//   - less: A comparison function (bst.LessFunc[K]) that defines the ordering of keys.
//
// Behavior:
//   - Initializes an empty Red-Black Tree.
//   - Uses the provided less function to maintain BST ordering.
//   - Ensures the sentinel nil node is properly initialized as black.
//
// Returns:
//   - A pointer to a newly created Tree[K, V] instance.
func New[K, V any](less bst.LessFunc[K]) *Tree[K, V] {
	t := &Tree[K, V]{
		Tree: bst.New[K, V, Color](less),
	}
	t.Tree.MustSetMetadata(t.Root(), Black) // set sentinel nil to black
	return t
}
