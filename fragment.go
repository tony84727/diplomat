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

func (f Fragment) getPairsByLocale(locale string) {

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

type FragmentDifferentLocalesError struct {
	Basis string
	Key   string
	Diff  goset.Set
}

func (e FragmentDifferentLocalesError) Error() string {
	return fmt.Sprintf("%s locales is different to %s, diff %s", e.Key, e.Basis, e.Diff)
}
