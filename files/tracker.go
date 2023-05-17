package files

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type StatusTracker struct {
	files        []string
	totalFiles   int
	filesDone    int32
	filesErrored int32
}

// NewStatusTracker creates a new StatusTracker.
func NewStatusTracker(files []string) *StatusTracker {
	return &StatusTracker{
		files:        files,
		totalFiles:   len(files),
		filesDone:    0,
		filesErrored: 0,
	}
}

// ProcessFiles processes all the files and updates the status.
func (s *StatusTracker) ProcessFiles(mode string, dryRun bool) error {
	// Capture the start time
	startTime := time.Now()
	fmt.Printf("⏳⏳⏳ START Processing mode=%s dryRun=%v\n\n", mode, dryRun)

	var wg sync.WaitGroup
	for _, file := range s.files {
		wg.Add(1)
		// @todo: replace values with actual incoming flags (mode string, dryRun bool)
		go s.processFile(file, "encrypt", true, &wg)
	}
	wg.Wait()

	// Capture the end time
	endTime := time.Now()

	// Calculate the difference and convert it to a duration
	duration := endTime.Sub(startTime)

	// Log out the total time taken
	fmt.Printf("Processed all files in %v\n", duration)

	return nil
}

func (s *StatusTracker) processFile(file string, mode string, dryRun bool, wg *sync.WaitGroup) {
	defer wg.Done() // signal that this goroutine is done when the function exits

	// Capture the start time
	startTime := time.Now()

	// TODO: Replace this with actual file processing code
	// if mode is "encrypt", encrypt the files
	// if mode is "decrypt", decrypt the files
	// if dryRun is true, don't actually change the files but simulate the process
	fmt.Printf("⏳ Processing file=%s\n", file)
	time.Sleep(time.Duration(rand.Intn(3)) * time.Second)

	// Capture the end time
	endTime := time.Now()

	// Calculate the difference and convert it to a duration
	duration := endTime.Sub(startTime)

	// Log out the duration
	fmt.Printf("✅ Processed %s in %v\n", file, duration)

	// Update the count of files processed in a thread-safe manner
	// atomic.AddInt32 adds 1 to the int32 value at &s.filesDone
	// The atomic package provides functions for performing atomic operations
	// which are safe to use concurrently (from multiple goroutines) without additional locking
	atomic.AddInt32(&s.filesDone, 1)
}

// GetStatus returns the current status.
func (s *StatusTracker) GetStatus() string {
	return fmt.Sprintf("Processed %d of %d files. Failed: %d", s.filesDone, s.totalFiles, s.filesErrored)
}
