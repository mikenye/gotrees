package bst

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew(t *testing.T) {
	tree := New[int, struct{}, struct{}](func(a, b int) bool {
		return a < b
	})
	assert.NoError(t, tree.IsTreeValid(), "expected valid tree")
	assert.True(t, tree.IsNil(tree.Root()), "expected new tree to have nil root")
	assert.True(t, tree.IsNil(tree.Parent(tree.Root())), "expected tree root to have nil parent")
}

func TestTree_Insert(t *testing.T) {
	tree := New[int, int, int](func(a, b int) bool {
		return a < b
	})

	// insert unique keys
	keys := []int{12, 5, 2, 9, 18, 15, 19, 13, 17, 20}
	for _, key := range keys {
		node, inserted := tree.Insert(key, key)
		assert.True(t, inserted, "expected inserted to be true when inserting unique nodes")
		tree.SetMetadata(node, key)
		assert.Equal(t, key, tree.Key(node), "expected added node's key to match")
		assert.Equal(t, key, tree.Value(node), "expected added node's value to match")
		assert.Equal(t, key, tree.Metadata(node), "expected added node's metadata to match")
	}

	t.Logf("tree after insert:\n%s", tree)

	// expect tree to be valid
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")

	// update duplicate key
	node, inserted := tree.Insert(15, 1515)
	assert.False(t, inserted, "expected inserted to be false when inserting duplicate node")
	assert.Equal(t, 15, tree.Key(node), "expected added node's key to match")
	assert.Equal(t, 1515, tree.Value(node), "expected added node's value to match")
	assert.Equal(t, 15, tree.Metadata(node), "expected added node's metadata to match")

	t.Logf("tree after update:\n%s", tree)

	// expect tree to be valid
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")

	// check structure is completely correct

	root := tree.Root() // root should be node 12
	assert.Equal(t, 12, tree.Key(root), "expected root node key of 12")
	assert.True(t, tree.IsNil(tree.Parent(root)), "expected root node parent to be nil")
	assert.True(t, tree.IsFull(root), "root should be full node")
	assert.True(t, tree.IsInternal(root), "root should be internal node")
	assert.False(t, tree.IsLeaf(root), "root should not be leaf")
	assert.False(t, tree.IsNil(root), "root should not be nil")
	assert.False(t, tree.IsUnary(root), "root should not be unary")

	n5 := tree.Left(tree.Root()) // node 5 should be left child of root (12)
	assert.Equal(t, 5, tree.Key(n5), "expected node 5 to be left child of root (12)")
	assert.Equal(t, tree.Root(), tree.Parent(n5), "expected parent of node 5 to be root (12)")
	assert.True(t, tree.IsFull(n5), "n5 should be full node")
	assert.True(t, tree.IsInternal(n5), "n5 should be internal node")
	assert.False(t, tree.IsLeaf(n5), "n5 should not be leaf")
	assert.False(t, tree.IsNil(n5), "n5 should not be nil")
	assert.False(t, tree.IsUnary(n5), "n5 should not be unary")

	n2 := tree.Left(n5) // node 2 should be left child of 5
	assert.Equal(t, 2, tree.Key(n2), "expected node 2 to be left child of node 5")
	assert.Equal(t, n5, tree.Parent(n2), "expected parent of node 2 to be node 5")
	assert.False(t, tree.IsFull(n2), "n2 should not be full node")
	assert.False(t, tree.IsInternal(n2), "n2 should not be internal node")
	assert.True(t, tree.IsLeaf(n2), "n2 should be leaf")
	assert.False(t, tree.IsNil(n2), "n2 should not be nil")
	assert.False(t, tree.IsUnary(n2), "n2 should not be unary")
	assert.True(t, tree.IsNil(tree.Left(n2)), "n2 left child should be nil")
	assert.True(t, tree.IsNil(tree.Right(n2)), "n2 right child should be nil")

	n9 := tree.Right(n5) // node 9 should be right child of 5
	assert.Equal(t, 9, tree.Key(n9), "expected node 9 to be right child of node 5")
	assert.Equal(t, n5, tree.Parent(n9), "expected parent of node 9 to be node 5")
	assert.False(t, tree.IsFull(n9), "n9 should not be full node")
	assert.False(t, tree.IsInternal(n9), "n9 should not be internal node")
	assert.True(t, tree.IsLeaf(n9), "n9 should be leaf")
	assert.False(t, tree.IsNil(n9), "n9 should not be nil")
	assert.False(t, tree.IsUnary(n9), "n9 should not be unary")
	assert.True(t, tree.IsNil(tree.Left(n9)), "n9 left child should be nil")
	assert.True(t, tree.IsNil(tree.Right(n9)), "n9 right child should be nil")

	n18 := tree.Right(root) // node 18 should be right child of root (12)
	assert.Equal(t, 18, tree.Key(n18), "expected node 18 to be right child of root (12)")
	assert.Equal(t, root, tree.Parent(n18), "expected parent of node 18 to be root (12)")
	assert.True(t, tree.IsFull(n18), "n18 should be full node")
	assert.True(t, tree.IsInternal(n18), "n18 should be internal node")
	assert.False(t, tree.IsLeaf(n18), "n18 should not be leaf")
	assert.False(t, tree.IsNil(n18), "n18 should not be nil")
	assert.False(t, tree.IsUnary(n18), "n18 should not be unary")

	n15 := tree.Left(n18) // node 15 should be left child of 18
	assert.Equal(t, 15, tree.Key(n15), "expected node 15 to be left child of node 18")
	assert.Equal(t, n18, tree.Parent(n15), "expected parent of node 15 to be node 18")
	assert.True(t, tree.IsFull(n15), "n15 should be full node")
	assert.True(t, tree.IsInternal(n15), "n15 should be internal node")
	assert.False(t, tree.IsLeaf(n15), "n15 should not be leaf")
	assert.False(t, tree.IsNil(n15), "n15 should not be nil")
	assert.False(t, tree.IsUnary(n15), "n15 should not be unary")

	n13 := tree.Left(n15) // node 13 should be left child of 15
	assert.Equal(t, 13, tree.Key(n13), "expected node 13 to be left child of node 15")
	assert.Equal(t, n15, tree.Parent(n13), "expected parent of node 13 to be node 15")
	assert.False(t, tree.IsFull(n13), "n13 should not be full node")
	assert.False(t, tree.IsInternal(n13), "n13 should not be internal node")
	assert.True(t, tree.IsLeaf(n13), "n13 should be leaf")
	assert.False(t, tree.IsNil(n13), "n13 should not be nil")
	assert.False(t, tree.IsUnary(n13), "n13 should not be unary")
	assert.True(t, tree.IsNil(tree.Left(n13)), "n13 left child should be nil")
	assert.True(t, tree.IsNil(tree.Right(n13)), "n13 right child should be nil")

	n17 := tree.Right(n15) // node 17 should be right child of 15
	assert.Equal(t, 17, tree.Key(n17), "expected node 17 to be right child of node 15")
	assert.Equal(t, n15, tree.Parent(n17), "expected parent of node 17 to be node 15")
	assert.False(t, tree.IsFull(n17), "n17 should not be full node")
	assert.False(t, tree.IsInternal(n17), "n17 should not be internal node")
	assert.True(t, tree.IsLeaf(n17), "n17 should be leaf")
	assert.False(t, tree.IsNil(n17), "n17 should not be nil")
	assert.False(t, tree.IsUnary(n17), "n17 should not be unary")
	assert.True(t, tree.IsNil(tree.Left(n17)), "n17 left child should be nil")
	assert.True(t, tree.IsNil(tree.Right(n17)), "n17 right child should be nil")

	n19 := tree.Right(n18) // node 19 should be right child of 18
	assert.Equal(t, 19, tree.Key(n19), "expected node 19 to be right child of node 18")
	assert.Equal(t, n18, tree.Parent(n19), "expected parent of node 19 to be node 18")
	assert.False(t, tree.IsFull(n19), "n19 should not be full node")
	assert.True(t, tree.IsInternal(n19), "n19 should be internal node")
	assert.False(t, tree.IsLeaf(n19), "n19 should not be leaf")
	assert.False(t, tree.IsNil(n19), "n19 should not be nil")
	assert.True(t, tree.IsUnary(n19), "n19 should be unary")
	assert.True(t, tree.IsNil(tree.Left(n19)), "n19 left child should be nil")

	n20 := tree.Right(n19) // node 20 should be right child of 19
	assert.Equal(t, 20, tree.Key(n20), "expected node 20 to be right child of node 19")
	assert.Equal(t, n19, tree.Parent(n20), "expected parent of node 20 to be node 19")
	assert.False(t, tree.IsFull(n20), "n20 should not be full node")
	assert.False(t, tree.IsInternal(n20), "n20 should not be internal node")
	assert.True(t, tree.IsLeaf(n20), "n20 should be leaf")
	assert.False(t, tree.IsNil(n20), "n20 should not be nil")
	assert.False(t, tree.IsUnary(n20), "n20 should not be unary")
	assert.True(t, tree.IsNil(tree.Left(n20)), "n20 left child should be nil")
	assert.True(t, tree.IsNil(tree.Right(n20)), "n20 right child should be nil")
}

