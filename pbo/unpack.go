package pbo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type UnpackHook func(entry *FileEntry) error

func createDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0777); err != nil {
			return err
		}
	}
	return nil
}

func (f *File) Unpack(destination string, hook UnpackHook) error {
	if err := createDir(destination); err != nil {
		return err
	}

	cacheEnabled := f.CacheEnabled()
	if cacheEnabled {
		f.ToggleCache(false)
	}

	for _, entry := range f.Files {
		if hook != nil {
			if err := hook(entry); err != nil {
				return err
			}
		}

		var path string
		if runtime.GOOS == "linux" {
			path = filepath.Join(destination, strings.Replace(entry.Filename, "\\", "/", -1))
		} else {
			path = filepath.Join(destination, entry.Filename)
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

func Unpack(file, destination string, hook UnpackHook) error {
	p, err := Load(file)
	if err != nil {
		return err
	}
	defer p.Close()
	return p.Unpack(destination, hook)
}

func UnpackVerbose(file, destination string) error {
	return Unpack(file, destination, VerboseUnpack)
}

func VerboseUnpack(entry *FileEntry) error {
	fmt.Println("Unpacking", entry.Filename)
	return nil
}
