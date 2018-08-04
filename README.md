diplomat
---------
## Goal

1. Generate translation module
2. Auto chinese convertion (simplied <=> tranditional)

## Translation file format (see testdata/outline.yaml)
```
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
1. en.admin.js
2. zh-CN.admin.js
3. zh-TW.admin.js