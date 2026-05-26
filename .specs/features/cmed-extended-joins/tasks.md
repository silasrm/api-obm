# CMED Extended Joins Tasks

**Design**: `.specs/features/cmed-extended-joins/design.md`
**Status**: Done

---

## Execution Plan

### Phase 1: Foundation (Sequential)

T1 → T2

### Phase 2: Usecases + Tests (Parallel)

```
T2 ──┬→ T3 [P]
     └→ T4 [P]
```

### Phase 3: Handler + Routes + Wiring (Sequential)

```
T3, T4 ──→ T5 ──→ T6 ──→ T7
```

### Phase 4: Documentation (Sequential)

```
T7 ──→ T8
```

---

## Task Breakdown

### T1: Add `GetByCNPJ` to CMEDRepository interface

**What**: Adicionar método `GetByCNPJ` na interface `CMEDRepository` e helper `normalizeCNPJ`
**Where**: `internal/domain/repository/interfaces.go`
**Depends on**: None
**Reuses**: Padrão de `GetByEAN` / `GetByNuSanReg`
**Requirement**: SUP-06, CNPJ-01

**Done when**:

- [ ] `GetByCNPJ(ctx, cnpj string, dtReferencia string) ([]entity.CMEDConformidade, error)` adicionado à interface
- [ ] `normalizeCNPJ(s string) string` helper definido (exportado ou não)
- [ ] Build passa: `go build ./internal/...`

**Tests**: none (interface only)
**Gate**: build

---

### T2: Implement `GetByCNPJ` in CMEDRepo

**What**: Implementar `GetByCNPJ` no repositório PostgreSQL — normaliza CNPJ, busca por `REGEXP_REPLACE(nu_cnpj, '[^0-9]', '', 'g') = $1`, retorna `[]CMEDConformidade`
**Where**: `internal/infrastructure/persistence/postgres/cmed_repo.go`
**Depends on**: T1
**Reuses**: `scanCMEDRows`, `cmedColumns`, padrão de `GetByEAN`

**Done when**:

- [ ] `GetByCNPJ` implementado com normalização CNPJ
- [ ] Suporta `dtReferencia` opcional (filtra por data ou retorna mais recente por `nu_sanreg`)
- [ ] Filtra `st_registro_ativo = 'ACTIVE'`
- [ ] Ordena por `no_produto`
- [ ] `normalizeCNPJ` implementado no package
- [ ] Build passa: `go build ./internal/...`

**Tests**: none (PG repo — validação manual ou integration test futuro)
**Gate**: build

---

### T3: Create `AMPCMEDUsecase` + tests [P]

**What**: Criar usecase `AMPCMEDUsecase` com método `GetAMPWithCMED` + testes unitários
**Where**: `internal/usecase/amp_cmed.go`, `internal/usecase/amp_cmed_test.go`
**Depends on**: T2
**Reuses**: Padrão de `AMPPCMEDUsecase`

**Done when**:

- [ ] `AMPCMEDUsecase` struct com deps: `AMPRepository`, `VMPRepository`, `CMEDRepository`, `*CacheRepo`
- [ ] `GetAMPWithCMED(ctx, ampID, dtReferencia)` retorna `*AMPCMEDResponse`
- [ ] AMP não encontrado → erro "AMP not found"
- [ ] `nu_nreg` vazio/nulo → CMED nil (sem erro)
- [ ] VMP carregado via `AMP.COVpID` (se > 0)
- [ ] Cache Redis key `amp_cmed:{ampID}:{dtReferencia}`
- [ ] Graceful degradation: CMED falha → log + CMED nil
- [ ] Testes: AMP+CMED encontrado, AMP sem nu_nreg, AMP 404, CMED graceful, cache hit/miss
- [ ] Gate check passa: `go test ./internal/usecase/... -count=1`

**Tests**: unit
**Gate**: quick

---

### T4: Create `SupplierCMEDUsecase` + tests [P]

**What**: Criar usecase `SupplierCMEDUsecase` com método `GetSupplierWithCMED` + testes unitários
**Where**: `internal/usecase/supplier_cmed.go`, `internal/usecase/supplier_cmed_test.go`
**Depends on**: T2
**Reuses**: Padrão de `AMPPCMEDUsecase`

**Done when**:

- [ ] `SupplierCMEDUsecase` struct com deps: `SupplierRepository`, `CMEDRepository`, `*CacheRepo`
- [ ] `GetSupplierWithCMED(ctx, supplierID, dtReferencia)` retorna `*SupplierCMEDResponse`
- [ ] Supplier não encontrado → erro "Supplier not found"
- [ ] `nu_cnpj` vazio/nulo → CMED [] (sem erro)
- [ ] CNPJ normalizado via `normalizeCNPJ` antes do lookup
- [ ] Cache Redis key `supplier_cmed:{supplierID}:{dtReferencia}`
- [ ] Graceful degradation: CMED falha → log + CMED []
- [ ] Testes: Supplier+CMED lista, Supplier sem CNPJ, Supplier 404, CMED graceful, cache hit/miss, normalização CNPJ
- [ ] Gate check passa: `go test ./internal/usecase/... -count=1`

**Tests**: unit
**Gate**: quick

---

### T5: Extend CMEDHandler + handler tests

**What**: Adicionar `GetAMPWithCMED` e `GetSupplierWithCMED` ao handler + testes
**Where**: `internal/interface/http/handler/cmed_handler.go`, `internal/interface/http/handler/cmed_handler_test.go`
**Depends on**: T3, T4
**Reuses**: Padrão de `GetAMPPWithCMED`

