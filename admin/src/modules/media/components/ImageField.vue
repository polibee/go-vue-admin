<script setup lang="ts">
import { ImageIcon } from '@lucide/vue'

defineProps<{
  src?: string
  alt: string
}>()
</script>

<template>
  <div class="flex aspect-[4/3] items-center justify-center overflow-hidden rounded-md bg-muted">
    <img v-if="src" :src="src" :alt="alt" class="size-full object-cover" />
    <ImageIcon v-else class="size-10 text-muted-foreground" aria-hidden="true" />
  </div>
</template>
