import { resourceDefinition } from './resource';

// This contract check uses the existing Node type/runtime surface; the generator
// does not add a frontend test-runner dependency.
if (resourceDefinition.name !== "departments") throw new Error('generated resource name mismatch');
if (resourceDefinition.route !== "/departments") throw new Error('generated resource route mismatch');
if (JSON.stringify(resourceDefinition.actions) !== JSON.stringify([{ name: "view", label: "View", kind: "", permission: "admin.departments.view", batch: false, payload: "" }, { name: "create", label: "Create", kind: "", permission: "admin.departments.create", batch: false, payload: "" }, { name: "update", label: "Update", kind: "", permission: "admin.departments.update", batch: false, payload: "" }, { name: "delete", label: "Delete", kind: "", permission: "admin.departments.delete", batch: false, payload: "" }])) throw new Error('generated resource actions mismatch');
