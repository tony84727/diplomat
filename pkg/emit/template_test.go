package emit

import (
	"fmt"
	"github.com/tony84727/diplomat/pkg/parser/yaml"
)

const translationFile = `
admin:
  admin:
    zh-TW: '管理員'
    en: 'Admin'
  message:
    hello:
      zh-TW: '您好'
      en: 'Hello!'
`

func ExampleTemplateEmitter_Emit() {
	parser := yaml.NewParser([]byte(translationFile))
	translation, err := parser.GetTranslation()
	if err != nil {
		panic(err)
	}
	templateSource := `Translations:
{{-  range .Pairs }}
{{ JoinKeys .Key "." }} => {{ .Text }}
{{- end -}}`
	emitter, err := NewTemplateEmitter(templateSource)
	if err != nil {
		panic(err)
	}
	output, err := emitter.Emit(translation)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(output))
	// Output:
	//Translations:
	//admin.admin.zh-TW => 管理員
	//admin.admin.en => Admin
	//admin.message.hello.zh-TW => 您好
	//admin.message.hello.en => Hello!
}
