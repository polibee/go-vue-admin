<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { ChevronsUpDown } from '@lucide/vue'
import {
  Sidebar, SidebarContent, SidebarFooter, SidebarGroup, SidebarGroupContent,
  SidebarGroupLabel, SidebarHeader, SidebarInset, SidebarMenu, SidebarMenuButton,
  SidebarMenuItem, SidebarProvider, SidebarRail, SidebarTrigger,
} from '@/components/ui/sidebar'
import {
  Breadcrumb, BreadcrumbItem, BreadcrumbLink, BreadcrumbList,
  BreadcrumbPage, BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { Button } from '@/components/ui/button'
import { useAuth } from '@/core/auth'
import { navigationItems } from '@/core/navigation'
import { useRouter } from 'vue-router'

const themeLabel = '主题：系统'
const auth = useAuth()
const router = useRouter()

async function handleLogout() {
  await auth.logout()
  await router.replace('/login')
}
</script>

<template>
  <SidebarProvider storage-key="go-vue-admin-sidebar">
    <Sidebar>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg">
              <div class="flex size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                <ChevronsUpDown />
              </div>
              <div class="grid flex-1 text-left text-sm leading-tight">
                <span class="truncate font-semibold">Go Vue Admin</span>
                <span class="truncate text-xs">Platform</span>
              </div>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>平台</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem v-for="item in navigationItems" :key="item.id">
                <SidebarMenuButton as-child>
                  <RouterLink :to="item.route">
                    <component :is="item.icon" />
                    <span>{{ item.label }}</span>
                  </RouterLink>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem><SidebarMenuButton>{{ themeLabel }}</SidebarMenuButton></SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
    <SidebarInset>
      <header class="flex h-16 shrink-0 items-center gap-2 border-b px-4">
        <SidebarTrigger class="-ml-1" />
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem>
              <BreadcrumbLink as-child><RouterLink to="/admin/dashboard">平台</RouterLink></BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator />
            <BreadcrumbItem><BreadcrumbPage>管理后台</BreadcrumbPage></BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
        <span v-if="auth.user" class="hidden text-sm text-muted-foreground md:inline">{{ auth.user.email }}</span>
        <Button class="ml-auto" size="sm" variant="ghost" @click="handleLogout">退出</Button>
      </header>
      <main class="flex flex-1 flex-col gap-4 p-4 md:p-6"><RouterView /></main>
    </SidebarInset>
  </SidebarProvider>
</template>
