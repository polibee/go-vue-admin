export default {
  apiErrors: { csrfUnavailable: 'Unable to obtain a security token. Please sign in again.', csrfInvalid: 'The security token is invalid. Please sign in again.', network: 'Unable to connect to the backend. Confirm that the API port is running.', requestFailed: 'Request failed. Please try again later.' },
  app: {
    title: 'Go Vue Admin', platform: 'Platform', admin: 'Admin Console', login: 'Admin Console Login',
    logout: 'Log out', dashboard: 'Dashboard', settings: 'Settings', media: 'Media', audit: 'Audit log',
    apiDocs: 'API docs', users: 'Users', roles: 'Roles', permissions: 'Permissions', modules: 'Business modules',
    plugins: 'Platform plugins', about: 'About', language: 'Language', theme: 'Theme', light: 'Light',
    dark: 'Dark', system: 'System', shadcn: 'shadcn default', semi: 'Semi Design', tdesign: 'TDesign',
    wechat: 'WeChat',
  },
} as const
