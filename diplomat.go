package diplomat

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-yaml/yaml"
	"github.com/siongui/gojianfan"
)

type Diplomat struct {
	outline           Outline
	outlinePath       string
	outputPath        string
	messengerHandlers map[string]MessengerHandler
	watch             bool
	watcherEvents     <-chan fsnotify.Event
}

func (d *Diplomat) applyChineseConvertor(mode, from, to string) error {
	convertor, err := NewChineseConvertor(mode, from, to)
	if err != nil {
		return err
	}
	for _, f := range d.outline.Fragments {
		for _, t := range f.Translations {
			convertor.Apply(t)
		}
	}
	return nil
}

func (d Diplomat) GetOutline() Outline {
	return d.outline
}

func (d Diplomat) Output() error {
	var wg sync.WaitGroup
	for fragmentName, f := range d.outline.Fragments {
		locales := f.GetLocaleMap()
		for _, locale := range locales.GetLocales() {
			for _, outConfig := range d.outline.Output.Fragments {
				messengerHandler, exist := d.hasMessenger(outConfig.Type)
				if exist {
					dir, err := d.dirForMessenger(outConfig.Type)
					if err != nil {
						return err
					}
					wg.Add(1)
					go func(fragmentName, locale, name string, content LocaleTranslations, path string) {
						messengerHandler(
							fragmentName,
							locale,
							name,
							content,
							path,
						)
						wg.Done()
					}(fragmentName, locale, outConfig.Name, *locales.data[locale], dir)
				}
			}
		}
	}
	wg.Wait()

	return nil
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

func (d Diplomat) getAllLocales() LocaleMap {
	maps := make([]LocaleMap, len(d.outline.Fragments))
	c := 0
	for _, f := range d.outline.Fragments {
		maps[c] = *f.GetLocaleMap()
		c++
	}
	return MergeLocaleMap(maps...)
}

func (d Diplomat) applyTransformers() {
	if d.outline.Settings.Chinese != nil {
		d.applyChineseConvertor(
			d.outline.Settings.Chinese.Convert.Mode,
			d.outline.Settings.Chinese.Convert.From,
			d.outline.Settings.Chinese.Convert.To,
		)
	}
	d.applyCopyConvertor()
}

func (d *Diplomat) applyCopyConvertor() {
	if len(d.outline.Settings.Copy) <= 0 {
		return
	}
	for _, f := range d.outline.Fragments {
		for _, cp := range d.outline.Settings.Copy {
			for _, translation := range f.Translations {
				translated, exist := translation.Get(cp.From)
				if exist {
					translation.Set(cp.To, translated)
				}
			}
		}
	}
}

type MessengerHandler func(fragmentName, locale, name string, content LocaleTranslations, path string)

func (d *Diplomat) RegisterMessenger(name string, messenger MessengerHandler) {
	d.messengerHandlers[name] = messenger
}

func (d *Diplomat) Watch() {
	for range throttle(20*time.Millisecond, d.watcherEvents) {
		data, err := ioutil.ReadFile(d.outlinePath)
		if err != nil {
			log.Println("outline read error", err)
			continue
		}
		var newOutline Outline
		err = yaml.Unmarshal(data, &newOutline)
		if err != nil {
			log.Println("outline parse error", err)
			continue
		}
		d.outline = newOutline
		d.applyTransformers()
		err = d.Output()
		if err != nil {
			log.Println("output error", err)
		}
	}
}

func NewDiplomatForFile(path string, outputPath string) (*Diplomat, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var outline Outline
	err = yaml.Unmarshal(data, &outline)
	if err != nil {
		return nil, err
	}
	dip := NewDiplomat(outline, outputPath)
	dip.outlinePath = path
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

func NewDiplomat(outline Outline, outputPath string) Diplomat {
	d := Diplomat{
		outline:           outline,
		outputPath:        outputPath,
		messengerHandlers: make(map[string]MessengerHandler, 1),
	}
	d.applyTransformers()
	return d
}

type ChineseConvertor struct {
	From          string
	To            string
	transformFunc func(string) string
}

func (cc ChineseConvertor) Apply(t Translation) {
	if cc.appliable(t) {
		from, _ := t.Get(cc.From)
		t.Set(cc.To, cc.transformFunc(from))
	}
}

func (cc ChineseConvertor) appliable(t Translation) bool {
	_, fromExist := t.Get(cc.From)
	_, toExist := t.Get(cc.To)
	return fromExist && !toExist
}

func NewChineseConvertor(mode, from, to string) (*ChineseConvertor, error) {
	var transformFunc func(string) string
	switch mode {
	case "s2f":
		transformFunc = gojianfan.S2T
		break
	case "t2s":
		transformFunc = gojianfan.T2S
		break
	default:
		return nil, fmt.Errorf("chinese convertor: unknown mode %s", mode)
	}
	return &ChineseConvertor{
		transformFunc: transformFunc,
		From:          from,
		To:            to,
	}, nil
}
