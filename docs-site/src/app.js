(() => {
  const root = document.documentElement
  const theme = document.getElementById('theme')
  const stored = localStorage.getItem('go-vue-admin-docs-theme')
  if (stored) root.dataset.theme = stored
  theme?.addEventListener('click', () => {
    const next = root.dataset.theme === 'dark' ? 'light' : 'dark'
    root.dataset.theme = next
    localStorage.setItem('go-vue-admin-docs-theme', next)
  })
  const input = document.getElementById('search')
  const results = document.getElementById('search-results')
  input?.addEventListener('input', () => {
    const query = input.value.trim().toLowerCase()
    if (!query) { results.hidden = true; results.innerHTML = ''; return }
    const matches = (window.DOCS_SEARCH || []).filter((item) => (item.title + ' ' + item.text).toLowerCase().includes(query)).slice(0, 8)
    results.innerHTML = matches.length ? matches.map((item) => '<a href="' + item.slug + '.html"><strong>' + item.title + '</strong><small>' + item.section + '</small></a>').join('') : '<small>没有找到匹配的文档。</small>'
    results.hidden = false
  })
})()
