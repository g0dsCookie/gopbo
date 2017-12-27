package pbo

import (
	"time"
)

// FileEntry contains informations for a PBO file entry.
type FileEntry struct {
	Filename     string        // Filename contains the name of the file.
	Packing      PackingMethod // Packing contains the packing method.
	OriginalSize uint32        // OriginalSize contains the filesize before it has been packed.
	Reserved     uint32        // Reserved contains no useful information.
	Timestamp    time.Time     // Timestamp contains the last modification time.
	DataSize     uint32        // DataSize contains the actual size within the pbo.

	parent *File
	start  int64
}

// IsEmpty returns true if everything is set to 0/empty/default.
func (e *FileEntry) IsEmpty() bool {
	return e.Filename == "" && e.Packing == PackingMethodUncompressed &&
		e.OriginalSize == 0 && e.Reserved == 0 &&
		e.Timestamp.Unix() == 0 && e.DataSize == 0
}

// IsProductEntry returns true if Packing is equals PackingMethodProductEntry
func (e *FileEntry) IsProductEntry() bool { return e.Packing == PackingMethodProductEntry }

// Data returns the current files data.
func (e *FileEntry) Data() ([]byte, error) { return e.parent.file.readData(e) }
