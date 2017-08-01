package watcher

import (
	"os"
	"reflect"
	"testing"
)

func TestAddFile(t *testing.T) {
	w := New()
	expectedMissing := &os.PathError{}
	actualMissing := w.Add("sdasd")
	if reflect.TypeOf(actualMissing) != reflect.TypeOf(expectedMissing) {
		t.Fatal("Added missing file without complaints")
	}
	var expectedExist error
	actualExist := w.Add("./watcher.go")
	if actualExist != expectedExist {
		t.Fatal("Existing file could not be added")
	}
	checkListSize(t, w, 1)
}

func TestRemoveFile(t *testing.T) {
	w := New()
	w.Add("./watcher.go")
	expectedMissing := &os.PathError{}
	actualMissing := w.Remove("asd")
	if reflect.TypeOf(actualMissing) != reflect.TypeOf(expectedMissing) {
		t.Fatal("Removed missing file without complaints")
	}
	var expectedExist error
	actualExist := w.Remove("./watcher.go")
	if actualExist != expectedExist {
		t.Fatal("Existing file could not be removed")
	}
	checkListSize(t, w, 0)
}

func TestRun(t *testing.T) {
	w := New()
	w.Add("./watcher.go")
	w.Run(100)
	t.Fatal("Run is not tested")
}

func checkListSize(t *testing.T, w *Watcher, expected int) {
	fileCount := len(w.files)
	watchCount := len(w.watchList)
	if fileCount != expected || watchCount != expected {
		t.Fatalf("Lists was not properly updated. files: %d, watchList: %d", fileCount, watchCount)
	}
}
