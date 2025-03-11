package rbtree

import (
	"github.com/mikenye/gotrees/bst"
	"testing"
)

// BenchmarkTree_SearchDelete creates a very large tree (10M nodes),
// then deletes items from said tree in the benchmarking loop.
func BenchmarkTree_SearchDelete(b *testing.B) {
	var n *bst.Node[int, struct{}, Color]

	// create a tree with integer key & no value,
	tree := New[int, struct{}](func(a, b int) bool {
		return a < b
	})

	// create large tree to delete from
	for i := 0; i <= 10_000_000; i++ {
		tree.Insert(i, struct{}{})
	}

	i := 0
	b.ResetTimer()
	for b.Loop() {
		n, _ = tree.Search(i)
		tree.Delete(n)
		i++
	}
}

// BenchmarkTree_Insert creates inserts items into a tree in the benchmarking loop.
func BenchmarkTree_Insert(b *testing.B) {
	tree := New[int, struct{}](func(a, b int) bool {
		return a < b
	})
	i := 0
	b.ResetTimer()
	for b.Loop() {
		tree.Insert(i, struct{}{})
		i++
	}
}
