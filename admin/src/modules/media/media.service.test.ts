import { describe, expect, it, vi } from 'vitest'

import { MediaService } from './media.service'

describe('MediaService', () => {
  it('lists media and forwards upload/delete operations to the generated client', async () => {
    const client = {
      listMedia: vi.fn().mockResolvedValue({ data: [{ id: 'media-1', original_name: 'avatar.png' }] }),
      uploadMedia: vi.fn().mockResolvedValue({ data: { id: 'media-2' } }),
      deleteMedia: vi.fn().mockResolvedValue({ data: { deleted: true } }),
    }
    const service = new MediaService(client as never)
    const file = new File(['image'], 'avatar.png', { type: 'image/png' })

    await expect(service.list()).resolves.toEqual([{ id: 'media-1', original_name: 'avatar.png' }])
    await expect(service.upload(file)).resolves.toEqual({ id: 'media-2' })
    await service.remove('media-2')

    expect(client.uploadMedia).toHaveBeenCalledWith(file)
    expect(client.deleteMedia).toHaveBeenCalledWith('media-2')
  })
})
