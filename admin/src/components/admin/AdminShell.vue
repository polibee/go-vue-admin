<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink, RouterView, useRouter } from 'vue-router'
import { ChevronsUpDown, ChevronRight } from '@lucide/vue'
import {
  Sidebar, SidebarContent, SidebarFooter, SidebarGroup, SidebarGroupContent,
  SidebarGroupLabel, SidebarHeader, SidebarInset, SidebarMenu, SidebarMenuButton,
  SidebarMenuItem, SidebarMenuSub, SidebarMenuSubButton, SidebarMenuSubItem,
  SidebarProvider, SidebarRail, SidebarTrigger,
} from '@/components/ui/sidebar'
import {
  Breadcrumb, BreadcrumbItem, BreadcrumbLink, BreadcrumbList,
  BreadcrumbPage, BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { Button } from '@/components/ui/button'
import { useAuth } from '@/core/auth'
import { navigationRegistry, NavigationService } from '@/core/navigation'
import { applyTheme, initializeTheme, type ThemeMode, type ThemePalette } from '@/core/theme'
import { setLocale, type SupportedLocale } from '@/core/i18n'

const { t, locale } = useI18n()
const auth = useAuth()
const router = useRouter()
const navigationSections = computed(() => navigationRegistry.sections(auth.user?.permissions ?? []))
const collapsedGroups = ref(new Set<string>())
const theme = ref(initializeTheme())
const navigationLabels: Record<string, string> = {
  dashboard: 'app.dashboard', settings: 'app.settings', media: 'app.media', audit: 'app.audit',
  'api-docs': 'app.apiDocs', users: 'app.users', roles: 'app.roles', permissions: 'app.permissions',
  modules: 'app.modules', plugins: 'app.plugins', about: 'app.about',
}

function navigationLabel(id: string, fallback: string): string {
  const key = navigationLabels[id]
  return key ? t(key) : fallback
}
function isGroupExpanded(id: string): boolean { return !collapsedGroups.value.has(id) }
function toggleGroup(id: string): void {
  const next = new Set(collapsedGroups.value)
  if (next.has(id)) next.delete(id)
  else next.add(id)
  collapsedGroups.value = next
}
function cycleLocale(): void {
  const next = locale.value === 'zh-CN' ? 'en' : 'zh-CN'
  setLocale(next as SupportedLocale)
}
function cycleMode(): void {
  const modes: ThemeMode[] = ['system', 'light', 'dark']
  const next = modes[(modes.indexOf(theme.value.mode) + 1) % modes.length]
  theme.value = { ...theme.value, mode: next }
  applyTheme(theme.value)
}
function cyclePalette(): void {
  const palettes: ThemePalette[] = ['shadcn', 'semi', 'wechat']
  const next = palettes[(palettes.indexOf(theme.value.palette) + 1) % palettes.length]
  theme.value = { ...theme.value, palette: next }
  applyTheme(theme.value)
}
function modeLabel(): string {
  return theme.value.mode === 'light' ? t('app.light') : theme.value.mode === 'dark' ? t('app.dark') : t('app.system')
}
function paletteLabel(): string {
  return theme.value.palette === 'semi' ? t('app.semi') : theme.value.palette === 'wechat' ? t('app.wechat') : t('app.shadcn')
}

onMounted(async () => {
  try { navigationRegistry.registerMany(await new NavigationService().list()) } catch { /* static metadata remains usable */ }
})
async function handleLogout() {
  await auth.logout()
  await router.replace('/login')
}
</script>

<template>
  <SidebarProvider storage-key="go-vue-admin-sidebar">
    <Sidebar>
      <SidebarHeader>
        <SidebarMenu><SidebarMenuItem><SidebarMenuButton size="lg">
          <div class="flex size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground"><ChevronsUpDown /></div>
          <div class="grid flex-1 text-left text-sm leading-tight"><span class="truncate font-semibold">{{ t('app.title') }}</span><span class="truncate text-xs">{{ t('app.platform') }}</span></div>
        </SidebarMenuButton></SidebarMenuItem></SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup><SidebarGroupLabel>{{ t('app.platform') }}</SidebarGroupLabel><SidebarGroupContent><SidebarMenu>
          <SidebarMenuItem v-for="section in navigationSections" :key="section.kind === 'item' ? section.item.id : section.group.id">
            <SidebarMenuButton v-if="section.kind === 'item'" as-child><RouterLink :to="section.item.route"><component :is="section.item.icon" /><span>{{ navigationLabel(section.item.id, section.item.label) }}</span></RouterLink></SidebarMenuButton>
            <SidebarMenuButton v-else type="button" @click="toggleGroup(section.group.id)" :aria-expanded="isGroupExpanded(section.group.id)"><component :is="section.group.icon" /><span>{{ section.group.label }}</span><ChevronRight class="ml-auto transition-transform" :class="{ 'rotate-90': isGroupExpanded(section.group.id) }" /></SidebarMenuButton>
            <SidebarMenuSub v-if="section.kind === 'group' && isGroupExpanded(section.group.id)"><SidebarMenuSubItem v-for="item in section.group.items" :key="item.id"><SidebarMenuSubButton as-child><RouterLink :to="item.route"><component :is="item.icon" /><span>{{ navigationLabel(item.id, item.label) }}</span></RouterLink></SidebarMenuSubButton></SidebarMenuSubItem></SidebarMenuSub>
          </SidebarMenuItem>
        </SidebarMenu></SidebarGroupContent></SidebarGroup>
      </SidebarContent>
      <SidebarFooter><SidebarMenu><SidebarMenuItem class="flex flex-col gap-2">
        <Button variant="ghost" size="sm" class="justify-start" aria-label="language" @click="cycleLocale">{{ t('app.language') }} · {{ locale === 'en' ? 'English' : '简体中文' }}</Button>
        <Button variant="ghost" size="sm" class="justify-start" aria-label="theme-mode" @click="cycleMode">{{ t('app.theme') }} · {{ modeLabel() }}</Button>
        <Button variant="ghost" size="sm" class="justify-start" aria-label="theme-palette" @click="cyclePalette">{{ t('app.theme') }} · {{ paletteLabel() }}</Button>
      </SidebarMenuItem></SidebarMenu></SidebarFooter>
      <SidebarRail />
    </Sidebar>
    <SidebarInset>
      <header class="flex h-16 shrink-0 items-center gap-2 border-b px-4"><SidebarTrigger class="-ml-1" /><Breadcrumb><BreadcrumbList><BreadcrumbItem><BreadcrumbLink as-child><RouterLink to="/admin/dashboard">{{ t('app.platform') }}</RouterLink></BreadcrumbLink></BreadcrumbItem><BreadcrumbSeparator /><BreadcrumbItem><BreadcrumbPage>{{ t('app.admin') }}</BreadcrumbPage></BreadcrumbItem></BreadcrumbList></Breadcrumb><span v-if="auth.user" class="hidden text-sm text-muted-foreground md:inline">{{ auth.user.email }}</span><Button class="ml-auto" size="sm" variant="ghost" @click="handleLogout">{{ t('app.logout') }}</Button></header>
      <main class="flex flex-1 flex-col gap-4 p-4 md:p-6"><RouterView /></main>
    </SidebarInset>
  </SidebarProvider>
</template>
