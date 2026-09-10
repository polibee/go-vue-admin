<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { apiClient } from '@/core/api/client'
import { createGeneratedApiClient, type ExtensionResource } from '@/generated/api'
const client = createGeneratedApiClient(apiClient)
const extensions = ref<ExtensionResource[]>([])
const error = ref('')
const busy = ref('')
async function load() { error.value = ''; try { extensions.value = (await client.listExtensions()).data } catch (cause) { error.value = cause instanceof Error ? cause.message : '扩展列表加载失败' } }
async function setState(item: ExtensionResource, state: 'enabled' | 'disabled') { busy.value = item.id; error.value = ''; try { extensions.value = (await client.setPluginState(item.id, state)).data } catch (cause) { error.value = cause instanceof Error ? cause.message : '插件状态更新失败' } finally { busy.value = '' } }
onMounted(load)
</script>
<template><section class="flex flex-col gap-6"><div><h1 class="text-2xl font-semibold tracking-tight">模块与插件</h1><p class="text-sm text-muted-foreground">查看已注册模块，并管理内置插件生命周期。</p></div><Alert v-if="error" variant="destructive"><AlertTitle>操作失败</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert><div class="grid gap-4 md:grid-cols-2"><Card v-for="item in extensions" :key="item.id"><CardHeader><CardTitle>{{ item.name }}</CardTitle><CardDescription>{{ item.kind }} · {{ item.id }}</CardDescription></CardHeader><CardContent>状态：{{ item.state === 'enabled' ? '已启用' : '已停用' }}</CardContent><CardFooter class="gap-2"><Button v-if="item.kind === 'plugin' && item.state === 'disabled'" :disabled="busy === item.id" @click="setState(item, 'enabled')">启用</Button><Button v-if="item.kind === 'plugin' && item.state === 'enabled'" variant="outline" :disabled="busy === item.id" @click="setState(item, 'disabled')">停用</Button><Button variant="outline" as-child><RouterLink :to="`/admin/extensions/${item.id}`">查看</RouterLink></Button></CardFooter></Card></div></section></template>
