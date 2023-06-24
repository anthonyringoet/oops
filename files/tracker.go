package files

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/anthonyringoet/oops/crypto"
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
		go s.processFile(file, mode, dryRun, &wg)
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

	var processingError bool
	// Capture the start time
	startTime := time.Now()

	fmt.Printf("⏳ Processing file=%s\n", file)

	if dryRun {
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
		fmt.Printf("👉 Dry run, not actually processing file=%s\n", file)
	} else {
		if mode == "encrypt" {
			// ideally, dont have to load keys again for every processFile call
			pubKey, err := crypto.LoadPublicKey(crypto.PublicKeyPath)
			if err != nil {
				fmt.Printf("Error loading public key: %s\n", err)
			}
			err = crypto.EncryptFile(file, pubKey)
			if err != nil {
				fmt.Printf("Error decrypting file: %s\n", err)
				atomic.AddInt32(&s.filesErrored, 1)
				processingError = true
			}
		} else if mode == "decrypt" {
			privKey, err := crypto.LoadPrivateKey(crypto.PrivateKeyPath)
			if err != nil {
				fmt.Printf("Error loading private key: %s\n", err)
			}
			err = crypto.DecryptFile(file, privKey)
			if err != nil {
				fmt.Printf("Error decrypting file: %s\n", err)
				atomic.AddInt32(&s.filesErrored, 1)
				processingError = true
			}
		}
	}

	// Capture the end time
	endTime := time.Now()

	// Calculate the difference and convert it to a duration
	duration := endTime.Sub(startTime)

	// Log out the duration
	fmt.Printf("✅ Processed %s in %v\n", file, duration)

	if !processingError {
		// Update the count of files processed in a thread-safe manner
		// atomic.AddInt32 adds 1 to the int32 value at &s.filesDone
		// The atomic package provides functions for performing atomic operations
		// which are safe to use concurrently (from multiple goroutines) without additional locking
		atomic.AddInt32(&s.filesDone, 1)
	}
}

// GetStatus returns the current status.
func (s *StatusTracker) GetStatus() string {
	return fmt.Sprintf("Processed %d of %d files. Failed: %d", s.filesDone, s.totalFiles, s.filesErrored)
}
