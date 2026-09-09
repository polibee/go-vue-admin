import { defineResource } from '@/resource-engine/core/ResourceDefinition'

export interface RoleResourceRecord {
  id: string
  name: string
  description: string
  permissions: string
}

export const roleResource = defineResource<RoleResourceRecord>({
  name: 'roles',
  label: '角色',
  permissions: {
    list: 'roles.view',
    get: 'roles.view',
    create: 'roles.create',
    update: 'roles.update',
    delete: 'roles.delete',
    bulkDelete: 'roles.delete',
  },
  columns: [
    { key: 'name', label: '名称', sortable: true },
    { key: 'description', label: '描述' },
    { key: 'permissions', label: '权限' },
  ],
  fields: [
    { name: 'name', label: '名称', required: true, placeholder: '请输入角色名称' },
    { name: 'description', label: '描述', type: 'textarea', placeholder: '请输入角色描述' },
    { name: 'permissions', label: '权限', type: 'textarea', required: true, placeholder: '每行一个权限，例如 users.view' },
  ],
})

export const roleResourceSeed: RoleResourceRecord[] = [
  { id: 'platform-admin', name: '平台管理员', description: '访问平台资源的开发管理员', permissions: 'dashboard.view\nsettings.view\nusers.*\nroles.*\npermissions.*' },
]
