# CMED Conformidade — Tasks

**Design**: `.specs/features/cmed-conformidade/design.md`
**Status**: Draft

---

## Execution Plan

### Phase 1: Infrastructure (Sequential)

Foundation — Redis, Go dependencies, DB migration, entity.

```
T1 → T2 → T3 → T4
```

### Phase 2: Persistence Layer (Parallel OK)

Postgres repo, Redis cache client, SyncRepo extension — all independent after T4.

```
       ┌→ T5 [P] ─┐
T4 ────┼→ T6 [P] ─┼──→ T9
       └→ T7 [P] ─┘
       T8 (Meilisearch indexer) ────→ T9
```

### Phase 3: Business Layer (Sequential after T5+T6+T7+T8)

Use cases — depend on repos and cache.

```
T9 → T10 → T11
```

### Phase 4: HTTP Layer + Integration (Sequential)

Handlers, DTOs, routes, main.go wiring.

```
T11 → T12 → T13
```

### Phase 5: Import CLI (Sequential)

Standalone CLI — depends on repo + cache + meilisearch.

```
T5 + T6 + T8 → T14
```

### Phase 6: Integration Updates (Sequential)

Update existing components: reindex, search, ROADMAP, STATE.

```
T13 + T14 → T15 → T16 → T17
```

---

## Task Breakdown

### T1: Add Redis to docker-compose + config + env

**What**: Add Redis service to docker-compose, RedisConfig to config.go, env vars to .env.example
**Where**: `docker-compose.yml`, `internal/infrastructure/config/config.go`, `.env.example`
**Depends on**: None
**Reuses**: Existing docker-compose service pattern (postgres, meilisearch)
**Requirement**: CMED-18, CMED-19

**Tools**:
- MCP: NONE
- Skill: NONE

**Done when**:
- [ ] `docker-compose.yml` has `redis` service (redis:7-alpine, port 6380:6379, healthcheck, volume)
- [ ] `config.go` has `RedisConfig` struct (Host, Port, Password, DB, CacheTTL) loaded from env
- [ ] `.env.example` has REDIS_* vars
- [ ] `docker-compose.yml` api service depends_on redis
- [ ] `docker-compose.yml` api environment has REDIS_HOST=redis

**Tests**: none
**Gate**: build — `go build ./...`

**Verify**:
```bash
docker compose config | grep redis
go build ./...
```

---

### T2: Add go-redis + excelize dependencies

**What**: Install `github.com/redis/go-redis/v9` and `github.com/xuri/excelize/v2`
**Where**: `go.mod`, `go.sum`
**Depends on**: T1 (config needs RedisConfig before go-redis is used)
**Reuses**: NONE

**Tools**:
- MCP: NONE
- Skill: NONE

**Done when**:
- [ ] `go get github.com/redis/go-redis/v9` executed
- [ ] `go get github.com/xuri/excelize/v2` executed
- [ ] `go mod tidy` passes
- [ ] `go build ./...` passes

**Tests**: none
**Gate**: build — `go build ./...`

**Verify**:
```bash
grep 'go-redis/v9' go.mod
grep 'excelize/v2' go.mod
go build ./...
```

---

### T3: Create migration SQL for tb_cmed_conformidade

**What**: Write SQL migration creating `tb_cmed_conformidade` with all columns, constraints, and indexes
**Where**: `migrations/postgres/002_cmed_conformidade.sql`
**Depends on**: None (independent of Go code)
**Reuses**: Existing OBM table patterns from `obm_09-05-2026_19-35.sql`
**Requirement**: CMED-22

**Tools**:
- MCP: NONE
- Skill: NONE

**Done when**:
- [ ] SQL file creates `tb_cmed_conformidade` with all columns from design
- [ ] UNIQUE constraint on `(nu_sanreg, dt_referencia)`
- [ ] All 4 indexes created (idx_cmed_ean1, idx_cmed_referencia, idx_cmed_sanreg_ativo, idx_cmed_produto)
- [ ] SQL is idempotent (IF NOT EXISTS)
- [ ] File can be executed against PostgreSQL without errors

**Tests**: none
**Gate**: build — `go build ./...`

