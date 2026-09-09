<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toTypedSchema } from '@vee-validate/zod'
import { z } from 'zod'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { useAuth } from '@/core/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuth()
const submitting = ref(false)

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
</script>

<template>
  <main class="flex min-h-svh items-center justify-center bg-muted/30 px-4 py-10">
    <Card class="w-full max-w-md">
      <CardHeader class="space-y-1">
        <CardTitle class="text-2xl">登录管理后台</CardTitle>
        <CardDescription>使用平台账号进入 Go Vue Admin。</CardDescription>
      </CardHeader>
      <CardContent>
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
