---
title: "Где используются миграции"
---

Миграции БД используются ТОЛЬКО в двух сервисах:
- `order/migrations` — embed в `order/internal/app/di.go` (`//go:embed ../migrations/*.sql`), env `MIGRATION_DIRECTORY=./order/migrations` (deploy/compose/order/.env)
- `iam/migrations` — embed в `iam/internal/app/di.go` (`//go:embed migrations/*.sql`), env `IAM_MIGRATION_DIRECTORY=./iam/migrations` (deploy/env/.env)

В payment, inventory, assembly, notification, platform миграций нет.
