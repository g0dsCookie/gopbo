package pbo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	humanize "github.com/dustin/go-humanize"
)

func createDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0777); err != nil {
			return err
		}
	}
	return nil
}

// Unpack unpacks the PBO file.
func (f *File) Unpack(destination string, verbose bool) error {
	if err := createDir(destination); err != nil {
		return err
	}

	cacheEnabled := f.CacheEnabled()
	if cacheEnabled {
		f.ToggleCache(false)
	}

	for _, entry := range f.Files {
		path := filepath.Join(destination, entry.Filename)

		if verbose {
			fmt.Printf("Unpacking %s to %s with size %s\n", entry.Filename, path, humanize.Bytes(uint64(entry.DataSize)))
		}

		createDir(filepath.Dir(path))

		data, err := entry.Data()
		if err != nil {
			return err
		}

		file, err := os.Create(path)
		if err != nil {
			return err
		}

		if _, err := file.Write(data); err != nil {
			return err
		}
		file.Close()

		if err = os.Chtimes(path, entry.Timestamp, entry.Timestamp); err != nil {
			return err
		}
	}

	for key, value := range f.Headers {
		if err := ioutil.WriteFile(filepath.Join(destination, "$"+strings.ToUpper(key)+"$"), []byte(value), 0666); err != nil {
			return err
		}
	}

	if cacheEnabled {
		f.ToggleCache(true)
	}

	return nil
}

// Unpack loads the PBO and unpacks it.
func Unpack(file, destination string, verbose bool) error {
	p, err := Load(file)
	if err != nil {
		return err
	}
	defer p.Close()
	return p.Unpack(destination, verbose)
}
