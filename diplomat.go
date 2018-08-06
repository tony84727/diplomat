package diplomat

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

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

func (d Diplomat) GetOutline() Outline {
	return d.outline
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
