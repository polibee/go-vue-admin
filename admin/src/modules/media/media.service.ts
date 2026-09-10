import type { GeneratedApiClient, MediaResource } from '@/generated/api'

export class MediaService {
  constructor(private readonly client: Pick<GeneratedApiClient, 'listMedia' | 'uploadMedia' | 'deleteMedia'>) {}

  async list(): Promise<MediaResource[]> {
    return (await this.client.listMedia()).data
  }

  async upload(file: File): Promise<MediaResource> {
    return (await this.client.uploadMedia(file)).data
  }

  async remove(id: string): Promise<void> {
    await this.client.deleteMedia(id)
  }
}
