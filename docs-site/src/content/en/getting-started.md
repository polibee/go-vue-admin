# Getting started

## Requirements

- Node.js 24 or newer
- pnpm 10
- Go 1.27
- MySQL 8 or PostgreSQL

## Local development

Install admin dependencies from the repository root:

    pnpm install

Provide a local application key and development account:

    APP_ENV=local
    APP_KEY=replace-with-a-32-character-dev-key
    AUTH_BOOTSTRAP_EMAIL=admin@example.com
    AUTH_BOOTSTRAP_PASSWORD=test-only-password
    RESOURCE_PROVIDER=memory

Start the backend and admin:

    go -C backend run .
    pnpm --dir admin run dev -- --host 0.0.0.0

Open http://127.0.0.1:5173/login. The admin uses the HTTP Resource Provider by default, so CRUD data survives browser reloads through the backend.

## Language and themes

The admin defaults to Simplified Chinese and can switch to English from the sidebar. Language preferences are stored locally; feature modules should use vue-i18n message keys instead of language branches in components.

Theme mode supports light, dark and system. Palette presets include shadcn default, Semi Design and WeChat. Presets only override semantic CSS variables and do not replace shadcn-vue components.

## Production warning

Never commit demo credentials, application keys or database passwords. Disable the bootstrap credential endpoint in production.
