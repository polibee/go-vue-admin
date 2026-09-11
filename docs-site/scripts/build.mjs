import fs from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const scriptDir = path.dirname(fileURLToPath(import.meta.url))
const siteDir = path.resolve(scriptDir, '..')
const contentDir = path.join(siteDir, 'src', 'content')
const outputDir = path.join(siteDir, 'dist')
const assetsDir = path.join(outputDir, 'assets')
const pages = [
  { slug: 'index', title: '开发者文档', file: 'index.md', section: '开始使用' },
  { slug: 'getting-started', title: '快速开始', file: 'getting-started.md', section: '开始使用' },
  { slug: 'architecture', title: '架构与边界', file: 'architecture.md', section: '核心概念' },
  { slug: 'modules', title: '开发业务模块', file: 'modules.md', section: '扩展开发' },
  { slug: 'plugins', title: '开发平台插件', file: 'plugins.md', section: '扩展开发' },
  { slug: 'resources', title: '资源与代码生成', file: 'resources.md', section: '扩展开发' },
  { slug: 'api-contracts', title: 'API 契约与客户端', file: 'api-contracts.md', section: '集成开发' },
  { slug: 'database', title: '数据库与生产配置', file: 'database.md', section: '部署运维' },
  { slug: 'production-readiness', title: '生产级检查清单', file: 'production-readiness.md', section: '部署运维' },
  { slug: 'contributing', title: '贡献与发布', file: 'contributing.md', section: '项目协作' },
]
const groups = [...new Set(pages.map((page) => page.section))]
const searchItems = pages.map((page) => {
  const source = fs.readFileSync(path.join(contentDir, page.file), 'utf8')
  return { slug: page.slug, title: page.title, section: page.section, text: source.replace(/\s+/g, ' ').slice(0, 1200) }
})

fs.rmSync(outputDir, { recursive: true, force: true })
fs.mkdirSync(assetsDir, { recursive: true })
fs.copyFileSync(path.join(siteDir, 'src', 'style.css'), path.join(assetsDir, 'style.css'))
fs.copyFileSync(path.join(siteDir, 'src', 'app.js'), path.join(assetsDir, 'app.js'))
for (const page of pages) {
  const source = fs.readFileSync(path.join(contentDir, page.file), 'utf8')
  fs.writeFileSync(path.join(outputDir, page.slug + '.html'), renderDocument(page, renderMarkdown(source)))
}

function escapeHtml(value) {
  return value.replaceAll('&', '&amp;').replaceAll('<', '&lt;').replaceAll('>', '&gt;').replaceAll('"', '&quot;')
}
function inlineMarkdown(value) {
  let html = escapeHtml(value)
  const tick = String.fromCharCode(96)
  html = html.replace(new RegExp(tick + '([^' + tick + ']*)' + tick, 'g'), '<code>$1</code>')
  html = html.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2">$1</a>')
  return html
}
function renderMarkdown(source) {
  const lines = source.replace(/\r\n/g, '\n').split('\n')
  const result = []
  let paragraph = []
  let list = []
  let code = []
  const flushParagraph = () => { if (paragraph.length) result.push('<p>' + paragraph.map(inlineMarkdown).join(' ') + '</p>'); paragraph = [] }
  const flushList = () => { if (list.length) result.push('<ul>' + list.map((item) => '<li>' + inlineMarkdown(item) + '</li>').join('') + '</ul>'); list = [] }
  const flushCode = () => { if (code.length) result.push('<pre><code>' + escapeHtml(code.join('\n')) + '</code></pre>'); code = [] }
  for (const line of lines) {
    if (line.startsWith('    ')) { flushParagraph(); flushList(); code.push(line.slice(4)); continue }
    flushCode()
    if (!line.trim()) { flushParagraph(); flushList(); continue }
    if (line.startsWith('# ')) { flushParagraph(); flushList(); result.push('<h1>' + inlineMarkdown(line.slice(2)) + '</h1>') }
    else if (line.startsWith('## ')) { flushParagraph(); flushList(); result.push('<h2>' + inlineMarkdown(line.slice(3)) + '</h2>') }
    else if (line.startsWith('### ')) { flushParagraph(); flushList(); result.push('<h3>' + inlineMarkdown(line.slice(4)) + '</h3>') }
    else if (line.startsWith('- ')) { flushParagraph(); list.push(line.slice(2)) }
    else { flushList(); paragraph.push(line) }
  }
  flushParagraph(); flushList(); flushCode()
  return result.join('\n')
}
function renderDocument(page, body) {
  const navigation = groups.map((group) => {
    const items = pages.filter((item) => item.section === group).map((item) => {
      const active = item.slug === page.slug ? ' class="active"' : ''
      return '<a' + active + ' href="' + item.slug + '.html">' + item.title + '</a>'
    }).join('')
    return '<div class="nav-group"><div class="nav-label">' + group + '</div>' + items + '</div>'
  }).join('')
  return '<!doctype html><html lang="zh-CN"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><meta name="description" content="Go Vue Admin 开发者文档"><title>' + page.title + ' · Go Vue Admin</title><link rel="stylesheet" href="assets/style.css"></head><body><header class="topbar"><a class="brand" href="index.html"><span class="brand-mark">G</span><span>Go Vue Admin<small>Developer Docs</small></span></a><div class="top-actions"><input id="search" type="search" placeholder="搜索文档"><button id="theme" type="button" aria-label="切换主题">◐</button><a href="https://github.com/polibee/go-vue-admin">GitHub</a></div></header><div id="search-results" class="search-results" hidden></div><div class="layout"><aside class="sidebar">' + navigation + '</aside><main class="content"><div class="eyebrow">' + page.section + '</div>' + body + '<footer>Go Vue Admin · 面向开发者的独立文档站</footer></main></div><script>window.DOCS_SEARCH=' + JSON.stringify(searchItems) + '</script><script src="assets/app.js"></script></body></html>'
}
