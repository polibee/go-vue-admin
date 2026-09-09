import { expect, test } from '@playwright/test'

const email = process.env.AUTH_BOOTSTRAP_EMAIL ?? 'admin@example.com'
const password = process.env.AUTH_BOOTSTRAP_PASSWORD ?? 'test-only-password'

test('generic resource pages support list, create, edit, and show flows', async ({ page }) => {
  await page.goto('/login')
  await page.getByLabel('邮箱').fill(email)
  await page.getByLabel('密码', { exact: true }).fill(password)
  await page.getByRole('button', { name: '登录' }).click()
  await expect(page).toHaveURL(/\/admin\/dashboard/)

  await page.goto('/admin/resources/demo')
  await expect(page.getByRole('heading', { name: '示例资源' })).toBeVisible()
  await expect(page.getByText('资源引擎示例')).toBeVisible()

  await page.getByRole('button', { name: '新建' }).click()
  await expect(page).toHaveURL(/\/admin\/resources\/demo\/create/)
  await page.getByLabel('名称').fill('浏览器创建的记录')
  await page.getByRole('combobox').click()
  await page.getByRole('option', { name: '启用' }).click()
  await page.getByLabel('负责人').fill('QA')
  await page.getByRole('button', { name: '创建' }).click()
  await expect(page).toHaveURL(/\/admin\/resources\/demo\/demo-/)
  await expect(page.getByText('浏览器创建的记录', { exact: true }).last()).toBeVisible()
})
