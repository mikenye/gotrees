# Red-Black Binary Search Tree - Go Implementation

[![codecov](https://codecov.io/gh/mikenye/gotrees/graph/badge.svg?token=SXQZZUMRAX)](https://codecov.io/gh/mikenye/gotrees)

## Overview

The `rbtree` package provides a **pointer-based Red-Black Binary Search Tree** implementation in Go. It is designed to be:

- **Generic**: Supports Go generics (`K`, `V`) for flexible key-value storage.
- **Efficient**: Balances itself to maintain **O(log n)** insertions, deletions, and lookups.
- **Extensible**: Built on top of the `bst` package, allowing modifications for custom balancing strategies.

## Why This Package Exists

1. **Generic Support** – Prior to Go generics, existing BST implementations relied on `interface{}`, leading to runtime type assertions and reduced type safety.
2. **Performance** – By leveraging generics, this implementation avoids the overhead of type assertions and provides efficient operations.
3. **Learning & Exploration** – This package was developed to deepen the understanding of balanced tree structures, balancing efficiency with usability.

## Generics: `K` and `V`

- `K` (**Key type**) – Defines the ordering of nodes. The ordering must be specified by a user-defined **comparison function**.
- `V` (**Value type**) – The data stored in each node. If no value is needed, `struct{}` can be used for **zero memory overhead**.

## Installation

```sh
# Using Go modules
go get github.com/mikenye/gotrees/rbtree
```
## Basic Usage

### Creating a New Tree

```go
tree := rbtree.New[int, string](func(a, b int) bool { return a < b })
```

### Inserting & Deleting Nodes

```go
tree.Insert(10, "ten")
tree.Insert(20, "twenty")
node, found := tree.Search(10)
if found {
    tree.Delete(node)
}
```

### Traversing the Tree

#### Recursive In-Order Traversal

```go
tree.TraverseInOrder(tree.Root(), func(n *bst.Node[int, string, struct{}]) bool {
    fmt.Println(n.Key())
    return true
})
```

**Note:** `TraverseInOrder` uses recursion. If the tree is deep and highly unbalanced, this could lead to a **stack overflow**. In such cases, consider an **iterative traversal** using `Successor`:

#### Iterative In-Order Traversal
```go
for node := tree.Min(tree.Root()); !tree.IsNil(node); node = tree.Successor(node) {
    fmt.Println(tree.Key(node))
}
```

## Limitations
- **Not Thread-Safe** – Requires external synchronization for concurrent use.
- **No Duplicate Keys** – Keys must be unique.
