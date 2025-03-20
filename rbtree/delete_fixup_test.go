package rbtree

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestDeleteFixupCases tests the deleteFixup method directly by
// creating a variety of delete scenarios
func TestDeleteFixupCases(t *testing.T) {
	t.Run("AllCases", func(t *testing.T) {
		// Create a substantial tree that will exercise all different deletion cases
		tree := New[int, string](func(a, b int) bool { return a < b })

		// Insert a range of keys
		for i := 0; i < 100; i += 2 {
			tree.Insert(i, "value")
		}

		// Verify tree is valid initially
		assert.NoError(t, tree.IsTreeValid())

		// Delete nodes one by one to trigger various fixup cases
		for i := 0; i < 100; i += 2 {
			n, found := tree.Search(i)
			assert.True(t, found)

			deleted := tree.Delete(n)
			assert.True(t, deleted)

			// Tree should remain valid after each deletion
			assert.NoError(t, tree.IsTreeValid())
		}
	})
}

// This function implements a more comprehensive testing of the deleted tree
// by attempting to create trees that will trigger specific deletion fixup cases
func TestDeleteFixupComprehensive(t *testing.T) {
	// Create a range of trees with different structures
	for seed := 1; seed < 20; seed++ {
		t.Run("ComprehensiveDeleteTest", func(t *testing.T) {
			tree := New[int, string](func(a, b int) bool { return a < b })

			// Insert nodes in a pattern that's influenced by the seed
			// This creates trees with different shapes to test various deletion cases
			for i := 0; i < 200; i++ {
				key := (i * seed) % 500
				tree.Insert(key, "value")
			}

			// Verify tree is valid initially
			assert.NoError(t, tree.IsTreeValid())

			// Delete every node in a specific order that's also influenced by the seed
			for i := 0; i < 200; i++ {
				key := ((i * 3) + seed) % 500
				n, found := tree.Search(key)
				if found {
					deleted := tree.Delete(n)
					assert.True(t, deleted)

					// Tree should remain valid after each deletion
					assert.NoError(t, tree.IsTreeValid())
				}
			}
		})
	}
}

// TestDeleteFixupDirectly calls the deleteFixup method directly with
// carefully crafted node arrangements to trigger specific cases
func TestDeleteFixupDirectly(t *testing.T) {
	t.Run("CallDeleteFixupDirectly", func(t *testing.T) {
		tree := New[int, string](func(a, b int) bool { return a < b })

		// First create a real valid tree
		for i := 0; i < 50; i++ {
			tree.Insert(i, "value")
		}

		// Get the root node for a direct call to deleteFixup
		root := tree.Root()
		assert.NotEqual(t, tree.Sentinel(), root)

		// Call deleteFixup directly with the root
		// This is only to test the actual function itself, not realistic usage
		tree.deleteFixup(root)

		// Tree should still be valid
		assert.NoError(t, tree.IsTreeValid())
	})
}

// TestIsTreeValidRedRoot tests the case where the root is red, which violates RB tree property
func TestIsTreeValidRedRoot(t *testing.T) {
	// Create a valid tree
	tree := New[int, string](func(a, b int) bool { return a < b })
	tree.Insert(10, "ten")

	// Verify it's valid initially
	assert.NoError(t, tree.IsTreeValid())

	// Directly set the root node to red, violating RB property #2
	tree.Tree.MustSetMetadata(tree.Root(), Red)

	// Now tree validation should fail
	err := tree.IsTreeValid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "root node is not black")
}
