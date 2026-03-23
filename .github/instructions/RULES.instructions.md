---
description: Describe when these instructions should be loaded by the agent based on task context
applyTo: '**'
# applyTo: 'Describe when these instructions should be loaded by the agent based on task context' # when provided, instructions will automatically be added to the request context when the pattern matches an attached file
---

<!-- Tip: Use /create-instructions in chat to generate content with agent assistance -->

Provide project context and coding guidelines that AI should follow when generating code, answering questions, or reviewing changes.

# RULES.md

You are generating code for a Go admin backend based on Gin.

Tech stack:
- Gin
- GORM
- MySQL
- Redis
- JWT
- Casbin
- Viper
- Zap
- Lumberjack
- Cron
- Excelize

Follow this architecture strictly:
router -> middleware -> handler -> service -> repository -> db

Rules:
1. Keep clear boundaries between layers.
2. router only registers routes.
3. middleware only handles cross-cutting concerns.
4. handler only handles request binding, validation, service calls, and unified response.
5. service only contains business logic and transaction orchestration.
6. repository only contains database access logic.
7. do not access DB directly in handler.
8. do not put business logic in repository.
9. do not use `*gin.Context` in service or repository; use `context.Context`.
10. separate `model` and `dto` strictly.
11. never bind request directly to DB model.
12. always use unified response format.
13. do not expose raw internal errors to clients.
14. log important operations and failures with zap.
15. never log passwords, secrets, or full tokens.
16. use viper for config; never hardcode secrets.
17. use bcrypt for password hashing.
18. use Casbin for authorization; do not hardcode role checks everywhere.
19. use pagination for list APIs; enforce max page size.
20. use stable ordering for all list queries.
21. use SQL migrations for production schema changes.
22. put transaction logic in service layer only.
23. validate all inputs, including query params and upload params.
24. do not create random directories or over-engineered abstractions.
25. prefer small, focused functions and files.
26. prefer constructor-based dependency injection.
27. reuse existing structure and naming style.
28. when generating a feature, generate all necessary layers consistently.
29. do not write everything in main.go.
30. optimize for readability, maintainability, and production safety.
31. Always follow this layering:  router -> middleware -> handler -> service -> repository -> db
32. prefer ID-based relations over database foreign keys.
33. do not use foreign key constraints unless explicitly required.
34. enforce relationships at application (service) layer, not database layer.
35. avoid ON DELETE CASCADE in database; handle cascading logic in service layer.
36. design tables to be loosely coupled for better scalability and flexibility.
37. always validate related entity existence in service before insert/update.
38. use indexes on foreign key fields (e.g., user_id, role_id) for performance.

Project structure:
- cmd/server
- configs
- internal/bootstrap
- internal/router
- internal/middleware
- internal/handler
- internal/service
- internal/repository
- internal/model
- internal/dto
- internal/response
- internal/pkg
- internal/cron
- migrations