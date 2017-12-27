package pbo

import (
	"io"
)

// File contains methods to handle PBO files.
type File struct {
	Path    string            // Path contains the path to the current PBO file.
	Files   []*FileEntry      // Files contains all files within the PBO file.
	Headers map[string]string // Headers contains all header fields of the PBO file.

	file *pboStream
}

// Load loads the PBO file.
func Load(path string) (file *File, err error) {
	file = &File{
		Path: path,
	}
	err = file.Load()
	return
}

// Load (re)loads the PBO file.
func (f *File) Load() (err error) {
	if f.file != nil {
		f.Close()
	}

	if f.file, err = openPBO(f.Path); err != nil {
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

// CacheEnabled returns true if file caching is enabled.
func (f *File) CacheEnabled() bool { return f.file.cache != nil }

// ToggleCache enables/disables the file caching mechanism according to enable.
func (f *File) ToggleCache(enable bool) {
	if enable {
		if f.CacheEnabled() {
			f.file.cache = make(map[string][]byte)
		}
	} else {
		f.file.cache = nil
	}
}

// ClearCache clears the file cache.
func (f *File) ClearCache() {
	if f.file.cache != nil {
		f.file.cache = make(map[string][]byte)
	}
}

// Close closes the file stream and clears the file cache.
// You can't read any files after this anymore!
func (f *File) Close() {
	f.file.Close()
	f.file = nil
}
