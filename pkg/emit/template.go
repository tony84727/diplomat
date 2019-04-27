package emit

import (
	"bytes"
	"github.com/tony84727/diplomat/pkg/data"
	"strings"
	"text/template"
)

type TemplateEmitter struct {
	template *template.Template
}

type Pair struct {
	Key  []string
	Text string
}

type templateContext struct {
	data.Translation
}

func (t templateContext) Pairs() []Pair {
	walker := data.NewTranslationWalker(t)
	pairs := make([]Pair, 0)
	_ = walker.ForEachTextNode(func(paths []string, textNode data.Translation) error {
		pairs = append(pairs, Pair{paths, *textNode.GetText()})
		return nil
	})
	return pairs
}

func (t TemplateEmitter) Emit(translation data.Translation) ([]byte, error) {
	var buffer bytes.Buffer
	err := t.template.Execute(&buffer, templateContext{translation})
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (t *TemplateEmitter) SetTemplate(content string) error {
	tpl, err := template.New("main").Funcs(template.FuncMap{
		"JoinKeys": strings.Join,
	}).Parse(content)
	if err != nil {
		return err
	}
	t.template = tpl
	return nil
}

func NewTemplateEmitter(templateSource string) (*TemplateEmitter, error) {
	emitter := &TemplateEmitter{}
	if err := emitter.SetTemplate(templateSource); err != nil {
		return nil, err
	}
	return emitter, nil
}
