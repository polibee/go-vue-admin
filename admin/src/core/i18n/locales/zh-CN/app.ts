export default {
  apiErrors: { csrfUnavailable: '无法获取安全令牌，请重新登录', csrfInvalid: '安全令牌无效，请重新登录', network: '无法连接后端服务，请确认 API 端口已启动。', requestFailed: '请求失败，请稍后重试' },
  app: {
    title: 'Go Vue Admin', platform: '平台', admin: '管理后台', login: '登录管理后台',
    logout: '退出', dashboard: '仪表盘', settings: '设置', media: '媒体库', audit: '审计日志',
    apiDocs: 'API 文档', users: '用户', roles: '角色', permissions: '权限', modules: '业务模块',
    plugins: '平台插件', about: '关于项目', language: '语言', theme: '主题', light: '浅色',
    dark: '深色', system: '跟随系统', shadcn: 'shadcn 默认', semi: 'Semi Design', tdesign: 'TDesign',
    wechat: '微信风格',
  },
} as const
