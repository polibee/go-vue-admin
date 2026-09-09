<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toTypedSchema } from '@vee-validate/zod'
import { z } from 'zod'
import { Check, Copy, Eye, EyeOff, LoaderCircle } from '@lucide/vue'

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
const showPassword = ref(false)
const authService = new AuthService()

async function loadBootstrapCredentials() {
  credentialsError.value = null
  try {
    bootstrapCredentials.value = await authService.bootstrapCredentials()
  } catch {
    credentialsError.value = '开发环境凭据暂时无法加载，请检查后端服务。'
  }
}

onMounted(loadBootstrapCredentials)

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
        <Form v-slot="{ handleSubmit, setValues }" :validation-schema="formSchema" :initial-values="{ email: '', password: '' }" as="div">
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
            <Button type="button" variant="secondary" class="w-full" @click="setValues({ email: bootstrapCredentials.email, password: bootstrapCredentials.password })">
              一键填入登录表单
            </Button>
          </div>
          </Alert>
          <div v-else-if="credentialsError" class="mb-5 flex items-center justify-between gap-3 rounded-lg border border-dashed p-3 text-sm text-muted-foreground" role="status">
            <span>{{ credentialsError }}</span>
            <Button type="button" variant="outline" size="sm" @click="loadBootstrapCredentials">重试</Button>
          </div>
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
                <div class="relative">
                  <FormControl><Input :type="showPassword ? 'text' : 'password'" autocomplete="current-password" placeholder="请输入密码" class="pr-10" v-bind="componentField" /></FormControl>
                  <Button type="button" variant="ghost" size="icon" class="absolute right-0 top-0 h-full px-3" :aria-label="showPassword ? '隐藏密码' : '显示密码'" @click="showPassword = !showPassword">
                    <EyeOff v-if="showPassword" data-icon="inline-start" />
                    <Eye v-else data-icon="inline-start" />
                  </Button>
                </div>
                <FormMessage />
              </FormItem>
            </FormField>
            <p v-if="auth.error" class="text-sm text-destructive" role="alert">{{ auth.error }}</p>
            <Button type="submit" class="w-full" :disabled="submitting">
              <LoaderCircle v-if="submitting" class="animate-spin" data-icon="inline-start" />
              {{ submitting ? '登录中…' : '登录' }}
            </Button>
          </form>
        </Form>
      </CardContent>
    </Card>
  </main>
</template>
