# Plugin System 设计预留

## 目标

插件系统用于扩展 Admin 平台，不用于替代业务 Module。

第一阶段只设计和预留接口，不实现：

- 插件安装；
- 插件动态发现；
- 插件执行；
- 插件启停；
- 插件卸载；
- 插件市场；
- 动态加载 Go `.so`/`.dll`。

## Module 与 Plugin

| 项目 | Module | Plugin |
|---|---|---|
| 来源 | 项目源码 | 外部扩展包 |
| 生命周期 | 编译期 | 运行时管理 |
| 信任级别 | 项目内可信代码 | 必须经过能力和安全检查 |
| 前端 | 构建时模块 | 未来由插件契约提供 |
| 后端 | 直接使用 Core | 通过 Extension Contract |
| 第一阶段 | 实现 | 只预留 |

## Manifest

插件必须提供显式 Manifest：

```yaml
id: example.blog
name: Blog
version: 1.0.0
apiVersion: admin.v1
description: Blog management plugin

dependencies:
  - id: admin.core
    version: ">=1.0.0"

capabilities:
  - navigation.register
  - resource.register
  - permission.register
  - event.subscribe
  - locale.register

permissions:
  - blog.posts.view
  - blog.posts.create

navigation:
  - id: blog.posts
    path: /blog/posts
    permission: blog.posts.view

resources:
  - id: blog.posts

locales:
  - zh-CN
  - en-US

frontend:
  type: build-time-module

backend:
  type: process
```

插件只能使用 Manifest 中声明的能力，不允许通过扫描目录隐式加载文件或配置。

## 生命周期状态

```text
discovered
    ↓
available
    ↓
installing
    ↓
installed
    ↓
enabled
    ↓
disabled
    ↓
uninstalling
    ↓
uninstalled
```

异常状态：

```text
failed
dependency_blocked
quarantined
```

## Plugin Manager

未来的 Plugin Manager 负责：

- 发现插件；
- 校验 Manifest；
- 校验插件 API 版本；
- 解析依赖和冲突；
- 校验签名或校验和；
- 安装插件；
- 启用插件；
- 停用插件；
- 卸载插件；
- 记录生命周期日志；
- 管理插件语言包。

发现来源预留：

```text
Built-in Registry
Local Plugin Directory
Configured Registry
Package Archive
```

## 启用

启用插件前必须：

1. 检查 Admin API 版本；
2. 检查依赖插件；
3. 检查冲突插件；
4. 检查 Capability；
5. 注册权限；
6. 注册导航；
7. 注册 Resource；
8. 注册语言包；
9. 注册事件；
10. 启动插件运行时。

依赖必须按拓扑顺序启用。

## 停用

```text
拒绝新的插件请求
    ↓
停止定时任务和后台任务
    ↓
注销导航、Resource、权限、事件和语言包
    ↓
关闭插件运行时
    ↓
保留配置和业务数据
    ↓
记录 disabled 状态
```

停用不能删除插件数据。

## 卸载

卸载必须先停用，并检查：

- 其他插件是否依赖它；
- 是否存在未完成任务；
- 是否存在插件拥有的数据；
- 是否需要保留审计记录；
- 是否需要备份配置。

默认卸载行为：

```text
删除插件代码
删除注册信息
保留业务数据
保留审计记录
保留迁移记录
```

删除数据必须是独立的、明确授权的危险操作。

## Plugin Context

未来插件通过稳定的上下文接口使用平台能力：

```go
type PluginContext interface {
    Manifest() Manifest
    Config() Config
    RegisterPermission(Permission) error
    RegisterNavigation(Navigation) error
    RegisterResource(Resource) error
    DeclareLocale(LocaleDeclaration) error
    Subscribe(event string, handler Handler) Disposable
    Logger() Logger
}
```

每项注册都应返回可撤销句柄：

```go
type Disposable interface {
    Dispose() error
}
```

停用插件时由宿主统一撤销其注册内容。

## 运行时安全

第一阶段不在 Goravel 主进程中动态加载任意 Go 插件。未来优先采用隔离进程：

```text
Admin Host
├── Plugin Manager
├── Plugin Process A
├── Plugin Process B
└── gRPC / HTTP Extension Protocol
```

Capability 使用白名单：

```text
navigation.register
resource.register
permission.register
event.subscribe
locale.register
storage.read
storage.write
http.outbound
schedule.register
```

默认禁止：

```text
process.exec
filesystem.root
database.raw
auth.override
permission.override
```

## 管理页面预留

未来页面：

```text
/plugins
/plugins/:id
/plugins/:id/install
/plugins/:id/settings
/plugins/:id/logs
```

第一阶段只预留导航、API Contract 和权限命名，不提供真实安装和执行。

## 设计来源

借鉴 Reasonix 的显式 Manifest、声明式能力、扩展与运行时分离、dry-run 和能力边界思想；不复制其实现。[Reasonix 文档](https://reasonix.io/docs/)

借鉴 Cordis 的上下文、依赖解析、生命周期、事件订阅和可撤销注册思想；不直接依赖 Cordis。Cordis 官方仓库说明其 API 仍处于活跃开发阶段且尚未稳定。[Cordis](https://github.com/cordiverse/cordis)
