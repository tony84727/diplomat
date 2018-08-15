package diplomat

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"

	yaml "gopkg.in/yaml.v2"
)

type YAMLOption struct {
	data interface{}
}

func (y YAMLOption) Get(paths ...interface{}) (interface{}, error) {
	current := y.data
	for i, p := range paths {
		switch v := p.(type) {
		case string:
			cv, ok := current.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("%v is not a map", paths[:i])
			}
			current = cv[v]
			break
		case int:
			cv, ok := current.([]interface{})
			if !ok {
				return nil, fmt.Errorf("%v is not a slice", paths[:i])
			}
			current = cv[v]
		}
	}
	return current, nil
}

func (y *YAMLOption) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var varient interface{}
	err := unmarshal(&varient)
	if err != nil {
		return err
	}
	switch v := varient.(type) {
	case map[interface{}]interface{}:
		y.data = interfaceMapToStringMap(v)
		break

	case []interface{}:
		y.data = checkInterfaceSlice(v)
		break
	default:
		return fmt.Errorf("YAMLOption should be either map or slice")
	}
	return nil
}

func interfaceMapToStringMap(in map[interface{}]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, i := range in {
		switch v := i.(type) {
		case map[interface{}]interface{}:
			out[k.(string)] = interfaceMapToStringMap(v)
			break
		case []interface{}:
			out[k.(string)] = checkInterfaceSlice(v)
			break
		default:
			out[k.(string)] = i
		}
	}
	return out
}

func checkInterfaceSlice(in []interface{}) []interface{} {
	out := make([]interface{}, len(in))
	for i, e := range in {
		switch v := e.(type) {
		case map[interface{}]interface{}:
			out[i] = interfaceMapToStringMap(v)
			break
		case []interface{}:
			out[i] = checkInterfaceSlice(v)
			break
		default:
			out[i] = e
		}
	}
	return out
}

func (y YAMLOption) MarshalYAML() (interface{}, error) {
	switch v := y.data.(type) {
	case []interface{}:
		return yaml.Marshal(v)
	case map[string]interface{}:
		return yaml.Marshal(v)
	}
	return nil, fmt.Errorf("unknown type %v", y)
	// return yaml.Marshal(y.data)
}

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

func Read(path string) (*Outline, error) {
	var content Outline
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &content)
	if err != nil {
		return nil, err
	}
	return &content, nil
}

// NestedKeyValue is a tree node to store nested translations.
type NestedKeyValue struct {
	data map[string]interface{}
}

func nkvDataFromStringMap(m map[string]interface{}) (map[string]interface{}, error) {
	for k, p := range m {
		switch v := p.(type) {
		case int:
			m[k] = strconv.Itoa(v)
			break
		case string:
			m[k] = v
			break
		case map[interface{}]interface{}:
			stringMap := make(map[string]interface{}, len(v))
			for i, j := range v {
				stringMap[i.(string)] = j
			}
			anotherNkv, err := nkvDataFromStringMap(stringMap)
			if err != nil {
				return m, err
			}
			m[k] = &NestedKeyValue{anotherNkv}
			break
		default:
			return m, fmt.Errorf("unexcepted type: %T at %s", v, k)
		}
	}
	return m, nil
}

func (nkv *NestedKeyValue) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var root map[string]interface{}
	err := unmarshal(&root)
	if err != nil {
		return err
	}
	d, err := nkvDataFromStringMap(root)
	if err != nil {
		return err
	}
	nkv.data = d
	return nil
}

func (nkv NestedKeyValue) GetKey(paths ...string) (value interface{}, exist bool) {
	if len(paths) <= 0 {
		return nkv, true
	}
	d, exist := nkv.data[paths[0]]
	if !exist {
		return nil, false
	}
	v, ok := d.(string)
	if ok {
		return v, true
	}
	return d.(*NestedKeyValue).GetKey(paths[1:]...)
}

func (nkv NestedKeyValue) GetKeys() [][]string {
	keys := make([][]string, 0, 1)
	for k, v := range nkv.data {
		switch i := v.(type) {
		case string:
			keys = append(keys, []string{k})
			break
		case *NestedKeyValue:
			nKeys := i.GetKeys()
			for _, s := range nKeys {
				keys = append(keys, append([]string{k}, s...))
			}
			break
		default:
			continue
		}
	}
	return keys
}

func (nkv NestedKeyValue) filterBySelectorOnBase(base []string, s Selector) NestedKeyValue {
	filtered := NestedKeyValue{
		data: make(map[string]interface{}),
	}
	for k, i := range nkv.data {
		key := make([]string, len(base)+1)
		for b, bk := range base {
			key[b] = bk
		}
		key[len(base)] = k
		switch v := i.(type) {
		case string:
			if s.IsValid(key) {
				filtered.data[k] = v
			}
			break
		case *NestedKeyValue:
			if s.IsValid(key) {
				filtered.data[k] = v
			} else {
				child := v.filterBySelectorOnBase(key, s)
				if len(child.data) > 0 {
					filtered.data[k] = child
				}
			}
		}
	}
	return filtered
}

func (nkv NestedKeyValue) FilterBySelector(s Selector) NestedKeyValue {
	return nkv.filterBySelectorOnBase([]string{}, s)
}

func (nkv NestedKeyValue) HasKey(keys ...string) bool {
	_, exist := nkv.GetKey(keys...)
	return exist
}

func (nkv NestedKeyValue) LanguageHasKey(language string, keys ...string) bool {
	keys = append(keys, language)
	return nkv.HasKey(keys...)
}

func (nkv *NestedKeyValue) Set(path []string, value string) error {
	if len(path) < 1 {
		return errors.New("except at least on path")
	}
	var previous *NestedKeyValue
	var current interface{} = nkv
	lastID := len(path) - 1
	for i, p := range path {
		switch v := current.(type) {
		case *NestedKeyValue:
			n, exist := v.data[p]
			if !exist {
				if i != lastID {
					n = &NestedKeyValue{
						data: make(map[string]interface{}),
					}
					v.data[p] = n
				}
			}
			previous = v
			current = n
			break
		case string:
			n := &NestedKeyValue{
				data: map[string]interface{}{p: ""},
			}
			previous.data[path[i-1]] = n
			previous = n
			current = n.data[p]
			break
		}
	}
	previous.data[path[len(path)-1]] = value
	return nil
}

type PartialTranslation struct {
	path string
	data map[string]NestedKeyValue
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
	var outline *Outline
	err = yaml.Unmarshal(data, outline)
	if err != nil {
		return nil, err
	}
	return outline, nil
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
