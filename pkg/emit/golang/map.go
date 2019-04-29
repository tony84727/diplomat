package golang

import (
	"github.com/tony84727/diplomat/pkg/data"
	"github.com/tony84727/diplomat/pkg/emit"
	"github.com/tony84727/diplomat/pkg/emit/types"
)

const (
	mapTemplate = `// {{ .DoNotEditWarning }}
package {{ .Options.PackageName }}
var {{ .Options.VariableName }} = map[string]string{
	{{ range .Pairs -}}
		"{{ JoinKeys .Key "." }}": "{{ .Text }}",
	{{ end }}
}
`
)
type MapEmitterOption struct {
	data.TemplateOption
}

func (o MapEmitterOption) VariableName() string {
	name, exist := o.TemplateOption.GetMapElement()["variableName"]
	if exist {
		if str, ok := name.(string); ok {
			return str
		}
	}
	return "Translations"
}

func (o MapEmitterOption) PackageName() string {
	name, exist := o.TemplateOption.GetMapElement()["packageName"]
	if exist {
		if str, ok := name.(string); ok {
			return str
		}
	}
	return "i18n"
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
	emit.GlobalRegistry.Registry("go-map",types.FactoryWrapper{Constructor:NewMapEmitter})
}