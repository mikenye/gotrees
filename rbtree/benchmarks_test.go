package rbtree

import (
	"github.com/emirpasic/gods/trees/redblacktree"
	"github.com/mikenye/gotrees/bst"
	"testing"
)

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

	// delete
	i := 0
	for b.Loop() {
		n, _ = tree.Search(i)
		tree.Delete(n)
		i++
	}
}

func BenchmarkGoDSRedBlackTree_SearchDelete(b *testing.B) {
	tree := redblacktree.NewWithIntComparator()

	// create large tree to delete from
	for i := 0; i <= 10_000_000; i++ {
		tree.Put(i, struct{}{})
	}

	// delete
	i := 0
	for b.Loop() {
		tree.Remove(i)
		i++
	}
}

func BenchmarkTree_Insert(b *testing.B) {
	// create a tree with integer key & no value,
	tree := New[int, struct{}](func(a, b int) bool {
		return a < b
	})
	i := 0
	for b.Loop() {
		tree.Insert(i, struct{}{})
		i++
	}
}

func BenchmarkGoDSRedBlackTree_Insert(b *testing.B) {
	tree := redblacktree.NewWithIntComparator()
	i := 0
	for b.Loop() {
		tree.Put(i, struct{}{})
		i++
	}
}
