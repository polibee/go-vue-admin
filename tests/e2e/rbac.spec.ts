import { expect, test } from '@playwright/test'

const email = process.env.AUTH_BOOTSTRAP_EMAIL ?? 'admin@example.com'
const password = process.env.AUTH_BOOTSTRAP_PASSWORD ?? 'test-only-password'

test.describe('RBAC navigation', () => {
  test('loads the permission-aware menu for the authenticated session', async ({ page }) => {
    await page.goto('/login')
    await page.getByLabel('邮箱').fill(email)
    await page.getByLabel('密码').fill(password)
    await page.getByRole('button', { name: '登录' }).click()

    await expect(page).toHaveURL(/\/admin\/dashboard/)
    await expect(page.getByRole('link', { name: '仪表盘' })).toBeVisible()
    await expect(page.getByRole('link', { name: '设置' })).toBeVisible()
    await expect(page.getByText('开发进度', { exact: true })).toBeVisible()
    await expect(page.getByText('Task 7 · Resource Engine Core')).toBeVisible()
  })
})
