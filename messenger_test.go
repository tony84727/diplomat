package diplomat

import (
	"bytes"
	"fmt"
)

func ExampleJsModuleMessengerSend() {
	config := BasicMessengerConfig{
		messengerType: "js",
		name:          "{{.FragmentName}}.{{.Locale}}.js",
		pairs: []TranslationPair{
			TranslationPair{
				Key:        "admin_user",
				Translated: "管理員",
			},
		},
		path:         "",
		locale:       "zh-TW",
		fragmentName: "admin",
	}
	messenger := NewJsModuleMessenger(config)
	var buffer bytes.Buffer
	messenger.Send(&buffer)
	fmt.Println(string(buffer.Bytes()))
	// Output:
	// export default {
	//     admin_user: "管理員",
	// }
}
