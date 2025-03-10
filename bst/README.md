# Binary Search Tree (BST) - Go Implementation

## Overview

The `bst` package provides a **pointer-based Binary Search Tree (BST)** implementation in Go. It is designed to be:

- **Generic**: Supports Go generics (`K`, `V`, `M` types) for flexible key-value and metadata storage.
- **Extensible**: Serves as a foundation for other tree structures (e.g., Red-Black Trees, AVL Trees).
- **Efficient**: Implements core BST operations while maintaining clarity and performance.

## Why This Package Exists

1. **Generic Support** – Prior to Go generics, existing BST implementations relied on `interface{}`, leading to runtime type assertions and reduced type safety.
2. **Performance** – By leveraging generics, this implementation avoids the overhead of type assertions and provides efficient operations.
3. **Learning & Exploration** – This package was developed to deepen the understanding of balanced tree structures, balancing efficiency with usability.

## Generics: `K`, `V` and `M`

- `K` (**Key type**) – Defines the ordering of nodes. The ordering must be specified by a user-defined **comparison function**.
- `V` (**Value type**) – The data stored in each node. If no value is needed, `struct{}` can be used for **zero memory overhead**.
- `M` (**Metadata type**) - The metadata stored in each node. This exists to allow the tree to be extended to other tree types. For example, it could store the color in a Red-Black Tree implementation. If no value is needed, `struct{}` can be used for **zero memory overhead**.

## Metadata

Each node in the BST contains an optional metadata field, which is intended to be used when extending the base `Tree` type.

Examples of metadata use cases include:

- **Red-Black Trees**: Store the node’s color (Red or Black).
- **AVL Trees**: Store the node’s balance factor.
- **Augmented Trees**: Store additional computed values (e.g., subtree sizes, sums).

If metadata is not needed, `struct{}` can be used as the metadata type, ensuring zero memory overhead.

## Sentinel Nil Node

This BST uses a **[sentinel nil node](https://en.wikipedia.org/wiki/Sentinel_node)** to represent the absence of a valid node. **Do not compare nodes to `nil` directly** - instead, **always use `tree.IsNil(n)` to check whether a node is nil.**

## Installation

```sh
# Using Go modules
go get github.com/mikenye/gotrees/bst
```

## Basic Usage

### Creating a New Tree

```go
tree := bst.New[int, string, struct{}](func(a, b int) bool { return a < b })
```

### Inserting & Deleting Nodes

```go
tree.Insert(10, "ten")
tree.Insert(20, "twenty")
tree.Delete(tree.Search(10))
```

### Traversing the Tree

```go
tree.TraverseInOrder(tree.Root(), func(n *bst.Node[int, string, struct{}]) bool {
    fmt.Println(n.Key())
    return true
})
```

Note: `TraverseInOrder` uses recursion. If the tree is deep and highly unbalanced, this could lead to a stack overflow. Consider using `Tree.Successor` and `Tree.Predecessor` in these cases:

```go
for node := tree.Min(tree.Root); !tree.IsNil(node); node = tree.Successor(node) {
    fmt.Println(tree.Key(node))
}
```

## Limitations
- **Not Thread-Safe** – Requires external synchronization for concurrent use.
- **No Duplicate Keys** – Keys must be unique.

## Future Enhancements
- Implement a `BST` interface for swappable tree backends.
- Add an array-backed BST alternative for optimized lookup performance.
- Extend support for AVL and Red-Black Trees.
