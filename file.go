package filecacher

import (
	"bytes"
	"io"
	"os"
	"sync"

	"github.com/PathDNA/atoms"
	"github.com/missionMeteora/journaler"

	"github.com/Hatch1fy/errors"
	"github.com/fsnotify/fsnotify"
)

// NewFile will return a new file
func NewFile(filename string) (fp *File, err error) {
	var f File
	if f.w, err = fsnotify.NewWatcher(); err != nil {
		return
	}

	if err = f.w.Add(filename); err != nil {
		return
	}

	f.out = journaler.New("FileCacher", filename)
	f.b = bytes.NewBuffer(nil)
	f.filename = filename

	if err = f.refreshBuffer(); err != nil {
		return
	}

	go f.watch()
	fp = &f
	return
}

// File represents a file
type File struct {
	mu  sync.RWMutex
	out *journaler.Journaler

	w *fsnotify.Watcher
	b *bytes.Buffer

	filename string

	closed atoms.Bool
}

func (f *File) watch() {
	for {
		select {
		case event, ok := <-f.w.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				if err := f.refreshBuffer(); err != nil {
					f.out.Error("Error refreshing file: %v", err)
				}
			}

		case err, ok := <-f.w.Errors:
			if !ok {
				return
			}

			f.out.Error("Error event encountered: %v", err)
			f.Close()
		}
	}
}

func (f *File) refreshBuffer() (err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	var tgt *os.File
	if tgt, err = os.Open(f.filename); err != nil {
		return
	}
	defer tgt.Close()

	f.b.Reset()
	if _, err = io.Copy(f.b, tgt); err != nil {
		return
	}

	return
}

// Read will read a file
func (f *File) Read(fn func(io.Reader) error) (err error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if f.closed.Get() {
		return errors.ErrIsClosed
	}

	r := bytes.NewReader(f.b.Bytes())

	return fn(r)
}

// Close will close a file
func (f *File) Close() (err error) {
	if !f.closed.Set(true) {
		return errors.ErrIsClosed
	}

	return f.w.Close()
}
