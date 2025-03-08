package rbtree

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func FuzzTree_Insert(f *testing.F) {
	f.Add(1, 11, 12, 69, 4, 14, 82, 50, 77, 3)
	f.Fuzz(func(t *testing.T, k1, k2, k3, k4, k5, k6, k7, k8, k9, k10 int) {
		// build tree
		keys := []int{k1, k2, k3, k4, k5, k6, k7, k8, k9, k10}
		t.Logf("input: %v", keys)
		tree := New[int, struct{}](func(a, b int) bool {
			return a < b
		})
		for _, k := range keys {
			t.Logf("inserting node: %d", k)
			tree.Insert(k, struct{}{})
		}
		t.Logf("rbtree:\n%s", tree)
		// check
		err := tree.IsTreeValid()
		if err != nil {
			t.Error(err)
		}
	})
}

func TestTree_Insert_update(t *testing.T) {
	keys := []int{11, 2, 14, 1, 7, 15, 5, 8, 4}
	tree := New[int, string](func(a, b int) bool {
		return a < b
	})
	for _, k := range keys {
		tree.Insert(k, fmt.Sprintf("%d", k))
	}
	t.Logf("rbtree:\n%s", tree)

	// underlying bst should be valid
	require.NoError(t, tree.IsTreeValid(), "tree should be valid")

	n4, _ := tree.Search(4)
	require.Equal(t, "4", tree.Value(n4))

	// update node 4
	tree.Insert(4, "updated")
	assert.Equal(t, "updated", tree.Value(n4))
}

func TestTree_Insert_fixup_cases(t *testing.T) {

	tests := map[string]struct {
		keys   []int // in order of insert
		checks func(t *testing.T, tree *Tree[int, struct{}])
	}{
		"case 1, z's parent is a left child": {
			keys:   []int{11, 2, 14, 1},
			checks: func(t *testing.T, tree *Tree[int, struct{}]) { return },
		},
		"case 1, z's parent is a right child": {
			keys:   []int{1, 11, 12, 69},
			checks: func(t *testing.T, tree *Tree[int, struct{}]) { return },
		},
		"case 2 & 3, z's parent is a left child": {
			keys:   []int{11, 2, 14, 1, 7, 15, 5, 8, 4},
			checks: func(t *testing.T, tree *Tree[int, struct{}]) { return },
		},
		"case 2 & 3, z's parent is a right child": {
			keys:   []int{1, 11, 12, 69, 4, 14},
			checks: func(t *testing.T, tree *Tree[int, struct{}]) { return },
		},
		"case 3, z's parent is a right child": {
			keys:   []int{1, 11, 12},
			checks: func(t *testing.T, tree *Tree[int, struct{}]) { return },
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// make tree
			tree := New[int, struct{}](func(a, b int) bool {
				return a < b
			})
			for _, k := range tc.keys {
				t.Logf("inserting node: %d", k)
				tree.Insert(k, struct{}{})
				t.Logf("rbtree after insert:\n%s", tree)
			}
			require.NoError(t, tree.IsTreeValid(), "tree should be valid")

			// other checks
			tc.checks(t, tree)
		})

	}
}

