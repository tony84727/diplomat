package diplomat

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	DefaultDirectoryPerm = 0755
	DefaultFilePerm = 644
)

type Output interface {
	WriteFile(filename string, data []byte) error
}

type OutputDirectory struct {
	directory string
}

func NewOutputDirectory(root string) *OutputDirectory {
	return &OutputDirectory{root}
}

func (o OutputDirectory) WriteFile(filename string, data []byte) error {
	actualPath := o.absPath(filename)
	if err := o.ensureDirExistsForPath(actualPath); err != nil {
		return err
	}
	return ioutil.WriteFile(actualPath, data, DefaultFilePerm)
}

func (o OutputDirectory) ensureDirExists(dirPath string) error {
	return os.MkdirAll(dirPath, DefaultDirectoryPerm)
}

func (o OutputDirectory) ensureDirExistsForPath(filePath string) error {
	return o.ensureDirExists(filepath.Dir(filePath))
}

func (o OutputDirectory) absPath(relative string) string {
	return filepath.Join(o.directory, relative)
}