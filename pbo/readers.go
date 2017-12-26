package pbo

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"
)

func readString(reader io.Reader) (string, error) {
	buf := bytes.Buffer{}
	b := make([]byte, 1)
	for {
		if _, err := reader.Read(b); err != nil {
			return buf.String(), err
		}
		if b[0] == 0x00 {
			break
		}
		buf.Write(b)
	}
	return buf.String(), nil
}

func readUInt32(reader io.Reader) (uint32, error) {
	var v uint32
	if err := binary.Read(reader, binary.LittleEndian, &v); err != nil {
		return 0, err
	}
	return v, nil
}

func readPackingMethod(reader io.Reader) (PackingMethod, error) {
	v, err := readUInt32(reader)
	if err != nil {
		return PackingMethodUncompressed, err
	}
	switch PackingMethod(v) {
	case PackingMethodUncompressed, PackingMethodPacked, PackingMethodProductEntry:
		return PackingMethod(v), nil
	default:
		return PackingMethodUncompressed, &InvalidPackingMethod{Packing: PackingMethod(v)}
	}
}

func readTimestamp(reader io.Reader) (time.Time, error) {
	v, err := readUInt32(reader)
	if err != nil {
		return time.Unix(0, 0), err
	}
	return time.Unix(int64(v), 0), nil
}

func loadFileEntry(reader io.Reader, parent *File) (entry FileEntry, err error) {
	entry.parent = parent
	if entry.Filename, err = readString(reader); err != nil {
		return
	}
	if entry.Packing, err = readPackingMethod(reader); err != nil {
		return
	}
	if entry.OriginalSize, err = readUInt32(reader); err != nil {
		return
	}
	if entry.Reserved, err = readUInt32(reader); err != nil {
		return
	}
	if entry.Timestamp, err = readTimestamp(reader); err != nil {
		return
	}
	if entry.DataSize, err = readUInt32(reader); err != nil {
		return
	}
	return
}

func loadFileHeader(reader io.Reader) (key string, value string, err error) {
	if key, err = readString(reader); err != nil {
		return
	}
	if value, err = readString(reader); err != nil {
		return
	}
	return
}

func loadHeaders(reader io.ReadSeeker) (map[string]string, error) {
	values := map[string]string{}
	b := make([]byte, 1)
	for {
		if _, err := reader.Read(b); err != nil {
			return values, err
		}
		if b[0] == 0x00 {
			break
		}
		reader.Seek(-1, io.SeekCurrent)
		key, value, err := loadFileHeader(reader)
		if err != nil {
			return values, err
		}
		values[key] = value
	}
	return values, nil
}

func loadEntries(reader io.Reader, parent *File) ([]*FileEntry, error) {
	entries := make([]*FileEntry, 0, 1)
	var start int64
	for {
		entry, err := loadFileEntry(reader, parent)
		if err != nil {
			return entries, err
		}

		if entry.Filename == "" && entry.Packing == PackingMethodUncompressed &&
			entry.OriginalSize == 0 && entry.Reserved == 0 &&
			entry.Timestamp.Unix() == 0 && entry.DataSize == 0 {
			break
		}

		entry.start = start
		start += int64(entry.DataSize)

		entries = append(entries, &entry)
	}
	return entries, nil
}
