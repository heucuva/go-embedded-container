# go-embedded-container
Embedded Containers for Go

## What are they?
Various containers supporting Go Generics that don't own their entities. Instead, they facilitate container operations on values that might be owned by something else.

## What containers are available?

| Name | Description |
|------|-------------|
| `embedded.Hash` | A map-style container with hashed value (of `int` type) lookup |
| `embedded.List` | A list-style container with a doubly-linked interface |
| `embedded.HashList` | A container with a map combined with a doubly-linked list interface referencing values by a hashed (of `int` type) lookup |
| `embedded.HashListMap` | A container with a map combined with a doubly-linked list interface. Internally, item keys are hashed (using FNV 64-bit hashing) so the `embedded.HashList` mechanisms can be reused |
| `embedded.Map` | A map-style container with red-black tree internally |
