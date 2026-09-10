import { describe, expect, it, vi } from 'vitest'

import { AuditService } from './audit.service'

describe('AuditService', () => {
  it('lists entries and loads one diff record through generated client', async () => {
    const client = {
      listAudit: vi.fn().mockResolvedValue({ data: [{ id: 'audit-1', action: 'settings.update' }] }),
      getAudit: vi.fn().mockResolvedValue({ data: { id: 'audit-1', before: { value: 'old' }, after: { value: 'new' } } }),
    }
    const service = new AuditService(client as never)

    await expect(service.list()).resolves.toEqual([{ id: 'audit-1', action: 'settings.update' }])
    await expect(service.get('audit-1')).resolves.toEqual({ id: 'audit-1', before: { value: 'old' }, after: { value: 'new' } })
    expect(client.getAudit).toHaveBeenCalledWith('audit-1')
  })
})
