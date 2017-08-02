package main

import (
	"fmt"
	"github.com/daavidjobe/knights_watch"
	"log"
	"os"
	"os/signal"
)

func main() {
	w := watcher.New()
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			select {
			case event := <-w.Event:
				fmt.Println(event.Describe())

			case err := <-w.Error:
				if err == watcher.ErrWatchedFileDeleted {
					fmt.Println(err)
					continue
				}
				log.Fatalln(err)
			}
		}

	}()

	// for testing
	if err := w.Add("README.md"); err != nil {
		log.Fatalln(err)
	}

	for path, f := range w.FileList() {
		fmt.Printf("%s: %s\n", f.Name(), path)
	}
	fmt.Println()
	fmt.Printf("Watching %d wildlings\n", len(w.FileList()))

	closed := make(chan struct{})
	c := make(chan os.Signal)
	signal.Notify(c, os.Kill, os.Interrupt)

	go func() {
		<-c
		w.Close()
		<-done
		fmt.Println("knights watch ended")
		close(closed)
	}()

	// Start the watching process
	if err := w.Run(200); err != nil {
		log.Fatalln(err)
	}

	<-closed
}
