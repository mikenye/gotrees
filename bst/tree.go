// Package bst provides a generic binary search tree (BST) implementation in Go.
//
// This package allows users to create, manipulate, and traverse binary search trees
// with customizable key-value pairs and node metadata. The BST maintains order using a user-supplied
// comparison function, ensuring efficient insertions, deletions, and lookups.
//
// # Tree Behavior
//
// This implementation **does not** balance itself. If self-balancing behavior is required,
// consider using an AVL Tree or Red-Black Tree, which can be implemented by extending bst.Tree
// (e.g., for a Red-Black Tree, the node color can be stored in the node metadata).
//
// Keys must have strict weak ordering. If keys do not have a strict weak ordering, the behavior is undefined.
// Strict weak ordering means that the LessFunc function must define a consistent and transitive ordering.
// That is, for any values a, b, and c:
//
//   - If a < b and b < c, then a < c must be true.
//   - If a < b is true, then b < a must be false.
//
// # Metadata: Node-Persistent Data
//
// Each node in the BST contains an optional metadata field (M), which persists with the node across operations.
// This metadata is transferred to replacement nodes during deletions, ensuring that critical information
// (such as balancing factors, colors, or aggregated values) is retained as the tree structure evolves.
//
// Examples of metadata use cases include:
//
//   - **Red-Black Trees**: Store the node’s color (Red or Black).
//   - **AVL Trees**: Store the node’s balance factor.
//   - **Augmented Trees**: Store additional computed values (e.g., subtree sizes, sums).
//
// If metadata is not needed, `struct{}` can be used as the metadata type, ensuring zero memory overhead.
//
// # Development Notes
//
// This package was developed as part of a learning exercise to explore data structures
// and algorithmic efficiency in Go. It is designed to be extendable for different tree types
// while maintaining a clean and reusable BST foundation.
package bst

import (
	"fmt"
	"strings"
)

// These "connectors" are used for the Tree.String method when drawing the BST.
const (
	connectorLeft     = " ╭── "
	connectorRight    = " ╰── "
	connectorVertical = " │   "
	connectorSpace    = "     "
)

// LessFunc is a comparison function used to define the ordering of keys in the BST.
//
// It should return true if 'a' is less than 'b', and false otherwise.
//
// For example, in a Tree where the key type is int:
//
//	lessFunc := func(a, b int) bool { return a < b }
//
// This function must define a consistent and transitive ordering to ensure correct BST behavior.
type LessFunc[K any] func(a, b K) bool

// TraversalFunc defines a function type used for processing nodes during a tree traversal.
//
// This function receives each node in the tree and performs an operation on it.
// The traversal continues as long as the function returns true. If the function returns false,
// the traversal stops early, allowing for early exits.
//
// Parameters:
//   - node: A pointer to the current node being processed.
//
// Returns:
//   - A boolean indicating whether to continue traversal (true) or stop early (false).
type TraversalFunc[K, V, M any] func(node *Node[K, V, M]) bool

// Tree represents a generic binary search tree (BST).
//
// It stores Nodes containing key-value pairs and maintains order based on the provided
// LessFunc function, which defines how keys are compared.
//
// ⚠️Important: This implementation does not perform automatic re-balancing.
// If the tree becomes skewed (e.g., inserting keys in sorted order),
// operations will degrade to O(n) complexity.
type Tree[K, V, M any] struct {
	root *Node[K, V, M] // Root node of the tree.
	less LessFunc[K]    // Function to compare keys and maintain order.
	nil  *Node[K, V, M]
}

// New creates and returns a new empty binary search tree (BST).
//
// The Tree is initialized with a user-provided LessFunc function, which defines
// how keys are compared for ordering.
//
// The user must specify the key (K), value (V) and metadata (M) types when creating the tree.
//
// If value and/or metadata are not needed, use struct{} as the type for zero memory overhead.
//
// Parameters:
//   - less: A comparison function that determines the ordering of keys.
//
// Returns:
//   - A pointer to an empty Tree.
//
// Example Usage:
//
//	// Creating a BST with integer keys and string values.
//	tree := New[int, string, struct{}](func(a, b int) bool { return a < b })
//	tree.Insert(10, "ten")
func New[K, V, M any](less LessFunc[K]) *Tree[K, V, M] {
	t := &Tree[K, V, M]{
		less: less,
		nil:  &Node[K, V, M]{},
	}
	t.root = t.nil
	t.root.parent = t.nil
	return t
}

