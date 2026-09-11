import { describe, expect, it } from 'vitest'

import source from './ApiDocsPage.vue?raw'

describe('ApiDocsPage', () => {
  it('uses the local OpenAPI contract in read-only mode', () => {
    expect(source).toContain("fetch('/api/docs/openapi.json'")
    expect(source).toContain("apiDocs.all")
    expect(source).toContain("apiDocs.extensions")
    expect(source).toContain('hideClientButton: true')
    expect(source).toContain('hideDownloadButton: true')
    expect(source).toContain('hideTestRequestButton: true')
  })
})
