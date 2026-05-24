# API OBM — Ontologia Brasileira de Medicamentos

## Vision

REST API and CLI tooling for Brazil's Medicines Ontology (OBM) — an adaptation of the NHS dm+d model. Provides structured access to virtual medicinal products (VMPs), actual medicinal products (AMPs), suppliers, and 32 domain types, with full-text search via Meilisearch.

## Goals

1. **Reliable data access** — Stable, paginated REST API for VMPs, AMPs, VTMs, VMPPs, AMPPs, DCBs, Ingredients, Suppliers, and Domains
2. **Fast search** — Meilisearch-powered global search with filters on name, code, manufacturer, description, and active status
3. **Simple imports** — CLI tool to ingest OBM data dumps (ZIP/SQL/MySQL), convert MySQL→PostgreSQL, import, validate, and reindex — one command
4. **Self-documenting** — Swagger UI, Postman collection, comprehensive PT-BR README

## Non-Goals

- Scheduled/automatic imports (on-demand only)
- Incremental/delta updates (always full replace)
- Expanding Meilisearch indexes beyond VMP/AMP/Supplier for now
- Downloading from portal URL (access is restricted)
- CSV/XML/JSON input formats

## Stack

| Layer | Tech |
|-------|------|
| Language | Go 1.25 |
| HTTP | Gin |
| DB | PostgreSQL 16 (pgx/v5) |
| Search | Meilisearch v1.8 |
| Auth | JWT HS256 (golang-jwt/jwt/v5) |
| Docs | swaggo/swag |
| Infra | Docker Compose |
| Env | godotenv |

## Architecture

Clean architecture — domain entities at core, repositories as interfaces, HTTP handlers at edge.

```
cmd/
  api/          → REST API entrypoint
  import/       → Import CLI entrypoint (upcoming)
internal/
  domain/       → Entities + repository interfaces
  usecase/      → Business logic (search, reindex, CRUD)
  infrastructure/
    config/     → Env-based config loading
    persistence/
      postgres/ → pgxpool repositories
      meilisearch/ → Search client + indexer
  interface/
    http/       → Handlers, DTOs, middleware, router
  importer/     → Import pipeline (upcoming)
scripts/        → Standalone utilities (seed, convert_sql, gen_postman)
```

## Data Model

67 tables from OBM MySQL dump:
- 7 `tb_*` — main entities (vmp, amp, vtm, vmpp, ampp, ingredient, supplier)
- 31+ `td_*` — domain tables (form, route, atc class, control category, etc.)
- 25+ `rl_*`/`rt_*` — relationship tables
- 3 Laravel tables (roles, role_users, sessions — skipped on import)
