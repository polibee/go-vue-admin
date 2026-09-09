import { describe, expect, it } from 'vitest'

import { resolveBackendUrl } from '../vite-backend'

describe('resolveBackendUrl', () => {
  it('selects the first reachable backend candidate', async () => {
    const attempts: string[] = []
    const url = await resolveBackendUrl({
      VITE_BACKEND_PORT: '3000',
      VITE_BACKEND_URL: '',
    }, async (candidate) => {
      attempts.push(candidate)
      return candidate.endsWith(':3003')
    })

    expect(url).toBe('http://127.0.0.1:3003')
    expect(attempts).toEqual([
      'http://127.0.0.1:3000',
      'http://127.0.0.1:3003',
    ])
  })

  it('honors an explicitly configured backend URL without probing', async () => {
    const probe = async () => {
      throw new Error('should not probe an explicit URL')
    }

    await expect(resolveBackendUrl({
      VITE_BACKEND_PORT: '',
      VITE_BACKEND_URL: 'http://127.0.0.1:3999/',
    }, probe)).resolves.toBe('http://127.0.0.1:3999')
  })

  it('falls back to the conventional backend port when no candidate is ready', async () => {
    await expect(resolveBackendUrl({
      VITE_BACKEND_PORT: '',
      VITE_BACKEND_URL: '',
    }, async () => false)).resolves.toBe('http://127.0.0.1:3003')
  })
})
