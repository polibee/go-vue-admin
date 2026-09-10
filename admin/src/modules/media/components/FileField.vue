<script setup lang="ts">
import { Field, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'

defineProps<{
  id: string
  label: string
  accept?: string
  disabled?: boolean
}>()

const emit = defineEmits<{ change: [file: File | null] }>()

function handleChange(event: Event) {
  const input = event.target as HTMLInputElement
  emit('change', input.files?.[0] ?? null)
}
</script>

<template>
  <Field>
    <FieldLabel :for="id">{{ label }}</FieldLabel>
    <Input :id="id" type="file" :accept="accept" :disabled="disabled" @change="handleChange" />
  </Field>
</template>
