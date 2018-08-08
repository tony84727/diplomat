package diplomat

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type Diplomat struct {
	outputPath        string
	messengerHandlers map[string]MessengerHandler
}

func (d Diplomat) dirForMessenger(messengerType string) (string, error) {
	path := filepath.Join(d.outputPath, messengerType)
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(path, 0755)
			if err != nil {
				return path, err
			}
			return path, nil
		}
		return path, err
	}
	if !info.IsDir() {
		return path, fmt.Errorf("output dir for [%s](%s) already exist, but is not a directory", messengerType, path)
	}
	return path, err
}

func (d Diplomat) hasMessenger(messengerType string) (MessengerHandler, bool) {
	m, exist := d.messengerHandlers[messengerType]
	return m, exist
}

type MessengerHandler func(fragmentName, locale, name string, content LocaleTranslations, path string)

func (d *Diplomat) RegisterMessenger(name string, messenger MessengerHandler) {
	d.messengerHandlers[name] = messenger
}

func NewDiplomatForDir(path string, outputPath string) (*Diplomat, error) {
	dip := NewDiplomat(outputPath)
	return &dip, nil
}

func NewDiplomatWatchFile(path string, outputPath string) (*Diplomat, error) {
	d, err := NewDiplomatForFile(path, outputPath)
	if err != nil {
		return d, err
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	watcher.Add(path)
	d.watcherEvents = watcher.Events

	return d, nil
}

func NewDiplomat(outputPath string) Diplomat {
	d := Diplomat{
		outputPath:        outputPath,
		messengerHandlers: make(map[string]MessengerHandler, 1),
	}
	return d
}
