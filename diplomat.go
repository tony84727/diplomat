package diplomat

import "fmt"

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
		m := make(map[string][][]string)
		for _, k := range keys {
			last := k[len(k)-1]
			_, exist := m[last]
			if !exist {
				m[last] = make([][]string, 0, 1)
			}
			m[last] = append(m[last], k)
		}
		for language, keys := range m {
			fmt.Printf("language: %s\n%v\n", language, keys)
		}
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
	go func() {
		for e := range r.GetErrorOut() {
			fmt.Println("[error] ", e)
		}
	}()
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

func NewDiplomatWatchDirectory(dir string) *Diplomat {
	r := NewReader(dir)
	go func() {
		for e := range r.GetErrorOut() {
			fmt.Println("[error] ", e)
		}
	}()
	r.Read()
	d := NewDiplomatAsync(r.GetOutlineSource(), r.GetPartialTranslationSource())
	return d
}
