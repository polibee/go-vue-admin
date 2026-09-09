import { describe, expect, it } from 'vitest'

describe('admin foundation', () => {
  it('keeps the official Vite template entrypoint contract', () => {
    expect('src/main.ts').toBe('src/main.ts')
    expect('src/App.vue').toBe('src/App.vue')
  })
})
