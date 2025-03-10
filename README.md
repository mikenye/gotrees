# gotrees - Generic Binary Search Trees in Go

[![codecov](https://codecov.io/gh/mikenye/gotrees/graph/badge.svg?token=SXQZZUMRAX)](https://codecov.io/gh/mikenye/gotrees)

## Overview

**gotrees** is a pure Go implementation of **generic binary search trees**, designed for flexibility and efficiency. It provides:

- **`bst`:** A **basic, non-self-balancing** Binary Search Tree (BST).
- **`rbtree`:** A **self-balancing Red-Black Tree** (extends `bst`).

Both implementations are **written entirely in Go** (**no Cgo**), ensuring **portability** and **easy integration** into any Go project.

## Why gotrees Exists

1. **Generic Support:** Before Go introduced generics, BSTs often relied on `interface{}`, leading to runtime type assertions and reduced type safety.
2. **Extensibility:** `bst` provides a foundation for creating custom tree structures (e.g., Red-Black Trees, AVL Trees).
3. **Learning & Exploration:** Developed to deepen understanding of balanced tree structures while prioritizing usability.

## Installation

```sh
go get github.com/mikenye/gotrees
```

## Package Overview

### **[bst - Binary Search Tree](./bst/)**

A **generic, pointer-based Binary Search Tree (BST)** that supports:
- **Ordered key-value storage** (via a user-defined comparison function).
- **Efficient traversal and search operations**.
- **Manual structure modification** (for creating custom balanced trees).
- **Extensibility** (for extending to create Red-Black Trees, AVL Trees, etc).

> ⚠️ **`bst` does not self-balance** – for a balanced BST, use `rbtree`.

### **[rbtree - Red-Black Tree](./rbtree/)**

A **self-balancing Red-Black Tree**, extending `bst`. It ensures:
- **Automatic balancing** for O(log n) performance.
- **Preservation of Red-Black Tree properties** (no consecutive red nodes, balanced black height).
- **Safe insertions and deletions without manual balancing**.

## Features
- **✅ Well documented** – Every function documented.
- **✅ 100% Go Implementation** – No Cgo dependencies.
- **✅ Fully Generic** – Supports any key and value types.
- **✅ Extensible** – `bst` can be used to build other trees.

## Limitations
- **Not Thread-Safe** – External synchronization is required for concurrent access.
- **No Duplicate Keys** – Each key must be unique.
