package rbtree

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// FuzzTree inserts 10 nodes and deletes between 1 and 10 of them.
// Tree structure and validity is checked after each insert and delete.
func FuzzTree(f *testing.F) {
	f.Add(1, 11, 12, 69, 4, 14, 82, 50, 77, 3, 10)
	f.Fuzz(func(t *testing.T, k1, k2, k3, k4, k5, k6, k7, k8, k9, k10, deleteKeys int) {
		if deleteKeys < 0 || deleteKeys > 9 {
			return
		}

		// create tree
		tree := New[int, struct{}](func(a, b int) bool {
			return a < b
		})

		// insert nodes
		keys := []int{k1, k2, k3, k4, k5, k6, k7, k8, k9, k10}
		t.Logf("input: %v", keys)
		for _, k := range keys {

			// insert node
			t.Logf("inserting node: %d", k)
			tree.Insert(k, struct{}{})

			// check
			t.Logf("rbtree after insert of node %d:\n%s", k, tree)
			err := tree.IsTreeValid()
			if err != nil {
				t.Error(err)
			}
		}

		// delete nodes
		deletedNodes := map[int]struct{}{}
		for i := 0; i <= deleteKeys; i++ {
			t.Logf("deleting node: %d", keys[i])

			// has the node already been deleted?
			_, alreadyDeleted := deletedNodes[keys[i]]

			// search for node
			n, found := tree.Search(keys[i])
			if !found && !alreadyDeleted {
				// if node not found and hasn't already been deleted, something is wrong
				t.Errorf("node %d not found", keys[i])
			}

			// delete node
			deleted := tree.Delete(n)
			if !deleted && !alreadyDeleted {
				// if node not deleted and hasn't already been deleted, something is wrong
				t.Errorf("node %d not deleted", keys[i])
			}

			// check validity of tree
			if !alreadyDeleted {
				t.Logf("rbtree after delete of node %d:\n%s", keys[i], tree)
				err := tree.IsTreeValid()
				if err != nil {
					t.Error(err)
				}
			}

			// add deleted node to map set
			deletedNodes[keys[i]] = struct{}{}
		}
	})
}

