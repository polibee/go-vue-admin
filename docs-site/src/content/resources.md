# 资源与代码生成

## Resource Manifest

Resource Manifest 是数据库结构和后台表单之间的稳定中间结构。数据库反向读取只能提供默认推断，后续生成应只读取 Manifest。

最小资源定义包括：

    id: products
    table: products
    label: 产品
    primary_key: id
    permissions:
      list: products.view
      create: products.create

## 单表资源生成

使用 Go 生成器发现数据库表并落成 Manifest：

    go -C backend run ./cmd/admin-gen resource --module catalog --table products

生成的后端 CRUD 使用通用 GORM Repository，不包含 MySQL 专用实现。生成器必须对字段、排序和写入列使用白名单，不能把用户输入直接拼接进 SQL。

## 生成后的职责

生成器负责标准 CRUD 骨架。复杂查询、外键选择器、枚举策略、软删除、审计和业务动作属于模块扩展代码。
