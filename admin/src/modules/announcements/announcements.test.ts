import { resourceDefinition } from './resource'

if (resourceDefinition.name !== 'announcements') throw new Error('generated resource name mismatch')
if (resourceDefinition.admin_route !== '/admin/announcements') throw new Error('generated resource route mismatch')
if (resourceDefinition.actions.length !== 4 || resourceDefinition.actions.some((action) => action.batch)) throw new Error('generated resource actions mismatch')
