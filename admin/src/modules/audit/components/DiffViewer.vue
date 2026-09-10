<script setup lang="ts">
defineProps<{
  before: unknown
  after: unknown
}>()

function formatSnapshot(value: unknown) {
  if (value === null || value === undefined) return '—'
  return JSON.stringify(value, null, 2)
}
</script>

<template>
  <div class="grid gap-4 lg:grid-cols-2">
    <section class="overflow-hidden rounded-lg border border-destructive/20 bg-destructive/5">
      <header class="border-b border-destructive/20 px-4 py-3 text-sm font-medium text-destructive">变更前</header>
      <pre class="max-h-80 overflow-auto p-4 text-xs leading-6 text-foreground">{{ formatSnapshot(before) }}</pre>
    </section>
    <section class="overflow-hidden rounded-lg border border-emerald-500/20 bg-emerald-500/5">
      <header class="border-b border-emerald-500/20 px-4 py-3 text-sm font-medium text-emerald-700 dark:text-emerald-300">变更后</header>
      <pre class="max-h-80 overflow-auto p-4 text-xs leading-6 text-foreground">{{ formatSnapshot(after) }}</pre>
    </section>
  </div>
</template>
