package pbo

import (
	"io"
	"os"
)

type File struct {
	Path    string
	Files   []*FileEntry
	Headers map[string]string

	file *pboStream
}

func Load(path string) (file *File, err error) {
	file = &File{
		Path: path,
	}
	err = file.Load()
	return
}

func (f *File) Load() (err error) {
	if f.file, err = newPboStream(f.Path); err != nil {
		return
	}

	if err = f.file.validateFile(); err != nil {
		return
	}

	if f.Headers, err = f.file.readHeaders(); err != nil {
		return
	}
	if f.Files, err = f.file.readFileEntries(f); err != nil {
		return
	}
	f.file.dataStart, err = f.file.Seek(0, io.SeekCurrent)
	return
}

func (f *File) Save() error {
	if f.file != nil {
		f.Close()
	}

	file, err := os.Create(f.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	// TODO

	return nil
}

func (f *File) CacheEnabled() bool { return f.file.cache != nil }

func (f *File) ToggleCache(enable bool) {
	if enable {
		if f.CacheEnabled() {
			f.file.cache = make(map[string][]byte)
		}
	} else {
		f.file.cache = nil
	}
}

func (f *File) ClearCache() {
	if f.file.cache != nil {
		f.file.cache = make(map[string][]byte)
	}
}

func (f *File) Close() {
	f.file.Close()
	f.file = nil
}
