package diplomat

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"

	yaml "gopkg.in/yaml.v2"
)

type PreprocessorConfig struct {
	Type    string
	Options YAMLOption
}

type OutputConfig struct {
	Selectors []string
	Templates []MessengerConfig
}

// Outline is the struct of translation file.
type Outline struct {
	Version       string
	Preprocessors []PreprocessorConfig
	Output        []OutputConfig
}

type PartialTranslation struct {
	path string
	data YAMLMap
}

func NewReader(dir string) *Reader {
	return &Reader{
		dir: dir,
	}
}

type Reader struct {
	dir string
}

func (r Reader) Read() (*Outline, []*PartialTranslation, error) {
	outlineChan, translationChan, errorChan := r.doRead(true)
	var wg sync.WaitGroup
	wg.Add(2)
	var outline *Outline
	go func() {
		outline = <-outlineChan
		wg.Done()
	}()
	translations := make([]*PartialTranslation, 0, 1)
	go func() {
		for t := range translationChan {
			translations = append(translations, t)
		}
		wg.Done()
	}()
	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()
	for {
		select {
		case <-done:
			return outline, translations, nil
		case err := <-errorChan:
			return nil, nil, err
		}
	}
}

type asyncErrorSink struct {
	errorChan chan error
}

func (a asyncErrorSink) push(err error) {
	go func() {
		select {
		case a.errorChan <- err:
			return
		default:
			log.Println("[error-sink]an error dropped ", err)
		}
	}()
}

func newAsyncErrorSink() asyncErrorSink {
	return asyncErrorSink{
		errorChan: make(chan error),
	}
}

func (r Reader) doRead(closeChannels bool) (<-chan *Outline, <-chan *PartialTranslation, <-chan error) {
	outlineChan := make(chan *Outline)
	translationChan := make(chan *PartialTranslation)
	errorSink := newAsyncErrorSink()
	go func() {
		o, err := parseOutline(filepath.Join(r.dir, "diplomat.yaml"))
		if err != nil {
			errorSink.push(err)
			return
		}
		outlineChan <- o
		if closeChannels {
			close(outlineChan)
		}
	}()
	go func() {
		var wg sync.WaitGroup
		paths, err := filepath.Glob(filepath.Join(r.dir, "*.yaml"))
		if err != nil {
			errorSink.push(err)
			return
		}
		for _, p := range paths {
			if isOutlineFile(p) {
				continue
			}
			wg.Add(1)
			go func(path string) {
				t, err := parsePartialTranslation(path)
				if err != nil {
					errorSink.push(err)
					return
				}
				translationChan <- t
				wg.Done()
			}(p)
		}
		wg.Wait()
		if closeChannels {
			close(translationChan)
		}
	}()
	return outlineChan, translationChan, errorSink.errorChan
}

func doWatch(events <-chan fsnotify.Event) (<-chan *Outline, <-chan *PartialTranslation, <-chan error) {
	outlineChan := make(chan *Outline)
	partialTranslationChan := make(chan *PartialTranslation)
	errorSink := newAsyncErrorSink()
	for e := range events {
		if isOutlineFile(e.Name) {
			go func(path string) {
				o, err := parseOutline(path)
				if err != nil {
					errorSink.push(err)
					return
				}
				outlineChan <- o
			}(e.Name)
		} else {
			go func(path string) {
				t, err := parsePartialTranslation(path)
				if err != nil {
					errorSink.push(err)
					return
				}
				partialTranslationChan <- t
			}(e.Name)
		}
	}
	return outlineChan, partialTranslationChan, errorSink.errorChan
}

func (r Reader) Watch() (<-chan *Outline, <-chan *PartialTranslation, <-chan error) {
	errorSink := newAsyncErrorSink()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		errorSink.push(err)
		return nil, nil, errorSink.errorChan
	}
	watcher.Add(r.dir)
	return doWatch(watcher.Events)
}

func isOutlineFile(name string) bool {
	return strings.TrimRight(filepath.Base(name), filepath.Ext(name)) == "diplomat"
}

func parseOutline(name string) (*Outline, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	var outline Outline
	err = yaml.Unmarshal(data, &outline)
	if err != nil {
		return nil, err
	}
	return &outline, nil
}

func parsePartialTranslation(path string) (*PartialTranslation, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var t YAMLMap = make(YAMLMap)
	err = yaml.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}
	return &PartialTranslation{
		path: path,
		data: t,
	}, nil
}
