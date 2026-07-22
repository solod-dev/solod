package maps

import (
	"testing"
)

// FuzzMap checks the map against Go's builtin map as an oracle: a random
// sequence of Set/Get/Delete operations must leave both in agreement.
func FuzzMap(f *testing.F) {
	f.Add([]byte{0, 1, 10, 0, 2, 20, 1, 1, 2, 1})
	f.Fuzz(func(t *testing.T, ops []byte) {
		m := New[int, int](nil, 8)
		defer m.Free()
		oracle := map[int]int{}

		// Each op is a triple: (kind, key, value).
		for i := 0; i+2 < len(ops); i += 3 {
			key := int(ops[i+1])
			val := int(ops[i+2])
			switch ops[i] % 3 {
			case 0: // Set
				m.Set(key, val)
				oracle[key] = val
			case 1: // Delete
				m.Delete(key)
				delete(oracle, key)
			case 2: // Get
				if got := m.Get(key); got != oracle[key] {
					t.Fatalf("Get(%d) = %d, want %d", key, got, oracle[key])
				}
			}
		}

		if m.Len() != len(oracle) {
			t.Fatalf("Len() = %d, want %d", m.Len(), len(oracle))
		}
		for key, want := range oracle {
			if !m.Has(key) {
				t.Fatalf("Has(%d) = false, want true", key)
			}
			if got := m.Get(key); got != want {
				t.Fatalf("Get(%d) = %d, want %d", key, got, want)
			}
		}
	})
}
