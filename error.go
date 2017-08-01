package watcher

import "errors"

var (
	// ErrFileMissing occurs when files are missing.
	ErrFileMissing = errors.New("error: file missing")
	// ErrUnknownOperation occurs if a file event is not handeled.
	// By default we are watching for create, delete, modify and move.
	ErrUnknownOperation = errors.New("error: unknown operation")
	// ErrDuration if the polling interval is 1 nanosecond or below.
	ErrDuration = errors.New("error: duration is to short")
	// ErrAlreadyRunning if an instance is already running.
	ErrAlreadyRunning = errors.New("error: watcher is already running")
	// ErrWatchedFileDeleted if the file is deleted.
	ErrWatchedFileDeleted = errors.New("error: watched file has been deleted")
)
