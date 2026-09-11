import { describe, expect, it } from 'vitest'

import source from './ApiDocsPage.vue?raw'

describe('ApiDocsPage', () => {
  it('uses the local OpenAPI contract in read-only mode', () => {
    expect(source).toContain("fetch('/api/docs/openapi.json'")
    expect(source).toContain('全部 API')
    expect(source).toContain('模块与插件')
    expect(source).toContain('hideClientButton: true')
    expect(source).toContain('hideDownloadButton: true')
    expect(source).toContain('hideTestRequestButton: true')
  })
})
