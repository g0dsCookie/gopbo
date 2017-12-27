package pbo

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"io"
	"os"
	"time"
)

const (
	bufferLength = 8192
)

type pboStream struct {
	dataStart int64
	cache     map[string][]byte

	*os.File
}

func newPboStream(path string) (*pboStream, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &pboStream{0, nil, file}, nil
}

func (p *pboStream) readString() (string, error) {
	buf := bytes.Buffer{}
	b := make([]byte, 1)
	for {
		if _, err := p.Read(b); err != nil {
			return buf.String(), err
		}
		if b[0] == 0x00 {
			break
		}
		buf.Write(b)
	}
	return buf.String(), nil
}

func (p *pboStream) readUInt32() (val uint32, err error) {
	err = binary.Read(p, binary.LittleEndian, &val)
	return
}

func (p *pboStream) readPackingMethod() (PackingMethod, error) {
	v, err := p.readUInt32()
	if err != nil {
		return PackingMethodUncompressed, err
	}
	pack := PackingMethod(v)
	switch pack {
	case PackingMethodUncompressed, PackingMethodPacked, PackingMethodProductEntry:
		return pack, nil
	default:
		return PackingMethodUncompressed, &InvalidPackingMethod{pack}
	}
}

func (p *pboStream) readTimestamp() (time.Time, error) {
	v, err := p.readUInt32()
	return time.Unix(int64(v), 0), err
}

func (p *pboStream) readFileEntry(parent *File) (entry *FileEntry, err error) {
	entry = &FileEntry{parent: parent}
	if entry.Filename, err = p.readString(); err != nil {
		return
	}
	if entry.Packing, err = p.readPackingMethod(); err != nil {
		return
	}
	if entry.OriginalSize, err = p.readUInt32(); err != nil {
		return
	}
	if entry.Reserved, err = p.readUInt32(); err != nil {
		return
	}
	if entry.Timestamp, err = p.readTimestamp(); err != nil {
		return
	}
	if entry.DataSize, err = p.readUInt32(); err != nil {
		return
	}
	return
}

func (p *pboStream) readFileEntries(parent *File) (entries []*FileEntry, err error) {
	var start int64
	var entry *FileEntry
	entries = make([]*FileEntry, 0, 1)
	for {
		if entry, err = p.readFileEntry(parent); err != nil {
			return
		}
		if entry.IsEmpty() {
			break
		}

		entry.start = start
		start += int64(entry.DataSize)

		entries = append(entries, entry)
	}
	return
}

func (p *pboStream) readHeader() (key, value string, err error) {
	if key, err = p.readString(); err != nil {
		return
	}
	if value, err = p.readString(); err != nil {
		return
	}
	return
}

func (p *pboStream) readHeaders() (headers map[string]string, err error) {
	var key, value string
	b := make([]byte, 1)
	headers = make(map[string]string)
	for {
		if _, err = p.Read(b); err != nil || b[0] == 0x00 {
			return
		}
		if _, err = p.Seek(-1, io.SeekCurrent); err != nil {
			return
		}
		if key, value, err = p.readHeader(); err != nil {
			return
		}
		headers[key] = value
	}
}

func (p *pboStream) readData(entry *FileEntry) ([]byte, error) {
	if p.cache != nil {
		if buf, ok := p.cache[entry.Filename]; ok {
			return buf, nil
		}
	}

	if _, err := p.Seek(p.dataStart+entry.start, io.SeekStart); err != nil {
		return nil, err
	}

	buf := make([]byte, entry.DataSize)
	if _, err := p.Read(buf); err != nil {
		return buf, err
	}

	// only cache if we fetched valid data without errors
	if p.cache != nil {
		p.cache[entry.Filename] = buf
	}
	return buf, nil
}

func (p *pboStream) validateProductEntry() error {
	if _, err := p.Seek(0, io.SeekStart); err != nil {
		return err
	}

	entry, err := p.readFileEntry(nil)
	if err != nil {
		return err
	}

	if !entry.IsProductEntry() {
		return ErrInvalidProductEntry
	}
	return nil
}

func (p *pboStream) validateFile() (err error) {
	var currentPos, length int64
	if length, err = p.Seek(-20, io.SeekEnd); err != nil {
		return
	}
	length--

	currentHash := make([]byte, 20)
	if _, err = p.Read(currentHash); err != nil {
		return
	}
	if _, err = p.Seek(0, io.SeekStart); err != nil {
		return
	}

	buf := make([]byte, bufferLength)
	hasher := sha1.New()
	for currentPos < length {
		if (currentPos + bufferLength) > length {
			buf = make([]byte, length-currentPos)
		}
		if _, err = p.Read(buf); err != nil {
			return
		}
		hasher.Write(buf)
		currentPos += bufferLength
	}

	calculatedHash := hasher.Sum(nil)
	if bytes.Compare(currentHash, calculatedHash) != 0 {
		err = ErrFileCorrupted
	}

	if err = p.validateProductEntry(); err != nil {
		return
	}
	return
}
