# 开发平台插件

## 插件边界

v1 插件是静态编译、受信任的 builtin 插件。当前不支持网络下载、任意外部 JavaScript 执行、Marketplace 或签名验证。

插件通过 defineAdminPlugin 声明 Manifest，并在 setup 中注册菜单、路由和资源。

    const plugin = defineAdminPlugin({
      manifest: {
        id: 'example-plugin',
        name: 'Example Plugin',
        version: '1.0.0',
        runtime: 'builtin',
        uiCompatibility: 'shadcn-vue'
      },
      setup(context) {
        context.registerMenu({
          id: 'example',
          label: 'Example',
          route: '/admin/example',
          permission: 'example.view'
        })
      }
    })

## 生命周期

启用插件时先由后端校验权限、状态和依赖，然后前端同步本地 PluginRegistry 并挂载 RuntimeRegistrationBridge。停用插件时清理菜单、动态路由和插件资源。

## 配置和业务页面

/admin/plugins 是平台控制面，负责启停、版本、依赖、权限和配置状态。支付密钥、证书和回调参数应放在插件自己的配置页面，由后端加密保存，页面不回显完整密钥。

支付订单、渠道和退款等业务功能应使用独立业务菜单；不要把所有插件页面堆在插件管理页中。
