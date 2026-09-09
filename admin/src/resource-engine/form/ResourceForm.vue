<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'

import { Button } from '@/components/ui/button'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { Textarea } from '@/components/ui/textarea'

import type { ResourceDefinition, ResourceField } from '../core/ResourceDefinition'
import { createResourceSchema } from './fieldSchema'

const props = withDefaults(defineProps<{
  definition: ResourceDefinition<object>
  initialValues?: Record<string, unknown>
  submitLabel?: string
}>(), {
  initialValues: () => ({}),
  submitLabel: '保存',
})

const emit = defineEmits<{
  submit: [values: Record<string, unknown>]
}>()

const fields = () => props.definition.fields ?? []
const validationSchema = toTypedSchema(createResourceSchema(fields()))

function inputType(field: ResourceField<object>): string {
  if (field.type === 'number' || field.type === 'date' || field.type === 'datetime') return field.type === 'datetime' ? 'datetime-local' : field.type
  return 'text'
}
</script>

<template>
  <Form v-slot="{ handleSubmit }" :validation-schema="validationSchema" :initial-values="initialValues" as="div">
    <form class="flex flex-col gap-5" @submit="handleSubmit($event, (values) => emit('submit', values as Record<string, unknown>))">
      <FormField v-for="field in fields()" :key="String(field.name)" v-slot="{ componentField }" :name="String(field.name)">
        <FormItem>
          <FormLabel>{{ field.label }}</FormLabel>
          <FormControl>
            <Textarea v-if="field.type === 'textarea'" :placeholder="field.placeholder" v-bind="componentField" />
            <Select v-else-if="field.type === 'select'" v-bind="componentField">
              <SelectTrigger><SelectValue :placeholder="field.placeholder ?? `请选择${field.label}`" /></SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectItem v-for="option in field.options ?? []" :key="option.value" :value="option.value">{{ option.label }}</SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
            <Switch v-else-if="field.type === 'switch'" v-bind="componentField" />
            <Input v-else :type="inputType(field)" :placeholder="field.placeholder" v-bind="componentField" />
          </FormControl>
          <FormMessage />
        </FormItem>
      </FormField>
      <Button type="submit">{{ submitLabel }}</Button>
    </form>
  </Form>
</template>
