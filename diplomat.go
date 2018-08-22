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
	all := make(YAMLMap)
	for _, n := range d.translations {
		for _, k := range n.data.GetKeys() {
			v, _ := n.data.GetKey(k...)
			all.Set(k, v.(string))
		}
	}
	return all
}

func (d Diplomat) Output(outDir string) error {
	all := d.getMergedYamlMap()
	ps, err := d.GetPreprocessors()
	if err != nil {
		return err
	}
	p := combinePreprocessor(ps...)
	p(all)
	for _, oc := range d.outline.Output {
		prefixSelectors := make([]Selector, len(oc.Selectors))
		for i, prefix := range oc.Selectors {
			prefixSelectors[i] = NewPrefixSelector(prefix)
		}
		selector := NewCombinedSelector(prefixSelectors...)
		selected := all.FilterBySelector(selector)
		fmt.Println(selected)
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

func NewDiplomatForDirectory(dir string) *Diplomat {
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