func TestTree_Min(t *testing.T) {
	tree := New[int, struct{}, struct{}](func(a, b int) bool {
		return a < b
	})
	tree.Insert(100, struct{}{})
	tree.Insert(50, struct{}{})
	tree.Insert(10, struct{}{})
	tree.Insert(20, struct{}{})
	tree.Insert(65, struct{}{})
	n150, _ := tree.Insert(150, struct{}{})
	tree.Insert(125, struct{}{})
	t.Logf("tree after creation:\n%s", tree)

	// expect tree to be valid
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")

	// check minimum in entire tree (from root) is 10
	n := tree.Min(tree.Root())
	assert.Equal(t, 10, tree.Key(n), "unexpected minimum from root")

	// check minimum from node with key = 150 is 125
	n = tree.Min(n150)
	assert.Equal(t, 125, tree.Key(n), "unexpected minimum from node 150")
}

func TestTree_Max(t *testing.T) {
	tree := New[int, struct{}, struct{}](func(a, b int) bool {
		return a < b
	})
	tree.Insert(100, struct{}{})
	n50, _ := tree.Insert(50, struct{}{})
	tree.Insert(10, struct{}{})
	tree.Insert(20, struct{}{})
	tree.Insert(65, struct{}{})
	tree.Insert(150, struct{}{})
	tree.Insert(125, struct{}{})
	t.Logf("tree after creation:\n%s", tree)

	// expect tree to be valid
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")

	// check maximum in entire tree (from root) is 150
	n := tree.Max(tree.Root())
	assert.Equal(t, 150, tree.Key(n), "unexpected minimum from root")

	// check maximum from node with key = 50 is 65
	n = tree.Max(n50)
	assert.Equal(t, 65, tree.Key(n), "unexpected minimum from node 50")
}