func TestTree_Delete(t *testing.T) {
	// todo: add structure checks
	tests := map[string]struct {
		keys     []int // in order of insert
		deletion func(t *testing.T, tree *Tree[int, struct{}])
		checks   func(t *testing.T, tree *Tree[int, struct{}])
	}{
		"nil node": {
			keys: []int{20, 10, 30},
			deletion: func(t *testing.T, tree *Tree[int, struct{}]) {
				deleted := tree.Delete(nil)
				require.False(t, deleted, "expected nil node to not be deleted")
				deleted = tree.Delete(tree.Sentinel())
				require.False(t, deleted, "expected nil node to not be deleted")
			},
			checks: func(t *testing.T, tree *Tree[int, struct{}]) {
				assert.Equal(t, tree.Sentinel(), tree.Parent(tree.Root()), "unexpected structure after delete")
				assert.Equal(t, 20, tree.Key(tree.Root()), "unexpected structure after delete")
				assert.Equal(t, 10, tree.Key(tree.Left(tree.Root())), "unexpected structure after delete")
				assert.Equal(t, 30, tree.Key(tree.Right(tree.Root())), "unexpected structure after delete")
			},
		},
		"left child delete, no fixup cases": {
			keys: []int{14, 11, 69, 3, 12, 50, 82, 1, 4, 77},
			deletion: func(t *testing.T, tree *Tree[int, struct{}]) {
				n1, _ := tree.Search(1)
				ok := tree.Delete(n1)
				require.True(t, ok)
			},
			checks: func(t *testing.T, tree *Tree[int, struct{}]) {
				n3, _ := tree.Search(3)
				n4, _ := tree.Search(4)
				assert.Equal(t, Black, tree.Metadata(n3), "expected node 3 to remain black")
				assert.Equal(t, tree.Sentinel(), tree.Left(n3), "expected left child of node 3 to be sentinel after delete")
				assert.Equal(t, n4, tree.Right(n3), "expected right child of node 3 to be node 4")
				assert.Equal(t, Red, tree.Metadata(n4), "expected node 4 to remain red")
			},
		},
		"successor transplant, fixup cases 3 & 4": {
			keys: []int{14, 11, 69, 3, 12, 50, 82, 1, 4, 77},
			deletion: func(t *testing.T, tree *Tree[int, struct{}]) {
				n1, _ := tree.Search(1)
				tree.Delete(n1)
				// no assertions for above deletions as this follows on from previous case(s) above
				n11, _ := tree.Search(11)
				ok := tree.Delete(n11)
				require.True(t, ok)
			},
			checks: func(t *testing.T, tree *Tree[int, struct{}]) {
				n3, _ := tree.Search(3)
				n4, _ := tree.Search(4)
				n12, _ := tree.Search(12)

				assert.Equal(t, n4, tree.Left(tree.Root()), "expected node 4 to be root left child")
				assert.Equal(t, Red, tree.Metadata(n4), "expected node 4 to remain red")
				assert.Equal(t, n3, tree.Left(n4), "expected left child of node 4 to be node 3")
				assert.Equal(t, Black, tree.Metadata(n3), "expected node 3 to remain black")
				assert.Equal(t, n12, tree.Right(n4), "expected right child of node 4 to be node 12")
				assert.Equal(t, Black, tree.Metadata(n12), "expected node 12 to remain black")
				assert.True(t, tree.IsLeaf(n3), "expected node 3 to be leaf")
				assert.True(t, tree.IsLeaf(n12), "expected node 12 to be leaf")
			},
		},
		"left child replacement, fixup case 2": {
			keys: []int{14, 11, 69, 3, 12, 50, 82, 1, 4, 77},
			deletion: func(t *testing.T, tree *Tree[int, struct{}]) {
				n1, _ := tree.Search(1)
				tree.Delete(n1)
				n11, _ := tree.Search(11)
				tree.Delete(n11)
				// no assertions for above deletions as this follows on from previous case(s) above
				n12, _ := tree.Search(12)
				ok := tree.Delete(n12)
				require.True(t, ok)
			},
			checks: func(t *testing.T, tree *Tree[int, struct{}]) {
				n3, _ := tree.Search(3)
				n4, _ := tree.Search(4)

				assert.Equal(t, n4, tree.Left(tree.Root()), "expected node 4 to be root left child")
				assert.Equal(t, Black, tree.Metadata(n4), "expected node 4 to change to black")
				assert.Equal(t, n3, tree.Left(n4), "expected left child of node 4 to be node 3")
				assert.Equal(t, Red, tree.Metadata(n3), "expected node 3 to change to red")
				assert.Equal(t, tree.Sentinel(), tree.Right(n4), "expected right child of node 4 to be nil")
				assert.True(t, tree.IsLeaf(n3), "expected node 3 to be leaf")
			},
		},
		"successor transplant, no fixup": {
			keys: []int{14, 11, 69, 3, 12, 50, 82, 1, 4, 77},
			deletion: func(t *testing.T, tree *Tree[int, struct{}]) {
				n1, _ := tree.Search(1)
				tree.Delete(n1)
				n11, _ := tree.Search(11)
				tree.Delete(n11)
				n12, _ := tree.Search(12)
				tree.Delete(n12)
				// no assertions for above deletions as this follows on from previous case(s) above
				n69, _ := tree.Search(69)
				ok := tree.Delete(n69)
				require.True(t, ok)
			},
			checks: func(t *testing.T, tree *Tree[int, struct{}]) {
				n50, _ := tree.Search(50)
				n77, _ := tree.Search(77)
				n82, _ := tree.Search(82)

				assert.Equal(t, n77, tree.Right(tree.Root()), "expected node 77 to be root right child")
				assert.Equal(t, Red, tree.Metadata(n77), "expected node 77 to be red")
				assert.Equal(t, n50, tree.Left(n77), "expected left child of node 77 to be node 50")
				assert.Equal(t, Black, tree.Metadata(n50), "expected node 50 to be black")
				assert.Equal(t, n82, tree.Right(n77), "expected right child of node 77 to be node 82")
				assert.Equal(t, Black, tree.Metadata(n82), "expected node 82 to be black")
				assert.True(t, tree.IsLeaf(n50), "expected node 50 to be leaf")
				assert.True(t, tree.IsLeaf(n82), "expected node 82 to be leaf")
			},
		},
		"right child replacement, no fixup": {
			keys: []int{14, 11, 69, 3, 12, 50, 82, 1, 4, 77},
			deletion: func(t *testing.T, tree *Tree[int, struct{}]) {
				n1, _ := tree.Search(1)
				tree.Delete(n1)
				n11, _ := tree.Search(11)
				tree.Delete(n11)
				n12, _ := tree.Search(12)
				tree.Delete(n12)
				n69, _ := tree.Search(69)
				tree.Delete(n69)
				// no assertions for above deletions as this follows on from previous case(s) above
				n4, _ := tree.Search(4)
				ok := tree.Delete(n4)
				require.True(t, ok)
			},
			checks: func(t *testing.T, tree *Tree[int, struct{}]) {
				n3, _ := tree.Search(3)

				assert.Equal(t, n3, tree.Left(tree.Root()), "expected node 3 to be root left child")
				assert.Equal(t, Black, tree.Metadata(n3), "expected node 3 to be black")
				assert.True(t, tree.IsLeaf(n3), "expected node 3 to be leaf")
			},
		},
		"root node with two children": {
			keys: []int{14, 11, 69, 3, 12, 50, 82, 1, 4, 77},
			deletion: func(t *testing.T, tree *Tree[int, struct{}]) {
				n1, _ := tree.Search(1)
				tree.Delete(n1)
				n11, _ := tree.Search(11)
				tree.Delete(n11)
				n12, _ := tree.Search(12)
				tree.Delete(n12)
				n69, _ := tree.Search(69)
				tree.Delete(n69)
				n4, _ := tree.Search(4)
				tree.Delete(n4)
				// no assertions for above deletions as this follows on from previous case(s) above
				n14, _ := tree.Search(14)
				ok := tree.Delete(n14)
				require.True(t, ok)
			},
			checks: func(t *testing.T, tree *Tree[int, struct{}]) {
				n50, _ := tree.Search(50)
				n3, _ := tree.Search(3)
				n77, _ := tree.Search(77)
				n82, _ := tree.Search(82)

				assert.Equal(t, tree.Root(), n50, "expected node 50 to be new tree root")
				assert.Equal(t, n3, tree.Left(tree.Root()), "expected node 3 to be root left child")
				assert.Equal(t, Black, tree.Metadata(n3), "expected node 3 to be black")
				assert.True(t, tree.IsLeaf(n3), "expected node 3 to be leaf")
				assert.Equal(t, n77, tree.Right(tree.Root()), "expected node 77 to be root right child")
				assert.Equal(t, Black, tree.Metadata(n77), "expected node 77 to be black")
				assert.Equal(t, tree.Sentinel(), tree.Left(n77), "expected node 77 left child to be nil")
				assert.Equal(t, n82, tree.Right(n77), "expected node 77 right child to be node 82")
				assert.True(t, tree.IsLeaf(n82), "expected node 82 to be leaf")
				assert.Equal(t, Red, tree.Metadata(n82), "expected node 77 to be black")
			},
		},
		"right child delete, no fixup": {
			keys: []int{14, 11, 69, 3, 12, 50, 82, 1, 4, 77},
			deletion: func(t *testing.T, tree *Tree[int, struct{}]) {
				n1, _ := tree.Search(1)
				tree.Delete(n1)
				n11, _ := tree.Search(11)
				tree.Delete(n11)
				n12, _ := tree.Search(12)
				tree.Delete(n12)
				n69, _ := tree.Search(69)
				tree.Delete(n69)
				n4, _ := tree.Search(4)
				tree.Delete(n4)
				n14, _ := tree.Search(14)
				tree.Delete(n14)
				// no assertions for above deletions as this follows on from previous case(s) above
				n82, _ := tree.Search(82)
				ok := tree.Delete(n82)
				require.True(t, ok)
			},
			checks: func(t *testing.T, tree *Tree[int, struct{}]) {
				n50, _ := tree.Search(50)
				n3, _ := tree.Search(3)
				n77, _ := tree.Search(77)

				assert.Equal(t, tree.Root(), n50, "expected node 50 to be tree root")
				assert.Equal(t, n3, tree.Left(tree.Root()), "expected node 3 to be root left child")
				assert.Equal(t, Black, tree.Metadata(n3), "expected node 3 to be black")
				assert.True(t, tree.IsLeaf(n3), "expected node 3 to be leaf")
				assert.Equal(t, n77, tree.Right(tree.Root()), "expected node 77 to be root right child")
				assert.Equal(t, Black, tree.Metadata(n77), "expected node 77 to be black")
				assert.True(t, tree.IsLeaf(n77), "expected node 77 to be leaf")
			},
		},
		"root delete, fixup case 2": {
			keys: []int{14, 11, 69, 3, 12, 50, 82, 1, 4, 77},
			deletion: func(t *testing.T, tree *Tree[int, struct{}]) {
				n1, _ := tree.Search(1)
				tree.Delete(n1)
				n11, _ := tree.Search(11)
				tree.Delete(n11)
				n12, _ := tree.Search(12)
				tree.Delete(n12)
				n69, _ := tree.Search(69)
				tree.Delete(n69)
				n4, _ := tree.Search(4)
				tree.Delete(n4)
				n14, _ := tree.Search(14)
				tree.Delete(n14)
				n82, _ := tree.Search(82)
				tree.Delete(n82)
				// no assertions for above deletions as this follows on from previous case(s) above
				n50, _ := tree.Search(50)
				ok := tree.Delete(n50)
				require.True(t, ok)
			},
			checks: func(t *testing.T, tree *Tree[int, struct{}]) {
				n3, _ := tree.Search(3)
				n77, _ := tree.Search(77)

				assert.Equal(t, tree.Root(), n77, "expected node 77 to be tree root")
				assert.Equal(t, n3, tree.Left(tree.Root()), "expected node 3 to be root left child")
				assert.Equal(t, Red, tree.Metadata(n3), "expected node 3 to be black")
				assert.True(t, tree.IsLeaf(n3), "expected node 3 to be leaf")
				assert.Equal(t, tree.Sentinel(), tree.Right(tree.Root()), "expected root right child to be nil")
			},
		},
		"root node with one child, no fixup": {
			keys: []int{14, 11, 69, 3, 12, 50, 82, 1, 4, 77},
			deletion: func(t *testing.T, tree *Tree[int, struct{}]) {
				n1, _ := tree.Search(1)
				tree.Delete(n1)
				n11, _ := tree.Search(11)
				tree.Delete(n11)
				n12, _ := tree.Search(12)
				tree.Delete(n12)
				n69, _ := tree.Search(69)
				tree.Delete(n69)
				n4, _ := tree.Search(4)
				tree.Delete(n4)
				n14, _ := tree.Search(14)
				tree.Delete(n14)
				n82, _ := tree.Search(82)
				tree.Delete(n82)
				n50, _ := tree.Search(50)
				tree.Delete(n50)
				// no assertions for above deletions as this follows on from previous case(s) above
				n77, _ := tree.Search(77)
				ok := tree.Delete(n77)
				require.True(t, ok)
			},
			checks: func(t *testing.T, tree *Tree[int, struct{}]) {
				n3, _ := tree.Search(3)

				assert.Equal(t, tree.Root(), n3, "expected node 77 to be tree root")
				assert.True(t, tree.IsLeaf(n3), "expected node 3 to be leaf")
			},
		},
		"root node with no children, no fixup": {
			keys: []int{14, 11, 69, 3, 12, 50, 82, 1, 4, 77},
			deletion: func(t *testing.T, tree *Tree[int, struct{}]) {
				n1, _ := tree.Search(1)
				tree.Delete(n1)
				n11, _ := tree.Search(11)
				tree.Delete(n11)
				n12, _ := tree.Search(12)
				tree.Delete(n12)
				n69, _ := tree.Search(69)
				tree.Delete(n69)
				n4, _ := tree.Search(4)
				tree.Delete(n4)
				n14, _ := tree.Search(14)
				tree.Delete(n14)
				n82, _ := tree.Search(82)
				tree.Delete(n82)
				n50, _ := tree.Search(50)
				tree.Delete(n50)
				n77, _ := tree.Search(77)
				tree.Delete(n77)
				// no assertions for above deletions as this follows on from previous case(s) above
				n3, _ := tree.Search(3)
				ok := tree.Delete(n3)
				require.True(t, ok)
			},
			checks: func(t *testing.T, tree *Tree[int, struct{}]) {
				assert.Equal(t, tree.Sentinel(), tree.Root(), "expected empty tree")
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// build tree from keys
			tree := New[int, struct{}](func(a, b int) bool { return a < b })
			for _, k := range tc.keys {
				tree.Insert(k, struct{}{})
			}
			t.Logf("rbtree before delete:\n%s", tree)
			require.NoError(t, tree.IsTreeValid(), "tree should be valid")

			// perform deletion
			tc.deletion(t, tree)
			t.Logf("rbtree after delete:\n%s", tree)
			require.NoError(t, tree.IsTreeValid(), "tree should be valid")

			// remaining checks
			tc.checks(t, tree)
		})
	}
}

func TestTree_Insert_fixup_cases(t *testing.T) {
	// todo: add structure checks
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
	assert.Panics(t, func() {
		tree.MustSetMetadata()
	})
	assert.Panics(t, func() {
		tree.SetMetadata()
	})
	assert.Panics(t, func() {
		tree.RotateLeft()
	})
	assert.Panics(t, func() {
		tree.RotateRight()
	})
	assert.Panics(t, func() {
		tree.SetLeft()
	})
	assert.Panics(t, func() {
		tree.SetParent()
	})
	assert.Panics(t, func() {
		tree.SetRight()
	})
	assert.Panics(t, func() {
		tree.Transplant()
	})
}

func TestTree_Size(t *testing.T) {
	tree := New[int, struct{}](func(a, b int) bool { return a < b })
	assert.Equal(t, 0, tree.Size(), "expected empty tree")
	tree.Insert(10, struct{}{})
	tree.Insert(5, struct{}{})
	tree.Insert(15, struct{}{})
	tree.Insert(14, struct{}{})
	assert.Equal(t, 4, tree.Size(), "expected 4 nodes in tree")
}
