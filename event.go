package watcher

import (
	"fmt"
	"os"
)

// Event holds the type of event e.g "Create" or "Write".
type Event struct {
	EventType
	Path string
	os.FileInfo
}

// EventType the type of event
type EventType uint32

const (
	// Create event
	Create EventType = iota
	// Write event
	Write
	// Remove event
	Remove
	// Rename event
	Rename
	// Chmod event
	Chmod
	// Move event
	Move
)

var eventTypes = map[EventType]string{
	Create: "CREATE",
	Write:  "WRITE",
	Remove: "REMOVE",
	Rename: "RENAME",
	Chmod:  "CHMOD",
	Move:   "MOVE",
}

// Describe The string representation of the event
func (e Event) Describe() (string, error) {
	if e.FileInfo != nil {
		return fmt.Sprintf("%s -> %s", eventTypes[e.EventType], e.Path), nil
	}
	return "", ErrFileMissing
}