**Done when**:

- [ ] `NewCMEDHandler` recebe `ampCmedUC` e `supplierCmedUC` adicionais
- [ ] `GetAMPWithCMED` parse `:id` + `dt_referencia` query
- [ ] `GetSupplierWithCMED` parse `:id` + `dt_referencia` query
- [ ] Respostas 200 (sucesso), 400 (id inválido), 404 (não encontrado)
- [ ] Testes handler: AMP+CMED sucesso, AMP 404, Supplier+CMED sucesso, Supplier 404, dt_referencia query param
- [ ] Gate check passa: `go test ./internal/interface/http/handler/... -count=1`

**Tests**: unit
**Gate**: quick

---

### T6: Add routes + wire in main.go

**What**: Adicionar rotas `/amp/:id/cmed` e `/suppliers/:id/cmed` no router + instanciar usecases no main.go
**Where**: `internal/interface/http/router/router.go`, `cmd/api/main.go`
**Depends on**: T5
**Reuses**: Padrão de `/ampp/:id/cmed`

**Done when**:

- [ ] Rota `api.GET("/amp/:id/cmed", cmedHandler.GetAMPWithCMED)` adicionada
- [ ] Rota `api.GET("/suppliers/:id/cmed", cmedHandler.GetSupplierWithCMED)` adicionada
- [ ] `AMPCMEDUsecase` e `SupplierCMEDUsecase` instanciados em `main.go`
- [ ] `NewCMEDHandler` chamado com novos usecases
- [ ] Build passa: `go build ./cmd/api/...`

**Tests**: none (wiring — validado por build + T5 handler tests)
**Gate**: build

---

### T7: Full verification

**What**: Rodar build completo + todos os testes para verificar zero regressões
**Where**: N/A
**Depends on**: T6
**Reuses**: N/A

**Done when**:

- [ ] `go build ./internal/...` passa
- [ ] `go build ./cmd/api/...` passa
- [ ] `go build ./cmd/cmed_import/...` passa
- [ ] `go test ./internal/usecase/... ./internal/interface/http/handler/... -count=1` passa
- [ ] Zero regressões em testes existentes

**Tests**: none
**Gate**: full

---

### T8: Update documentation

**What**: Atualizar README, VERSION.md e criar release notes v1.4.0
**Where**: `README.md`, `docs/VERSION.md`, `docs/release-notes/v1.4.0.md`
**Depends on**: T7
**Reuses**: Formato do README e release notes existentes

**Done when**:

- [ ] README: endpoints `/amp/:id/cmed` e `/suppliers/:id/cmed` na tabela de rotas
- [ ] README: exemplos 13-14 com curl + response JSON
- [ ] README: seção explicando joins AMP↔CMED e Supplier↔CMED
- [ ] VERSION.md: atualizado para v1.4.0 com changelog + histórico
- [ ] `docs/release-notes/v1.4.0.md` criado com mudanças, endpoints, exemplos, compatibilidade

**Tests**: none
**Gate**: build

---

## Parallel Execution Map

```
Phase 1 (Sequential):
T1 ──→ T2

Phase 2 (Parallel):
T2 ──┬→ T3 [P]
     └→ T4 [P]

Phase 3 (Sequential):
T3, T4 ──→ T5 ──→ T6 ──→ T7 ──→ T8
```

---

## Task Granularity Check

| Task | Scope | Status |
| --- | --- | --- |
| T1: Add GetByCNPJ interface | 1 method + 1 helper | ✅ Granular |
| T2: Implement GetByCNPJ | 1 method implementation | ✅ Granular |
| T3: AMPCMEDUsecase + tests | 1 usecase + test file | ✅ Granular |
| T4: SupplierCMEDUsecase + tests | 1 usecase + test file | ✅ Granular |
| T5: Handler methods + tests | 2 handler methods + test update | ✅ Granular |
| T6: Routes + wiring | 2 files modify | ✅ Granular |
| T7: Full verification | Build + test run | ✅ Granular |
| T8: Documentation | 3 files | ⚠️ OK — cohesive (all docs) |

---

## Diagram-Definition Cross-Check

| Task | Depends On (body) | Diagram Shows | Status |
| --- | --- | --- | --- |
| T1 | None | No incoming arrows | ✅ Match |
| T2 | T1 | T1 → T2 | ✅ Match |
| T3 | T2 | T2 → T3 | ✅ Match |
| T4 | T2 | T2 → T4 | ✅ Match |
| T5 | T3, T4 | T3,T4 → T5 | ✅ Match |
| T6 | T5 | T5 → T6 | ✅ Match |
| T7 | T6 | T6 → T7 | ✅ Match |
| T8 | T7 | T7 → T8 | ✅ Match |

---

## Test Co-location Validation

| Task | Code Layer | Matrix Requires | Task Says | Status |
| --- | --- | --- | --- | --- |
| T1 | Interface | none | none | ✅ OK |
| T2 | Repo implementation | none | none | ✅ OK |
| T3 | Usecase | unit | unit | ✅ OK |
| T4 | Usecase | unit | unit | ✅ OK |
| T5 | Handler | unit | unit | ✅ OK |
| T6 | Router + main | none | none | ✅ OK |
| T7 | Verification | none | none | ✅ OK |
| T8 | Documentation | none | none | ✅ OK |
