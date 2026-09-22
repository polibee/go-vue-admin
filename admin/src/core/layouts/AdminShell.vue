<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ClipboardList, Languages, LayoutDashboard, LogOut, Search, ShieldCheck, Unplug } from '@lucide/vue'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Breadcrumb, BreadcrumbItem, BreadcrumbLink, BreadcrumbList, BreadcrumbPage, BreadcrumbSeparator } from '@/components/ui/breadcrumb'
import { Button } from '@/components/ui/button'
import { CommandDialog, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from '@/components/ui/command'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import { Separator } from '@/components/ui/separator'
import { Sidebar, SidebarContent, SidebarFooter, SidebarGroup, SidebarGroupContent, SidebarGroupLabel, SidebarHeader, SidebarInset, SidebarMenu, SidebarMenuButton, SidebarMenuItem, SidebarProvider, SidebarTrigger } from '@/components/ui/sidebar'
import { useAuthStore } from '@/stores/auth'
import { generatedApi, type GlobalSearchResult, type ResourceManifest } from '@/generated/api'
import { dashboardResourceRoute, visibleDashboardResources } from '@/lib/dashboard-resources'
import { groupResourceNavigation } from '@/lib/resource-navigation'

const { t, locale } = useI18n()
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const resourceManifests = ref<ResourceManifest[]>([])
const searchOpen = ref(false)
const searchResults = computed(() => visibleDashboardResources(resourceManifests.value, auth.user?.permissions || []))
const resourceNavigationGroups = computed(() => groupResourceNavigation(resourceManifests.value, auth.user?.permissions || []))
const globalSearchResults = ref<GlobalSearchResult[]>([])
const globalSearchLoading = ref(false)
let globalSearchTimer: ReturnType<typeof setTimeout> | undefined
let globalSearchRequest = 0
const breadcrumbResource = computed(() => {
  const routeResource = typeof route.params.resource === 'string' ? route.params.resource : ''
  if (routeResource) return routeResource
  const routeName = String(route.name || '')
  if (routeName.includes('user')) return 'users'
  if (routeName.includes('role')) return 'roles'
  if (routeName.includes('permission')) return 'permissions'
  const generatedResource = routeName.match(/^(.+)-resource-/)?.[1]
  return generatedResource || ''
})
const breadcrumbLabel = computed(() => {
  if (route.name === 'home') return t('auth.dashboard')
  if (route.name === 'rbac') return t('rbac.title')
  if (route.name === 'audit-logs') return t('auth.auditLogs')
  const resource = resourceManifests.value.find((item) => item.name === breadcrumbResource.value)
  return resource?.label || breadcrumbResource.value || t('auth.dashboard')
})

function toggleLocale() {
  locale.value = locale.value === 'zh-CN' ? 'en-US' : 'zh-CN'
  localStorage.setItem('locale', locale.value)
}

async function logout() {
  await auth.logout()
  await router.replace({ name: 'login' })
}

async function logoutAll() {
  await auth.logoutAll()
  await router.replace({ name: 'login' })
}

function openResource(resource: { name: string; route: string }) {
  searchOpen.value = false
  void router.push(dashboardResourceRoute(resource))
}

function openSearchResult(result: GlobalSearchResult) {
  searchOpen.value = false
  void router.push(result.route)
}

function resourceGroupLabel(name: string) {
  if (name === 'system') return locale.value === 'zh-CN' ? '系统管理' : 'System'
  if (name === 'business') return locale.value === 'zh-CN' ? '业务管理' : 'Business'
  return name
}

function handleSearchInput(value: string) {
  const query = value.trim()
  globalSearchRequest += 1
  const request = globalSearchRequest
  if (globalSearchTimer) clearTimeout(globalSearchTimer)
  if (query.length < 2 || !auth.token) {
    globalSearchResults.value = []
    globalSearchLoading.value = false
    return
  }
  globalSearchLoading.value = true
  globalSearchTimer = setTimeout(async () => {
    try {
      const results = await generatedApi.globalSearch(query, auth.token as string)
      if (request === globalSearchRequest) globalSearchResults.value = results
    } catch {
      if (request === globalSearchRequest) globalSearchResults.value = []
    } finally {
      if (request === globalSearchRequest) globalSearchLoading.value = false
    }
  }, 250)
}

function handleSearchShortcut(event: KeyboardEvent) {
  if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === 'k') {
    event.preventDefault()
    searchOpen.value = true
  }
}

onMounted(async () => {
  window.addEventListener('keydown', handleSearchShortcut)
  if (!auth.token) return
  try { resourceManifests.value = await generatedApi.resourceRegistry(auth.token) } catch { resourceManifests.value = [] }
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleSearchShortcut)
  if (globalSearchTimer) clearTimeout(globalSearchTimer)
})
</script>