// keyEq is a helper function that performs an equality check between keys using the LessFunc function.
func (t *Tree[K, V, M]) keyEq(a, b K) bool {
	return !t.less(a, b) && !t.less(b, a)
}

// transplant replaces node "toReplace" with node "replacement"
func (t *Tree[K, V, M]) transplant(toReplace, replacement *Node[K, V, M]) {

	// perform transplant
	if toReplace.parent == t.nil {

		// if old node has nil parent, then it was the root
		t.root = replacement // update the root to replacement

	} else if toReplace == toReplace.parent.left {

		// if the old node is a left child, replace its parent's left child
		toReplace.parent.left = replacement // update the node to replacement

	} else {

		// else, it was a right child, replace its parent's right child
		toReplace.parent.right = replacement // update the node to replacement

	}

	// set replacement's parent to node's parent
	if replacement != t.nil {
		replacement.parent = toReplace.parent
	}
}

// Contains checks whether the given node n is present in the tree.
//
// The function searches for n's key in the tree and verifies that the
// returned node is the same as n. This ensures that the node belongs
// to this specific tree instance and is not an external or detached node.
//
// Returns:
//   - true if n is in the tree.
//   - false if n is not found or belongs to a different tree.
func (t *Tree[K, V, M]) Contains(n *Node[K, V, M]) bool {
	n2, found := t.Search(t.Key(n))
	return found && n == n2
}

// Delete removes the specified node `n` from the tree.
//
// If the deletion is successful, it returns the replacement node (if any) and true.
// If the node does not exist or is nil, it returns the tree's sentinel nil node and false.
//
// Metadata from the deleted node is transferred to its replacement node to ensure
// consistency when extending the tree (e.g., for Red-Black Trees or AVL Trees).
//
// The deletion process follows standard BST deletion rules:
//   - If n has no left child, it is replaced by its right child.
//   - If n has no right child, it is replaced by its left child.
//   - If n has two children, it is replaced by its in-order successor.
//
// Returns:
//   - (*Node[K, V, M], true) if the node was successfully deleted.
//   - (t.nil, false) if the node was not found or nil.
func (t *Tree[K, V, M]) Delete(n *Node[K, V, M]) (*Node[K, V, M], bool) {

	// if nil input, don't delete anything and give nil output
	if t.IsNil(n) || n == nil {
		return t.nil, false
	}

	if n.left == t.nil {
		replacement := n.right
		t.transplant(n, n.right)
		if replacement != t.nil {
			t.SetMetadata(replacement, n.metadata)
		}
		return replacement, true

	} else if n.right == t.nil {
		replacement := n.left
		t.transplant(n, n.left)
		if replacement != t.nil {
			t.SetMetadata(replacement, n.metadata)
		}
		return replacement, true

	} else {
		successor := t.Min(n.right)
		replacement := successor
		if successor.parent != n {
			t.transplant(successor, successor.right)
			successor.right = n.right
			successor.right.parent = successor
		}
		t.transplant(n, successor)
		successor.left = n.left
		successor.left.parent = successor
		if replacement != t.nil {
			t.SetMetadata(replacement, n.metadata)
		}
		return replacement, true
	}
}

// Depth returns the depth of node n.
//
// The depth of a node is the number of edges from the root to the node.
// The root node has a depth of 0.
//
// ⚠️ Important: This function does not validate whether node actually belongs to the tree.
// Calling it on an arbitrary node could lead to undefined behavior. See Tree.Contains.
func (t *Tree[K, V, M]) Depth(n *Node[K, V, M]) int {
	h := 0
	for !t.IsNil(n.parent) {
		h++
		n = n.parent
	}
	return h
}

