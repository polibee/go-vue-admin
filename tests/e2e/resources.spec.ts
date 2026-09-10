import { expect, test } from '@playwright/test'

const email = process.env.AUTH_BOOTSTRAP_EMAIL ?? 'admin@example.com'
const password = process.env.AUTH_BOOTSTRAP_PASSWORD ?? 'test-only-password'

test('generic resource pages support list, create, edit, and show flows', async ({ page }) => {
  await page.goto('/login')
  await page.getByLabel('邮箱').fill(email)
  await page.getByLabel('密码', { exact: true }).fill(password)
  await page.getByRole('button', { name: '登录', exact: true }).click()
  await expect(page).toHaveURL(/\/admin\/dashboard/)

  await page.goto('/admin/resources/demo')
  await expect(page.getByRole('heading', { name: '示例资源' })).toBeVisible()
  await expect(page.getByText('资源引擎示例')).toBeVisible()

  await page.getByRole('button', { name: '新建' }).click()
  await expect(page).toHaveURL(/\/admin\/resources\/demo\/create/)
  const createdName = `浏览器创建的记录-${Date.now()}`
  await page.getByLabel('名称').fill(createdName)
  await page.getByRole('combobox').click()
  await page.getByRole('option', { name: '启用' }).click()
  await page.getByLabel('负责人').fill('QA')
  await page.getByRole('button', { name: '创建' }).click()
  await expect(page).toHaveURL(/\/admin\/resources\/demo\/demo-/)
  await expect(page.getByText(createdName, { exact: true }).last()).toBeVisible()

  await page.getByRole('button', { name: '编辑' }).click()
  await expect(page).toHaveURL(/\/admin\/resources\/demo\/demo-.*\/edit/)
  const updatedName = `${createdName}-已更新`
  await page.getByLabel('名称').fill(updatedName)
  await page.getByRole('combobox').click()
  await page.getByRole('option', { name: '启用' }).click()
  await page.getByLabel('负责人').fill('QA')
  await page.getByRole('button', { name: '保存' }).click()
  await expect(page.getByText(updatedName, { exact: true }).last()).toBeVisible()

  await page.getByRole('button', { name: '删除' }).click()
  await expect(page.getByRole('alertdialog')).toBeVisible()
  await page.getByRole('alertdialog').getByRole('button', { name: '确认删除' }).click()
  await expect(page).toHaveURL(/\/admin\/resources\/demo$/)
  await expect(page.getByText(updatedName, { exact: true })).not.toBeVisible()
})

test('generic resource list supports filters, sorting, and bulk delete', async ({ page }) => {
  await page.goto('/login')
  await page.getByLabel('邮箱').fill(email)
  await page.getByLabel('密码', { exact: true }).fill(password)
  await page.getByRole('button', { name: '登录', exact: true }).click()
  await expect(page).toHaveURL(/\/admin\/dashboard/)
  await page.goto('/admin/resources/demo')

  await page.getByLabel('状态筛选').click()
  await page.getByRole('option', { name: '启用' }).click()
  await expect(page.getByText('资源引擎示例')).toBeVisible()
  await expect(page.getByText('可编辑记录')).not.toBeVisible()

  await page.getByRole('button', { name: '按名称排序' }).click()
  await expect(page.getByText('page=1')).toBeVisible()

  await page.getByLabel('选择 demo-1').click()
  await page.getByRole('button', { name: '批量删除' }).click()
  await expect(page.getByRole('alertdialog')).toBeVisible()
  await page.getByRole('alertdialog').getByRole('button', { name: '确认删除' }).click()
  await expect(page.getByText('资源引擎示例')).not.toBeVisible()
})

test('core user, role, and permission resources use the same generic routes', async ({ page }) => {
  await page.goto('/login')
  await page.getByLabel('邮箱').fill(email)
  await page.getByLabel('密码', { exact: true }).fill(password)
  await page.getByRole('button', { name: '登录', exact: true }).click()
  await expect(page).toHaveURL(/\/admin\/dashboard/)
  await expect(page.getByRole('link', { name: '用户' })).toBeVisible()

  await page.goto('/admin/resources/users')
  await expect(page.getByRole('heading', { name: '用户' })).toBeVisible()
  await expect(page.getByRole('cell', { name: email })).toBeVisible()

  await page.goto('/admin/resources/roles')
  await expect(page.getByRole('heading', { name: '角色' })).toBeVisible()
  await expect(page.getByText('平台管理员')).toBeVisible()

  await page.goto('/admin/resources/permissions')
  await expect(page.getByRole('heading', { name: '权限' })).toBeVisible()
  await expect(page.getByText('dashboard.view')).toBeVisible()
})
