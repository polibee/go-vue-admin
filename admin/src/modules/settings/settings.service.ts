import type { GeneratedApiClient, SettingResource } from '@/generated/api'

export type SettingsByNamespace = Record<string, SettingResource[]>

export class SettingsService {
  constructor(private readonly client: Pick<GeneratedApiClient, 'listSettings' | 'upsertSetting' | 'deleteSetting'>) {}

  async list(namespace?: string): Promise<SettingsByNamespace> {
    const { data } = await this.client.listSettings(namespace ? { namespace } : undefined)
    return data.reduce<SettingsByNamespace>((groups, item) => {
      ;(groups[item.namespace] ??= []).push(item)
      return groups
    }, {})
  }

  async upsert(input: SettingResource): Promise<SettingResource> {
    return (await this.client.upsertSetting(input)).data
  }

  async remove(namespace: string, key: string): Promise<void> {
    await this.client.deleteSetting(namespace, key)
  }
}
