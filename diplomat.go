package diplomat

import "log"

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

func NewDiplomat(outline Outline, translations map[string]*PartialTranslation) *Diplomat {
	d := &Diplomat{}
	d.SetOutline(&outline)
	d.SetTranslations(translations)
	return d
}

func NewDiplomatAsync(outlineSource <-chan *Outline, transitionSource <-chan *PartialTranslation) *Diplomat {
	d := &Diplomat{}
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

func newReaderAndDiplomat(dir string) (*Reader, *Diplomat) {
	r := NewReader(dir)
	d := NewDiplomatAsync(r.GetOutlineSource(), r.GetPartialTranslationSource())
	go func() {
		for e := range r.GetErrorOut() {
			log.Println("reader:", e)
		}
	}()
	return r, d
}

func NewDiplomatForDirectory(dir string) *Diplomat {
	r, d := newReaderAndDiplomat(dir)
	r.Read()
	return d
}

func NewDiplomatWatchDirectory(dir string) *Diplomat {
	r, d := newReaderAndDiplomat(dir)
	r.Watch()
	return d
}
