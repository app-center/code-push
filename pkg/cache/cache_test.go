package cache

import (
	assert "github.com/stretchr/testify/require"
	"testing"
)

func TestCache(t *testing.T) {
	allocatedFunc := func(name KeyType) (value ValueType, ok bool) {
		return name, true
	}

	unAllocatedFunc := func(name KeyType) (value ValueType, ok bool) {
		return nil, false
	}

	getAllocatedCache := func(optionalCapacity ...int) *Cache {
		capacity := 3

		if len(optionalCapacity) >= 1 {
			capacity = optionalCapacity[0]
		}

		return New(CtorConfig{
			Capacity:  capacity,
			AllocFunc: allocatedFunc,
		})
	}

	getUnAllocatedCache := func(optionalCapacity ...int) *Cache {
		capacity := 3

		if len(optionalCapacity) >= 1 {
			capacity = optionalCapacity[0]
		}

		return New(CtorConfig{
			Capacity:  capacity,
			AllocFunc: unAllocatedFunc,
		})
	}

	t.Run("length of initial cache", func(t *testing.T) {
		assert.Equal(t, 0, getAllocatedCache().Len())
		assert.Equal(t, 0, getUnAllocatedCache().Len())
	})

	t.Run("query non existed key", func(t *testing.T) {
		var queryOk, containOk bool

		key := "not existed"

		allocatedCache := getAllocatedCache()
		unAllocatedCache := getUnAllocatedCache()

		_, queryOk = allocatedCache.Query(key)
		assert.False(t, queryOk)

		_, queryOk = unAllocatedCache.Query(key)
		assert.False(t, queryOk)

		containOk = allocatedCache.Contains(key)
		assert.False(t, queryOk)

		containOk = unAllocatedCache.Contains(key)
		assert.False(t, containOk)
	})

	t.Run("get non existed key", func(t *testing.T) {
		key := "not existed"

		allocatedCache := getAllocatedCache()
		unAllocatedCache := getUnAllocatedCache()

		if val, ok := allocatedCache.Get(key); ok {
			allocVal, _ := allocatedFunc(key)
			assert.Equal(t, allocVal, val)
		} else {
			assert.Fail(t, "alloc-able cache should return valid value in getting")
		}

		if val, ok := unAllocatedCache.Get(key); ok || val != nil {
			assert.Fail(t, "alloc-unable cache should return nil and false in getting")
		}
	})

	t.Run("existed key shall be queryable and deletable", func(t *testing.T) {
		key := "code-push-key"
		val := "code-push-value"
		val2 := "code-push-value2"

		cache := getUnAllocatedCache()

		cache.Set(key, val)

		t.Run("query after setting", func(t *testing.T) {
			getVal, getOk := cache.Get(key)
			assert.True(t, getOk)
			assert.Equal(t, val, getVal)

			queryVal, queryOk := cache.Query(key)
			assert.True(t, queryOk)
			assert.Equal(t, val, queryVal)

			assert.True(t, cache.Contains(key))
		})

		cache.Set(key, val2)
		t.Run("query after override setting", func(t *testing.T) {
			getVal, getOk := cache.Get(key)
			assert.True(t, getOk)
			assert.Equal(t, val2, getVal)

			queryVal, queryOk := cache.Query(key)
			assert.True(t, queryOk)
			assert.Equal(t, val2, queryVal)

			assert.True(t, cache.Contains(key))
		})

		assert.True(t, cache.Remove(key), "delete existed key")

		t.Run("query after deleting", func(t *testing.T) {
			queryVal, queryOk := cache.Query(key)
			assert.False(t, queryOk)
			assert.Nil(t, queryVal)

			assert.False(t, cache.Contains(key))
		})

		assert.False(t, cache.Remove(key), "delete non existed key")
	})

	t.Run("length of cache always under capacity", func(t *testing.T) {
		cache := getUnAllocatedCache(3)

		assert.LessOrEqual(t, cache.Len(), cache.Capacity())

		cache.Set(1, 1)
		assert.Less(t, cache.Len(), cache.Capacity())

		cache.Set(2, 1)
		assert.Less(t, cache.Len(), cache.Capacity())

		cache.Set(3, 1)
		assert.Equal(t, cache.Len(), cache.Capacity())

		assert.True(t, cache.Set(4, 4), "there is items evicted when len oversize")
		assert.Equal(t, cache.Len(), cache.Capacity())
	})

	t.Run("clean cache", func(t *testing.T) {
		cache := getUnAllocatedCache(3)

		cache.Set(1, 1)
		cache.Set(2, 2)
		cache.Set(3, 3)

		assert.Equal(t, cache.Len(), 3, "len before clean")

		cache.Purge()
		assert.Equal(t, cache.Len(), 0, "len after clean")
		assert.False(t, cache.Contains(1))
	})
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
