export default {
  app: {
    title: 'Go Vue Admin', platform: '平台', admin: '管理后台', login: '登录管理后台',
    logout: '退出', dashboard: '仪表盘', settings: '设置', media: '媒体库', audit: '审计日志',
    apiDocs: 'API 文档', users: '用户', roles: '角色', permissions: '权限', modules: '业务模块',
    plugins: '平台插件', about: '关于项目', language: '语言', theme: '主题', save: '保存',
    cancel: '取消', light: '浅色', dark: '深色', system: '跟随系统', shadcn: 'shadcn 默认',
    semi: 'Semi Design', wechat: '微信风格',
  },
  about: {
    title: '关于 Go Vue Admin', description: '面向模块、资源和插件开发的通用后台管理面板。',
    foundation: 'Foundation 基础能力', foundationDescription: '认证、RBAC、GORM 资源引擎、OpenAPI 契约和运行时扩展边界。',
    docs: '开发者文档', docsDescription: '查看架构、模块、插件、资源生成、API 契约和生产检查清单。',
    repository: '代码仓库', version: '版本与状态', status: 'Foundation 14 / 14',
  },
  dashboard: {
    title: '仪表盘', description: '平台基础能力已准备好，业务模块可按契约接入。',
    components: '组件基础', backend: '后端基础', data: '数据基础', progress: '开发进度',
    progressDescription: 'Foundation 计划与数据库兼容验收已完成；业务资源可按配置切换 MySQL / PostgreSQL。',
    done: '已完成', next: '下一步', language: '开发语言', repository: '代码仓库',
  },
  login: {
    description: '使用平台账号进入 Go Vue Admin。', credentials: '开发环境登录凭据',
    credentialsDescription: '凭据从后端实时读取，修改 AUTH_BOOTSTRAP_* 后刷新页面即可同步。',
    account: '账号', password: '密码', fill: '一键填入登录表单', email: '邮箱',
    passwordPlaceholder: '请输入密码', submit: '登录', submitting: '登录中…', retry: '重试',
    copyAccount: '复制账号', copyPassword: '复制密码', accountCopied: '账号已复制', passwordCopied: '密码已复制', showPassword: '显示密码', hidePassword: '隐藏密码',
  },
} as const
