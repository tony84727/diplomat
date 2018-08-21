package diplomat

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

// func (d Diplomat) getMergedNKV() *NestedKeyValue {
// 	for _, n := range d.translations {
// 	}
// }

func (d Diplomat) Output() error {
	return nil
}

func NewDiplomat(outline Outline, translations map[string]*PartialTranslation) *Diplomat {
	d := &Diplomat{}
	d.SetOutline(&outline)
	d.SetTranslations(translations)
	return d
}

func NewDiplomatAsync(outlineSource <-chan *Outline, transitionSource <-chan *PartialTranslation) *Diplomat {
	d := &Diplomat{}
	d.SetOutline(<-outlineSource)
	go func() {
		for o := range outlineSource {
			d.SetOutline(o)
		}
	}()

	go func() {
		for t := range transitionSource {
			d.SetTranslation(t)
		}
	}()
	return d
}

func NewDiplomatForDirectory(dir string) *Diplomat {
	r := NewReader(dir)
	r.Read()
	d := NewDiplomatAsync(r.GetOutlineSource(), r.GetPartialTranslationSource())
	return d
}

func NewDiplomatWatchDirectory(dir string) *Diplomat {
	r := NewReader(dir)
	r.Read()
	d := NewDiplomatAsync(r.GetOutlineSource(), r.GetPartialTranslationSource())
	return d
}
