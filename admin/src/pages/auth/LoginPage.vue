<script setup lang="ts">
import { ref } from 'vue'
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toTypedSchema } from '@vee-validate/zod'
import { z } from 'zod'
import { Check, Copy } from '@lucide/vue'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { AuthService, type BootstrapCredentials } from '@/core/auth/auth.service'
import { useAuth } from '@/core/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuth()
const submitting = ref(false)
const bootstrapCredentials = ref<BootstrapCredentials | null>(null)
const credentialsError = ref<string | null>(null)
const copiedField = ref<'email' | 'password' | null>(null)
const authService = new AuthService()

onMounted(async () => {
  try {
    bootstrapCredentials.value = await authService.bootstrapCredentials()
  } catch {
    credentialsError.value = '开发环境凭据暂时无法加载，请检查后端服务。'
  }
})

const formSchema = toTypedSchema(z.object({
  email: z.string().email('请输入有效邮箱'),
  password: z.string().min(1, '请输入密码'),
}))

function redirectTarget(): string {
  const redirect = route.query.redirect
  return typeof redirect === 'string' && redirect.startsWith('/') ? redirect : '/admin/dashboard'
}

async function onSubmit(values: Record<string, unknown>) {
  submitting.value = true
  try {
    await auth.login({ email: String(values.email ?? ''), password: String(values.password ?? '') })
    await router.replace(redirectTarget())
  } finally {
    submitting.value = false
  }
}

async function copyCredential(field: 'email' | 'password') {
  const value = bootstrapCredentials.value?.[field]
  if (!value) return

  try {
    await navigator.clipboard.writeText(value)
  } catch {
    const fallback = document.createElement('textarea')
    fallback.value = value
    fallback.setAttribute('readonly', '')
    fallback.style.position = 'fixed'
    fallback.style.opacity = '0'
    document.body.appendChild(fallback)
    fallback.select()
    document.execCommand('copy')
    fallback.remove()
  }
  copiedField.value = field
  window.setTimeout(() => {
    if (copiedField.value === field) copiedField.value = null
  }, 1500)
}
</script>

<template>
  <main class="flex min-h-svh items-center justify-center bg-muted/30 px-4 py-10">
    <Card class="w-full max-w-md">
      <CardHeader class="space-y-1">
        <CardTitle class="text-2xl">登录管理后台</CardTitle>
        <CardDescription>使用平台账号进入 Go Vue Admin。</CardDescription>
      </CardHeader>
      <CardContent>
        <Alert v-if="bootstrapCredentials" class="mb-5">
          <AlertTitle>开发环境登录凭据</AlertTitle>
          <AlertDescription>凭据从后端实时读取，修改 AUTH_BOOTSTRAP_* 后刷新页面即可同步。</AlertDescription>
          <div class="mt-3 flex flex-col gap-3">
            <div class="flex flex-col gap-2">
              <Label for="bootstrap-email">账号</Label>
              <div class="flex gap-2">
                <Input id="bootstrap-email" :model-value="bootstrapCredentials.email" readonly aria-label="开发账号" />
                <Button type="button" variant="outline" size="icon" :aria-label="copiedField === 'email' ? '账号已复制' : '复制账号'" @click="copyCredential('email')">
                  <Check v-if="copiedField === 'email'" data-icon="inline-start" />
                  <Copy v-else data-icon="inline-start" />
                </Button>
              </div>
            </div>
            <div class="flex flex-col gap-2">
              <Label for="bootstrap-password">密码</Label>
              <div class="flex gap-2">
                <Input id="bootstrap-password" :model-value="bootstrapCredentials.password" readonly aria-label="开发密码" />
                <Button type="button" variant="outline" size="icon" :aria-label="copiedField === 'password' ? '密码已复制' : '复制密码'" @click="copyCredential('password')">
                  <Check v-if="copiedField === 'password'" data-icon="inline-start" />
                  <Copy v-else data-icon="inline-start" />
                </Button>
              </div>
            </div>
          </div>
        </Alert>
        <p v-else-if="credentialsError" class="mb-5 text-sm text-muted-foreground" role="status">{{ credentialsError }}</p>
        <Form v-slot="{ handleSubmit }" :validation-schema="formSchema" as="div">
          <form class="grid gap-5" @submit="handleSubmit($event, onSubmit)">
            <FormField v-slot="{ componentField }" name="email">
              <FormItem>
                <FormLabel>邮箱</FormLabel>
                <FormControl><Input type="email" autocomplete="email" placeholder="admin@example.com" v-bind="componentField" /></FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
            <FormField v-slot="{ componentField }" name="password">
              <FormItem>
                <FormLabel>密码</FormLabel>
                <FormControl><Input type="password" autocomplete="current-password" placeholder="请输入密码" v-bind="componentField" /></FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
            <p v-if="auth.error" class="text-sm text-destructive" role="alert">{{ auth.error }}</p>
            <Button type="submit" class="w-full" :disabled="submitting">
              {{ submitting ? '登录中…' : '登录' }}
            </Button>
          </form>
        </Form>
      </CardContent>
    </Card>
  </main>
</template>
