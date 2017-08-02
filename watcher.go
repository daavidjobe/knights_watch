package watcher

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Watcher represents the watcher
type Watcher struct {
	Event     chan Event
	Error     chan error
	Closed    chan struct{}
	close     chan struct{}
	files     map[string]os.FileInfo
	watchList map[string]bool
	running   bool
	events    map[EventType]struct{}
	waitGroup *sync.WaitGroup
	mu        *sync.Mutex
}

// New creates a new Watcher instance
func New() *Watcher {
	var wg sync.WaitGroup
	wg.Add(1)
	return &Watcher{
		Event:     make(chan Event),
		Error:     make(chan error),
		mu:        new(sync.Mutex),
		waitGroup: &wg,
		files:     make(map[string]os.FileInfo),
		watchList: make(map[string]bool),
	}
}

// Add adds a file or folder to watch
func (w *Watcher) Add(filename string) (err error) {
	w.lock()
	filename, _ = filepath.Abs(filename) // Ignore the error. We catch it later
	fileList, err := w.list(filename)
	if err != nil {
		return err
	}
	for k, v := range fileList {
		w.files[k] = v
	}
	w.watchList[filename] = false
	w.unlock()
	return nil
}

// Remove removes a file or folder from the list of watched files
func (w *Watcher) Remove(filename string) (err error) {
	w.lock()
	filename, _ = filepath.Abs(filename) // Ignore the error. We catch it later
	_, found := w.files[filename]
	if !found {
		return &os.PathError{}
	}
	delete(w.watchList, filename)
	delete(w.files, filename)
	w.unlock()
	return nil
}

// Run this method runs the program and polls changes with a specified interval
func (w *Watcher) Run(d time.Duration) error {
	// Return an error if d is less than 1 nanosecond.
	if d < time.Nanosecond {
		return ErrDuration
	}

	// Make sure the Watcher is not already running.
	w.lock()
	if w.running {
		w.unlock()
		return ErrAlreadyRunning
	}
	w.running = true
	w.unlock()
	w.wait()
	for {
		done := make(chan struct{})
		event := make(chan Event)
		fileList := w.fetchList()
		cancel := make(chan struct{})

		// Look for events.
		go func() {
			w.pollEvents(fileList, event, cancel)
			done <- struct{}{}
		}()

	inner:
		for {
			select {
			case event := <-event:
				if len(w.events) > 0 {
					_, found := w.events[event.EventType]
					if !found {
						continue
					}
				}
				w.Event <- event
			case <-done: // Current cycle is finished.
				break inner
			}
		}
		// Update the file's list.
		w.lock()
		w.files = fileList
		w.unlock()
		time.Sleep(d)
	}
}

// FileList return the list with watched files
func (w *Watcher) FileList() map[string]os.FileInfo {
	w.lock()
	defer w.unlock()
	return w.files
}

// Close performs cleanup on the watcher instance
func (w *Watcher) Close() {
	w.mu.Lock()
	if !w.running {
		w.mu.Unlock()
		return
	}
	w.running = false
	w.files = make(map[string]os.FileInfo)
	w.watchList = make(map[string]bool)
	w.mu.Unlock()
	// Send a close signal to the Run method.
	w.close <- struct{}{}
}

func (w *Watcher) list(filename string) (map[string]os.FileInfo, error) {
	fileList := make(map[string]os.FileInfo)

	stat, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	fileList[filename] = stat
	return fileList, nil
}

func (w *Watcher) pollEvents(files map[string]os.FileInfo,
	event chan Event,
	cancel chan struct{}) {
	w.lock()
	defer w.unlock()
	creates := make(map[string]os.FileInfo)
	removes := make(map[string]os.FileInfo)

	// Get all "Removes"
	for path, info := range w.files {
		if _, found := files[path]; !found {
			removes[path] = info
		}
	}

	// Get all "Creates" & "Chmods"
	for path, info := range files {
		oldInfo, found := w.files[path]
		if !found {
			creates[path] = info
			continue
		}
		if oldInfo.ModTime() != info.ModTime() {
			select {
			case <-cancel:
				return
			case event <- Event{Write, path, info}:
			}
		}
		if oldInfo.Mode() != info.Mode() {
			select {
			case <-cancel:
				return
			case event <- Event{Chmod, path, info}:
			}
		}
	}
	// Get all "Renames" & "Moves"
	for path1, info1 := range removes {
		for path2, info2 := range creates {
			if os.SameFile(info1, info2) {
				e := Event{
					EventType: Move,
					Path:      fmt.Sprintf("%s -> %s", path1, path2),
					FileInfo:  info1,
				}
				// If they are from the same directory, it's a rename
				// instead of a move event.
				if filepath.Dir(path1) == filepath.Dir(path2) {
					e.EventType = Rename
				}

				delete(removes, path1)
				delete(creates, path2)

				select {
				case <-cancel:
					return
				case event <- e:
				}
			}
		}
	}
}

func (w *Watcher) fetchList() map[string]os.FileInfo {
	w.lock()
	defer w.unlock()
	files := make(map[string]os.FileInfo)
	var list map[string]os.FileInfo
	var err error
	for name := range w.watchList {
		list, err = w.list(name)
		if err != nil {
			if os.IsNotExist(err) {
				w.Error <- ErrWatchedFileDeleted
				w.unlock()
				w.Remove(name)
				w.lock()
			} else {
				w.Error <- err
			}
		}
	}
	// Add the file's to the file list.
	for key, value := range list {
		files[key] = value
	}
	return files
}

func (w *Watcher) wait() {
	w.waitGroup.Done()
}

func (w *Watcher) lock() {
	w.mu.Lock()
}

func (w *Watcher) unlock() {
	w.mu.Unlock()
}