**Verify**:
```bash
# Verify SQL syntax by running against test DB
psql -h localhost -p 5433 -U obm -d dbportalobm -f migrations/postgres/002_cmed_conformidade.sql
```

---

### T4: Create CMEDConformidade entity + CMEDFilterParams + CMEDRepository interface

**What**: Define Go entity struct, filter params, and repository interface
**Where**: `internal/domain/entity/cmed.go` (new), `internal/domain/repository/interfaces.go` (modify)
**Depends on**: T3 (entity must match migration columns)
**Reuses**: Entity pattern from `entities.go`, FilterParams from `interfaces.go`
**Requirement**: CMED-23, CMED-08 thru CMED-11

**Tools**:
- MCP: NONE
- Skill: NONE

**Done when**:
- [ ] `entity/cmed.go` defines `CMEDConformidade` struct with all fields from design
- [ ] `repository/interfaces.go` has `CMEDFilterParams` struct
- [ ] `repository/interfaces.go` has `CMEDRepository` interface with all methods
- [ ] `SyncRepository` interface updated with `GetAllCMED` method
- [ ] `go build ./...` passes

**Tests**: none
**Gate**: build — `go build ./...`

**Verify**:
```bash
go build ./...
go vet ./...
```

---

### T5: Implement Postgres CMEDRepo [P]

**What**: Implement `CMEDRepository` in PostgreSQL — CRUD, List with filters, GetByNuSanReg, GetByEAN, GetHistorico, UpsertBatch
**Where**: `internal/infrastructure/persistence/postgres/cmed_repo.go` (new)
**Depends on**: T4 (interface definition)
**Reuses**: `GenericRepo` patterns from `generic_repo.go` (ILIKE, cursor pagination, scan)
**Requirement**: CMED-08 thru CMED-11, CMED-17

**Tools**:
- MCP: NONE
- Skill: NONE

**Done when**:
- [ ] `NewCMEDRepo(pool)` constructor
- [ ] `GetByID` — select by co_seq_id
- [ ] `List` — ILIKE on no_produto, ds_substancia; exact on nu_sanreg, nu_ean1; filter on tarja, tp_produto, dt_referencia
- [ ] `GetByNuSanReg` — select by nu_sanreg + dt_referencia (latest if dt_referencia empty)
- [ ] `GetByEAN` — select by nu_ean1/2/3 + dt_referencia
- [ ] `GetHistorico` — select all versions by nu_sanreg ordered by dt_referencia DESC
- [ ] `UpsertBatch` — INSERT ON CONFLICT (nu_sanreg, dt_referencia) DO UPDATE, batch 500
- [ ] Cursor-based pagination following existing pattern
- [ ] `go build ./...` passes

**Tests**: none
**Gate**: build — `go build ./...`

**Verify**:
```bash
go build ./...
go vet ./...
```

---

### T6: Implement Redis cache client [P]

**What**: Create Redis cache client with Get, Set, DeleteByPattern, HealthCheck
**Where**: `internal/infrastructure/persistence/redis/client.go` (new)
**Depends on**: T4 (needs CMEDConformidade type for JSON marshaling — actually just string cache, so T1+T2)
**Reuses**: NONE
**Requirement**: CMED-15, CMED-16

**Tools**:
- MCP: NONE
- Skill: NONE

**Done when**:
- [ ] `CacheRepo` struct wrapping `redis.Client`
- [ ] `NewCacheRepo(cfg RedisConfig)` constructor with connection
- [ ] `Get(ctx, key) (string, error)` — returns empty string on cache miss (not error)
- [ ] `Set(ctx, key, value, ttl)` — set with TTL from config
- [ ] `DeleteByPattern(ctx, pattern)` — SCAN + DEL for pattern-based invalidation
- [ ] `HealthCheck(ctx) error` — PING
- [ ] Graceful handling: constructor doesn't fail if Redis unavailable (log warning, return nil-safe client)
- [ ] `go build ./...` passes

**Tests**: none
**Gate**: build — `go build ./...`

**Verify**:
```bash
go build ./...
go vet ./...
```

---

### T7: Extend SyncRepo with GetAllCMED [P]

