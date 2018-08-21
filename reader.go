package diplomat

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"

	yaml "gopkg.in/yaml.v2"
)

type PreprocessorConfig struct {
	Type    string
	Options YAMLOption
}

type OutputConfig struct {
	Selectors []string
	Template  MessengerConfig
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
		dir:                    dir,
		outlineChan:            make(chan *Outline, 1),
		partialTranslationChan: make(chan *PartialTranslation, 1),
		errChan:                make(chan error, 10),
	}
}

type Reader struct {
	dir                    string
	outlineChan            chan *Outline
	partialTranslationChan chan *PartialTranslation
	errChan                chan error
}

func (r Reader) GetOutlineSource() <-chan *Outline {
	return r.outlineChan
}

func (r Reader) GetPartialTranslationSource() <-chan *PartialTranslation {
	return r.partialTranslationChan
}

func (r Reader) GetErrorOut() <-chan error {
	return r.errChan
}

func (r Reader) pushError(e error) {
	go func() {
		select {
		case r.errChan <- e:
			return
		default:
			log.Println("an error drop by reader", e)
		}
	}()
}

func (r Reader) Read() {
	r.doRead(true)
}

func (r Reader) doRead(closeChannel bool) {
	var mainWg sync.WaitGroup
	mainWg.Add(1)
	go func() {
		o, err := parseOutline(filepath.Join(r.dir, "diplomat.yaml"))
		if err != nil {
			r.pushError(err)
			mainWg.Done()
			return
		}
		r.outlineChan <- o
		if closeChannel {
			close(r.outlineChan)
		}
		mainWg.Done()
	}()
	mainWg.Add(1)
	go func() {
		var wg sync.WaitGroup
		paths, err := filepath.Glob(filepath.Join(r.dir, "**", "*.yaml"))
		if err != nil {
			r.pushError(err)
			mainWg.Done()
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
					r.pushError(err)
					return
				}
				r.partialTranslationChan <- t
				wg.Done()
			}(p)
		}
		wg.Wait()
		if closeChannel {
			close(r.partialTranslationChan)
		}
		mainWg.Done()
	}()

	mainWg.Wait()
}

func (r Reader) Watch() {
	r.doRead(false)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		r.pushError(err)
	}
	watcher.Add(r.dir)
	for e := range nameBaseThrottler(watcher.Events) {
		if isOutlineFile(e.Name) {
			go func(path string) {
				o, err := parseOutline(path)
				if err != nil {
					r.pushError(err)
					return
				}
				r.outlineChan <- o
			}(e.Name)
		} else {
			go func(path string) {
				t, err := parsePartialTranslation(path)
				if err != nil {
					r.pushError(err)
					return
				}
				r.partialTranslationChan <- t
			}(e.Name)
		}
	}
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
	var t *PartialTranslation
	err = yaml.Unmarshal(data, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

type watchThrottler struct {
	source     <-chan fsnotify.Event
	out        chan<- fsnotify.Event
	throttlers map[string]chan<- fsnotify.Event
}

func (wt watchThrottler) loop() {
	for e := range wt.source {
		c, exist := wt.throttlers[e.Name]
		if !exist {
			nc := make(chan fsnotify.Event, 1)
			go func() {
				for e := range throttle(time.Second, nc) {
					wt.out <- e
				}
			}()
			wt.throttlers[e.Name] = nc
			c = nc
		}
		c <- e
	}
}

func (wt watchThrottler) close() {
	for _, c := range wt.throttlers {
		close(c)
	}
	close(wt.out)
}

func nameBaseThrottler(source <-chan fsnotify.Event) <-chan fsnotify.Event {
	c := make(chan fsnotify.Event, 1)
	w := watchThrottler{
		source,
		c,
		make(map[string]chan<- fsnotify.Event),
	}
	go w.loop()
	return c
}
