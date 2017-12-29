package pbo

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func listFiles(dir string) ([]*FileEntry, map[string]string, error) {
	files := make([]*FileEntry, 0, 1)
	headers := make(map[string]string)
	return files, headers, filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if strings.HasPrefix(info.Name(), "$") && strings.HasSuffix(info.Name(), "$") {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			headers[strings.ToLower(info.Name()[1:len(info.Name())-1])] = string(data)
		} else {
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			files = append(files, &FileEntry{
				Filename:     relPath,
				Packing:      PackingMethodUncompressed,
				OriginalSize: uint32(info.Size()),
				Reserved:     0,
				Timestamp:    info.ModTime(),
				DataSize:     uint32(info.Size()),
			})
		}
		return nil
	})
}

// Pack packs dir into destination. The destination will be overwritten if it exists.
func Pack(dir, destination string, verbose bool) error {
	stream, err := createPBO(destination)
	if err != nil {
		return err
	}
	defer stream.Close()

	if verbose {
		fmt.Printf("Building file list for %s\n", dir)
	}
	files, headers, err := listFiles(dir)
	if err != nil {
		return err
	}

	if err = stream.writeFileEntry(NewProductEntry()); err != nil {
		return err
	}
	if err = stream.writeHeaders(headers); err != nil {
		return err
	}
	if err = stream.writeFileEntries(files); err != nil {
		return err
	}
	if err = stream.writeFileEntry(NewEmptyEntry()); err != nil {
		return err
	}

	for _, entry := range files {
		path := filepath.Join(dir, entry.Filename)
		if verbose {
			fmt.Printf("Packing file %s\n", entry.Filename)
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		if _, err = io.Copy(stream, file); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}

	return stream.writeHash()
}
