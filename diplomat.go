package diplomat

import (
	"fmt"
	"log"
)

type Diplomat struct {
	outline         *Outline
	translations    map[string]*PartialTranslation
	outlineChan     chan *Outline
	translationChan chan *PartialTranslation
	changeListeners *fanoutHub
}

func (d *Diplomat) SetOutline(o *Outline) {
	d.outlineChan <- o
}

func (d *Diplomat) SetTranslation(translation *PartialTranslation) {
	d.translationChan <- translation
}

func (d Diplomat) GetOutputSettings() []OutputConfig {
	return d.outline.Output
}

func (d Diplomat) GetPreprocessors() ([]PreprocesserFunc, error) {
	return preprocessorManagerInstance.buildPreprocessors(d.outline.Preprocessors)
}

func (d Diplomat) getMergedYamlMap() YAMLMap {
	maps := make([]YAMLMap, len(d.translations))
	for _, p := range d.translations {
		maps = append(maps, p.data)
	}
	return MergeYAMLMaps(maps...)
}

func (d Diplomat) Watch(outDir string) error {
	l := make(chan interface{})
	d.changeListeners.addListener(l)
	log.Println("output")
	for range l {
		if d.outline != nil {
			err := d.Output(outDir)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (d Diplomat) Output(outDir string) error {
	all := d.getMergedYamlMap()
	ps, err := d.GetPreprocessors()
	if err != nil {
		return fmt.Errorf("building preprocessers: %s", err)
	}
	p := combinePreprocessor(ps...)
	err = p(all)
	if err != nil {
		return fmt.Errorf("preprocessor: %s", err)
	}
	for _, oc := range d.outline.Output {
		prefixSelectors := make([]Selector, len(oc.Selectors))
		for i, prefix := range oc.Selectors {
			prefixSelectors[i] = NewPrefixSelector(prefix)
		}
		selector := NewCombinedSelector(prefixSelectors...)
		selected := all.FilterBySelector(selector)
		keys := selected.GetKeys()
		m := make(map[string]YAMLMap)
		for _, k := range keys {
			last := k[len(k)-1]
			_, exist := m[last]
			if !exist {
				m[last] = make(YAMLMap)
			}
			v, _ := selected.GetKey(k...)
			m[last].Set(k, v.(string))
		}
		err := d.runMessengers(oc, m, outDir)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d Diplomat) runMessengers(oc OutputConfig, languages map[string]YAMLMap, outDir string) error {
	bs := newBuildSpace(outDir)
	errChan := make(chan *messengerError)
	running := 0
	for _, t := range oc.Templates {
		mf, exist := messengerRegistryInstance[t.Type]
		if exist {
			running++
			mbs := bs.ForMessenger(t.Type)
			go func(f MessengerFunc, messengerType string, option YAMLOption, b IBuildSpace) {
				err := f(languages, option, b)
				if err != nil {
					errChan <- newMessengerError(messengerType, err)
				} else {
					errChan <- nil
				}
			}(mf, t.Type, t.Options, mbs)
		}
	}
	for running > 0 {
		e := <-errChan
		if e != nil {
			return e
		}
		running--
	}
	return nil
}

func (d *Diplomat) startMaintenanceLoops() {
	go d.maintainOutline()
	go d.maintainTranslations()
}

func (d *Diplomat) maintainTranslations() {
	for t := range d.translationChan {
		d.translations[t.path] = t
		d.changeListeners.broadcast(true)
	}
}

func (d *Diplomat) maintainOutline() {
	for o := range d.outlineChan {
		d.outline = o
		d.changeListeners.broadcast(true)
	}
}

func New() *Diplomat {
	listeners := newFanoutHub()
	go listeners.run()
	// debug := make(chan interface{})
	// listeners.addListener(debug)
	// for e := range debug {
	// 	log.Println(e)
	// }
	d := &Diplomat{
		changeListeners: listeners,
		outlineChan:     make(chan *Outline, 10),
		translationChan: make(chan *PartialTranslation, 10),
		translations:    make(map[string]*PartialTranslation),
		outline:         nil,
	}
	d.startMaintenanceLoops()
	return d
}

func NewByValue(outline *Outline, translations []*PartialTranslation) *Diplomat {
	d := New()
	d.SetOutline(outline)
	for _, t := range translations {
		d.SetTranslation(t)
	}
	return d
}

func NewDiplomatAsync(outlineSource <-chan *Outline, translationSource <-chan *PartialTranslation) *Diplomat {
	d := New()
	go func() {
		for o := range outlineSource {
			d.SetOutline(o)
		}
	}()

	go func() {
		for t := range translationSource {
			d.SetTranslation(t)
		}
	}()
	return d
}

func NewDiplomatForDirectory(dir string) (*Diplomat, error) {
	r := NewReader(dir)
	outline, translations, err := r.Read()
	if err != nil {
		return nil, err
	}
	d := NewByValue(outline, translations)
	return d, nil
}

func NewDiplomatWatchDirectory(dir string) (d *Diplomat, errorChan <-chan error) {
	r := NewReader(dir)
	outlineChan, translationChan, ec := r.Watch()
	d = NewDiplomatAsync(outlineChan, translationChan)
	errorChan = ec
	return
}
