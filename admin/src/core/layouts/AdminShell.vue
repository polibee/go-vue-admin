<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ClipboardList, Languages, LayoutDashboard, LogOut, ShieldCheck, Unplug } from '@lucide/vue'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Breadcrumb, BreadcrumbItem, BreadcrumbLink, BreadcrumbList, BreadcrumbPage, BreadcrumbSeparator } from '@/components/ui/breadcrumb'
import { Button } from '@/components/ui/button'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import { Separator } from '@/components/ui/separator'
import { Sidebar, SidebarContent, SidebarFooter, SidebarGroup, SidebarGroupContent, SidebarGroupLabel, SidebarHeader, SidebarInset, SidebarMenu, SidebarMenuButton, SidebarMenuItem, SidebarProvider, SidebarTrigger } from '@/components/ui/sidebar'
import { useAuthStore } from '@/stores/auth'
import { generatedResourceDefinitions } from '@/core/resource/generated'

const { t, locale } = useI18n()
const router = useRouter()
const auth = useAuthStore()

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
              <template v-for="item in generatedResourceDefinitions" :key="item.name">
                <SidebarMenuItem v-if="auth.can(item.permission)">
                  <SidebarMenuButton as-child :is-active="$route.path.startsWith(item.route)" :tooltip="item.label">
                    <RouterLink :to="item.route"><LayoutDashboard /><span>{{ item.label }}</span></RouterLink>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              </template>
              <SidebarMenuItem v-if="auth.can('admin.users.view')">
                <SidebarMenuButton as-child :is-active="$route.name === 'resource-list'" :tooltip="t('resource.title')">
                  <RouterLink to="/users"><LayoutDashboard /><span>{{ t('resource.title') }}</span></RouterLink>
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
              <BreadcrumbPage>{{ t('auth.dashboard') }}</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
        <div class="ml-auto md:hidden">
          <Button variant="ghost" size="icon" @click="logout">
            <LogOut />
            <span class="sr-only">{{ t('auth.logout') }}</span>
          </Button>
        </div>
        <Button variant="ghost" size="sm" @click="toggleLocale"><Languages />{{ locale === 'zh-CN' ? 'EN' : '中文' }}</Button>
      </header>
      <div class="flex flex-1 flex-col gap-4 p-4 pt-6">
        <RouterView />
      </div>
    </SidebarInset>
  </SidebarProvider>
</template>