func TestTree_Delete(t *testing.T) {
	tests := map[string]struct {
		creation func() *Tree[int, string, struct{}]
		deletion func(*Tree[int, string, struct{}])
		checks   func(*Tree[int, string, struct{}])
	}{
		"nil node": {
			creation: func() *Tree[int, string, struct{}] {
				// create right-leaning tree
				tree := New[int, string, struct{}](func(a, b int) bool {
					return a < b
				})
				tree.Insert(20, "z")
				tree.Insert(10, "l")
				tree.Insert(30, "y")
				return tree
			},
			deletion: func(tree *Tree[int, string, struct{}]) {
				_, deleted := tree.Delete(nil)
				require.False(t, deleted, "expected nil node to not be deleted")
			},
			checks: func(tree *Tree[int, string, struct{}]) {
				assert.Equal(t, tree.nil, tree.Parent(tree.Root()), "unexpected structure after delete")
				assert.Equal(t, 20, tree.Key(tree.Root()), "unexpected structure after delete")
				assert.Equal(t, 10, tree.Key(tree.Left(tree.Root())), "unexpected structure after delete")
				assert.Equal(t, 30, tree.Key(tree.Right(tree.Root())), "unexpected structure after delete")
			},
		},
		"node is root and has no left child": {
			creation: func() *Tree[int, string, struct{}] {
				// create right-leaning tree
				tree := New[int, string, struct{}](func(a, b int) bool {
					return a < b
				})
				tree.Insert(10, "z")
				tree.Insert(20, "r")
				return tree
			},
			deletion: func(tree *Tree[int, string, struct{}]) {
				// delete root node
				n, found := tree.Search(10)
				require.True(t, found, "expected to find node to be deleted")
				_, deleted := tree.Delete(n)
				require.True(t, deleted, "expected node to be deleted")
			},
			checks: func(tree *Tree[int, string, struct{}]) {
				// ensure structure as expected
				assert.True(t, tree.IsNil(tree.Parent(tree.Root())), "expected root node parent to be nil")
				assert.True(t, tree.IsNil(tree.Left(tree.Root())), "expected root left child to be nil")
				assert.True(t, tree.IsNil(tree.Right(tree.Root())), "expected root right child to be nil")
				// ensure Transplant worked as expected
				assert.Equal(t, 20, tree.Key(tree.Root()), "unexpected root node key after deletion")
				assert.Equal(t, "r", tree.Value(tree.Root()), "unexpected root node value after deletion")
			},
		},
		"node is root and has a left child but no right child": {
			creation: func() *Tree[int, string, struct{}] {
				tree := New[int, string, struct{}](func(a, b int) bool {
					return a < b
				})
				tree.Insert(10, "z")
				tree.Insert(5, "l")
				return tree
			},
			deletion: func(tree *Tree[int, string, struct{}]) {
				n, found := tree.Search(10)
				require.True(t, found, "expected to find node to be deleted")
				_, deleted := tree.Delete(n)
				require.True(t, deleted, "expected node to be deleted")
			},
			checks: func(tree *Tree[int, string, struct{}]) {
				// ensure structure as expected
				assert.True(t, tree.IsNil(tree.Parent(tree.Root())), "expected root node parent to be nil")
				assert.True(t, tree.IsNil(tree.Left(tree.Root())), "expected root left child to be nil")
				assert.True(t, tree.IsNil(tree.Right(tree.Root())), "expected root right child to be nil")
				// ensure Transplant worked as expected
				assert.Equal(t, 5, tree.Key(tree.Root()), "unexpected root node key after deletion")
				assert.Equal(t, "l", tree.Value(tree.Root()), "unexpected root node value after deletion")
			},
		},
		"node is root and has two children, successor has right child": {
			creation: func() *Tree[int, string, struct{}] {
				tree := New[int, string, struct{}](func(a, b int) bool {
					return a < b
				})
				tree.Insert(20, "z")
				tree.Insert(10, "l")
				tree.Insert(30, "y")
				tree.Insert(40, "x")
				return tree
			},
			deletion: func(tree *Tree[int, string, struct{}]) {
				n, found := tree.Search(20)
				require.True(t, found, "expected to find node to be deleted")
				_, deleted := tree.Delete(n)
				require.True(t, deleted, "expected node to be deleted")
			},
			checks: func(tree *Tree[int, string, struct{}]) {
				// ensure structure as expected
				assert.True(t, tree.IsNil(tree.Parent(tree.Root())), "expected root node parent to be nil")
				assert.False(t, tree.IsNil(tree.Left(tree.Root())), "expected root left child to be non-nil")
				assert.False(t, tree.IsNil(tree.Right(tree.Root())), "expected root right child to be non-nil")
				// ensure transplants worked as expected
				assert.Equal(t, 30, tree.Key(tree.Root()), "unexpected root node key after deletion")
				assert.Equal(t, "y", tree.Value(tree.Root()), "unexpected root node key after deletion")
				// ensure structure is as expected - left child of root should be 10/l
				assert.Equal(t, 10, tree.Key(tree.Left(tree.Root())), "unexpected root left child key after deletion")
				assert.Equal(t, "l", tree.Value(tree.Left(tree.Root())), "unexpected root left child value after deletion")
				assert.Equal(t, tree.Root(), tree.Parent(tree.Left(tree.Root())), "expected parent of root's left child node to be root")
				// ensure structure is as expected - right child of root should be 40/x
				assert.Equal(t, 40, tree.Key(tree.Right(tree.Root())), "unexpected root right child key after deletion")
				assert.Equal(t, "x", tree.Value(tree.Right(tree.Root())), "unexpected root right child value after deletion")
				assert.Equal(t, tree.Root(), tree.Parent(tree.Right(tree.Root())), "expected parent of root's right child node to be root")
			},
		},
		"node is root and has two children, successor has left child": {
			creation: func() *Tree[int, string, struct{}] {
				tree := New[int, string, struct{}](func(a, b int) bool {
					return a < b
				})
				tree.Insert(20, "z")
				tree.Insert(10, "l")
				tree.Insert(30, "r")
				tree.Insert(25, "y")
				tree.Insert(27, "x")
				return tree
			},
			deletion: func(tree *Tree[int, string, struct{}]) {
				n, found := tree.Search(20)
				require.True(t, found, "expected to find node to be deleted")
				_, deleted := tree.Delete(n)
				require.True(t, deleted, "expected node to be deleted")
			},
			checks: func(tree *Tree[int, string, struct{}]) {
				// ensure structure as expected
				assert.True(t, tree.IsNil(tree.Parent(tree.Root())), "expected root node parent to be nil")
				assert.False(t, tree.IsNil(tree.Left(tree.Root())), "expected root left child to be non-nil")
				assert.False(t, tree.IsNil(tree.Right(tree.Root())), "expected root right child to be non-nil")
				// ensure transplants worked as expected
				assert.Equal(t, 25, tree.Key(tree.Root()), "unexpected root node key after deletion")
				assert.Equal(t, "y", tree.Value(tree.Root()), "unexpected root node key after deletion")
				// ensure structure is as expected - left child of root should be 10/l
				assert.Equal(t, 10, tree.Key(tree.Left(tree.Root())), "unexpected root left child key after deletion")
				assert.Equal(t, "l", tree.Value(tree.Left(tree.Root())), "unexpected root left child value after deletion")
				assert.Equal(t, tree.Root(), tree.Parent(tree.Left(tree.Root())), "expected parent of root's left child node to be root")
				// ensure structure is as expected - right child of root should be 30/r
				assert.Equal(t, 30, tree.Key(tree.Right(tree.Root())), "unexpected root right child key after deletion")
				assert.Equal(t, "r", tree.Value(tree.Right(tree.Root())), "unexpected root right child value after deletion")
				assert.Equal(t, tree.Root(), tree.Parent(tree.Right(tree.Root())), "expected parent of root's right child node to be root")
				// ensure structure is as expected - left child of root's right child
				assert.Equal(t, 27, tree.Key(tree.Left(tree.Right(tree.Root()))), "unexpected structure after deletion")
				assert.Equal(t, "x", tree.Value(tree.Left(tree.Right(tree.Root()))), "unexpected structure after deletion")
				assert.Equal(t, tree.Right(tree.Root()), tree.Parent(tree.Left(tree.Right(tree.Root()))), "unexpected structure after deletion")
			},
		},
		"node is right child of its parent": {
			creation: func() *Tree[int, string, struct{}] {
				tree := New[int, string, struct{}](func(a, b int) bool {
					return a < b
				})
				tree.Insert(10, "root")
				tree.Insert(20, "right")
				tree.Insert(30, "right-right")
				return tree
			},
			deletion: func(tree *Tree[int, string, struct{}]) {
				n, found := tree.Search(20)
				require.True(t, found, "expected to find node to be deleted")
				_, deleted := tree.Delete(n)
				require.True(t, deleted, "expected node to be deleted")
			},
			checks: func(tree *Tree[int, string, struct{}]) {
				assert.Equal(t, 10, tree.Key(tree.Root()), "unexpected root node key after deletion")
				assert.Equal(t, 30, tree.Key(tree.Right(tree.Root())), "unexpected right child key after deletion")
				assert.True(t, tree.IsNil(tree.Left(tree.Root())), "expected left child to be nil")
				assert.True(t, tree.IsNil(tree.Right(tree.Right(tree.Root()))), "expected right-right child to be nil")
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tree := tc.creation()
			t.Logf("tree after creation:\n%s", tree)

			// expect tree to be valid
			require.NoError(t, tree.IsTreeValid(), "expected valid tree")

			tc.deletion(tree)
			t.Logf("tree after deletion:\n%s", tree)

			// expect tree to be valid
			require.NoError(t, tree.IsTreeValid(), "expected valid tree")

			tc.checks(tree)
		})
	}
}