// Insert inserts a new node with the given key and value into the tree.
//
// If a node with the same key already exists, its value is updated,
// and the existing node is returned with true.
//
// Otherwise, a new node is created, inserted at the appropriate position,
// and returned with false.
//
// The function maintains BST ordering:
//   - If key is less than the current node's key, it is inserted in the left subtree.
//   - If key is greater, it is inserted in the right subtree.
//   - If key already exists, its value is updated instead of creating a duplicate.
//
// Returns:
//   - (*Node[K, V, M], true) if the key existed and the value was updated.
//   - (*Node[K, V, M], false) if a new node was inserted.
func (t *Tree[K, V, M]) Insert(key K, value V) (*Node[K, V, M], bool) {

	parent := t.nil      // trailing pointer - parent of current node
	currNode := t.Root() // current node

	// find nil leaf where new node will be inserted
	for !t.IsNil(currNode) {

		// update trailing pointer
		parent = currNode

		if t.keyEq(currNode.key, key) {

			// If key already exists, update the value
			currNode.value = value
			return currNode, true

		} else if t.less(key, currNode.key) {

			// If key is smaller, go left
			currNode = currNode.left

		} else {

			// If key is larger, go right
			currNode = currNode.right
		}
	}

	// Create a new node to insert
	newNode := &Node[K, V, M]{
		key:    key,
		value:  value,
		parent: parent,
		left:   t.nil,
		right:  t.nil,
	}

	if parent == t.nil {

		// If the tree was empty, set root
		t.root = newNode

	} else if t.less(key, parent.key) {

		// if the key is less than the parent key, insert new node as left child
		parent.left = newNode

	} else {

		// if the key is greater than the parent key, insert new node as right child
		parent.right = newNode
	}

	return newNode, false
}

// IsFull returns true if the given node `n` has both left and right children.
//
// A full node is one that has exactly two children.
func (t *Tree[K, V, M]) IsFull(n *Node[K, V, M]) bool {
	return n.left != t.nil && n.right != t.nil
}

// IsInternal returns true if the given node `n` is an internal node,
// meaning it has at least one child (left or right).
//
// Internal nodes are non-leaf nodes that contribute to the tree structure.
func (t *Tree[K, V, M]) IsInternal(n *Node[K, V, M]) bool {
	return n.left != t.nil || n.right != t.nil
}

// IsLeaf returns true if the given node `n` has no children,
// meaning both its left and right pointers are nil.
//
// A leaf node is a terminal node in the tree.
func (t *Tree[K, V, M]) IsLeaf(n *Node[K, V, M]) bool {
	return n.left == t.nil && n.right == t.nil
}

// IsNil returns true if the given node `n` is the tree's sentinel nil node.
//
// The nil node is used to represent the absence of a real node in the tree.
func (t *Tree[K, V, M]) IsNil(n *Node[K, V, M]) bool {
	return n == t.nil
}

// IsUnary returns true if the given node `n` has exactly one child
// (either left or right, but not both).
//
// This is determined using a logical XOR operation on the child checks.
func (t *Tree[K, V, M]) IsUnary(n *Node[K, V, M]) bool {
	return (n.left == t.nil) != (n.right == t.nil) // Logical XOR
}

// IsTreeValid performs structural validation of the tree.
//
// This function verifies:
//   - The sentinel nil node has not been altered in a way that could break functionality.
//   - The root node’s parent is correctly set to the sentinel nil node.
//   - The tree maintains proper in-order traversal order (keys are correctly sorted).
//   - Parent-child relationships are correctly maintained.
//
// The validation is performed using an in-order traversal to ensure that
// all nodes follow the correct key ordering and structural constraints.
//
// Returns:
//   - nil if the tree is valid.
//   - An error describing the first encountered issue if the tree is invalid.
func (t *Tree[K, V, M]) IsTreeValid() error {
	// check sentinel has not been changed
	if t.nil.parent != t.nil {
		return fmt.Errorf("sentinel nil node parent not sentinel nil node")
	}

	// check root node has nil parent
	if t.root.parent != t.nil {
		return fmt.Errorf("root node parent not sentinel nil node")
	}

	// Recurse the tree in order. Check:
	//  - node keys are in order
	//  - node parent/child relationships are correct
	var (
		err              error
		currKey, prevKey K
	)
	first := true
	t.TraverseInOrder(t.root, func(node *Node[K, V, M]) bool {
		prevKey = currKey
		currKey = node.key

		if first {
			first = false
		} else {

			// if not first node, currKey should be greater than prevKey
			if !t.less(prevKey, currKey) {
				err = fmt.Errorf("traversal error: out of order keys at node: %v", node.key)
				return false
			}
		}

		// check parent/child
		parentNil := node.parent == t.nil
		leftChild := node == node.parent.left && node != node.parent.right
		rightChild := node != node.parent.left && node == node.parent.right
		if !parentNil && !(leftChild || rightChild) {
			err = fmt.Errorf("traversal error: parent/child mismatch for node: %v", node.key)
			return false
		}

		return true
	})
	if err != nil {
		return err
	}
	return nil
}

