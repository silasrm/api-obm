# Roadmap — API OBM

## v1.0.0 — Released (2026-05-24)

- [x] REST API with all endpoints (VMP, AMP, VTM, VMPP, AMPP, DCB, Ingredient, Supplier, Domains)
- [x] Meilisearch global search with filters
- [x] JWT authentication
- [x] Cursor-based pagination
- [x] Swagger UI + Postman collection
- [x] Docker Compose infrastructure
- [x] Manual MySQL→PostgreSQL conversion script (`scripts/convert_sql.go`)

## v1.1.0 — Import CLI (Done)

- [x] T1: Refactor converter into reusable package (`internal/importer/converter.go`)
- [x] T2: Source resolver (ZIP/SQL/MySQL → io.Reader)
- [x] T3: PostgreSQL importer (streaming SQL execution)
- [x] T4: Post-import validator (row counts, FK integrity)
- [x] T5: Metadata manager (`obm_metadata` table)
- [x] T6: Pipeline orchestrator (streaming via io.Pipe)
- [x] T7: CLI entrypoint (`cmd/import/main.go`)
- [x] T8: Backward-compat wrapper (`scripts/convert_sql.go`)
- [x] T9: README update with import docs
- [ ] E2E test with real dump

## v1.2.0 — Future

- [ ] Expand Meilisearch indexes (VTM, VMPP, AMPP, DCB, Ingredient)
- [ ] Incremental/delta import support
- [ ] Import scheduling (cron or API trigger)
- [ ] API endpoint for import triggering
- [ ] Monitoring / observability (structured logging, metrics)
