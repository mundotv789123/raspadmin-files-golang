package watch

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/farmergreg/rfsnotify"
	"github.com/mundotv789123/raspadmin/internal/config"
	"github.com/mundotv789123/raspadmin/internal/validations"
)

type FileHandlerFunc func(filePath string)

type EventQueueWorker struct {
	queue   chan string
	handler FileHandlerFunc
	set     map[string]bool
	mu      sync.Mutex
}

func NewEventQueueWorker(bufferSize int, handler FileHandlerFunc) *EventQueueWorker {
	qw := &EventQueueWorker{
		queue:   make(chan string, bufferSize),
		handler: handler,
		set:     make(map[string]bool),
	}
	go qw.start()
	return qw
}

func (qw *EventQueueWorker) Push(filePath string) {
	qw.mu.Lock()
	if qw.set[filePath] {
		qw.mu.Unlock()
		return
	}
	qw.set[filePath] = true
	qw.mu.Unlock()

	qw.queue <- filePath
}

func (qw *EventQueueWorker) start() {
	for filePath := range qw.queue {
		qw.mu.Lock()
		delete(qw.set, filePath)
		qw.mu.Unlock()

		if qw.handler != nil {
			qw.handler(filePath)
		}
	}
}

type FileDebouncer struct {
	mu      sync.Mutex
	timers  map[string]*time.Timer
	delay   time.Duration
	handler FileHandlerFunc
}

func NewFileDebouncer(delay time.Duration, handler FileHandlerFunc) *FileDebouncer {
	return &FileDebouncer{
		timers:  make(map[string]*time.Timer),
		delay:   delay,
		handler: handler,
	}
}

func (d *FileDebouncer) AddEvent(path string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, exists := d.timers[path]; exists {
		t.Stop()
	}

	d.timers[path] = time.AfterFunc(d.delay, func() {
		d.mu.Lock()
		delete(d.timers, path)
		d.mu.Unlock()

		d.handler(path)
	})
}

func InitWatch(path string, handler FileHandlerFunc) {
	watcher, err := rfsnotify.NewWatcher()
	if err != nil {
		slog.Error(fmt.Sprintf("error creating watcher: %v", err))
		return
	}
	defer watcher.Close()

	done := make(chan bool)
	go watch(watcher, handler)
	err = watcher.AddRecursive(path)
	if err != nil {
		slog.Error(fmt.Sprintf("error adding recursive watch path: %v", err))
		return
	}
	slog.Info(fmt.Sprintf("started watching path: %s", path))
	<-done
}

func watch(watcher *rfsnotify.RWatcher, handler FileHandlerFunc) {
	queueWorker := NewEventQueueWorker(10000, handler)

	debouncer := NewFileDebouncer(15*time.Second, func(filePath string) {
		slog.Info(fmt.Sprintf("debounced timer expired, queuing file: %s", filePath))
		queueWorker.Push(filePath)
	})

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			filePath := filepath.Clean(event.Name)

			// Ignore events occurring in cache directory
			if strings.HasPrefix(filePath, config.CacheDirAds) {
				continue
			}

			// Ignore temporary or hidden files
			if validations.IsTempOrHiddenFile(filepath.Base(filePath)) {
				continue
			}

			slog.Debug(fmt.Sprintf("watch event received for %s: %s", filePath, event.Op.String()))
			debouncer.AddEvent(filePath)

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			slog.Error(fmt.Sprintf("watcher error: %v", err))
		}
	}
}