func TestTree_RotateLeft_root(t *testing.T) {

	tree := New[int, string, struct{}](func(a, b int) bool {
		return a < b
	})
	x, _ := tree.Insert(100, "x")
	a, _ := tree.Insert(50, "a")
	y, _ := tree.Insert(200, "y")
	b, _ := tree.Insert(150, "b")
	c, _ := tree.Insert(250, "c")
	t.Logf("tree after creation:\n%s", tree)

	// expect tree to be valid
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")

	tree.RotateLeft(x)

	t.Logf("tree after rotation:\n%s", tree)

	// expect tree to be valid
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")

	assert.Equal(t, y, tree.Root(), "expected node y to be new root")
	assert.Equal(t, c, tree.Right(y), "expected node c to be root's right child")
	assert.True(t, tree.IsLeaf(c), "expected node c to be leaf node")
	assert.Equal(t, x, tree.Left(y), "expected x to be root's left child")
	assert.Equal(t, a, tree.Left(x), "expected a to be x's left child")
	assert.True(t, tree.IsLeaf(a), "expected node a to be leaf node")
	assert.Equal(t, b, tree.Right(x), "expected b to be x's right child")
	assert.True(t, tree.IsLeaf(b), "expected node b to be leaf node")
}

func TestTree_RotateLeft_leftchild(t *testing.T) {

	tree := New[int, string, struct{}](func(a, b int) bool {
		return a < b
	})
	r, _ := tree.Insert(500, "root")
	x, _ := tree.Insert(250, "x")
	a, _ := tree.Insert(200, "a")
	y, _ := tree.Insert(300, "y")
	b, _ := tree.Insert(299, "b")
	c, _ := tree.Insert(301, "c")
	t.Logf("tree after creation:\n%s", tree)

	// expect tree to be valid
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")

	tree.RotateLeft(x)

	t.Logf("tree after rotation:\n%s", tree)

	// expect tree to be valid
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")

	assert.Equal(t, r, tree.Root(), "root should be unchanged")
	assert.Equal(t, y, tree.Left(r), "expected node y to be new root")
	assert.Equal(t, c, tree.Right(y), "expected node c to be root's right child")
	assert.True(t, tree.IsLeaf(c), "expected node c to be leaf node")
	assert.Equal(t, x, tree.Left(y), "expected x to be root's left child")
	assert.Equal(t, a, tree.Left(x), "expected a to be x's left child")
	assert.True(t, tree.IsLeaf(a), "expected node a to be leaf node")
	assert.Equal(t, b, tree.Right(x), "expected b to be x's right child")
	assert.True(t, tree.IsLeaf(b), "expected node b to be leaf node")
}