**What**: Add `GetAllCMED()` method to SyncRepo following existing pattern
**Where**: `internal/infrastructure/persistence/postgres/sync_repo.go` (modify)
**Depends on**: T4 (needs CMEDConformidade entity / SyncRepository interface update)
**Reuses**: `GetAllAMPs` pattern — query, scan, return `[]map[string]interface{}`
**Requirement**: CMED-12

**Tools**:
- MCP: NONE
- Skill: NONE

**Done when**:
- [ ] `GetAllCMED(ctx)` queries `tb_cmed_conformidade` WHERE `st_registro_ativo = 'ACTIVE'`
- [ ] Returns map with fields: co_seq_id, nu_sanreg, no_produto, ds_substancia, no_laboratorio, nu_ean1, ds_classe_terapeutica, ds_apresentacao, tp_produto, tp_regime_preco, ds_tarja, dt_referencia, vr_pf_sem_impostos, vr_pmc_sem_impostos
- [ ] `go build ./...` passes

**Tests**: none
**Gate**: build — `go build ./...`

**Verify**:
```bash
go build ./...
```

---

### T8: Extend Meilisearch indexer + client for CMED

**What**: Add `obm_cmed` index config, `IndexCMEDs` method, and hit mapping in Search
**Where**: `internal/infrastructure/persistence/meilisearch/indexer.go` (modify), `client.go` (modify)
**Depends on**: T4 (interface), T7 (SyncRepo provides CMED data format)
**Reuses**: `batchIndex`, Search hit mapping pattern
**Requirement**: CMED-12, CMED-13

**Tools**:
- MCP: NONE
- Skill: NONE

**Done when**:
- [ ] `ConfigureIndexes` updated with `cmed` index config (searchable, filterable, sortable per design)
- [ ] `IndexCMEDs(ctx, docs)` method added
- [ ] `Search` method: default entities now includes `cmed` (vmp, amp, supplier, cmed)
- [ ] Hit mapping for `cmed`: co_seq_id→ID, no_produto→Nome, nu_sanreg→Codigo, no_laboratorio→Fabricante, ds_apresentacao→Descricao
- [ ] `go build ./...` passes

**Tests**: none
**Gate**: build — `go build ./...`

**Verify**:
```bash
go build ./...
```

---

### T9: Create CMED Usecase

**What**: Implement CMEDUsecase with List, GetByID, GetByRegistro, GetByEAN, GetHistorico
**Where**: `internal/usecase/cmed.go` (new)
**Depends on**: T5 (CMEDRepo), T6 (CacheRepo)
**Reuses**: VMPUsecase pattern
**Requirement**: CMED-08 thru CMED-11, CMED-17

**Tools**:
- MCP: NONE
- Skill: NONE

**Done when**:
- [ ] `NewCMEDUsecase(cmedRepo, cacheRepo)` constructor
- [ ] `List(ctx, filter)` — delegates to repo, returns CursorPage
- [ ] `GetByID(ctx, id)` — delegates to repo
- [ ] `GetByRegistro(ctx, nuSanReg, dtReferencia)` — delegates to repo
- [ ] `GetByEAN(ctx, ean, dtReferencia)` — delegates to repo
- [ ] `GetHistorico(ctx, id)` — delegates to repo
- [ ] `go build ./...` passes

**Tests**: none
**Gate**: build — `go build ./...`

**Verify**:
```bash
go build ./...
```

---

### T10: Create AMPP-CMED Usecase (JOIN with cache)

**What**: Implement AMPPCMEDUsecase — GetAMPPWithCMED with Redis cache
**Where**: `internal/usecase/ampp_cmed.go` (new)
**Depends on**: T5 (CMEDRepo), T6 (CacheRepo)
**Reuses**: NONE (new pattern)
**Requirement**: CMED-14, CMED-15, CMED-16

**Tools**:
- MCP: NONE
- Skill: NONE

**Done when**:
- [ ] `NewAMPPCMEDUsecase(amppRepo, cmedRepo, cacheRepo)` constructor
- [ ] `GetAMPPWithCMED(ctx, amppID, dtReferencia)`:
  1. Get AMPP by ID from amppRepo
  2. If AMPP.nu_sanreg is null → return AMPP with cmed=null
  3. Check Redis cache `ampp_cmed:{id}:{dt_referencia}`
  4. Cache hit → unmarshal and return
  5. Cache miss → query CMED by nu_sanreg + dt_referencia
  6. Get AMP by AMPP.co_apid
  7. Get VMP by AMP.co_vpid
  8. Build AMPPCMEDResponse, cache it, return
