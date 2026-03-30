---
description: Core engineering rules for sys-admin-serve backend tasks (Go/Gin). Apply for API, service, repository, model, and migration work.
applyTo: 'cmd/server/**,configs/**,internal/**,migrations/**'
---

# Sys Admin Serve Backend Rules

## Context

Go admin backend tech stack:
Gin, GORM, MySQL, Redis, JWT, Casbin, Viper, Zap, Lumberjack, Cron, Excelize.

Project layout:
`cmd/server`, `configs`, `internal/bootstrap`, `internal/router`, `internal/middleware`, `internal/handler`, `internal/service`, `internal/repository`, `internal/model`, `internal/dto`, `internal/response`, `internal/pkg`, `internal/cron`, `migrations`.

## Architecture (strict)

Always follow:
`router -> middleware -> handler -> service -> repository -> db`

- `router`: route registration only.
- `middleware`: cross-cutting concerns only.
- `handler`: request binding, validation, service call, unified response.
- `service`: business logic, orchestration, transactions.
- `repository`: data access only.

## Engineering Rules

1. Keep layer boundaries clear; do not bypass layers.
2. Never access DB directly in handler.
3. Do not place business logic in repository.
4. Do not use `*gin.Context` in service/repository; use `context.Context`.
5. Keep `model` and `dto` strictly separated.
6. Never bind request payload directly into DB model.
7. Always return the unified response format in handlers.
8. Do not expose raw internal errors to clients.
9. Log important operations and failures with `zap`.
10. Never log passwords, secrets, or full tokens.
11. Use `viper` for config; never hardcode secrets.
12. Use `bcrypt` for password hashing.
13. Use `Casbin` for authorization; avoid scattered hardcoded role checks.
14. List APIs must support pagination and enforce max page size.
15. List queries must have stable ordering.
16. Use SQL migrations for production schema changes.
17. Put transaction boundaries in service layer only.
18. Validate all inputs (body/query/path/upload params).
19. Prefer constructor-based dependency injection.
20. Reuse existing naming/style and keep functions focused.
21. Prefer ID-based relations; avoid DB foreign key constraints unless explicitly required.
22. Handle cascading behavior in service layer instead of `ON DELETE CASCADE`.
23. Validate related entity existence in service before insert/update.
24. Add indexes for high-frequency relation fields (e.g., `user_id`, `role_id`).
