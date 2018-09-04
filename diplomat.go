package diplomat

import (
	"fmt"
)

type Diplomat struct {
	outline      *Outline
	translations map[string]*PartialTranslation
}

func (d *Diplomat) SetOutline(o *Outline) {
	d.outline = o
}

func (d *Diplomat) SetTranslations(translations map[string]*PartialTranslation) {
	d.translations = translations
}

func (d *Diplomat) SetTranslation(translation *PartialTranslation) {
	if d.translations == nil {
		d.translations = make(map[string]*PartialTranslation)
	}
	d.translations[translation.path] = translation
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

func NewDiplomat(outline Outline, translations map[string]*PartialTranslation) *Diplomat {
	d := &Diplomat{}
	d.SetOutline(&outline)
	d.SetTranslations(translations)
	return d
}

func NewDiplomatAsync(outlineSource <-chan *Outline, translationSource <-chan *PartialTranslation) *Diplomat {
	d := &Diplomat{}
	d.SetOutline(<-outlineSource)
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
	d := &Diplomat{}
	d.SetOutline(outline)
	translationMap := make(map[string]*PartialTranslation, len(translations))
	for _, t := range translations {
		translationMap[t.path] = t
	}
	d.SetTranslations(translationMap)
	return d, nil
}

func NewDiplomatWatchDirectory(dir string) (d *Diplomat, errorChan <-chan error, changeListener <-chan bool) {
	r := NewReader(dir)
	outlineChan, translationChan, ec := r.Watch()
	proxiedOutlineChan := make(chan *Outline)
	proxiedTranslationChan := make(chan *PartialTranslation)
	c := make(chan bool)
	go func() {
		for o := range outlineChan {
			select {
			case c <- true:
			default:
			}
			proxiedOutlineChan <- o
		}
	}()

	go func() {
		for t := range translationChan {
			select {
			case c <- true:
			default:
			}
			proxiedTranslationChan <- t
		}
	}()

	d = NewDiplomatAsync(proxiedOutlineChan, proxiedTranslationChan)
	changeListener = c
	errorChan = ec
	return
}