- [ ] Graceful Redis degradation: if cacheRepo is nil, skip cache
- [ ] `go build ./...` passes

**Tests**: none
**Gate**: build — `go build ./...`

**Verify**:
```bash
go build ./...
```

---

### T11: Create CMED Handler + DTOs

**What**: HTTP handlers for CMED endpoints + request/response DTOs
**Where**: `internal/interface/http/handler/cmed_handler.go` (new), `internal/interface/http/dto/dto.go` (modify)
**Depends on**: T9 (CMEDUsecase), T10 (AMPPCMEDUsecase)
**Reuses**: VMPHandler pattern, SearchHandler pattern for filter binding
**Requirement**: CMED-08 thru CMED-11, CMED-14, CMED-17

**Tools**:
- MCP: NONE
- Skill: NONE

**Done when**:
- [ ] `CMEDHandler` struct with cmedUC and amppCmedUC
- [ ] `NewCMEDHandler(cmedUC, amppCmedUC)` constructor
- [ ] `List(c)` — binds query: nome, registro, ean, tarja, tipo_produto, regime_preco, dt_referencia, limit, cursor
- [ ] `GetByID(c)` — param :id
- [ ] `GetByRegistro(c)` — param :registro
- [ ] `GetByEAN(c)` — param :ean
- [ ] `GetHistorico(c)` — param :id
- [ ] `GetAMPPWithCMED(c)` — param :id, query dt_referencia
- [ ] DTOs: CMEDListRequest, CMEDResponse in dto.go
- [ ] Swagger annotations on all endpoints
- [ ] `go build ./...` passes

**Tests**: none
**Gate**: build — `go build ./...`

**Verify**:
```bash
go build ./...
```

---

### T12: Register CMED routes in router

**What**: Add CMED routes to the API router
**Where**: `internal/interface/http/router/router.go` (modify)
**Depends on**: T11 (CMEDHandler)
**Reuses**: Existing route registration pattern
**Requirement**: CMED-08 thru CMED-11, CMED-14, CMED-17

**Tools**:
- MCP: NONE
- Skill: NONE

**Done when**:
- [ ] `GET /api/v1/cmed` → cmedHandler.List
- [ ] `GET /api/v1/cmed/:id` → cmedHandler.GetByID
- [ ] `GET /api/v1/cmed/registro/:registro` → cmedHandler.GetByRegistro
- [ ] `GET /api/v1/cmed/ean/:ean` → cmedHandler.GetByEAN
- [ ] `GET /api/v1/cmed/:id/historico` → cmedHandler.GetHistorico
- [ ] `GET /api/v1/ampp/:id/cmed` → cmedHandler.GetAMPPWithCMED
- [ ] `SetupRouter` signature updated to accept cmedHandler
- [ ] `go build ./...` passes

**Tests**: none
**Gate**: build — `go build ./...`

**Verify**:
```bash
go build ./...
```

---

### T13: Wire CMED components in main.go

**What**: Instantiate all CMED components (repo, usecase, handler) and pass to router
**Where**: `cmd/api/main.go` (modify)
**Depends on**: T5, T6, T9, T10, T11, T12
**Reuses**: Existing wiring pattern in main.go
**Requirement**: CMED-08 thru CMED-17

**Tools**:
- MCP: NONE
- Skill: NONE

**Done when**:
- [ ] Redis client created with config.Redis
- [ ] CMEDRepo, CMEDUsecase, AMPPCMEDUsecase, CMEDHandler instantiated
- [ ] cmedHandler passed to router.SetupRouter
- [ ] CacheRepo passed to AMPPCMEDUsecase
- [ ] Graceful: if Redis connection fails, log warning, set cacheRepo=nil (degrade without cache)
- [ ] `go build ./...` passes

**Tests**: none
**Gate**: build — `go build ./...`

**Verify**:
```bash
go build ./...
```

---

### T14: Create CMED Import CLI

