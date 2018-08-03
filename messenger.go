package diplomat

import (
	"bytes"
	"html/template"
	"io"
	"log"
)

type TranslationPair struct {
	Key        string
	Translated string
}

type MessengerConfig interface {
	GetType() string
	GetName() string
	GetPairs() []TranslationPair
	GetPath() string
	GetLocale() string
	GetFragmentName() string
}

type Messenger interface {
	GetFolder() string
	Send(config MessengerConfig) error
}

var JsMessengerTemplate *template.Template

func init() {
	var err error
	JsMessengerTemplate, err = template.New("jsMessengerTemplate").Parse(`
        export default {
            {{ range . }}
                {{.Key}}:"{{.Translated}}",
            {{ end }}
        }
    `)
	if err != nil {
		log.Fatal(err)
	}
}

func NewJsModuleMessenger(config MessengerConfig) *JsModuleMessenger {
	return &JsModuleMessenger{config}
}

type JsModuleMessenger struct {
	config MessengerConfig
}

func (j JsModuleMessenger) GetFolder() string {
	return "js"
}

func (j JsModuleMessenger) Send(writer io.Writer) error {
	return JsMessengerTemplate.Execute(writer, j.config.GetPairs())
}

func (j JsModuleMessenger) getFileName(fragmentName, locale string) (string, error) {
	var buffer bytes.Buffer
	t, err := template.New("js").Parse(j.config.GetName())
	if err != nil {
		return "", err
	}
	t.Execute(
		&buffer,
		struct {
			Locale       string
			FragmentName string
		}{
			locale,
			fragmentName,
		},
	)
	return buffer.String(), nil
}
