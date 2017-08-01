package watcher

import "errors"

var (
	// ErrFileMissing occurs when files are missing
	ErrFileMissing = errors.New("error: file missing")
	// ErrUnknownOperation occurs if a file event is not handeled
	ErrUnknownOperation = errors.New("error: unknown operation")
	// ErrDuration if the polling interval is to low
	ErrDuration = errors.New("error: duration is to short")
	// ErrAlreadyRunning if an instance is already running
	ErrAlreadyRunning = errors.New("error: watcher is already running")
	// ErrWatchedFileDeleted if the file is deleted
	ErrWatchedFileDeleted = errors.New("error: watched file has been deleted")
)