// Key returns the key of the given node n.
func (t *Tree[K, V, M]) Key(n *Node[K, V, M]) K {
	return n.key
}

// Left returns the left child of the given node n.
//
// If the node has no left child, it returns the tree's sentinel nil node.
func (t *Tree[K, V, M]) Left(n *Node[K, V, M]) *Node[K, V, M] {
	return n.left
}

// Max returns the node with the maximum key in the subtree rooted at n.
//
// This function traverses to the rightmost node of the subtree.
// If n is nil or the subtree is empty, it returns n.
func (t *Tree[K, V, M]) Max(n *Node[K, V, M]) *Node[K, V, M] {
	for n.right != nil && n.right != t.nil {
		n = n.right
	}
	return n
}

// Metadata returns the metadata associated with the given node n.
//
// The metadata field can be used to store auxiliary information such as
// node color (for Red-Black Trees), balance factor (for AVL Trees), or
// any other user-defined data.
func (t *Tree[K, V, M]) Metadata(n *Node[K, V, M]) M {
	return n.metadata
}

// Min returns the node with the minimum key in the subtree rooted at n.
//
// This function traverses to the leftmost node of the subtree.
// If n is nil or the subtree is empty, it returns n.
func (t *Tree[K, V, M]) Min(n *Node[K, V, M]) *Node[K, V, M] {
	for n.left != nil && n.left != t.nil {
		n = n.left
	}
	return n
}

// MustSetMetadata updates the metadata of the given node n.
//
// This function allows modification of the metadata field, which can be used
// for storing auxiliary information such as node color (for Red-Black Trees),
// balance factors (for AVL Trees), or other user-defined data.
//
// No nil checks are performed. If the node is nil, the function will panic.
func (t *Tree[K, V, M]) MustSetMetadata(n *Node[K, V, M], metadata M) {
	n.metadata = metadata
}

// Parent returns the parent of the given node n.
//
// If n is the root, it returns the tree's sentinel nil node.
func (t *Tree[K, V, M]) Parent(n *Node[K, V, M]) *Node[K, V, M] {
	return n.parent
}

// Predecessor returns the in-order predecessor of the given node n.
//
// The predecessor is the largest node in n's left subtree.
// If n has no left subtree, it moves up the tree until it finds a parent
// where n is in the right subtree. If no predecessor exists, it returns the sentinel nil node.
func (t *Tree[K, V, M]) Predecessor(n *Node[K, V, M]) *Node[K, V, M] {
	if n.left != t.nil {
		return t.Max(n.left)
	}
	p := n.parent
	for p != t.nil && n != p.right {
		n = p
		p = p.parent
	}
	return p
}

// Right returns the right child of the given node n.
//
// If the node has no right child, it returns the tree's sentinel nil node.
func (t *Tree[K, V, M]) Right(n *Node[K, V, M]) *Node[K, V, M] {
	return n.right
}

// Root returns the root node of the tree.
//
// If the tree is empty, it returns the sentinel nil node.
func (t *Tree[K, V, M]) Root() *Node[K, V, M] {
	return t.root
}

// RotateLeft performs a left rotation on the given node within the tree.
//
// A left rotation moves the node down while promoting its right child.
// This is commonly used in self-balancing trees (e.g., AVL or Red-Black Trees).
//
// Rotation steps:
//  1. The right child of the node becomes the new parent of the node.
//  2. The left child of the node's right subtree becomes the new right child of the node.
//  3. The node's right subtree replaces the node in the tree structure.
//
// Preconditions:
//   - The given node must have a non-nil right child (node.right != nil).
//
// ⚠️ Important: This function does not validate whether node actually belongs to the tree.
// Calling it on an arbitrary node could lead to undefined behavior.
func (t *Tree[K, V, M]) RotateLeft(node *Node[K, V, M]) {
	if node == nil || node == t.nil || node.right == t.nil {
		return // No rotation possible if node is nil or has no right child
	}

	rightSubtree := node.right
	node.right = rightSubtree.left
	if rightSubtree.left != t.nil {
		rightSubtree.left.parent = node
	}

	rightSubtree.parent = node.parent
	if node.parent == t.nil {
		t.root = rightSubtree
	} else if node.parent.left == node {
		node.parent.left = rightSubtree
	} else {
		node.parent.right = rightSubtree
	}

	rightSubtree.left, node.parent = node, rightSubtree
}