func TestTree_RotateLeft_rightchild(t *testing.T) {

	tree := New[int, string, struct{}](func(a, b int) bool {
		return a < b
	})
	r, _ := tree.Insert(0, "root")
	x, _ := tree.Insert(250, "x")
	a, _ := tree.Insert(200, "a")
	y, _ := tree.Insert(300, "y")
	b, _ := tree.Insert(299, "b")
	c, _ := tree.Insert(301, "c")
	t.Logf("tree after creation:\n%s", tree)

	// expect tree to be valid
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")

	tree.RotateLeft(x)

	t.Logf("tree after rotation:\n%s", tree)

	// expect tree to be valid
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")

	assert.Equal(t, r, tree.Root(), "root should be unchanged")
	assert.Equal(t, y, tree.Right(r), "expected node y to be new root")
	assert.Equal(t, c, tree.Right(y), "expected node c to be root's right child")
	assert.True(t, tree.IsLeaf(c), "expected node c to be leaf node")
	assert.Equal(t, x, tree.Left(y), "expected x to be root's left child")
	assert.Equal(t, a, tree.Left(x), "expected a to be x's left child")
	assert.True(t, tree.IsLeaf(a), "expected node a to be leaf node")
	assert.Equal(t, b, tree.Right(x), "expected b to be x's right child")
	assert.True(t, tree.IsLeaf(b), "expected node b to be leaf node")
}

