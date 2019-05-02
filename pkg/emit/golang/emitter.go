package golang

import (
	"github.com/tony84727/diplomat/pkg/data"
	"github.com/tony84727/diplomat/pkg/emit"
	"github.com/tony84727/diplomat/pkg/emit/types"
)

const (
	mapTemplate = `// {{ .DoNotEditWarning }}
package {{ .Options.PackageName }}
{{ $dataName := .Options.VariableName }}
{{ with .Options.Interface }}
	import "strings"
{{ end }}
var {{ .Options.VariableName }} = map[string]string{
	{{ range .Pairs -}}
		"{{ JoinKeys .Key "." }}": "{{ .Text }}",
	{{ end }}
}
{{ with .Options.Interface }}
type {{ .Name }} interface {
	Translate(keys ...string) string
}
type impl struct {}

func (i impl) Translate(keys ...string) string {
	return {{ $dataName }}[strings.Join(keys, ".")]
}

func New{{ $dataName }}() {{ .Name }} {
	return &impl{}
}
{{ end }}
`
)
type MapEmitterOption struct {
	data.TemplateOption
}

func (o MapEmitterOption) VariableName() string {
	return o.getString("variableName", "Translations")
}

func (o MapEmitterOption) PackageName() string {
	return o.getString("packageName", "i18n")
}

type interfaceOption struct {
	Name string
}

func (o MapEmitterOption) Interface() *interfaceOption {
	i, exist := o.GetMapElement()["interface"]
	if exist {
		if m, ok := i.(map[interface{}]interface{}); ok {
			if name, exist := m["name"]; exist {
				if s, ok := name.(string); ok {
					return &interfaceOption{Name: s}
				}
			}
			return nil
		}
	}
	return nil
}

func (o MapEmitterOption) getString(key string, defaultValue string) string {
	val, exist := o.TemplateOption.GetMapElement()[key]
	if exist {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

type MapEmitter struct {
	*emit.TemplateEmitter
}

func NewMapEmitter() types.Emitter {
	templateEmitter, err := emit.NewTemplateEmitter(mapTemplate)
	if err != nil {
		panic(err)
	}
	return &MapEmitter{templateEmitter}
}

func (m MapEmitter) Emit(translations data.Translation, options data.TemplateOption) ([]byte, error) {
	return m.TemplateEmitter.Emit(translations, &MapEmitterOption{options})
}

func init() {
	emit.GlobalRegistry.Registry("go",types.FactoryWrapper{Constructor:NewMapEmitter})
}