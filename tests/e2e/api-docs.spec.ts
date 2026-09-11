import { expect, test } from '@playwright/test'

test('API docs are separated into functional pages', async ({ page }) => {
  await page.goto('/login')
  await page.getByLabel('邮箱').fill(process.env.AUTH_BOOTSTRAP_EMAIL ?? 'admin@example.com')
  await page.getByLabel('密码', { exact: true }).fill(process.env.AUTH_BOOTSTRAP_PASSWORD ?? 'test-only-password')
  await page.getByRole('button', { name: '登录', exact: true }).click()
  await expect(page).toHaveURL(/\/admin\/dashboard/)
  await page.getByRole('link', { name: 'API 文档', exact: true }).click()
  await expect(page).toHaveURL(/\/admin\/api-docs$/)
  await expect(page.getByRole('tab', { name: /全部 API/ })).toBeVisible()
  await expect(page.getByRole('tab', { name: /模块与插件/ })).toBeVisible()
  await expect(page.getByRole('tab', { name: /业务资源/ })).toBeVisible()
  await page.getByRole('tab', { name: /模块与插件/ }).click()
  await expect(page.getByRole('tab', { name: /模块与插件/ })).toHaveAttribute('data-state', 'active')
})
