export const DEMO_CREDENTIALS = {
  email: 'admin@example.com',
  password: 'Admin123!',
} as const

export type DemoCredentialKey = keyof typeof DEMO_CREDENTIALS

export function demoCredentialLabel(key: DemoCredentialKey) {
  return key === 'email' ? 'auth.demoEmail' : 'auth.demoPassword'
}
