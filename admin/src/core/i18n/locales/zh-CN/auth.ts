export default {
  login: {
    description: '使用平台账号进入 Go Vue Admin。', credentials: '开发环境登录凭据',
    credentialsDescription: '凭据从后端实时读取，修改 AUTH_BOOTSTRAP_* 后刷新页面即可同步。',
    account: '账号', password: '密码', fill: '一键填入登录表单', email: '邮箱',
    passwordPlaceholder: '请输入密码', submit: '登录', submitting: '登录中…', retry: '重试',
    copyAccount: '复制账号', copyPassword: '复制密码', accountCopied: '账号已复制',
    passwordCopied: '密码已复制', showPassword: '显示密码', hidePassword: '隐藏密码',
  },
} as const
