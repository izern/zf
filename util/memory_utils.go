package util

import (
	"fmt"
	"runtime"
)

// MemoryPool provides object pooling to reduce allocations
type MemoryPool struct {
	stringMaps chan map[string]interface{}
	interfaces chan []interface{}
}

// NewMemoryPool creates a new memory pool with the specified capacity
func NewMemoryPool(capacity int) *MemoryPool {
	return &MemoryPool{
		stringMaps: make(chan map[string]interface{}, capacity),
		interfaces: make(chan []interface{}, capacity),
	}
}

// GetStringMap returns a reusable map[string]interface{} from the pool
func (p *MemoryPool) GetStringMap() map[string]interface{} {
	select {
	case m := <-p.stringMaps:
		// Clear the map before reuse
		for k := range m {
			delete(m, k)
		}
		return m
	default:
		return make(map[string]interface{})
	}
}

// PutStringMap returns a map to the pool for reuse
func (p *MemoryPool) PutStringMap(m map[string]interface{}) {
	if m == nil || len(m) > 1000 { // Don't pool very large maps
		return
	}
	
	select {
	case p.stringMaps <- m:
	default:
		// Pool is full, discard
	}
}

// GetInterfaceSlice returns a reusable []interface{} from the pool
func (p *MemoryPool) GetInterfaceSlice() []interface{} {
	select {
	case s := <-p.interfaces:
		return s[:0] // Reset length but keep capacity
	default:
		return make([]interface{}, 0, 16) // Start with reasonable capacity
	}
}

// PutInterfaceSlice returns a slice to the pool for reuse
func (p *MemoryPool) PutInterfaceSlice(s []interface{}) {
	if s == nil || cap(s) > 1000 { // Don't pool very large slices
		return
	}
	
	select {
	case p.interfaces <- s:
	default:
		// Pool is full, discard
	}
}

// Global memory pool instance
var defaultPool = NewMemoryPool(100)

// GetPooledStringMap gets a map from the default pool
func GetPooledStringMap() map[string]interface{} {
	return defaultPool.GetStringMap()
}

// ReturnPooledStringMap returns a map to the default pool
func ReturnPooledStringMap(m map[string]interface{}) {
	defaultPool.PutStringMap(m)
}

// GetPooledInterfaceSlice gets a slice from the default pool
func GetPooledInterfaceSlice() []interface{} {
	return defaultPool.GetInterfaceSlice()
}

// ReturnPooledInterfaceSlice returns a slice to the default pool
func ReturnPooledInterfaceSlice(s []interface{}) {
	defaultPool.PutInterfaceSlice(s)
}

// OptimizedConvertMap2String uses memory pooling for better performance
func OptimizedConvertMap2String(m map[interface{}]interface{}) map[string]interface{} {
	if len(m) == 0 {
		return GetPooledStringMap()
	}

	res := GetPooledStringMap()
	if len(res) == 0 && len(m) > 16 {
		// If pooled map is empty and we need a larger map, create a new one
		ReturnPooledStringMap(res)
		res = make(map[string]interface{}, len(m))
	}

	for k, v := range m {
		key := fmt.Sprint(k)
		switch val := v.(type) {
		case map[interface{}]interface{}:
			res[key] = ConvertMap2String(val)
		case map[string]interface{}:
			res[key] = val // Already correct type
		case []interface{}:
			res[key] = ConvertArray2String(val)
		case []map[interface{}]interface{}:
			result := make([]map[string]interface{}, len(val))
			for i, item := range val {
				result[i] = ConvertMap2String(item)
			}
			res[key] = result
		case []map[string]interface{}:
			res[key] = val // Already correct type
		default:
			res[key] = v // Primitive types don't need conversion
		}
	}
	return res
}

// ForceGC forces garbage collection when processing large files
func ForceGC() {
	runtime.GC()
	runtime.GC() // Call twice for better cleanup
}

// EstimateMemoryUsage estimates memory usage of a data structure
func EstimateMemoryUsage(v interface{}) int64 {
	var size int64
	
	switch val := v.(type) {
	case map[string]interface{}:
		size += int64(len(val)) * 32 // Rough estimate for map overhead
		for k, v := range val {
			size += int64(len(k)) + EstimateMemoryUsage(v)
		}
	case map[interface{}]interface{}:
		size += int64(len(val)) * 32
		for k, v := range val {
			size += EstimateMemoryUsage(k) + EstimateMemoryUsage(v)
		}
	case []interface{}:
		size += int64(len(val)) * 8 // 8 bytes per pointer on 64-bit
		for _, item := range val {
			size += EstimateMemoryUsage(item)
		}
	case string:
		size += int64(len(val))
	case int, int32, int64, uint, uint32, uint64:
		size += 8
	case float32, float64:
		size += 8
	case bool:
		size += 1
	default:
		size += 16 // Default estimate
	}
	
	return size
}

// ShouldUseStreaming determines if streaming should be used based on data size
func ShouldUseStreaming(data []byte) bool {
	const streamingThreshold = 10 * 1024 * 1024 // 10MB
	return len(data) > streamingThreshold
}