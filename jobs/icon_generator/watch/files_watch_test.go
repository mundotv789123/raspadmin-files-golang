package watch

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)


func TestFileDebouncer_SingleEvent(t *testing.T) {
	var count int32
	debouncer := NewFileDebouncer(100*time.Millisecond, func(filePath string) {
		atomic.AddInt32(&count, 1)
	})

	debouncer.AddEvent("/test/file1.png")

	// Wait less than delay
	time.Sleep(50 * time.Millisecond)
	if atomic.LoadInt32(&count) != 0 {
		t.Fatalf("expected 0 calls before delay expired, got %d", count)
	}

	// Wait for delay to expire
	time.Sleep(100 * time.Millisecond)
	if atomic.LoadInt32(&count) != 1 {
		t.Fatalf("expected 1 call after delay expired, got %d", count)
	}
}

func TestFileDebouncer_MultipleEventsResetsTimer(t *testing.T) {
	var count int32
	var lastPath string
	debouncer := NewFileDebouncer(150*time.Millisecond, func(filePath string) {
		atomic.AddInt32(&count, 1)
		lastPath = filePath
	})

	// Fire event at T=0
	debouncer.AddEvent("/test/file1.png")

	// Fire event at T=80ms (before 150ms expires)
	time.Sleep(80 * time.Millisecond)
	debouncer.AddEvent("/test/file1.png")

	// Fire event at T=160ms (before 80+150=230ms expires)
	time.Sleep(80 * time.Millisecond)
	debouncer.AddEvent("/test/file1.png")

	// Check at T=200ms (160 + 40ms) - should NOT have executed yet
	time.Sleep(40 * time.Millisecond)
	if atomic.LoadInt32(&count) != 0 {
		t.Fatalf("expected 0 calls due to timer resets, got %d", count)
	}

	// Wait until T=350ms (> 160 + 150 = 310ms)
	time.Sleep(150 * time.Millisecond)
	if atomic.LoadInt32(&count) != 1 {
		t.Fatalf("expected exactly 1 call after reset timer expired, got %d", count)
	}
	if lastPath != "/test/file1.png" {
		t.Fatalf("expected path /test/file1.png, got %s", lastPath)
	}
}

func TestFileDebouncer_MultipleFilesIndependent(t *testing.T) {
	var countFile1 int32
	var countFile2 int32

	debouncer := NewFileDebouncer(100*time.Millisecond, func(filePath string) {
		if filePath == "/test/file1.png" {
			atomic.AddInt32(&countFile1, 1)
		} else if filePath == "/test/file2.png" {
			atomic.AddInt32(&countFile2, 1)
		}
	})

	debouncer.AddEvent("/test/file1.png")
	time.Sleep(30 * time.Millisecond)
	debouncer.AddEvent("/test/file2.png")

	time.Sleep(90 * time.Millisecond) // file1 should be done (120ms total), file2 not yet (90ms since file2)
	if atomic.LoadInt32(&countFile1) != 1 {
		t.Fatalf("expected file1 to be 1, got %d", countFile1)
	}
	if atomic.LoadInt32(&countFile2) != 0 {
		t.Fatalf("expected file2 to be 0, got %d", countFile2)
	}

	time.Sleep(50 * time.Millisecond)
	if atomic.LoadInt32(&countFile2) != 1 {
		t.Fatalf("expected file2 to be 1, got %d", countFile2)
	}
}

func TestEventQueueWorker_SequentialExecution(t *testing.T) {
	var executed []string
	var mu sync.Mutex

	worker := NewEventQueueWorker(10, func(filePath string) {
		time.Sleep(50 * time.Millisecond)
		mu.Lock()
		executed = append(executed, filePath)
		mu.Unlock()
	})

	worker.Push("/file1.mp3")
	worker.Push("/file2.mp3")
	worker.Push("/file3.mp3")

	time.Sleep(250 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(executed) != 3 {
		t.Fatalf("expected 3 items executed, got %d", len(executed))
	}
	if executed[0] != "/file1.mp3" || executed[1] != "/file2.mp3" || executed[2] != "/file3.mp3" {
		t.Fatalf("expected sequential order [/file1.mp3, /file2.mp3, /file3.mp3], got %v", executed)
	}
}

