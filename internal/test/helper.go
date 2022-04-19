package test

import "testing"

type BenchFunc func(b *testing.B)

type BenchDefinition struct {
	Name string
	Func BenchFunc
}

func RunBenchmarks(b *testing.B, withAllocs bool, defs []BenchDefinition) {
	for _, def := range defs {
		b.Run(def.Name, func(b *testing.B) {
			if withAllocs {
				b.ReportAllocs()
			}
			b.ResetTimer()
			def.Func(b)
		})
	}
}

type TestFunc func(t *testing.T)

type TestDefinition struct {
	Name string
	Func TestFunc
}

func RunTests(t *testing.T, defs []TestDefinition) {
	for _, def := range defs {
		t.Run(def.Name, def.Func)
	}
}
