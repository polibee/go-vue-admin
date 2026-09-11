export default {
  extensions: {
    moduleTitle: '业务模块', moduleDescription: '业务模块负责产品能力和资源功能。',
    pluginTitle: '平台插件', pluginDescription: '平台插件负责可插拔的基础能力和运行时扩展。',
    status: '状态', enabled: '已启用', disabled: '已停用', version: '版本',
    configure: '配置', details: '详情', view: '查看', enable: '启用', disable: '停用',
    loadFailed: '加载失败', actionFailed: '操作失败', emptyModules: '暂无已注册业务模块。',
    emptyPlugins: '暂无已注册平台插件。', moduleKind: '业务模块', pluginKind: '平台插件',
    configTitle: '扩展配置', configDescription: '业务配置由插件页面负责，平台管理页只维护生命周期和配置状态。', configId: '配置标识', noSchema: '未声明配置模型', configured: '已配置', notConfigured: '未配置', secretFields: '敏感字段', secretHint: '仅提交，不回显明文', configFormHint: '具体配置表单和保存接口由对应业务插件实现。', detailTitle: '扩展详情',
  },
} as const