func TestTree_RotateLeft_nil(t *testing.T) {
	tree := New[int, string, struct{}](func(a, b int) bool {
		return a < b
	})
	root, _ := tree.Insert(100, "root")
	lc, _ := tree.Insert(50, "left child")
	rc, _ := tree.Insert(150, "right child")

	t.Logf("tree after creation:\n%s", tree)

	// expect tree to be valid
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")

	tree.RotateLeft(nil)

	t.Logf("tree after rotation:\n%s", tree)

	// expect tree to be valid
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")

	assert.Equal(t, root, tree.Root(), "tree should be unchanged")
	assert.Equal(t, lc, tree.Left(root), "expected node lc to be left child of root")
	assert.Equal(t, rc, tree.Right(root), "expected node rc to be right child of root")
	assert.True(t, tree.IsLeaf(lc), "expected node lc to be leaf node")
	assert.True(t, tree.IsLeaf(rc), "expected node rc to be leaf node")
}

func TestTree_RotateRight_root(t *testing.T) {
	tree := New[int, string, struct{}](func(a, b int) bool {
		return a < b
	})
	y, _ := tree.Insert(100, "y")
	c, _ := tree.Insert(200, "c")
	x, _ := tree.Insert(50, "x")
	b, _ := tree.Insert(75, "b")
	a, _ := tree.Insert(25, "a")
	t.Logf("tree after creation:\n%s", tree)

	// expect tree to be valid
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")

	tree.RotateRight(y)

	t.Logf("tree after rotation:\n%s", tree)

	// expect tree to be valid
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")

	assert.Equal(t, x, tree.Root(), "expected node x to be new root")
	assert.Equal(t, a, tree.Left(x), "expected node a to be root's left child")
	assert.True(t, tree.IsLeaf(a), "expected node a to be leaf node")
	assert.Equal(t, y, tree.Right(x), "expected y to be root's right child")
	assert.Equal(t, b, tree.Left(y), "expected b to be y's left child")
	assert.True(t, tree.IsLeaf(b), "expected node b to be leaf node")
	assert.Equal(t, c, tree.Right(y), "expected c to be y's right child")
	assert.True(t, tree.IsLeaf(c), "expected node c to be leaf node")
}

func TestTree_RotateRight_leftchild(t *testing.T) {
	tree := New[int, string, struct{}](func(a, b int) bool {
		return a < b
	})
	r, _ := tree.Insert(500, "root")
	y, _ := tree.Insert(100, "y")
	c, _ := tree.Insert(200, "c")
	x, _ := tree.Insert(50, "x")
	b, _ := tree.Insert(75, "b")
	a, _ := tree.Insert(25, "a")
	t.Logf("tree after creation:\n%s", tree)

	// expect tree to be valid
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")

	tree.RotateRight(y)

	t.Logf("tree after rotation:\n%s", tree)

	// expect tree to be valid
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")

	assert.Equal(t, r, tree.Root(), "root should be unchanged")
	assert.Equal(t, x, tree.Left(r), "expected node x to be left child of root")
	assert.Equal(t, a, tree.Left(x), "expected node a to be x's left child")
	assert.True(t, tree.IsLeaf(a), "expected node a to be leaf node")
	assert.Equal(t, y, tree.Right(x), "expected y to be x's right child")
	assert.Equal(t, b, tree.Left(y), "expected b to be y's left child")
	assert.True(t, tree.IsLeaf(b), "expected node b to be leaf node")
	assert.Equal(t, c, tree.Right(y), "expected c to be y's right child")
	assert.True(t, tree.IsLeaf(c), "expected node c to be leaf node")
}

