import { expect, test } from '@playwright/test'

test('switches language and semantic theme palettes from the admin shell', async ({ page }) => {
  await page.goto('/login')
  await page.getByLabel('邮箱').fill(process.env.AUTH_BOOTSTRAP_EMAIL ?? 'admin@example.com')
  await page.getByLabel('密码', { exact: true }).fill(process.env.AUTH_BOOTSTRAP_PASSWORD ?? 'test-only-password')
  await page.getByRole('button', { name: '登录', exact: true }).click()
  await expect(page).toHaveURL(/\/admin\/dashboard/)

  await page.getByLabel('language').click()
  await expect(page.getByRole('link', { name: 'Dashboard' })).toBeVisible()

  await page.getByLabel('theme-palette').click()
  await expect(page.locator('html')).toHaveAttribute('data-palette', 'semi')

  await page.reload()
  await expect(page.getByRole('link', { name: 'Dashboard' })).toBeVisible()
  await expect(page.locator('html')).toHaveAttribute('data-palette', 'semi')
  await page.evaluate(() => localStorage.clear())
})
