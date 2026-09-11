export default {
  app: {
    title: 'Go Vue Admin', platform: 'Platform', admin: 'Admin Console', login: 'Admin Console Login',
    logout: 'Log out', dashboard: 'Dashboard', settings: 'Settings', media: 'Media', audit: 'Audit log',
    apiDocs: 'API docs', users: 'Users', roles: 'Roles', permissions: 'Permissions', modules: 'Business modules',
    plugins: 'Platform plugins', about: 'About', language: 'Language', theme: 'Theme', save: 'Save',
    cancel: 'Cancel', light: 'Light', dark: 'Dark', system: 'System', shadcn: 'shadcn default',
    semi: 'Semi Design', wechat: 'WeChat',
  },
  about: {
    title: 'About Go Vue Admin', description: 'A general-purpose admin console for modular resources and plugins.',
    foundation: 'Foundation capabilities', foundationDescription: 'Authentication, RBAC, GORM resources, OpenAPI contracts and runtime extension boundaries.',
    docs: 'Developer documentation', docsDescription: 'Architecture, modules, plugins, resource generation, API contracts and production checks.',
    repository: 'Repository', version: 'Version and status', status: 'Foundation 14 / 14',
  },
  dashboard: {
    title: 'Dashboard', description: 'The platform foundation is ready for contract-driven business modules.',
    components: 'UI foundation', backend: 'Backend foundation', data: 'Data foundation', progress: 'Development progress',
    progressDescription: 'Foundation and database compatibility acceptance are complete; resources can select MySQL or PostgreSQL.',
    done: 'Complete', next: 'Next', language: 'Primary language', repository: 'Repository',
  },
  login: {
    description: 'Sign in to Go Vue Admin.', credentials: 'Development credentials',
    credentialsDescription: 'Credentials are loaded from the backend and reflect AUTH_BOOTSTRAP_* after refresh.',
    account: 'Account', password: 'Password', fill: 'Fill login form', email: 'Email',
    passwordPlaceholder: 'Enter your password', submit: 'Sign in', submitting: 'Signing in…', retry: 'Retry',
    copyAccount: 'Copy account', copyPassword: 'Copy password', accountCopied: 'Account copied', passwordCopied: 'Password copied', showPassword: 'Show password', hidePassword: 'Hide password',
  },
} as const
