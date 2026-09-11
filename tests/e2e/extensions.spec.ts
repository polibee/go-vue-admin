import { expect, test } from '@playwright/test'

test('modules and plugins have separate admin pages', async ({ page }) => {
  await page.goto('/login')
  await page.getByLabel('邮箱').fill(process.env.AUTH_BOOTSTRAP_EMAIL ?? 'admin@example.com')
  await page.getByLabel('密码', { exact: true }).fill(process.env.AUTH_BOOTSTRAP_PASSWORD ?? 'test-only-password')
  await page.getByRole('button', { name: '登录', exact: true }).click()
  await expect(page).toHaveURL(/\/admin\/dashboard/)
  await page.getByRole('link', { name: '业务模块', exact: true }).click()
  await expect(page).toHaveURL(/\/admin\/modules$/)
  await expect(page.getByRole('heading', { name: '业务模块' })).toBeVisible()
  await expect(page.getByText('示例模块', { exact: true })).toBeVisible()
  await page.getByRole('link', { name: '查看', exact: true }).click()
  await expect(page).toHaveURL(/\/admin\/modules\/example$/)
  await expect(page.getByText('示例模块已注册并启动', { exact: true })).toBeVisible()

  await page.getByRole('link', { name: '平台插件', exact: true }).click()
  await expect(page).toHaveURL(/\/admin\/plugins$/)
  await expect(page.getByRole('heading', { name: '平台插件' })).toBeVisible()
  await expect(page.getByText('示例插件', { exact: true })).toBeVisible()
  const enableButton = page.getByRole('button', { name: '启用', exact: true })
  if (await enableButton.count() === 0) {
    await page.getByRole('button', { name: '停用', exact: true }).click()
    await expect(enableButton).toBeVisible()
  }
  await enableButton.click()
  await expect(page.getByText('状态：已启用', { exact: true })).toBeVisible()
  await page.getByRole('link', { name: '查看' }).last().click()
  await expect(page).toHaveURL(/\/admin\/plugins\/example-plugin$/)
  await expect(page.getByText('示例插件正在运行', { exact: true })).toBeVisible()
  await page.goto('/admin/plugins')
  await page.getByRole('button', { name: '停用', exact: true }).click()
  await expect(page.getByText('状态：已停用', { exact: true })).toBeVisible()
})
