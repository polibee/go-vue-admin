import { expect, test } from '@playwright/test'

const email = process.env.AUTH_BOOTSTRAP_EMAIL ?? 'admin@example.com'
const password = process.env.AUTH_BOOTSTRAP_PASSWORD ?? 'test-only-password'

test.describe('cookie session authentication', () => {
  test('guards the admin, keeps the session after refresh, and revokes it on logout', async ({ page }) => {
    await page.goto('/admin/dashboard')
    await expect(page).toHaveURL(/\/login/, { timeout: 15000 })
    await expect(page.getByLabel('开发账号')).toHaveValue(email)
    await expect(page.getByLabel('开发密码')).toHaveValue(password)
    await page.getByRole('button', { name: '复制账号' }).click()
    await expect(page.getByRole('button', { name: '账号已复制' })).toBeVisible()

    await page.getByLabel('邮箱').fill(email)
    await page.getByLabel('密码', { exact: true }).fill(password)
    await page.getByRole('button', { name: '登录' }).click()
    await expect(page).toHaveURL(/\/admin\/dashboard/)
    await expect(page.getByText('开发进度', { exact: true })).toBeVisible({ timeout: 15000 })

    await page.reload()
    await expect(page).toHaveURL(/\/admin\/dashboard/)
    await expect(page.getByText(email)).toBeVisible({ timeout: 15000 })

    await page.getByRole('button', { name: '退出' }).click()
    await expect(page).toHaveURL(/\/login/, { timeout: 15000 })
    await expect(page.getByRole('heading', { name: '登录管理后台' })).toBeVisible()

    await page.goto('/admin/dashboard')
    await expect(page).toHaveURL(/\/login/, { timeout: 15000 })
  })
})
