import { resourceDefinition } from './resource'

if (resourceDefinition.name !== 'announcements') throw new Error('generated resource name mismatch')
if (resourceDefinition.route !== '/announcements') throw new Error('generated resource route mismatch')
if (JSON.stringify(resourceDefinition.actions) !== JSON.stringify(['view', 'create', 'update', 'delete'])) throw new Error('generated resource actions mismatch')
