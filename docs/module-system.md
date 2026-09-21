# Module System

模块是功能边界和源码边界，不是进程插件。第一阶段只支持编译期模块注册，不实现第三方动态 `.so`/`.dll` 热插拔。

## 标准资源模块

```text
backend/app/modules/<name>/
├── resource/
├── permissions/
└── menu/
```

`admin:make-resource` 是唯一生成入口，从一份 ResourceSpec 生成后端资源边界、权限、菜单、前端页面、路由描述、迁移、测试和 README。生成器不会修改人工维护的 Provider、Routes、权限或菜单文件；它会更新带生成标记的 discovery 文件，使资源自动进入 Registry、路由和 Admin Shell。数据库迁移仍需人工审阅后执行。

模块可以提供自己的中英文语言包：前端由 vue-i18n 加载，后端由 Goravel Localization 加载。语言文件统一归档到前端 `locales/<locale>/<namespace>/` 和后端 `lang/<locale>/<namespace>/`。模块不得覆盖 `core` 或其他模块的命名空间。

## 依赖规则

- 模块可以依赖 Core；
- 模块之间不得直接修改彼此的模型；
- 跨模块调用使用公开 Service Contract 或 Goravel Event；
- `shared` 不是未知代码的收容目录；
- Core 不能导入 CMS、Trading、Payment 等业务包。

第一阶段平台模块：

```text
users
roles
permissions
menus
settings
audit
```

插件不是模块的别名。插件是未来的外部扩展包，必须遵循 [plugin-system.md](./plugin-system.md) 的 Manifest、Capability 和生命周期设计。
