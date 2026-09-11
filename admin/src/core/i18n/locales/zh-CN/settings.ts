export default {
  settings: {
    title: '设置', description: '按命名空间管理平台配置，值类型由契约校验。',
    unavailable: '设置暂时不可用', createTitle: '新增设置', createDescription: '使用 namespace + key 定位设置，修改不会改变数据库结构。',
    namespace: '命名空间', key: '键名', type: '值类型', chooseType: '选择值类型', value: '值', valuePlaceholder: '输入设置值',
    text: '文本', boolean: '布尔值', integer: '整数', number: '数字', json: 'JSON', note: '说明',
    notePlaceholder: '帮助管理员理解这个设置', noteHint: '建议使用稳定、可读的键名，敏感凭据不应直接写入此处。',
    create: '创建设置', configured: '已配置设置', configuredDescription: '当前显示 Memory 或 GORM Provider 返回的设置。',
    loading: '正在加载设置…', empty: '暂无设置', emptyDescription: '请先创建一个命名空间设置。',
    missingNote: '未填写说明', saving: '保存中…', save: '保存', loadFailed: '设置加载失败，请稍后重试。',
    saveFailed: '设置保存失败，请检查输入。', createFailed: '设置创建失败，请检查输入。', required: '请填写命名空间和键名。',
    saved: '已保存 {id}', created: '设置已创建',
  },
}