// RotateRight performs a right rotation on the given node within the tree.
//
// A right rotation moves the node down while promoting its left child.
// This is commonly used in self-balancing trees (e.g., AVL or Red-Black Trees).
//
// Rotation steps:
//  1. The left child of the node becomes the new parent of the node.
//  2. The right child of the node's left subtree becomes the new left child of the node.
//  3. The node's left subtree replaces the node in the tree structure.
//
// Preconditions:
//   - The given node must have a non-nil left child (node.left != nil).
//
// ⚠️ Important: This function does not validate whether node actually belongs to the tree.
// Calling it on an arbitrary node could lead to undefined behavior.
func (t *Tree[K, V, M]) RotateRight(node *Node[K, V, M]) {
	if node == nil || node == t.nil || node.left == t.nil {
		return // No rotation possible if node is nil or has no left child
	}

	leftSubtree := node.left
	node.left = leftSubtree.right
	if leftSubtree.right != t.nil {
		leftSubtree.right.parent = node
	}

	leftSubtree.parent = node.parent
	if node.parent == t.nil {
		t.root = leftSubtree
	} else if node.parent.left == node {
		node.parent.left = leftSubtree
	} else {
		node.parent.right = leftSubtree
	}

	leftSubtree.right, node.parent = node, leftSubtree
}

// Search looks for a node with the given key in the tree.
//
// The search follows standard BST lookup rules:
//   - If the key matches the current node, it is returned.
//   - If the key is smaller, the search continues in the left subtree.
//   - If the key is larger, the search continues in the right subtree.
//
// If the key is found, the corresponding node is returned along with true.
// If the key is not found, the tree's sentinel nil node is returned with false.
//
// Returns:
//   - (*Node[K, V, M], true) if the key exists in the tree.
//   - (*Node[K, V, M], false) if the key is not found.
func (t *Tree[K, V, M]) Search(key K) (*Node[K, V, M], bool) {
	currNode := t.root

	// if we arrive at a nil node, then node is not in tree
	for currNode != t.nil {

		// if we've found the matching node, return it
		if t.keyEq(currNode.key, key) {
			return currNode, true
		}

		// traverse the tree in the direction of key
		if t.less(key, currNode.key) {
			currNode = currNode.left
		} else {
			currNode = currNode.right
		}
	}
	return t.nil, false
}

// SetMetadata updates the metadata of the given node n.
//
// This function allows modification of the metadata field, which can be used
// for storing auxiliary information such as node color (for Red-Black Trees),
// balance factors (for AVL Trees), or other user-defined data.
func (t *Tree[K, V, M]) SetMetadata(n *Node[K, V, M], metadata M) {
	if n != nil && !t.IsNil(n) {
		n.metadata = metadata
	}
}

// Sibling returns the sibling of the given node n.
//
// If n has a parent:
//   - If n is the left child, its sibling is the right child.
//   - If n is the right child, its sibling is the left child.
//
// If n has no parent (i.e., it is the root), the sentinel nil node is returned.
//
// Returns:
//   - A pointer to the sibling node if one exists.
//   - The sentinel nil node if n is the root or has no sibling.
func (t *Tree[K, V, M]) Sibling(n *Node[K, V, M]) *Node[K, V, M] {
	if n.parent == t.nil {
		return t.nil
	}
	if n.parent.left == n {
		return n.parent.right
	}
	return n.parent.left
}

