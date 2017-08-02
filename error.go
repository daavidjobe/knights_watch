package watcher

import "errors"

var (
	// ErrFileMissing occurs when files are missing
	ErrFileMissing = errors.New("error: wildling is missing")
	// ErrUnknownOperation occurs if a file event is not handeled
	ErrUnknownOperation = errors.New("error: unknown dragonglass")
	// ErrDuration if the polling interval is to low
	ErrDuration = errors.New("error: unleashing arrows to quickly")
	// ErrAlreadyRunning if an instance is already running
	ErrAlreadyRunning = errors.New("error: watchers are already on the wall")
	// ErrWatchedFileDeleted if the file is deleted
	ErrWatchedFileDeleted = errors.New("error: wildling is vanquished")
)
