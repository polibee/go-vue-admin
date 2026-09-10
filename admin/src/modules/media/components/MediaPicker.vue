<script setup lang="ts">
import { ref } from 'vue'
import { Upload } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { FileField } from './index'

defineProps<{ busy?: boolean }>()
const emit = defineEmits<{ upload: [file: File] }>()
const selected = ref<File | null>(null)

function submit() {
  if (selected.value) emit('upload', selected.value)
}
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>上传媒体</CardTitle>
      <CardDescription>支持图片、文档和其他平台文件，单个文件不超过 10 MB。</CardDescription>
    </CardHeader>
    <CardContent>
      <FileField id="media-file" label="选择文件" accept="image/*,application/pdf,text/plain" :disabled="busy" @change="selected = $event" />
    </CardContent>
    <CardFooter class="justify-end">
      <Button :disabled="!selected || busy" @click="submit">
        <Upload class="mr-2 size-4" />
        {{ busy ? '上传中…' : '上传文件' }}
      </Button>
    </CardFooter>
  </Card>
</template>
