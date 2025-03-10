package rbtree

import (
	"fmt"
	"github.com/mikenye/gotrees/bst"
)

type Color bool

const (
	Red   Color = false
	Black Color = true
)

func (c Color) String() string {
	if c == Black {
		return "⬛"
	} else {
		return "🟥"
	}
}

type Tree[K, V any] struct {
	*bst.Tree[K, V, Color]
}

func (t *Tree[K, V]) isBlack(n *bst.Node[K, V, Color]) bool {
	if t.IsNil(n) || t.Metadata(n) != Red {
		return true
	}
	return false
}

func (t *Tree[K, V]) isRed(n *bst.Node[K, V, Color]) bool {
	if !t.IsNil(n) && t.Metadata(n) == Red {
		return true
	}
	return false
}

func (t *Tree[K, V]) setColor(n *bst.Node[K, V, Color], c Color) {
	if !t.IsNil(n) {
		t.Tree.SetMetadata(n, c)
	}
}

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
	return true
}

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

func (t *Tree[K, V]) Insert(key K, value V) (*bst.Node[K, V, Color], bool) {
	n, updated := t.Tree.Insert(key, value)
	if updated {
		return n, updated
	}
	t.setColor(n, Red)

	// fixup
	z := n // z is the inserted node, and will be reassigned during fixup
	for t.isRed(t.Parent(z)) {
		if t.Parent(z) == t.Left(t.Parent(t.Parent(z))) { // if z's parent a left child?
			y := t.Right(t.Parent(t.Parent(z))) // y is z's uncle
			if t.isRed(y) {                     // are z's parent and uncle both red?

				// case 1
				t.setColor(t.Parent(z), Black)
				t.setColor(y, Black)
				t.setColor(t.Parent(t.Tree.Parent(z)), Red)
				z = t.Parent(t.Parent(z))

			} else {
				if z == t.Right(t.Parent(z)) { // is z is a right child?

					// case 2
					z = t.Parent(z)
					t.Tree.RotateLeft(z)
				}

				// case 3
				t.setColor(t.Parent(z), Black)
				t.setColor(t.Parent(t.Parent(z)), Red)
				t.Tree.RotateRight(t.Parent(t.Parent(z)))
			}
		} else {

			// as above, but with right & left exchanged

			y := t.Left(t.Parent(t.Parent(z))) // y is z's uncle
			if t.isRed(y) {                    // are z's parent and uncle both red?

				// case 1
				t.setColor(t.Parent(z), Black)
				t.setColor(y, Black)
				t.setColor(t.Parent(t.Parent(z)), Red)
				z = t.Parent(t.Parent(z))

			} else {
				if z == t.Left(t.Parent(z)) { // is z is a right child?

					// case 2
					z = t.Parent(z)
					t.Tree.RotateRight(z)
				}

				// case 3
				t.setColor(t.Parent(z), Black)
				t.setColor(t.Parent(t.Parent(z)), Red)
				t.Tree.RotateLeft(t.Parent(t.Parent(z)))
			}

		}
	}
	t.setColor(t.Root(), Black)
	return n, true
}

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

// Deprecated: Should not be called on an rbtree.Tree, doing so may corrupt the tree.
func (t *Tree[K, V]) Transplant() {
	panic(fmt.Errorf("Transplant should not be called on an rbtree.Tree, doing so may corrupt the tree"))
}

func New[K, V any](less bst.LessFunc[K]) *Tree[K, V] {
	t := &Tree[K, V]{
		Tree: bst.New[K, V, Color](less),
	}
	t.Tree.MustSetMetadata(t.Root(), Black) // set sentinel nil to black
	return t
}
