package data

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var (
	configFilePattern = regexp.MustCompile(`^diplomat\.ya?ml`)
	translationFilePattern = regexp.MustCompile(`.*ya?ml`)
)

type SourceSet interface {
	GetConfigurationFile() (string,error)
	GetTranslationFiles() ([]string, error)
}

type fileSystemSourceSet struct {
	rootPath string
}

func (f fileSystemSourceSet) ensureRootIsDirectory() error {
	state, err := os.Stat(f.rootPath)
	if err != nil {
		return err
	}
	if !state.IsDir() {
		return fmt.Errorf("%s is not a directory", f.rootPath)
	}
	return nil
}

func (f fileSystemSourceSet) GetConfigurationFile() (string, error) {
	if err := f.ensureRootIsDirectory(); err != nil {
		return "", err
	}
	rootDir, err := os.Open(f.rootPath)
	if err != nil {
		return "", err
	}
	files, err := rootDir.Readdir(-1)
	if err != nil {
		return "", err
	}
	for _, c := range files {
		if configFilePattern.MatchString(filepath.Base(c.Name())) {
			return f.prefixRoot(c.Name())[0], nil
		}
	}
	return "", fmt.Errorf("cannot find configuration file in %s", f.rootPath)
}

func (f fileSystemSourceSet) prefixRoot(relativePath ...string) []string {
	paths := make([]string, len(relativePath))
	for i, p := range relativePath {
		paths[i] = filepath.Join(f.rootPath, p)
	}
	return paths
}

func (f fileSystemSourceSet) GetTranslationFiles() ([]string, error) {
	files := make([]string, 0)
	if err := filepath.Walk(f.rootPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		basename := filepath.Base(path)
		if !configFilePattern.MatchString(basename) && translationFilePattern.MatchString(basename) {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return files,nil
}

func NewFileSystemSourceSet(rootPath string) SourceSet {
	return &fileSystemSourceSet{rootPath}
}
