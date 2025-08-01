package util

import (
	"bufio"
	"io"
	"sync"
	"github.com/izern/zf/types"
)

// StreamProcessor handles large file processing with streaming
type StreamProcessor struct {
	chunkSize int
	pool      *MemoryPool
}

// NewStreamProcessor creates a new stream processor
func NewStreamProcessor(chunkSize int) *StreamProcessor {
	return &StreamProcessor{
		chunkSize: chunkSize,
		pool:      NewMemoryPool(10),
	}
}

// ProcessLargeFile processes large files in chunks to reduce memory usage
func (sp *StreamProcessor) ProcessLargeFile(reader io.Reader, processor func([]byte) error) error {
	scanner := bufio.NewScanner(reader)
	
	// Increase buffer size for large files
	buf := make([]byte, 0, sp.chunkSize)
	scanner.Buffer(buf, sp.chunkSize)
	
	for scanner.Scan() {
		if err := processor(scanner.Bytes()); err != nil {
			return err
		}
	}
	
	return scanner.Err()
}

// ConcurrentProcessor handles concurrent processing of data chunks
type ConcurrentProcessor struct {
	workerCount int
	jobQueue    chan ProcessJob
	wg          sync.WaitGroup
}

// ProcessJob represents a unit of work
type ProcessJob struct {
	Data     []byte
	Callback func([]byte) ([]byte, error)
	Result   chan ProcessResult
}

// ProcessResult contains the result of processing
type ProcessResult struct {
	Data  []byte
	Error error
}

// NewConcurrentProcessor creates a new concurrent processor
func NewConcurrentProcessor(workerCount int) *ConcurrentProcessor {
	return &ConcurrentProcessor{
		workerCount: workerCount,
		jobQueue:    make(chan ProcessJob, workerCount*2),
	}
}

// Start begins the concurrent processing workers
func (cp *ConcurrentProcessor) Start() {
	for i := 0; i < cp.workerCount; i++ {
		cp.wg.Add(1)
		go cp.worker()
	}
}

// Stop stops all workers and waits for completion
func (cp *ConcurrentProcessor) Stop() {
	close(cp.jobQueue)
	cp.wg.Wait()
}

// Submit submits a job for processing
func (cp *ConcurrentProcessor) Submit(data []byte, callback func([]byte) ([]byte, error)) <-chan ProcessResult {
	result := make(chan ProcessResult, 1)
	job := ProcessJob{
		Data:     data,
		Callback: callback,
		Result:   result,
	}
	
	select {
	case cp.jobQueue <- job:
	default:
		// Queue is full, process synchronously
		output, err := callback(data)
		result <- ProcessResult{Data: output, Error: err}
		close(result)
	}
	
	return result
}

// worker processes jobs from the queue
func (cp *ConcurrentProcessor) worker() {
	defer cp.wg.Done()
	
	for job := range cp.jobQueue {
		output, err := job.Callback(job.Data)
		job.Result <- ProcessResult{Data: output, Error: err}
		close(job.Result)
	}
}

// BatchProcessor processes multiple items efficiently
type BatchProcessor struct {
	batchSize int
	timeout   int // milliseconds
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor(batchSize int, timeoutMs int) *BatchProcessor {
	return &BatchProcessor{
		batchSize: batchSize,
		timeout:   timeoutMs,
	}
}

// ProcessBatch processes items in batches for better performance
func (bp *BatchProcessor) ProcessBatch(items []interface{}, processor func([]interface{}) error) error {
	for i := 0; i < len(items); i += bp.batchSize {
		end := i + bp.batchSize
		if end > len(items) {
			end = len(items)
		}
		
		batch := items[i:end]
		if err := processor(batch); err != nil {
			return err
		}
		
		// Force GC periodically to manage memory
		if i%1000 == 0 {
			ForceGC()
		}
	}
	
	return nil
}

// FastStringBuilder provides optimized string building for large outputs
type FastStringBuilder struct {
	buffer []byte
	pool   *MemoryPool
}

// NewFastStringBuilder creates a new fast string builder
func NewFastStringBuilder(initialCapacity int) *FastStringBuilder {
	return &FastStringBuilder{
		buffer: make([]byte, 0, initialCapacity),
		pool:   defaultPool,
	}
}

// Write appends data to the builder
func (fsb *FastStringBuilder) Write(data []byte) {
	fsb.buffer = append(fsb.buffer, data...)
}

// WriteString appends a string to the builder
func (fsb *FastStringBuilder) WriteString(s string) {
	fsb.buffer = append(fsb.buffer, s...)
}

// String returns the built string
func (fsb *FastStringBuilder) String() string {
	return string(fsb.buffer)
}

// Reset clears the builder for reuse
func (fsb *FastStringBuilder) Reset() {
	fsb.buffer = fsb.buffer[:0]
}

// Len returns the current length
func (fsb *FastStringBuilder) Len() int {
	return len(fsb.buffer)
}

// Cap returns the current capacity
func (fsb *FastStringBuilder) Cap() int {
	return cap(fsb.buffer)
}

// OptimizedArrayCopy provides faster array copying with bounds checking
func OptimizedArrayCopy(src, dest []interface{}, srcPos, destPos, length int) types.ZfError {
	// Bounds checking
	if srcPos < 0 || destPos < 0 || length < 0 {
		return types.NewUnSupportError("negative indices not allowed")
	}
	
	if srcPos+length > len(src) {
		return types.NewIndexOutOfBoundError(len(src), "src", srcPos+length-1)
	}
	
	if destPos+length > len(dest) {
		return types.NewIndexOutOfBoundError(len(dest), "dest", destPos+length-1)
	}
	
	// Use built-in copy for better performance
	copy(dest[destPos:destPos+length], src[srcPos:srcPos+length])
	
	return nil
}

// CacheAwareProcessor adapts processing based on available memory
type CacheAwareProcessor struct {
	memoryThreshold int64
	highMemoryMode  bool
}

// NewCacheAwareProcessor creates a processor that adapts to memory pressure
func NewCacheAwareProcessor(memoryThreshold int64) *CacheAwareProcessor {
	return &CacheAwareProcessor{
		memoryThreshold: memoryThreshold,
		highMemoryMode:  false,
	}
}

// ShouldUseHighMemoryMode determines if we should optimize for speed over memory
func (cap *CacheAwareProcessor) ShouldUseHighMemoryMode(dataSize int64) bool {
	estimatedMemory := dataSize * 3 // Rough estimate including temporary objects
	cap.highMemoryMode = estimatedMemory < cap.memoryThreshold
	return cap.highMemoryMode
}

// ProcessWithAdaptiveStrategy processes data with memory-aware optimizations
func (cap *CacheAwareProcessor) ProcessWithAdaptiveStrategy(data []byte, processor func([]byte, bool) ([]byte, error)) ([]byte, error) {
	useHighMemory := cap.ShouldUseHighMemoryMode(int64(len(data)))
	
	if !useHighMemory {
		// Force GC before processing large data
		ForceGC()
	}
	
	result, err := processor(data, useHighMemory)
	
	if !useHighMemory {
		// Force GC after processing large data
		ForceGC()
	}
	
	return result, err
}