<template>
  <SidebarProvider>
    <Sidebar collapsible="icon">
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" :tooltip="t('core.appName')">
              <span class="flex size-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
                <LayoutDashboard />
              </span>
              <span class="truncate font-semibold">{{ t('core.appName') }}</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>{{ t('core.navigation') }}</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton as-child :is-active="true" :tooltip="t('auth.dashboard')">
                  <RouterLink to="/">
                    <LayoutDashboard />
                    <span>{{ t('auth.dashboard') }}</span>
                  </RouterLink>
                </SidebarMenuButton>
              </SidebarMenuItem>
              <SidebarMenuItem v-if="auth.canAny(['admin.users.view', 'admin.roles.manage', 'admin.permissions.manage'])">
                <SidebarMenuButton as-child :is-active="$route.name === 'rbac'" :tooltip="t('rbac.title')">
                  <RouterLink to="/rbac"><ShieldCheck /><span>{{ t('rbac.title') }}</span></RouterLink>
                </SidebarMenuButton>
              </SidebarMenuItem>
              <SidebarMenuItem v-if="auth.can('admin.users.view')">
                <SidebarMenuButton as-child :is-active="$route.name === 'audit-logs'" :tooltip="t('auth.auditLogs')">
                  <RouterLink to="/audit-logs"><ClipboardList /><span>{{ t('auth.auditLogs') }}</span></RouterLink>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
        <SidebarGroup v-for="group in resourceNavigationGroups" :key="group.name">
          <SidebarGroupLabel>{{ resourceGroupLabel(group.name) }}</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem v-for="item in group.items" :key="item.name">
                <SidebarMenuButton as-child :is-active="$route.path.startsWith(dashboardResourceRoute(item))" :tooltip="item.label">
                  <RouterLink :to="dashboardResourceRoute(item)"><LayoutDashboard /><span>{{ item.label }}</span></RouterLink>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <DropdownMenu>
              <DropdownMenuTrigger as-child>
                <SidebarMenuButton size="lg" :tooltip="auth.user?.name">
                  <Avatar class="size-8 rounded-lg">
                    <AvatarFallback class="rounded-lg">{{ auth.user?.name?.slice(0, 1).toUpperCase() }}</AvatarFallback>
                  </Avatar>
                  <span class="truncate">{{ auth.user?.name }}</span>
                </SidebarMenuButton>
              </DropdownMenuTrigger>
              <DropdownMenuContent side="right" align="end" class="w-56">
                <DropdownMenuLabel>{{ auth.user?.email }}</DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem @click="logout">
                  <LogOut />
                  {{ t('auth.logout') }}
                </DropdownMenuItem>
                <DropdownMenuItem @click="logoutAll">
                  <Unplug />
                  {{ t('auth.logoutAll') }}
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>

    <SidebarInset>
      <header class="flex h-16 shrink-0 items-center gap-2 border-b px-4">
        <SidebarTrigger class="-ml-1" />
        <Separator orientation="vertical" class="mr-2 h-4" />
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem class="hidden md:block">
              <BreadcrumbLink href="#">{{ t('core.appName') }}</BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator class="hidden md:block" />
            <BreadcrumbItem>
              <BreadcrumbPage>{{ breadcrumbLabel }}</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
        <Button variant="outline" class="ml-auto hidden h-9 w-56 justify-start gap-2 font-normal text-muted-foreground sm:flex" @click="searchOpen = true">
          <Search data-icon="inline-start" />
          <span>{{ t('core.searchResources') }}</span>
          <kbd class="ml-auto rounded border bg-muted px-1.5 py-0.5 text-[10px]">⌘K</kbd>
        </Button>
        <Button variant="ghost" size="icon" class="ml-auto sm:hidden" :aria-label="t('core.searchResources')" @click="searchOpen = true">
          <Search />
        </Button>
        <div class="ml-auto md:hidden">
          <Button variant="ghost" size="icon" @click="logout">
            <LogOut />
            <span class="sr-only">{{ t('auth.logout') }}</span>
          </Button>
        </div>
        <Button variant="ghost" size="sm" @click="toggleLocale"><Languages />{{ locale === 'zh-CN' ? 'EN' : '中文' }}</Button>
      </header>
      <CommandDialog v-model:open="searchOpen" :title="t('core.searchResources')" :description="t('core.searchResources')">
        <CommandInput :placeholder="t('core.searchResources')" @update:model-value="handleSearchInput" />
        <CommandList>
          <CommandEmpty v-if="!globalSearchLoading">{{ t('core.noSearchResults') }}</CommandEmpty>
          <CommandGroup v-if="globalSearchResults.length" :heading="t('core.searchData')">
            <CommandItem v-for="result in globalSearchResults" :key="result.resource + ':' + result.id" :value="result.title + ' ' + (result.subtitle || '')" @select="openSearchResult(result)">
              <span>{{ result.title }}</span>
              <span v-if="result.subtitle" class="truncate text-xs text-muted-foreground">{{ result.subtitle }}</span>
              <span class="ml-auto text-xs text-muted-foreground">{{ result.label }}</span>
            </CommandItem>
          </CommandGroup>
          <CommandGroup :heading="t('core.searchResources')">
            <CommandItem v-for="resource in searchResults" :key="resource.name" :value="resource.name + ' ' + resource.label" @select="openResource(resource)">
              <span>{{ resource.label }}</span>
              <span class="ml-auto text-xs text-muted-foreground">{{ resource.name }}</span>
            </CommandItem>
          </CommandGroup>
        </CommandList>
      </CommandDialog>
      <div class="flex flex-1 flex-col gap-4 p-4 pt-6">
        <RouterView />
      </div>
    </SidebarInset>
  </SidebarProvider>
</template>
