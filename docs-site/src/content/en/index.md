# Go Vue Admin Developer Docs

Go Vue Admin is a modular administration platform built with Go, Goravel, Vue 3, shadcn-vue, GORM and OpenAPI contracts.

This site is for developers building modules, plugins, resources and external integrations. It is isolated from the admin application and backend runtime.

## Start here

- Follow Getting started to run the backend and admin locally.
- Read Architecture and boundaries before adding a module.
- Build business modules with a manifest, runtime registration and contract-driven APIs.
- Build platform plugins with lifecycle, dependency and configuration boundaries.

## Stable boundaries

OpenAPI documents and schemas are machine-readable contracts. Do not edit generated clients by hand. Business pages should use the generated client and ResourceDataProvider.

Read Production readiness before exposing an installation to real users.