// String returns a visual representation of the binary search tree (BST).
//
// The tree is displayed in a structured format, resembling its actual shape.
// Nodes are printed with connectors indicating their relationships, making it
// easy to understand the hierarchy of the tree.
//
// The tree is ordered in ascending order, with the minimum node on the first line.
//
// The nodes are printed using the Node.String method.
//
// If the tree is empty, the function returns "Empty Tree".
//
// Returns:
//   - A formatted string representing the BST structure.
//
// This function uses an in-order iterator to traverse the tree and builds
// the output using a string builder. It tracks vertical lines dynamically
// to create a structured visualization of the BST.
func (t *Tree[K, V, M]) String() string {

	// if tree is empty, return early
	if t.root == t.nil {
		return "Empty Tree"
	}

	// prepare string builder
	builder := strings.Builder{}

	// prepare map to hold which levels to draw vertical lines
	verticalLineHeights := make(map[int]bool)

	// ascend the tree. for each node:
	t.TraverseInOrder(t.root, func(node *Node[K, V, M]) bool {
		// get height of node
		h := t.Depth(node)

		// if we are at a height that needs a vertical line, draw it,
		// otherwise draw a space
		for j := 0; j < h-1; j++ {
			if verticalLineHeights[j+1] {
				builder.WriteString(connectorVertical)
			} else {
				builder.WriteString(connectorSpace)
			}
		}

		// draw "connector" based on node orientation
		if node.parent != t.nil && node.parent.left == node {
			builder.WriteString(connectorLeft)
		} else if node.parent != t.nil && node.parent.right == node {
			builder.WriteString(connectorRight)
		}

		// print node key
		builder.WriteString(node.String())
		builder.WriteString("\n")

		// turn on/off vertical lines

		// if node parent is in the "right" direction ("down" in this representation),
		// turn on vertical lines for this height.
		if node.parent != t.nil && node.parent.left == node {
			verticalLineHeights[h] = true
		}
		// if node parent is in "left" direction ("up" in this representation),
		// turn off vertical lines for this height.
		if node.parent != t.nil && node.parent.right == node {
			verticalLineHeights[h] = false
		}
		// if node has right child ("down in this representation),
		// turn on vertical lines for the next height (h+1).
		if node.right != t.nil {
			verticalLineHeights[h+1] = true
		} else {
			verticalLineHeights[h+1] = false
		}

		return true
	})

	// return the tree
	return builder.String()
}

// Successor returns the in-order successor of the given node n.
//
// The successor is the smallest node that is greater than n in the tree.
//   - If n has a right subtree, the successor is the leftmost node in that subtree.
//   - Otherwise, the function moves up the tree until it finds a parent
//     where n is in the left subtree. That parent is the successor.
//
// If no successor exists, the sentinel nil node is returned.
//
// Returns:
//   - A pointer to the successor node if one exists.
//   - The sentinel nil node if n has no successor.
func (t *Tree[K, V, M]) Successor(n *Node[K, V, M]) *Node[K, V, M] {
	if n.right != t.nil {
		return t.Min(n.right)
	}
	p := n.parent
	for p != t.nil && n != p.left {
		n = p
		p = p.parent
	}
	return p
}

// TraverseInOrder performs an in-order traversal of the tree starting from node n.
//
// TraverseInOrder uses recursion. If the tree is deep and highly unbalanced, this could lead to stack overflow.
// Consider using Tree.Successor and Tree.Predecessor in these cases.
//
// The traversal order is:
//  1. Recursively visit the left subtree.
//  2. Process the current node.
//  3. Recursively visit the right subtree.
//
// The function applies the user-provided function f to each visited node.
// If f returns false, the traversal stops early.
//
// Returns:
//   - true if the traversal completes successfully.
//   - false if f returns false, causing an early exit.
func (t *Tree[K, V, M]) TraverseInOrder(n *Node[K, V, M], f TraversalFunc[K, V, M]) bool {

	// Recurse the left children of n
	if n.left != nil && n.left != t.nil && !t.TraverseInOrder(n.left, f) {
		return false
	}

	// Process n
	if !f(n) {
		return false
	}

	// Recurse the right children of n
	if n.right != nil && n.right != t.nil && !t.TraverseInOrder(n.right, f) {
		return false
	}

	// Continue traversing
	return true
}

// Value returns the value associated with the given node n.
//
// This function retrieves the stored value for the node's key.
func (t *Tree[K, V, M]) Value(n *Node[K, V, M]) V {
	return n.value
}
