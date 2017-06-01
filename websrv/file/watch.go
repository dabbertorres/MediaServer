package file

import (
	"errors"
	"io/ioutil"
	"os"
	
	"github.com/fsnotify/fsnotify"
)

type Event struct {
	Data []byte
	Error error
}

var (
	ErrorRemoved = errors.New("Watched file was removed")
	ErrorPermissions = errors.New("Watched file's permissions changed to inaccessible")
	ErrorRenamed = errors.New("Watched file was renamed")
)

var watcher *fsnotify.Watcher

func WatchInit() (err error) {
	watcher, err = fsnotify.NewWatcher()
	return
}

func WatchStop() error {
	if watcher != nil {
		return watcher.Close()
	} else {
		return nil
	}
}

func Watch(file string, stop <-chan bool) (<-chan Event, error) {
	if err := tryOpen(file); err != nil {
		return nil, err
	}
	
	if err := watcher.Add(file); err != nil {
		return nil, err
	}
	
	event := make(chan Event)
	
	go func() {
		defer close(event)
		defer watcher.Remove(file)
		
		for {
			select {
			case _, ok := <-stop:
				if !ok {
					return
				}
			
			case ev := <-watcher.Events:
				if ev.Name == file {
					op := watchEvent(ev.Op)
					
					if op.is(fsnotify.Remove) {
						event <- Event{nil, ErrorRemoved}
						return
					}
					
					if op.is(fsnotify.Chmod) {
						if err := tryOpen(file); err != nil {
							event <- Event{nil, ErrorPermissions}
							return
						}
					}
					
					if op.is(fsnotify.Rename) {
						event <- Event{nil, ErrorRenamed}
						return
					}
					
					if op.is(fsnotify.Write) {
						data, err := ioutil.ReadFile(file)
						if err != nil {
							event <- Event{nil, err}
							return
						} else {
							event <- Event{data, nil}
						}
					}
				}
			}
		}
	}()
	
	return event, nil
}

func tryOpen(file string) error {
	if fd, err := os.OpenFile(file, os.O_RDONLY, 0); err != nil {
		return err
	} else if err = fd.Close(); err != nil {
		return err
	} else {
		return nil
	}
}

// helper type (and function) to make checking event type easy to read
type watchEvent fsnotify.Op

func (fe watchEvent) is(op fsnotify.Op) bool {
	return fsnotify.Op(fe) & op == op
}