func TestTree_IsTreeValid(t *testing.T) {
	tests := map[string]struct {
		creation func() *Tree[int, struct{}]
		mutation func(tree *Tree[int, struct{}])
		checks   func(t *testing.T, tree *Tree[int, struct{}])
	}{
		"valid tree": {
			creation: func() *Tree[int, struct{}] {
				tree := New[int, struct{}](func(a, b int) bool { return a < b })
				for i := -20; i <= 20; i++ {
					tree.Insert(i, struct{}{})
				}
				for i := -40; i <= -21; i++ {
					tree.Insert(i, struct{}{})
				}
				for i := 21; i <= 40; i++ {
					tree.Insert(i, struct{}{})
				}
				return tree
			},
			mutation: func(tree *Tree[int, struct{}]) { return },
			checks: func(t *testing.T, tree *Tree[int, struct{}]) {
				assert.NoError(t, tree.IsTreeValid(), "expected valid tree")
			},
		},
		"red root": {
			creation: func() *Tree[int, struct{}] {
				tree := New[int, struct{}](func(a, b int) bool { return a < b })
				tree.Insert(10, struct{}{})
				return tree
			},
			mutation: func(tree *Tree[int, struct{}]) {
				tree.Tree.MustSetMetadata(tree.Root(), Red)
			},
			checks: func(t *testing.T, tree *Tree[int, struct{}]) {
				assert.Error(t, tree.IsTreeValid(), "expected invalid tree")
			},
		},
		"nil leaf nodes are not black": {
			creation: func() *Tree[int, struct{}] {
				tree := New[int, struct{}](func(a, b int) bool { return a < b })
				tree.Insert(10, struct{}{})
				return tree
			},
			mutation: func(tree *Tree[int, struct{}]) {
				tree.Tree.MustSetMetadata(tree.Left(tree.Root()), Red)
			},
			checks: func(t *testing.T, tree *Tree[int, struct{}]) {
				assert.Error(t, tree.IsTreeValid(), "expected invalid tree")
			},
		},
		"node is red and has red left child": {
			creation: func() *Tree[int, struct{}] {
				tree := New[int, struct{}](func(a, b int) bool { return a < b })
				tree.Insert(10, struct{}{})
				tree.Insert(5, struct{}{})
				tree.Insert(15, struct{}{})
				tree.Insert(20, struct{}{})
				return tree
			},
			mutation: func(tree *Tree[int, struct{}]) {
				n, _ := tree.Search(5)
				tree.Tree.MustSetMetadata(n, Red)
				n, _ = tree.Search(15)
				tree.Tree.MustSetMetadata(n, Red)
			},
			checks: func(t *testing.T, tree *Tree[int, struct{}]) {
				assert.Error(t, tree.IsTreeValid(), "expected invalid tree")
			},
		},
		"node is red and has red right child": {
			creation: func() *Tree[int, struct{}] {
				tree := New[int, struct{}](func(a, b int) bool { return a < b })
				tree.Insert(10, struct{}{})
				tree.Insert(5, struct{}{})
				tree.Insert(15, struct{}{})
				tree.Insert(14, struct{}{})
				return tree
			},
			mutation: func(tree *Tree[int, struct{}]) {
				n, _ := tree.Search(5)
				tree.Tree.MustSetMetadata(n, Red)
				n, _ = tree.Search(15)
				tree.Tree.MustSetMetadata(n, Red)
			},
			checks: func(t *testing.T, tree *Tree[int, struct{}]) {
				assert.Error(t, tree.IsTreeValid(), "expected invalid tree")
			},
		},
		"node has black count mismatch": {
			creation: func() *Tree[int, struct{}] {
				tree := New[int, struct{}](func(a, b int) bool { return a < b })
				tree.Insert(10, struct{}{})
				tree.Insert(5, struct{}{})
				tree.Insert(15, struct{}{})
				tree.Insert(14, struct{}{})
				return tree
			},
			mutation: func(tree *Tree[int, struct{}]) {
				n, _ := tree.Search(14)
				tree.Tree.MustSetMetadata(n, Black)
			},
			checks: func(t *testing.T, tree *Tree[int, struct{}]) {
				assert.Error(t, tree.IsTreeValid(), "expected invalid tree")
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tree := tc.creation()
			t.Logf("initial rbtree:\n%s", tree)
			require.NoError(t, tree.IsTreeValid(), "tree should be valid")
			// break tree
			tc.mutation(tree)
			t.Logf("rbtree after mutation:\n%s", tree)
			// checks
			tc.checks(t, tree)
		})
	}
}

func TestTree_panics(t *testing.T) {
	tree := New[int, struct{}](func(a, b int) bool { return a < b })
	root, _ := tree.Insert(10, struct{}{})
	assert.Panics(t, func() {
		tree.MustSetMetadata(root, Red)
	})
	assert.Panics(t, func() {
		tree.SetMetadata(root, Red)
	})
	assert.Panics(t, func() {
		tree.RotateLeft(root)
	})
	assert.Panics(t, func() {
		tree.RotateRight(root)
	})
}