func TestTree_RotateRight_rightchild(t *testing.T) {
	tree := New[int, string, struct{}](func(a, b int) bool {
		return a < b
	})
	r, _ := tree.Insert(0, "root")
	y, _ := tree.Insert(100, "y")
	c, _ := tree.Insert(200, "c")
	x, _ := tree.Insert(50, "x")
	b, _ := tree.Insert(75, "b")
	a, _ := tree.Insert(25, "a")
	t.Logf("tree after creation:\n%s", tree)

	// expect tree to be valid
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")

	tree.RotateRight(y)

	t.Logf("tree after rotation:\n%s", tree)

	// expect tree to be valid
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")

	assert.Equal(t, r, tree.Root(), "root should be unchanged")
	assert.Equal(t, x, tree.Right(r), "expected node x to be right child of root")
	assert.Equal(t, a, tree.Left(x), "expected node a to be x's left child")
	assert.True(t, tree.IsLeaf(a), "expected node a to be leaf node")
	assert.Equal(t, y, tree.Right(x), "expected y to be x's right child")
	assert.Equal(t, b, tree.Left(y), "expected b to be y's left child")
	assert.True(t, tree.IsLeaf(b), "expected node b to be leaf node")
	assert.Equal(t, c, tree.Right(y), "expected c to be y's right child")
	assert.True(t, tree.IsLeaf(c), "expected node c to be leaf node")
}

func TestTree_RotateRight_nil(t *testing.T) {
	tree := New[int, string, struct{}](func(a, b int) bool {
		return a < b
	})
	root, _ := tree.Insert(100, "root")
	lc, _ := tree.Insert(50, "left child")
	rc, _ := tree.Insert(150, "right child")

	t.Logf("tree after creation:\n%s", tree)

	// expect tree to be valid
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")

	tree.RotateRight(nil)

	t.Logf("tree after rotation:\n%s", tree)

	// expect tree to be valid
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")

	assert.Equal(t, root, tree.Root(), "tree should be unchanged")
	assert.Equal(t, lc, tree.Left(root), "expected node lc to be left child of root")
	assert.Equal(t, rc, tree.Right(root), "expected node rc to be right child of root")
	assert.True(t, tree.IsLeaf(lc), "expected node lc to be leaf node")
	assert.True(t, tree.IsLeaf(rc), "expected node rc to be leaf node")
}

func TestTree_IsTreeValid(t *testing.T) {
	createTree := func() *Tree[int, struct{}, struct{}] {
		tree := New[int, struct{}, struct{}](func(a, b int) bool {
			return a < b
		})
		tree.Insert(100, struct{}{})
		tree.Insert(50, struct{}{})
		tree.Insert(25, struct{}{})
		tree.Insert(75, struct{}{})
		tree.Insert(150, struct{}{})
		tree.Insert(125, struct{}{})
		tree.Insert(175, struct{}{})
		require.NoError(t, tree.IsTreeValid(), "expected valid tree")
		return tree
	}

	// break sentinel node
	tree := createTree()
	tree.nil.parent = nil
	require.Error(t, tree.IsTreeValid(), "expected sentinel nil parent to return error")

	// break root node
	tree = createTree()
	tree.root.parent = nil
	require.Error(t, tree.IsTreeValid(), "expected root nil parent to return error")

	// break tree: out of order node
	tree = createTree()
	minNode := tree.Min(tree.Root())
	minNode.key = 51
	require.Error(t, tree.IsTreeValid(), "expected out of order node key to return error")

	// break tree: broken parent/child relationship
	tree = createTree()
	brokenNode, _ := tree.Search(75)
	brokenNode.parent = tree.Root()
	require.Error(t, tree.IsTreeValid(), "expected parent/child mismatch to return error")

}

func TestTree_Predecessor(t *testing.T) {
	tree := New[int, struct{}, struct{}](func(a, b int) bool {
		return a < b
	})
	tree.Insert(100, struct{}{})
	tree.Insert(50, struct{}{})
	tree.Insert(25, struct{}{})
	tree.Insert(75, struct{}{})
	tree.Insert(150, struct{}{})
	tree.Insert(125, struct{}{})
	tree.Insert(175, struct{}{})
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")

	expected := []int{175, 150, 125, 100, 75, 50, 25}
	actual := make([]int, 0, len(expected))

	n := tree.Max(tree.Root())
	for !tree.IsNil(n) {
		actual = append(actual, n.key)
		n = tree.Predecessor(n)
	}

	assert.Equal(t, expected, actual)

}

