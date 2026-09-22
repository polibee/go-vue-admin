import assert from 'node:assert/strict'
import { resolveResourcePageMode, validateCustomPageOverrides } from './page-resolver.ts'

assert.equal(resolveResourcePageMode({ name: 'departments' }), 'generic')
assert.equal(resolveResourcePageMode({ name: 'departments', pageMode: 'generic' }), 'generic')
assert.equal(resolveResourcePageMode({ name: 'orders', pageMode: 'custom' }), 'custom')

assert.doesNotThrow(() => validateCustomPageOverrides({
  list: 'OrderListPage',
  form: 'OrderFormPage',
  detail: 'OrderDetailPage',
}))
assert.throws(() => validateCustomPageOverrides({ list: 'OrderListPage', form: 'OrderFormPage' }), /detail/)
