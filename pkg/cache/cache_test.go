package cache

import (
	"fmt"
	"testing"
)

func TestCache(t *testing.T) {
	data := map[string]int{
		"foo":    1,
		"bar":    2,
		"zoo":    3,
		"github": 4,
		"gitlab": 5,
	}

	allocFunc := func(name KeyType) (value ValueType, ok bool) {
		return name, true
	}

	t.Run("basic", func(t *testing.T) {
		c := New(CtorConfig{
			Capacity:  5,
			AllocFunc: allocFunc,
		})

		for name := range data {
			if c.Contains(name) == true {
				t.Fatalf(`c.Contains(%s) = %v, expected %v`, name, true, false)
			}

			if c.Remove(name) == true {
				t.Fatalf(`c.Remove(%s) = %v, expected %v`, name, true, false)
			}

			expectVal, _ := allocFunc(name)
			getVal, getOk := c.Get(name)
			if getOk == false || getVal != expectVal {
				t.Fatalf(`c.Get(%s) = %v, %v, expected %v, true`, name, getVal, getOk, expectVal)
			}

			evicted := c.Set(name, data[name])
			if evicted == true {
				t.Fatalf(`c.Set(%s, %v) == true, expected false`, name, data[name])
			}

			expectVal, _ = data[name]
			getVal, getOk = c.Get(name)
			if getOk == false || getVal != expectVal {
				t.Fatalf(`c.Get(%s) = %v, %v, expected %v, true`, name, getVal, getOk, expectVal)
			}
		}

		if cLen := c.Len(); cLen != len(data) {
			t.Fatalf(`c.Len() = %v, expected %v`, cLen, len(data))
		}

		if c.Remove("foo") == true {
			if c.Contains("foo") {
				t.Fatalf(`c.Contains(%s) = true, expected false`, "foo")
			}

			if cLen := c.Len(); cLen != len(data)-1 {
				t.Fatalf(`c.Len() = %v, expected %v`, cLen, len(data)-1)
			}
		} else {
			t.Fatalf(`c.Remove(%s) = false, expected true`, "foo")
		}

		c.Purge()
		if cLen := c.Len(); cLen != 0 {
			t.Fatalf(`c.Len() = %v, expected %v`, cLen, 0)
		}
	})

	t.Run("without_alloc_func", func(t *testing.T) {
		c := New(CtorConfig{
			Capacity: 10,
		})

		if _, ok := c.Get("not_found"); ok {
			t.Fatalf(`c.Get("not_found") = _, true, expected _, false`)
		}
	})

	t.Run("capacity", func(t *testing.T) {
		c := New(CtorConfig{
			Capacity:  3,
			AllocFunc: allocFunc,
		})

		c.Set("a", 1)
		c.Set("b", 2)
		c.Set("c", 3)

		if c.Len() != 3 {
			t.Fatalf("c.Len() = %v, expected 3", c.Len())
		}

		evicted := c.Set("d", 4)
		if !evicted {
			t.Fatalf(`c.Set(%s, %v) == false, expected true`, `"d"`, 4)
		}

		if cLen := c.Len(); cLen != 3 {
			t.Fatalf(`c.Len() = %v, expected %v`, cLen, 3)
		}
	})
}

func TestParallelSafety(t *testing.T) {
	c := New(CtorConfig{
		Capacity: 10,
		AllocFunc: func(key KeyType) (value ValueType, ok bool) {
			return key.(int) * -1, true
		},
	})

	for i := 0; i < 3; i++ {
		i := i
		t.Run(fmt.Sprintf(`parallel#%d`, i), func(t *testing.T) {
			t.Parallel()

			for j := 1; j < 6; j++ {
				val := i * j * 10
				c.Set(j, val)
				if getVal, _ := c.Get(j); getVal != val {
					t.Fatalf(`c.Get("%d") == %v; expected %d`, j, getVal, val)
				}

				c.Remove(j)
				if getVal, _ := c.Get(j); getVal != -1*j {
					t.Fatalf(`c.Get("%d") == %v; expected %d`, j, getVal, -1*j)
				}
			}
		})
	}
}

func BenchmarkCacheCapacity_N(b *testing.B) {
	benchmarkCache(b.N, b)
}

func BenchmarkCacheCapacity_HalfOfN(b *testing.B) {
	benchmarkCache(b.N/2, b)
}

func benchmarkCache(capacity int, b *testing.B) {
	c := New(CtorConfig{
		Capacity: capacity,
		AllocFunc: func(key KeyType) (value ValueType, ok bool) {
			return key, true
		},
	})

	b.ResetTimer()

	b.Run("setter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.Set(i, i+1)
		}
	})

	b.Run("getter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.Get(i)
		}
	})
}
