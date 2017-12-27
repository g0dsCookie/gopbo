package pbo

import (
	"time"
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
}

func (e *FileEntry) IsEmpty() bool {
	return e.Filename == "" && e.Packing == PackingMethodUncompressed &&
		e.OriginalSize == 0 && e.Reserved == 0 &&
		e.Timestamp.Unix() == 0 && e.DataSize == 0
}

func (e *FileEntry) IsProductEntry() bool { return e.Packing == PackingMethodProductEntry }

func (e *FileEntry) Data() ([]byte, error) { return e.parent.file.readData(e) }
