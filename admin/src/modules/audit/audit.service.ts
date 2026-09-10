import type { AuditResource, GeneratedApiClient } from '@/generated/api'

export class AuditService {
  constructor(private readonly client: Pick<GeneratedApiClient, 'listAudit' | 'getAudit'>) {}

  async list(): Promise<AuditResource[]> {
    return (await this.client.listAudit()).data
  }

  async get(id: string): Promise<AuditResource> {
    return (await this.client.getAudit(id)).data
  }
}
