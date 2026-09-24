<script setup lang="ts">
import { computed } from 'vue'
import { getResourceDetailExtension } from '@/core/resource/detail-extensions'

const props = defineProps<{
  resource: string
  id: string
  record: Record<string, unknown>
}>()

const extension = computed(() => getResourceDetailExtension(props.resource))
</script>

<template>
  <component v-if="extension" :is="extension.component" :resource="resource" :id="id" :record="record" />
</template>
