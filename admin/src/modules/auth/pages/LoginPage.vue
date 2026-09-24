<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Check, Copy, Languages } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import { ApiError, errorMessageKey } from '@/lib/api'
import { DEMO_CREDENTIALS, type DemoCredentialKey } from '@/lib/login-demo'
import { useAuthStore } from '@/stores/auth'
import { safeAdminRedirect } from '@/router/admin-routing'

const { t, locale } = useI18n()
const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const email = ref(DEMO_CREDENTIALS.email)
const password = ref(DEMO_CREDENTIALS.password)
const errorMessage = ref('')
const isSubmitting = ref(false)
const copiedCredential = ref<DemoCredentialKey | null>(null)

function toggleLocale() {
  locale.value = locale.value === 'zh-CN' ? 'en-US' : 'zh-CN'
  localStorage.setItem('locale', locale.value)
}

async function copyCredential(key: DemoCredentialKey) {
  try {
    await navigator.clipboard.writeText(DEMO_CREDENTIALS[key])
    copiedCredential.value = key
    window.setTimeout(() => {
      if (copiedCredential.value === key) copiedCredential.value = null
    }, 1600)
  } catch {
    errorMessage.value = t('auth.copyFailed')
  }
}

async function fillAndSubmit() {
  email.value = DEMO_CREDENTIALS.email
  password.value = DEMO_CREDENTIALS.password
  await submit()
}

async function submit() {
  errorMessage.value = ''
  isSubmitting.value = true
  try {
    await auth.login(email.value, password.value)
    await router.replace(safeAdminRedirect(typeof route.query.redirect === 'string' ? route.query.redirect : undefined))
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? t(errorMessageKey(error.code)) : t('errors.unknown')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <main class="relative flex min-h-svh items-center justify-center bg-muted/40 p-6">
    <Button class="absolute right-6 top-6" variant="ghost" size="sm" :aria-label="t('auth.language')" @click="toggleLocale">
      <Languages data-icon="inline-start" />
      {{ locale === 'zh-CN' ? 'EN' : '中文' }}
    </Button>
    <Card class="w-full max-w-sm">
      <CardHeader>
        <CardTitle>{{ t('auth.loginTitle') }}</CardTitle>
        <CardDescription>{{ t('auth.loginDescription') }}</CardDescription>
      </CardHeader>
      <CardContent>
        <form class="flex flex-col gap-4" @submit.prevent="submit">
          <section class="flex flex-col gap-3 rounded-lg border bg-muted/30 p-4" :aria-label="t('auth.demoTitle')">
            <div class="flex flex-col gap-1">
              <p class="text-sm font-medium">{{ t('auth.demoTitle') }}</p>
              <p class="text-xs text-muted-foreground">{{ t('auth.demoDescription') }}</p>
            </div>
            <div class="flex items-center justify-between gap-3 text-sm">
              <span class="text-muted-foreground">{{ t('auth.demoEmail') }}</span>
              <code class="truncate">{{ DEMO_CREDENTIALS.email }}</code>
              <Button type="button" variant="outline" size="sm" @click="copyCredential('email')">
                <Check v-if="copiedCredential === 'email'" data-icon="inline-start" />
                <Copy v-else data-icon="inline-start" />
                {{ copiedCredential === 'email' ? t('auth.copied') : t('auth.copy') }}
              </Button>
            </div>
            <Separator />
            <div class="flex items-center justify-between gap-3 text-sm">
              <span class="text-muted-foreground">{{ t('auth.demoPassword') }}</span>
              <code class="truncate">{{ DEMO_CREDENTIALS.password }}</code>
              <Button type="button" variant="outline" size="sm" @click="copyCredential('password')">
                <Check v-if="copiedCredential === 'password'" data-icon="inline-start" />
                <Copy v-else data-icon="inline-start" />
                {{ copiedCredential === 'password' ? t('auth.copied') : t('auth.copy') }}
              </Button>
            </div>
            <Button type="button" variant="secondary" @click="fillAndSubmit">
              {{ t('auth.fillAndLogin') }}
            </Button>
          </section>
          <div class="flex flex-col gap-2">
            <Label for="email">{{ t('auth.email') }}</Label>
            <Input id="email" v-model="email" type="email" autocomplete="username" required />
          </div>
          <div class="flex flex-col gap-2">
            <Label for="password">{{ t('auth.password') }}</Label>
            <Input id="password" v-model="password" type="password" autocomplete="current-password" required />
          </div>
          <p v-if="errorMessage" class="text-sm text-destructive" role="alert">{{ errorMessage }}</p>
          <Button type="submit" :disabled="isSubmitting">
            {{ isSubmitting ? t('auth.loggingIn') : t('auth.login') }}
          </Button>
        </form>
      </CardContent>
    </Card>
  </main>
</template>
