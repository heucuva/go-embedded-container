# go-embedded-container

[![Go](https://github.com/heucuva/go-embedded-container/actions/workflows/go.yml/badge.svg)](https://github.com/heucuva/go-embedded-container/actions/workflows/go.yml)

Embedded Containers for Go

## What are they?
Various containers supporting Go Generics that don't own their entities. Instead, they facilitate container operations on values that might be owned by something else.

## What containers are available?

| Name | Description |
|------|-------------|
| `embedded.Hash` | A map-style container with hashed value (of `int` type) lookup |
| `embedded.HashList` | A container combining the mechanisms of `embedded.Hash` and `embedded.List` |
| `embedded.HashListMap` | A container with a map combined with a doubly-linked list interface. Internally, item keys are hashed (using FNV 64-bit hashing) so the `embedded.HashList` mechanisms can be reused |
| `embedded.HashMap` | A container combining the mechanisms of `embedded.Hash` and `embedded.Map` without incurring the performance concerns of `embedded.Map` |
| `embedded.List` | A list-style container with a doubly-linked interface |
| `embedded.Map` | A map-style container with red-black tree internally |
| `embedded.PriorityQueue` | A priority queue-style container with heap sorting internally |