func TestTree_Successor(t *testing.T) {
	tree := New[int, struct{}, struct{}](func(a, b int) bool {
		return a < b
	})
	tree.Insert(100, struct{}{})
	tree.Insert(50, struct{}{})
	tree.Insert(25, struct{}{})
	tree.Insert(75, struct{}{})
	tree.Insert(150, struct{}{})
	tree.Insert(125, struct{}{})
	tree.Insert(175, struct{}{})
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")

	expected := []int{25, 50, 75, 100, 125, 150, 175}
	actual := make([]int, 0, len(expected))

	n := tree.Min(tree.Root())
	for !tree.IsNil(n) {
		actual = append(actual, n.key)
		n = tree.Successor(n)
	}

	assert.Equal(t, expected, actual)

}

func TestTree_Search(t *testing.T) {
	tree := New[int, struct{}, struct{}](func(a, b int) bool {
		return a < b
	})
	tree.Insert(100, struct{}{})
	tree.Insert(50, struct{}{})
	tree.Insert(25, struct{}{})
	tree.Insert(75, struct{}{})
	tree.Insert(150, struct{}{})
	tree.Insert(125, struct{}{})
	tree.Insert(175, struct{}{})
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")

	// ensure all valid keys are found
	keys := []int{25, 50, 75, 100, 125, 150, 175}
	for _, key := range keys {
		n, found := tree.Search(key)
		assert.Truef(t, found, "expected to find key: %d", key)
		assert.Equalf(t, key, tree.Key(n), "expected found node key %d to match searched key %d", tree.Key(n), key)
	}

	// attempt to find key that does not exist
	n, found := tree.Search(500)
	assert.Falsef(t, found, "expected to not find key: %d", 500)
	assert.True(t, tree.IsNil(n), "expected tree.nil for node not found")
}

func TestTree_Sibling(t *testing.T) {
	tree := New[int, struct{}, struct{}](func(a, b int) bool {
		return a < b
	})

	assert.True(t, tree.IsNil(tree.Sibling(tree.Root())), "expected empty tree to return t.nil sibling")

	n100, _ := tree.Insert(100, struct{}{})
	n50, _ := tree.Insert(50, struct{}{})
	tree.Insert(25, struct{}{})
	tree.Insert(75, struct{}{})
	n150, _ := tree.Insert(150, struct{}{})
	n175, _ := tree.Insert(175, struct{}{})
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")

	assert.True(t, tree.IsNil(tree.Sibling(n100)), "expected root node to return t.nil sibling")

	assert.Equal(t, n150, tree.Sibling(n50), "expected node 50 node to return node 150 as sibling")
	assert.Equal(t, n50, tree.Sibling(n150), "expected node 150 node to return node 50 as sibling")

	assert.True(t, tree.IsNil(tree.Sibling(n175)), "expected node 175 to return t.nil sibling")

}

func TestTree_String(t *testing.T) {
	tree := New[int, uint8, struct{}](func(a, b int) bool {
		return a < b
	})

	assert.Equal(t, "Empty Tree", tree.String())

	tree.Insert(100, 100)
	tree.Insert(50, 50)
	tree.Insert(25, 25)
	tree.Insert(75, 75)
	tree.Insert(150, 150)
	tree.Insert(125, 125)
	tree.Insert(175, 175)

	expected := `      ╭── 25: 25 [{}]
 ╭── 50: 50 [{}]
 │    ╰── 75: 75 [{}]
100: 100 [{}]
 │    ╭── 125: 125 [{}]
 ╰── 150: 150 [{}]
      ╰── 175: 175 [{}]
`

	assert.Equal(t, expected, tree.String())

}

func TestTree_Height(t *testing.T) {
	tree := New[int, struct{}, struct{}](func(a, b int) bool {
		return a < b
	})
	tree.Insert(100, struct{}{})
	n50, _ := tree.Insert(50, struct{}{})
	n25, _ := tree.Insert(25, struct{}{})
	require.NoError(t, tree.IsTreeValid(), "expected valid tree")
	assert.Equal(t, 0, tree.Depth(tree.Root()))
	assert.Equal(t, 1, tree.Depth(n50))
	assert.Equal(t, 2, tree.Depth(n25))
}

func TestTree_Contains(t *testing.T) {

	// Make two trees with matching keys.

	treeA := New[int, struct{}, struct{}](func(a, b int) bool {
		return a < b
	})
	treeA.Insert(100, struct{}{})
	treeA.Insert(50, struct{}{})

	treeB := New[int, struct{}, struct{}](func(a, b int) bool {
		return a < b
	})
	treeB.Insert(100, struct{}{})
	nB, _ := treeB.Insert(50, struct{}{})

	assert.False(t, treeA.Contains(nB), "node from tree B should not exist in node A")
	assert.True(t, treeB.Contains(nB), "expected to find node B in tree B")
}
