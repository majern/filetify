package shared

import (
	"fmt"
	"github.com/msoft-dev/filetify/pkg/shared"
	"sync"
	"testing"
)

var cache = shared.StaticCache()

func setValue(idx int) {
	kv := getKV(idx)
	cache.GetOrSet(kv, kv)
}

func getValue(idx int) interface{} {
	kv := getKV(idx)
	return cache.GetOrSet(kv, nil)
}

func getKV(idx int) string {
	return fmt.Sprintf("Key/Value_#%d", idx)
}

func BenchmarkTestSetValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		setValue(i)
	}
}

func BenchmarkTestSetValueConcurrent(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(kv int) {
			setValue(kv)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func BenchmarkTestGetValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getValue(i)
	}
}

func BenchmarkTestGetValueConcurrent(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(kv int) {
			getValue(kv)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
