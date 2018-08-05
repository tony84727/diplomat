diplomat
---------
[![Build Status](https://travis-ci.org/MinecraftXwinP/diplomat.svg?branch=master)](https://travis-ci.org/MinecraftXwinP/diplomat)
## Goal

1. Generate translation module
2. Auto chinese convertion (simplied <=> tranditional)

## Translation file format (see testdata/outline.yaml)
```yaml
version: '1'
settings:
  chinese:
    convert:
      mode: t2s
      from: zh-TW
      to: zh-CN
  copy:
  - from: en
    to: fr
fragments:
  admin:
    description: translations for admin page
    translations:
      admin:
        zh-TW: 管理員
        en: Admin
output:
  fragments:
  - type: js
    name: "{{.Locale}}.{{.FragmentName}}.js"
```

Above configuration will generate three files:
```js
// en.admin.js
export default {
    admin: "Admin",
}
// zh-CN.admin.js
export default {
    admin: "管理员",
}
// zh-TW.admin.js
export default {
    admin: "管理員",
}
```