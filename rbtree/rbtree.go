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

func New[K, V any](less bst.LessFunc[K]) *Tree[K, V] {
	t := &Tree[K, V]{
		Tree: bst.New[K, V, Color](less),
	}
	t.Tree.MustSetMetadata(t.Root(), Black) // set sentinel nil to black
	return t
}

func (t *Tree[K, V]) Insert(key K, value V) (*bst.Node[K, V, Color], bool) {
	n, updated := t.Tree.Insert(key, value)
	if updated {
		return n, updated
	}
	t.Tree.SetMetadata(n, Red)

	// fixup
	z := n
	for t.isRed(t.Parent(z)) {
		if t.Parent(z) == t.Left(t.Parent(t.Parent(z))) { // if z's parent a left child?
			y := t.Right(t.Parent(t.Parent(z))) // y is z's uncle
			if t.isRed(y) {                     // are z's parent and uncle both red?

				// case 1
				t.Tree.SetMetadata(t.Parent(z), Black)
				t.Tree.SetMetadata(y, Black)
				t.Tree.SetMetadata(t.Parent(t.Tree.Parent(z)), Red)
				z = t.Parent(t.Parent(z))

			} else {
				if z == t.Right(t.Parent(z)) { // is z is a right child?

					// case 2
					z = t.Parent(z)
					t.Tree.RotateLeft(z)
				}

				// case 3
				t.Tree.SetMetadata(t.Parent(z), Black)
				t.Tree.SetMetadata(t.Parent(t.Parent(z)), Red)
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

func (t *Tree[K, V]) MustSetMetadata(_ *bst.Node[K, V, Color], _ Color) {
	panic(fmt.Errorf("MustSetMetadata should not be called on an rbtree, doing so may corrupt the tree"))
}

func (t *Tree[K, V]) SetMetadata(_ *bst.Node[K, V, Color], _ Color) {
	panic(fmt.Errorf("SetMetadata should not be called on an rbtree, doing so may corrupt the tree"))
}

func (t *Tree[K, V]) RotateLeft(_ *bst.Node[K, V, Color]) {
	panic(fmt.Errorf("RotateLeft should not be called on an rbtree, doing so may corrupt the tree"))
}

func (t *Tree[K, V]) RotateRight(_ *bst.Node[K, V, Color]) {
	panic(fmt.Errorf("RotateRight should not be called on an rbtree, doing so may corrupt the tree"))
}
