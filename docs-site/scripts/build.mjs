import fs from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const scriptDir = path.dirname(fileURLToPath(import.meta.url))
const siteDir = path.resolve(scriptDir, '..')
const contentDir = path.join(siteDir, 'src', 'content')
const outputDir = path.join(siteDir, 'dist')
const assetsDir = path.join(outputDir, 'assets')
const pages = [
  { slug: 'index', titles: ['Developer Docs', '开发者文档'], file: 'index.md', sections: ['Getting started', '开始使用'] },
  { slug: 'getting-started', titles: ['Getting started', '快速开始'], file: 'getting-started.md', sections: ['Getting started', '开始使用'] },
  { slug: 'architecture', titles: ['Architecture and boundaries', '架构与边界'], file: 'architecture.md', sections: ['Core concepts', '核心概念'] },
  { slug: 'modules', titles: ['Build business modules', '开发业务模块'], file: 'modules.md', sections: ['Extension development', '扩展开发'] },
  { slug: 'plugins', titles: ['Build platform plugins', '开发平台插件'], file: 'plugins.md', sections: ['Extension development', '扩展开发'] },
  { slug: 'resources', titles: ['Resources and code generation', '资源与代码生成'], file: 'resources.md', sections: ['Extension development', '扩展开发'] },
  { slug: 'api-contracts', titles: ['API contracts and clients', 'API 契约与客户端'], file: 'api-contracts.md', sections: ['Integration', '集成开发'] },
  { slug: 'database', titles: ['Database and production config', '数据库与生产配置'], file: 'database.md', sections: ['Operations', '部署运维'] },
  { slug: 'production-readiness', titles: ['Production readiness', '生产级检查清单'], file: 'production-readiness.md', sections: ['Operations', '部署运维'] },
  { slug: 'contributing', titles: ['Contributing and releases', '贡献与发布'], file: 'contributing.md', sections: ['Collaboration', '项目协作'] },
]
const locales = [
  { id: 'en', dir: 'en', label: 'English', lang: 'en' },
  { id: 'zh-CN', dir: '', label: '简体中文', lang: 'zh-CN' },
]

fs.rmSync(outputDir, { recursive: true, force: true })
fs.mkdirSync(assetsDir, { recursive: true })
fs.copyFileSync(path.join(siteDir, 'src', 'style.css'), path.join(assetsDir, 'style.css'))
fs.copyFileSync(path.join(siteDir, 'src', 'app.js'), path.join(assetsDir, 'app.js'))
for (const locale of locales) {
  const localeOutput = locale.id === 'en' ? outputDir : path.join(outputDir, locale.id)
  const localeContent = locale.dir ? path.join(contentDir, locale.dir) : contentDir
  fs.mkdirSync(localeOutput, { recursive: true })
  const searchItems = pages.map((page) => {
    const source = fs.readFileSync(path.join(localeContent, page.file), 'utf8')
    return { slug: page.slug, title: page.titles[locale.id === 'en' ? 0 : 1], section: page.sections[locale.id === 'en' ? 0 : 1], text: source.replace(/\s+/g, ' ').slice(0, 1200) }
  })
  for (const page of pages) {
    const source = fs.readFileSync(path.join(localeContent, page.file), 'utf8')
    const localizedPage = { ...page, title: page.titles[locale.id === 'en' ? 0 : 1], section: page.sections[locale.id === 'en' ? 0 : 1] }
    fs.writeFileSync(path.join(localeOutput, page.slug + '.html'), renderDocument(localizedPage, renderMarkdown(source), locale, searchItems))
  }
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
function renderDocument(page, body, locale, searchItems) {
  const languageIndex = locale.id === 'en' ? 0 : 1
  const groups = [...new Set(pages.map((item) => item.sections[languageIndex]))]
  const navigation = groups.map((group) => {
    const items = pages.filter((item) => item.sections[languageIndex] === group).map((item) => {
      const active = item.slug === page.slug ? ' class="active"' : ''
      return '<a' + active + ' href="' + item.slug + '.html">' + item.titles[languageIndex] + '</a>'
    }).join('')
    return '<div class="nav-group"><div class="nav-label">' + group + '</div>' + items + '</div>'
  }).join('')
  const localePrefix = locale.id === 'en' ? '' : '../'
  const languageLinks = locale.id === 'en' ? '<a href="zh-CN/index.html">简体中文</a>' : '<a href="../index.html">English</a>'
  return '<!doctype html><html lang="' + locale.lang + '"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><meta name="description" content="Go Vue Admin developer documentation"><title>' + page.title + ' · Go Vue Admin</title><link rel="stylesheet" href="' + localePrefix + 'assets/style.css"></head><body><header class="topbar"><a class="brand" href="' + localePrefix + 'index.html"><span class="brand-mark">G</span><span>Go Vue Admin<small>Developer Docs</small></span></a><div class="top-actions"><input id="search" type="search" placeholder="' + (locale.id === 'en' ? 'Search docs' : '搜索文档') + '"><button id="theme" type="button" aria-label="' + (locale.id === 'en' ? 'Toggle theme' : '切换主题') + '">◐</button>' + languageLinks + '<a href="https://github.com/polibee/go-vue-admin">GitHub</a></div></header><div id="search-results" class="search-results" hidden></div><div class="layout"><aside class="sidebar">' + navigation + '</aside><main class="content"><div class="eyebrow">' + page.section + '</div>' + body + '<footer>Go Vue Admin · Developer documentation</footer></main></div><script>window.DOCS_SEARCH=' + JSON.stringify(searchItems) + '</script><script src="' + localePrefix + 'assets/app.js"></script></body></html>'
}