**What**: Standalone CLI to import XLSX CMED data into PostgreSQL + Meilisearch
**Where**: `cmd/cmed_import/main.go` (new)
**Depends on**: T5 (CMEDRepo.UpsertBatch), T6 (CacheRepo.DeleteByPattern), T8 (MeilisearchRepo.IndexCMEDs)
**Reuses**: `cmd/import/main.go` flag/logging pattern, `converter.go` batch pattern
**Requirement**: CMED-01 thru CMED-07

**Tools**:
- MCP: NONE
- Skill: NONE

**Done when**:
- [ ] CLI flags: --source (required), --date (required, YYYY-MM-DD), --header-row (default 42), --skip-index (default false)
- [ ] Opens XLSX with excelize
- [ ] Reads header at specified row, maps columns by name (case-insensitive)
- [ ] Validates required columns: REGISTRO, PRODUTO, SUBSTÂNCIA — error with list if missing
- [ ] Parses data rows starting at header-row + 1
- [ ] Cleans data: "-" → nil, comma decimals → dot, numeric strings → int64/float64
- [ ] Builds JSONB js_precos_pf and js_precos_pmc with all PF/PMC alíquotas
- [ ] Calls UpsertBatch in batches of 500
- [ ] Invalidates Redis cache: DeleteByPattern("cmed:*") + DeleteByPattern("ampp_cmed:*")
- [ ] Reindexes Meilisearch obm_cmed (unless --skip-index)
- [ ] Prints summary: rows imported, rows skipped, errors
- [ ] `go build ./cmd/cmed_import/...` passes

**Tests**: none
**Gate**: build — `go build ./cmd/cmed_import/...`

**Verify**:
```bash
go build -o /dev/null ./cmd/cmed_import/
```

---

### T15: Update ReindexUsecase for CMED

**What**: Add CMED sync and index to the reindex pipeline
**Where**: `internal/usecase/reindex.go` (modify)
**Depends on**: T7 (SyncRepo.GetAllCMED), T8 (MeilisearchRepo.IndexCMEDs)
**Reuses**: Existing reindex pattern
**Requirement**: CMED-12

**Tools**:
- MCP: NONE
- Skill: NONE

**Done when**:
- [ ] Reindex method calls `syncRepo.GetAllCMED(ctx)`
- [ ] Calls `meiliRepo.IndexCMEDs(ctx, cmedDocs)`
- [ ] Logs indexed count
- [ ] `indexed["cmed"]` added to result map
- [ ] `go build ./...` passes

**Tests**: none
**Gate**: build — `go build ./...`

**Verify**:
```bash
go build ./...
```

---

### T16: Update SearchUsecase + SearchHandler for CMED

**What**: Add `cmed` to default search entities and filter mapping
**Where**: `internal/infrastructure/persistence/meilisearch/client.go` (default entities), `internal/interface/http/handler/search_handler.go` (filter mapping)
**Depends on**: T8 (Meilisearch CMED index + hit mapping)
**Reuses**: Existing search filter binding pattern
**Requirement**: CMED-13

**Tools**:
- MCP: NONE
- Skill: NONE

**Done when**:
- [ ] Default entities in `Search` method: `["vmp", "amp", "supplier", "cmed"]`
- [ ] SearchHandler: filter[tarja] maps to `ds_tarja` for cmed entity
- [ ] SearchHandler: filter[registro] maps to `nu_sanreg` for cmed entity
- [ ] Swagger docs updated: entity param description includes "cmed"
- [ ] `go build ./...` passes

**Tests**: none
**Gate**: build — `go build ./...`

**Verify**:
```bash
go build ./...
```

---

### T17: Update ROADMAP + STATE docs

**What**: Update project documentation to reflect CMED Conformidade feature
**Where**: `.specs/project/ROADMAP.md` (modify), `.specs/project/STATE.md` (modify)
**Depends on**: T13, T14, T15, T16 (all implementation done)
**Reuses**: NONE

**Tools**:
- MCP: NONE
- Skill: NONE

**Done when**:
- [ ] ROADMAP.md has v1.3.0 section with CMED Conformidade feature items (checked off)
- [ ] STATE.md updated: D10 (CMED table independent from OBM tables), D11 (Redis for cache with graceful degradation), D12 (nu_sanreg as join key)
- [ ] STATE.md todos updated

