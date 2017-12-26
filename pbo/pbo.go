package pbo

import (
	"bytes"
	"crypto/sha1"
	"io"
	"os"
	"time"
)

type PackingMethod uint32

const (
	PackingMethodUncompressed PackingMethod = 0x00000000
	PackingMethodPacked                     = 0x43707273
	PackingMethodProductEntry               = 0x56657273
)

type FileEntry struct {
	Filename     string
	Packing      PackingMethod
	OriginalSize uint32
	Reserved     uint32
	Timestamp    time.Time
	DataSize     uint32

	parent *File
	start  int64
	buf    []byte
}

func (e *FileEntry) Data() ([]byte, error) {
	var err error
	if e.buf == nil {
		if !e.parent.CacheData {
			return e.parent.getData(e.start, e.DataSize, e.Packing == PackingMethodPacked)
		}
		e.buf, err = e.parent.getData(e.start, e.DataSize, e.Packing == PackingMethodPacked)
	}
	return e.buf, err
}

type File struct {
	Path    string
	Files   []*FileEntry
	Headers map[string]string

	CacheData bool

	file      *os.File
	dataStart int64
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

func (f *File) getData(offset int64, length uint32, packed bool) ([]byte, error) {
	if _, err := f.file.Seek(f.dataStart+offset, io.SeekStart); err != nil {
		return nil, err
	}
	buf := make([]byte, length)
	_, err := f.file.Read(buf)
	return buf, err
}

func (f *File) Close() {
	f.file.Close()
	f.file = nil
	f.dataStart = 0
}

func (f *File) Dispose() {
	if f.file != nil {
		f.Close()
	}
	f.Files = nil
	f.Headers = nil
}

func Load(path string) (*File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	if err = validateFile(file); err != nil {
		return nil, err
	}
	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	if err = validateEntry(file); err != nil {
		return nil, err
	}

	parent := &File{
		Path:      path,
		file:      file,
		CacheData: true,
	}

	parent.Headers, err = loadHeaders(file)
	if err != nil {
		return nil, err
	}
	parent.Files, err = loadEntries(file, parent)
	if err != nil {
		return nil, err
	}

	parent.dataStart, err = file.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	return parent, nil
}

func validateEntry(reader io.ReadSeeker) error {
	entry, err := loadFileEntry(reader, nil)
	if err != nil {
		return err
	}
	if entry.Packing != PackingMethodProductEntry {
		return ErrInvalidProductEntry
	}
	return nil
}

func validateFile(reader io.ReadSeeker) error {
	var currentPos, lastPos int64
	var err error
	if lastPos, err = reader.Seek(-20, io.SeekEnd); err != nil {
		return err
	}
	lastPos--

	currentHash := make([]byte, 20)
	if _, err := reader.Read(currentHash); err != nil {
		return err
	}

	if _, err = reader.Seek(0, io.SeekStart); err != nil {
		return err
	}

	buf := make([]byte, 8192)
	hasher := sha1.New()
	for currentPos < lastPos {
		if (currentPos + 8192) > lastPos {
			buf = make([]byte, lastPos-currentPos)
		}
		if _, err = reader.Read(buf); err != nil {
			return err
		}
		hasher.Write(buf)
		currentPos += 8192
	}

	calculatedHash := hasher.Sum(nil)
	if bytes.Compare(currentHash, calculatedHash) != 0 {
		return ErrFileCorrupted
	}
	return nil
}
