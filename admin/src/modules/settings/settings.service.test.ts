import { describe, expect, it, vi } from 'vitest'

import { SettingsService } from './settings.service'

describe('SettingsService', () => {
  it('groups settings by namespace and forwards typed mutations', async () => {
    const client = {
      listSettings: vi.fn().mockResolvedValue({
        data: [
          { namespace: 'general', key: 'site_name', value: 'Go Vue Admin', value_type: 'string' },
          { namespace: 'auth', key: 'session_ttl', value: 120, value_type: 'integer' },
        ],
      }),
      upsertSetting: vi.fn().mockResolvedValue({ data: { namespace: 'general', key: 'site_name', value: 'Admin' } }),
      deleteSetting: vi.fn().mockResolvedValue({ data: { deleted: true } }),
    }
    const service = new SettingsService(client as never)

    await expect(service.list()).resolves.toEqual({
      general: [{ namespace: 'general', key: 'site_name', value: 'Go Vue Admin', value_type: 'string' }],
      auth: [{ namespace: 'auth', key: 'session_ttl', value: 120, value_type: 'integer' }],
    })
    await service.upsert({ namespace: 'general', key: 'site_name', value: 'Admin', value_type: 'string' })
    await service.remove('general', 'site_name')

    expect(client.upsertSetting).toHaveBeenCalledWith({
      namespace: 'general', key: 'site_name', value: 'Admin', value_type: 'string',
    })
    expect(client.deleteSetting).toHaveBeenCalledWith('general', 'site_name')
  })
})
