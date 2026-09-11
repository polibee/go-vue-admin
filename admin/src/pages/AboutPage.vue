<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { ExternalLink, BookOpen, ShieldCheck } from '@lucide/vue'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

const { t } = useI18n()
const repositoryUrl = 'https://github.com/polibee/go-vue-admin'
const docsUrl = 'https://polibee.github.io/go-vue-admin/'
</script>

<template>
  <section class="flex flex-col gap-6">
    <div>
      <h1 class="text-2xl font-semibold tracking-tight">{{ t('about.title') }}</h1>
      <p class="text-sm text-muted-foreground">{{ t('about.description') }}</p>
    </div>
    <div class="grid gap-4 lg:grid-cols-3">
      <Card>
        <CardHeader>
          <ShieldCheck class="size-5 text-primary" />
          <CardTitle>{{ t('about.foundation') }}</CardTitle>
          <CardDescription>{{ t('about.foundationDescription') }}</CardDescription>
        </CardHeader>
        <CardContent><p class="text-sm text-muted-foreground">{{ t('about.status') }}</p></CardContent>
      </Card>
      <Card>
        <CardHeader>
          <BookOpen class="size-5 text-primary" />
          <CardTitle>{{ t('about.docs') }}</CardTitle>
          <CardDescription>{{ t('about.docsDescription') }}</CardDescription>
        </CardHeader>
        <CardContent>
          <Button as-child variant="outline">
            <a :href="docsUrl" target="_blank" rel="noreferrer">
              {{ t('about.docs') }} <ExternalLink data-icon="inline-end" />
            </a>
          </Button>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>{{ t('about.repository') }}</CardTitle>
          <CardDescription>{{ t('about.version') }}</CardDescription>
        </CardHeader>
        <CardContent>
          <Button as-child variant="outline">
            <a :href="repositoryUrl" target="_blank" rel="noreferrer">
              polibee/go-vue-admin <ExternalLink data-icon="inline-end" />
            </a>
          </Button>
        </CardContent>
      </Card>
    </div>
  </section>
</template>
