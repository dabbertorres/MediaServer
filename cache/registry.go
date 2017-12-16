package cache

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type IgnoreFunc func(path string, info os.FileInfo) bool

type dataMap map[string][]byte
type listenerMap map[string]chan []byte

type File struct {
	Name string
	Data []byte
}

type Registry struct {
	BasePath   string
	filesMutex sync.RWMutex
	files      dataMap
	listeners  listenerMap
	watcher    *fsnotify.Watcher
}

func NewRegistry(basePath string) (reg *Registry, err error) {
	reg = &Registry{
		BasePath:  filepath.Clean(basePath),
		files:     make(dataMap),
		listeners: make(listenerMap),
	}

	reg.watcher, err = fsnotify.NewWatcher()
	if err == nil {
		go reg.watch()
	}

	return
}

func (reg *Registry) Close() {
	for _, c := range reg.listeners {
		close(c)
	}
	reg.watcher.Close()
}

func (reg *Registry) Paths() []string {
	ret := make([]string, len(reg.files))

	i := 0
	for k := range reg.files {
		ret[i] = k
		i++
	}

	return ret
}

// Filter returns the files with names matching the given regex
func (reg *Registry) Filter(filter string) []File {
	reg.filesMutex.RLock()
	defer reg.filesMutex.RUnlock()

	ret := make([]File, 0, len(reg.files))

	for name, data := range reg.files {
		if match, _ := regexp.MatchString(filter, name); match {
			ret = append(ret, File{name, data})
		}
	}

	return ret
}

// Walk walks down Registry.BasePath and adds all files found
func (reg *Registry) Walk() error {
	return filepath.Walk(reg.BasePath, reg.add)
}

// WalkButIgnore walks down Registry.BasePath and calls ignore
// for each child found.
// If ignore returns true, the cache/sub-directory is skipped.
// Otherwise, the cache is added (or the sub-directory is walked).
func (reg *Registry) WalkButIgnore(ignore IgnoreFunc) error {
	return filepath.Walk(reg.BasePath,
		func(path string, info os.FileInfo, err error) error {
			if ignore(path, info) {
				fi, err := os.Stat(path)
				if err != nil {
					log.Printf("Error accessing cache '%s': %v\n", path, err)
					return nil
				}

				if fi.IsDir() {
					return filepath.SkipDir
				} else {
					return nil
				}
			}

			return reg.add(path, info, err)
		})
}

func (reg *Registry) Get(file string) []byte {
	reg.filesMutex.RLock()
	defer reg.filesMutex.RUnlock()

	ret := reg.files[file]
	if ret == nil {
		return nil
	} else {
		return ret
	}
}

func (reg *Registry) ListenTo(file string) <-chan []byte {
	c, ok := reg.listeners[file]
	if !ok {
		c = make(chan []byte)
		reg.listeners[file] = c
	}

	return c
}

func (reg *Registry) set(file string, data []byte) {
	reg.filesMutex.Lock()
	reg.files[filepath.ToSlash(file)] = data
	reg.filesMutex.Unlock()
}

func (reg *Registry) add(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Printf("Error accessing cache '%s': %v\n", path, err)
		return nil
	}

	// not interested in directories or hidden files
	if info.IsDir() || info.Name()[0] == '.' {
		return nil
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Could not read file '%s': %v\n", path, err)
		return nil
	}

	reg.watcher.Add(path)

	// map the cache path relative to the Base path for ease of use
	path, _ = filepath.Rel(reg.BasePath, path)
	reg.set(path, data)

	return nil
}

func (reg *Registry) watch() {
	for {
		select {
		case ev := <-reg.watcher.Events:
			if is(ev.Op, fsnotify.Remove) {
				log.Printf("File '%s' was removed - no longer watching.\n", ev.Name)
				reg.watcher.Remove(ev.Name)
				break
			}

			if is(ev.Op, fsnotify.Chmod) {
				if err := tryOpen(ev.Name); err != nil {
					log.Printf("File '%s' change made it inaccessible - no longer watching.\n", ev.Name)
					reg.watcher.Remove(ev.Name)
					break
				}
			}

			if is(ev.Op, fsnotify.Write) || is(ev.Op, fsnotify.Rename) {
				data, err := ioutil.ReadFile(ev.Name)
				if err == nil {
					log.Printf("Reloading '%s'\n", ev.Name)
					path, _ := filepath.Rel(reg.BasePath, ev.Name)
					reg.set(path, data)
					go reg.notify(ev.Name)
				} else {
					log.Printf("Error (%s) reading modified cache '%s' - no longer watching.\n", err, ev.Name)
					reg.watcher.Remove(ev.Name)
					break
				}
			}

		case err := <-reg.watcher.Errors:
			log.Println("File watch error:", err)
		}
	}
}

func (reg *Registry) notify(file string) {
	c, ok := reg.listeners[file]
	if ok {
		reg.filesMutex.RLock()
		c <- reg.files[file]
		reg.filesMutex.RUnlock()
	}
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

func is(op1, op2 fsnotify.Op) bool {
	return op1&op2 == op2
}
