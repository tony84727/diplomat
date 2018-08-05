package diplomat

import (
	"fmt"

	goset "github.com/deckarep/golang-set"
)

// Fragment is a group of translations with additional information.
type Fragment struct {
	Description  string
	Translations Translations
}

// Lint check problems of Fragment
func (f Fragment) Lint() []error {
	errors := make([]error, 0)
	// check all translations have the same locales
	var basisKey string
	var locales goset.Set
	for k, ts := range f.Translations {
		if locales == nil {
			basisKey = k
			locales = goset.NewSetFromSlice(stringSliceToInterfaceSlice(ts.GetLocales()))
			continue
		}
		diff := locales.Difference(goset.NewSetFromSlice(stringSliceToInterfaceSlice(ts.GetLocales())))
		if len(diff.ToSlice()) > 0 {
			errors = append(errors, &FragmentDifferentLocalesError{
				Basis: basisKey,
				Key:   k,
				Diff:  diff,
			})
		}
	}
	return errors
}

type LocaleTranslations struct {
	Locale       string
	Translations map[string]string
}

func (lt *LocaleTranslations) Add(key, translated string) {
	lt.Translations[key] = translated
}

func NewLocaleTranslations(locale string) *LocaleTranslations {
	return &LocaleTranslations{
		Locale:       locale,
		Translations: make(map[string]string, 1),
	}
}

type LocaleMap struct {
	data map[string]*LocaleTranslations
}

func (l LocaleMap) Contain(locale string) bool {
	_, exist := l.data[locale]
	return exist
}

func (l *LocaleMap) Add(locale, key, translated string) {
	if !l.Contain(locale) {
		l.data[locale] = NewLocaleTranslations(locale)
	}
	l.data[locale].Add(key, translated)
}

func (l LocaleMap) Get(locale string) (t LocaleTranslations, exist bool) {
	m, exist := l.data[locale]
	if exist {
		t = *m
	}
	return
}

func (l LocaleMap) GetLocales() []string {
	locales := make([]string, len(l.data))
	c := 0
	for k := range l.data {
		locales[c] = k
		c++
	}
	return locales
}

func (l LocaleMap) IterateLocale(locale string) <-chan TranslationPair {
	t, ok := l.data[locale]
	if !ok {
		return nil
	}
	c := make(chan TranslationPair)
	go func() {
		for key, translated := range t.Translations {
			c <- TranslationPair{
				Key:        key,
				Translated: translated,
			}
		}
		close(c)
	}()
	return c
}

func (l LocaleMap) GetLocaleTranslationPairs(locale string) []TranslationPair {
	buffer := make([]TranslationPair, 0)
	c := 0
	for p := range l.IterateLocale(locale) {
		buffer[c] = p
		c++
	}
	return buffer
}

func MergeLocaleMap(maps ...LocaleMap) LocaleMap {
	merged := LocaleMap{
		data: make(map[string]*LocaleTranslations),
	}

	return merged
}

func (f Fragment) GetLocaleMap() *LocaleMap {
	localeMap := &LocaleMap{
		data: make(map[string]*LocaleTranslations, 1),
	}
	for k, ts := range f.Translations {
		for entry := range ts.Iterate() {
			locale := entry.Locale
			translated := entry.Translated
			localeMap.Add(locale, k, translated)
		}
	}
	return localeMap
}

type FragmentDifferentLocalesError struct {
	Basis string
	Key   string
	Diff  goset.Set
}

func (e FragmentDifferentLocalesError) Error() string {
	return fmt.Sprintf("%s locales is different to %s, diff %s", e.Key, e.Basis, e.Diff)
}
