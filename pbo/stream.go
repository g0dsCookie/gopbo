package pbo

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"io"
	"os"
	"runtime"
	"strings"
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

func openPBO(path string) (*pboStream, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &pboStream{0, nil, file}, nil
}

func createPBO(path string) (*pboStream, error) {
	file, err := os.Create(path)
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

func (p *pboStream) writeString(s string) (err error) {
	b := append([]byte(s), 0x00)
	if _, err = p.Write(b); err != nil {
		return
	}
	return
}

func (p *pboStream) readUInt32() (val uint32, err error) {
	err = binary.Read(p, binary.LittleEndian, &val)
	return
}

func (p *pboStream) writeUInt32(val uint32) error {
	return binary.Write(p, binary.LittleEndian, val)
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

func (p *pboStream) writePackingMethod(pack PackingMethod) error {
	return p.writeUInt32(uint32(pack))
}

func (p *pboStream) readTimestamp() (time.Time, error) {
	v, err := p.readUInt32()
	return time.Unix(int64(v), 0), err
}

func (p *pboStream) writeTimestamp(t time.Time) error {
	return p.writeUInt32(uint32(t.Unix()))
}

func (p *pboStream) readFileEntry(parent *File) (entry *FileEntry, err error) {
	entry = &FileEntry{parent: parent}
	if entry.Filename, err = p.readString(); err != nil {
		return
	}
	if runtime.GOOS != "windows" {
		entry.Filename = strings.Replace(entry.Filename, "\\", "/", -1)
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
	entry.DataSize, err = p.readUInt32()
	return
}

func (p *pboStream) writeFileEntry(entry *FileEntry) (err error) {
	filename := entry.Filename
	if runtime.GOOS != "windows" {
		filename = strings.Replace(filename, "/", "\\", -1)
	}
	if err = p.writeString(filename); err != nil {
		return
	}
	if err = p.writePackingMethod(entry.Packing); err != nil {
		return
	}
	if err = p.writeUInt32(entry.OriginalSize); err != nil {
		return
	}
	if err = p.writeUInt32(entry.Reserved); err != nil {
		return
	}
	if err = p.writeTimestamp(entry.Timestamp); err != nil {
		return
	}
	err = p.writeUInt32(entry.DataSize)
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

func (p *pboStream) writeFileEntries(entries []*FileEntry) error {
	for _, entry := range entries {
		if err := p.writeFileEntry(entry); err != nil {
			return err
		}
	}
	return nil
}

func (p *pboStream) readHeader() (s string, err error) {
	s, err = p.readString()
	return
}

func (p *pboStream) writeHeader(key, value string) (err error) {
	if err = p.writeString(key); err != nil {
		return
	}
	err = p.writeString(value)
	return
}

func (p *pboStream) readHeaders() (headers map[string]string, err error) {
	var header string
	headerBuf := make([]string, 0, 1)
	for {
		if header, err = p.readHeader(); err != nil {
			return
		}
		if header == "" {
			break
		}
		headerBuf = append(headerBuf, header)
	}
	headers = make(map[string]string)
	i := 0
	for i < len(headerBuf) { // ignore last entry
		headers[headerBuf[i]] = headerBuf[i+1]
		i += 2
	}
	return
}

func (p *pboStream) writeHeaders(headers map[string]string) error {
	for key, value := range headers {
		if err := p.writeHeader(key, value); err != nil {
			return err
		}
	}
	_, err := p.Write([]byte{0x00})
	return err
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

func (p *pboStream) writeData(b []byte) (err error) {
	if _, err = p.Write(b); err != nil {
		return
	}
	return
}

func (p *pboStream) writeHash() error {
	length, err := p.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}
	calculatedHash, err := p.calculateHash(length)
	if err != nil {
		return err
	}
	b := append([]byte{0x00}, calculatedHash...)
	_, err = p.Write(b)
	return err
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
	var length int64
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

	calculatedHash, err := p.calculateHash(length)
	if err != nil {
		return err
	}
	if bytes.Compare(currentHash, calculatedHash) != 0 {
		err = ErrFileCorrupted
	}

	if err = p.validateProductEntry(); err != nil {
		return
	}
	return
}

func (p *pboStream) calculateHash(length int64) ([]byte, error) {
	var currentPos int64
	buf := make([]byte, bufferLength)
	hasher := sha1.New()
	for currentPos < length {
		if (currentPos + bufferLength) > length {
			buf = make([]byte, length-currentPos)
		}
		read, err := p.Read(buf)
		if err == io.EOF {
			hasher.Write(buf[:read])
			break
		} else if err != nil {
			return nil, err
		}
		hasher.Write(buf)
		currentPos += bufferLength
	}
	return hasher.Sum(nil), nil
}
