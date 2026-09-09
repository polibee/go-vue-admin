import { MemoryResourceDataProvider } from './core/ResourceDataProvider'
import { defineResource } from './core/ResourceDefinition'
import { ResourceRegistry } from './core/ResourceRegistry'

export interface DemoResourceRecord {
  id: string
  name: string
  status: string
  owner: string
}

export const demoResource = defineResource<DemoResourceRecord>({
  name: 'demo',
  label: '示例资源',
  permissions: {
    list: 'dashboard.view',
    get: 'dashboard.view',
    create: 'dashboard.view',
    update: 'dashboard.view',
    delete: 'dashboard.view',
    bulkDelete: 'dashboard.view',
  },
  columns: [
    { key: 'name', label: '名称', sortable: true },
    { key: 'status', label: '状态', sortable: true },
    { key: 'owner', label: '负责人' },
  ],
  fields: [
    { name: 'name', label: '名称', required: true, placeholder: '请输入名称' },
    { name: 'status', label: '状态', type: 'select', required: true, options: [
      { label: '草稿', value: 'draft' },
      { label: '启用', value: 'active' },
    ] },
    { name: 'owner', label: '负责人', required: true, placeholder: '请输入负责人' },
  ],
})

export const demoProvider = new MemoryResourceDataProvider<DemoResourceRecord>([
  { id: 'demo-1', name: '资源引擎示例', status: 'active', owner: 'Platform Admin' },
  { id: 'demo-2', name: '可编辑记录', status: 'draft', owner: 'Platform Admin' },
])

export const resourceRegistry = new ResourceRegistry()
resourceRegistry.register(demoResource, demoProvider)
