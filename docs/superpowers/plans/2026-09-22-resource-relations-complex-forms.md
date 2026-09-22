# Resource 关系与复杂表单扩展实施计划

> 执行约束：Native 当前工作区；不新增依赖、不执行迁移、不修改 Laragon PostgreSQL/Redis 配置、不推送 GitHub。每个任务完成后先验证并保留本地提交，阶段全部验收后再决定是否推送。

## 1. Core Resource 契约与校验

**修改范围**：`backend/app/core/resource/`、对应单元测试、OpenAPI schema。

- 增加 `Relation`、`FormGroup`、`DetailSection` 和基础联动元数据；
- 为关系 kind、字段引用、布局列数、重复字段和目标资源建立 Registry 校验；
- 保留未声明关系/分组的现有默认行为；
- 增加契约单元测试和 OpenAPI schema 测试。

**验收**：非法关系、未知字段、重复字段、非法布局值被稳定拒绝；现有 Manifest 全部通过。

## 2. 后端关系选项与写入边界

**修改范围**：`backend/app/core/admin/controllers/`、`backend/app/services/` 对应领域目录、路由和测试。

- 增加统一关系选项接口；
- 通过服务端 Registry 解析目标资源和字段，不允许客户端传表名/字段名；
- 应用目标资源权限、数据范围、字段 readable 策略；
- 在创建/更新前验证 belongsTo ID，非法值返回稳定错误；
- 保持 Service 分目录约束，不把关系逻辑散落到 `app/core` 或资源目录。

**验收**：授权/越权/不存在目标、own 数据范围和不可读 label 字段均有边界测试；现有 CRUD 无回归。

## 3. 前端共享表单与详情扩展

**修改范围**：`admin/src/components/resource/`、`admin/src/core/resource/`、`admin/src/lib/`、前端测试。

- `ResourceFormView` 支持分组、受控列布局和 belongsTo Select；
- 关系选项按需加载，加载失败时禁止提交；
- `ResourceDetailPage` 支持 hasMany 只读摘要和详情区块；
- 基础联动只支持声明式固定条件，不执行任意脚本；
- 保持共享组件在 `admin/src/components`，资源业务元数据在 `admin/src/modules/<name>`。

**验收**：无关系资源、users 现有表单、关系表单和响应式小屏幕均通过类型/组件测试。

## 4. 生成器和示例资源

**修改范围**：`backend/app/generator/`、生成器 Golden、示例资源、README、生成客户端。

- 支持关系、分组、详情区块输入和确定性输出；
- 生成关系 API 客户端契约和 README 审阅清单；
- 冲突预检覆盖未知资源、字段重复和半套写入；
- 不生成业务 Handler，不自动执行迁移，不覆盖人工注册点。

**验收**：真实示例资源生成后人工审阅迁移，再由人工决定是否执行；生成产物能进入菜单、列表、表单、详情、关系选项和权限校验。

## 5. 集成验证与里程碑

执行：

```text
go test ./... -count=1
go build -o storage/codex-backend-relations-check.exe .
npx vue-tsc -b
npm run build
node --test tests/resource-actions.test.ts
git diff --check
```

删除临时构建文件，检查 PostgreSQL/Redis 仅用于需要的联调，不执行迁移。完成自审后提交：

`feat: complete resource relations and forms milestone`

本计划完成后再进入真实资源端到端验收，不在本阶段扩展 many-to-many、导入、审计脱敏或复杂业务流程。
