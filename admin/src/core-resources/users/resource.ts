import { defineResource } from '@/resource-engine/core/ResourceDefinition'

export interface UserResourceRecord {
  id: string
  email: string
  name: string
  active: boolean
  role_ids: string
}

export const userResource = defineResource<UserResourceRecord>({
  name: 'users',
  label: '用户',
  permissions: {
    list: 'users.view',
    get: 'users.view',
    create: 'users.create',
    update: 'users.update',
    delete: 'users.delete',
    bulkDelete: 'users.delete',
  },
  columns: [
    { key: 'name', label: '姓名', sortable: true },
    { key: 'email', label: '邮箱', sortable: true },
    { key: 'active', label: '状态' },
  ],
  fields: [
    { name: 'name', label: '姓名', required: true, placeholder: '请输入姓名' },
    { name: 'email', label: '邮箱', required: true, placeholder: '请输入邮箱' },
    { name: 'active', label: '启用', type: 'checkbox' },
    { name: 'role_ids', label: '角色 ID', placeholder: '多个角色用逗号分隔' },
  ],
})

export const userResourceSeed: UserResourceRecord[] = [
  { id: 'bootstrap-admin', email: 'admin@example.com', name: 'Platform Admin', active: true, role_ids: 'platform-admin' },
]