**Tests**: none
**Gate**: none

**Verify**:
```bash
cat .specs/project/ROADMAP.md
cat .specs/project/STATE.md
```

---

## Parallel Execution Map

```
Phase 1 (Sequential - Foundation):
T1 → T2 → T3 → T4

Phase 2 (Parallel - Persistence):
       ┌→ T5 [P] (CMEDRepo) ──────┐
T4 ────┼→ T6 [P] (Redis cache) ───┼──→ T9
       └→ T7 [P] (SyncRepo) ──────┘
T8 (Meilisearch) ──────────────────→ T9

Phase 3 (Sequential - Business):
T9 → T10

Phase 4 (Sequential - HTTP):
T10 → T11 → T12 → T13

Phase 5 (Parallel with Phase 4 - CLI):
T5 + T6 + T8 → T14

Phase 6 (Sequential - Integration):
T13 + T14 + T15 → T15 → T16 → T17
```

---

## Task Granularity Check

| Task | Scope | Status |
|---|---|---|
| T1: Redis docker + config + env | 3 files, 1 concept | ✅ Granular |
| T2: Go dependencies | 1 command, go.mod | ✅ Granular |
| T3: Migration SQL | 1 file | ✅ Granular |
| T4: Entity + interface | 2 files, related types | ✅ Granular |
| T5: Postgres CMEDRepo | 1 file, 1 component | ✅ Granular |
| T6: Redis cache client | 1 file, 1 component | ✅ Granular |
| T7: SyncRepo extension | 1 method addition | ✅ Granular |
| T8: Meilisearch extension | 2 files modified, 1 concept | ✅ Granular |
| T9: CMED Usecase | 1 file, 1 component | ✅ Granular |
| T10: AMPP-CMED Usecase | 1 file, 1 component | ✅ Granular |
| T11: Handler + DTOs | 2 files, 1 component | ✅ Granular |
| T12: Router registration | 1 file, route wiring | ✅ Granular |
| T13: main.go wiring | 1 file, DI wiring | ✅ Granular |
| T14: Import CLI | 1 file, 1 CLI tool | ✅ Granular |
| T15: Reindex update | 1 file, add 1 step | ✅ Granular |
| T16: Search update | 2 files, add entity | ✅ Granular |
| T17: Docs update | 2 files, docs | ✅ Granular |

---

## Diagram-Definition Cross-Check

| Task | Depends On (body) | Diagram Shows | Status |
|---|---|---|---|
| T1 | None | Start node | ✅ Match |
| T2 | T1 | T1→T2 | ✅ Match |
| T3 | None | Independent start | ✅ Match |
| T4 | T3 | T3→T4 | ✅ Match |
| T5 | T4 | T4→T5 [P] | ✅ Match |
| T6 | T4 | T4→T6 [P] | ✅ Match |
| T7 | T4 | T4→T7 [P] | ✅ Match |
| T8 | T4, T7 | T4+T7→T8 | ✅ Match |
| T9 | T5, T6, T7, T8 | T5+T6+T7+T8→T9 | ✅ Match |
| T10 | T5, T6 | T9→T10 | ✅ Match |
| T11 | T9, T10 | T10→T11 | ✅ Match |
| T12 | T11 | T11→T12 | ✅ Match |
| T13 | T5, T6, T9, T10, T11, T12 | T12→T13 | ✅ Match |
| T14 | T5, T6, T8 | T5+T6+T8→T14 (parallel with Phase 4) | ✅ Match |
| T15 | T7, T8 | After T14 | ✅ Match |
| T16 | T8 | After T15 | ✅ Match |
| T17 | T13, T14, T15, T16 | After T16 | ✅ Match |

---

## Test Co-location Validation

No TESTING.md exists for this project. Test types are set to `none` for all tasks since the codebase currently has no test commands defined. All gates are `build` level (`go build ./...`).

| Task | Code Layer | Matrix Requires | Task Says | Status |
|---|---|---|---|---|
| T1-T17 | Various | none (no TESTING.md) | none | ✅ OK |

**Note**: Integration tests for the CMED pipeline (import + search + JOIN) should be added as a follow-up, similar to the existing `E2E test with real dump` item in ROADMAP v1.1.0.
