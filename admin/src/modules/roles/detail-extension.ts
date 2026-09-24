import { defineComponent, h } from 'vue'
import RolePermissionsPanel from './components/RolePermissionsPanel.vue'
import { registerResourceDetailExtension } from '@/core/resource/detail-extensions'
import type { RolePermissionAssignment } from '@/generated/api'

const RoleDetailExtension = defineComponent({
  name: 'RoleDetailExtension',
  props: {
    resource: { type: String, required: true },
    id: { type: String, required: true },
    record: { type: Object, required: true },
  },
  setup(props) {
    return () => h(RolePermissionsPanel, {
      roleId: Number(props.id),
      roleName: String(props.record.name || ''),
      assignments: Array.isArray(props.record.permissions) ? props.record.permissions as RolePermissionAssignment[] : [],
    })
  },
})

export function registerRoleDetailExtension() {
  registerResourceDetailExtension({ resource: 'roles', component: RoleDetailExtension })
}
