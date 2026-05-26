# Centralizador de Versão — API OBM

> **ATENÇÃO**: Este é o arquivo de referência e auditoria da versão atual da API. Todas as builds, tags no git ou release notes devem referenciar a versão travada aqui.

## Versão Atual

```
Versão: 1.4.0
Release Date: 2026-05-25
Timezone: America/Sao_Paulo
Status: Released
```

### v1.4.0 (2026-05-25)

**Status:** Released

**Mudanças:**

- feat: Endpoint AMP+CMED (`GET /api/v1/amp/:id/cmed`) via `AMP.nu_nreg = CMED.nu_sanreg` — retorna AMP + VMP pai + preço CMED
- feat: Endpoint Supplier+CMED (`GET /api/v1/suppliers/:id/cmed`) via `Supplier.nu_cnpj = CMED.nu_cnpj` — retorna Fornecedor + lista de produtos CMED
- feat: CMEDRepository.GetByCNPJ com normalização de CNPJ (remove não-dígitos antes do JOIN)
- feat: AMPCMEDUsecase + SupplierCMEDUsecase com cache Redis e graceful degradation
- feat: 23 novos testes unitários e de handler (105 total)

**Arquivos novos:**

- `internal/usecase/amp_cmed.go` — Usecase JOIN AMP+CMED com cache Redis
- `internal/usecase/supplier_cmed.go` — Usecase JOIN Supplier+CMED com cache Redis
- `internal/usecase/amp_cmed_test.go` — Testes usecase AMP-CMED
- `internal/usecase/supplier_cmed_test.go` — Testes usecase Supplier-CMED
- `internal/interface/http/handler/amp_cmed_handler.go` — Handler HTTP AMP+CMED
- `internal/interface/http/handler/supplier_cmed_handler.go` — Handler HTTP Supplier+CMED
- `internal/interface/http/handler/amp_cmed_handler_test.go` — Testes handler AMP+CMED
- `internal/interface/http/handler/supplier_cmed_handler_test.go` — Testes handler Supplier+CMED

**Arquivos modificados:**

- `internal/domain/repository/interfaces.go` — Adicionados GetByRegistroSanitário e GetByCNPJ no CMEDRepository
- `internal/infrastructure/persistence/postgres/cmed_repo.go` — Implementação GetByCNPJ com normalização de CNPJ
- `internal/interface/http/router/router.go` — Adicionadas rotas `/amp/:id/cmed` e `/suppliers/:id/cmed`
- `cmd/api/main.go` — Wiring de AMPCMEDUsecase, SupplierCMEDUsecase e handlers
- `README.md` — Seções JOIN AMP+CMED, JOIN Fornecedor+CMED, Mecanismos de JOIN, exemplos 13-14, tabela de rotas atualizada

### v1.3.0 (2026-05-25)

**Status:** Released

**Mudanças:**

- feat: Integração CMED Conformidade de Preços — tabela `tb_cmed_conformidade` com versionamento por data de referência
- feat: Endpoints CMED — listagem com filtros (nome, registro, EAN, tarja, tipo, regime), busca por ID, Registro Sanitário, EAN e histórico de preços
- feat: JOIN AMPP+CMED com cache Redis (`GET /api/v1/ampp/:id/cmed`) — dados ontológicos + preço regulado em uma única chamada
- feat: CLI de importação CMED (`cmd/cmed_import/`) — parse XLSX com `--header-row` configurável (default 42), UPSERT batch, invalidação de cache e reindexação Meilisearch
- feat: Index Meilisearch `obm_cmed` — CMED incluído na busca global com filtros por tarja e registro sanitário
- feat: Redis no Docker Compose com graceful degradation — API funciona sem Redis (sem cache)
- feat: 32 testes unitários (CMED usecase, AMPP-CMED usecase, CMED handler)
- feat: Documentação completa — seção CMED no README com exemplos, CLI de Importação CMED, variáveis Redis no .env
- fix: Conversão `dt_referencia::text` no cmed_repo.go para compatibilidade pgx com tipo DATE do PostgreSQL
- fix: Indentação YAML do docker-compose.yml (SYNC_ON_STARTUP, GIN_MODE, REDIS_HOST, restart)

**Arquivos novos:**

- `migrations/postgres/002_cmed_conformidade.sql` — DDL da tabela CMED
- `internal/domain/entity/cmed.go` — Entidade CMEDConformidade
- `internal/infrastructure/persistence/postgres/cmed_repo.go` — Repositório PostgreSQL CMED
- `internal/infrastructure/persistence/redis/client.go` — Cliente Redis com cache
- `internal/usecase/cmed.go` — Usecase CMED
- `internal/usecase/ampp_cmed.go` — Usecase JOIN AMPP+CMED com cache
- `internal/interface/http/handler/cmed_handler.go` — Handler HTTP CMED
- `cmd/cmed_import/main.go` — CLI importação CMED
- `internal/usecase/cmed_test.go` — Testes usecase CMED
- `internal/usecase/ampp_cmed_test.go` — Testes usecase AMPP-CMED
- `internal/interface/http/handler/cmed_handler_test.go` — Testes handler CMED
- `.specs/features/cmed-conformidade/spec.md` — Spec com 23 requisitos rastreáveis
- `.specs/features/cmed-conformidade/design.md` — Design arquitetural
- `.specs/features/cmed-conformidade/tasks.md` — 17 tasks com validação cruzada

**Arquivos modificados:**

- `docker-compose.yml` — Adicionado serviço Redis, depends_on, variáveis de ambiente
- `internal/infrastructure/config/config.go` — Adicionado RedisConfig
- `.env.example` — Adicionadas variáveis REDIS_*
- `go.mod` / `go.sum` — Adicionados go-redis/v9 e excelize/v2
- `internal/domain/repository/interfaces.go` — Adicionados CMEDRepository, CMEDFilterParams, GetAllCMED
- `internal/infrastructure/persistence/postgres/sync_repo.go` — Adicionado GetAllCMED
- `internal/infrastructure/persistence/meilisearch/indexer.go` — Adicionado index cmed e IndexCMEDs
- `internal/infrastructure/persistence/meilisearch/client.go` — Adicionado cmed como entidade default e mapeamento de hits
- `internal/interface/http/dto/dto.go` — Adicionados DTOs CMED e filtros tarja/registro
- `internal/interface/http/handler/search_handler.go` — Adicionados filter[tarja] e filter[registro]
- `internal/interface/http/router/router.go` — Adicionadas rotas CMED e /ampp/:id/cmed
- `cmd/api/main.go` — Wiring de CMEDRepo, CacheRepo, usecases e handler
- `internal/usecase/reindex.go` — Adicionado step CMED ao reindex
- `README.md` — Seções CMED Conformidade, CLI Importação CMED, exemplos 8-12, variáveis Redis, tabela de rotas atualizada
- `.specs/project/ROADMAP.md` — Adicionada v1.3.0
- `.specs/project/STATE.md` — Adicionadas decisões D10-D13

### v1.0.0 (2026-05-24)

**Status:** Released

**Mudanças:**

- feat: API REST da Ontologia Brasileira de Medicamentos (OBM) — modelo dm+d do NHS adaptado para o Brasil
- feat: Endpoints de consulta para VMP, AMP, VTM, VMPP, AMPP, DCB, Ingredientes, Fornecedores e Domínios
- feat: Busca global via Meilisearch com filtros por nome, código, fabricante, descrição e status ativo
- feat: Endpoints de detalhe para VMP (com VTM, domínios, ingredientes, rotas, formas, classes ATC) e AMP (com VMP, fornecedor, ingredientes, rotas)
- feat: 32 tipos de domínio (forma farmacêutica, via, classe ATC, categoria de controle, etc.)
- feat: Autenticação JWT (HS256) com login/geração de token
- feat: Paginação cursor-based em todos os endpoints de listagem
- feat: Health check público (PostgreSQL + Meilisearch)
- feat: Reindexação administrativa do Meilisearch
- feat: Swagger UI integrado (`/swagger/index.html`)
- feat: Collection Postman com ambiente local
- feat: Infraestrutura Docker Compose (PostgreSQL 16, Meilisearch v1.8, API)
- feat: Seed de usuários iniciais (admin/admin123, viewer/viewer123)
- feat: Suite de testes unitários e de integração
- fix: Correção da descrição de "Observatório de Medicamentos" para "Ontologia Brasileira de Medicamentos" em todos os arquivos
- fix: Auth middleware aceita token com ou sem prefixo Bearer (compatibilidade Swagger UI)
- fix: Porta do servidor alterada de 8080 para 8094
- fix: Porta do PostgreSQL no host alterada de 5432 para 5433 (evita conflito)
- fix: Healthcheck do Meilisearch usando curl em vez de wget
- fix: Syntax SQL na migration e carregamento de .env no seed script

**Stack:**

- Go 1.25, Gin, pgx/v5, Meilisearch Go SDK, swaggo/swag
- PostgreSQL 16, Meilisearch v1.8
- Docker Compose, JWT HS256, bcrypt

**Arquivos novos:**

- `cmd/api/main.go` — Entry point da API
- `internal/domain/entity/entities.go` — Entidades de domínio
- `internal/domain/repository/interfaces.go` — Contratos dos repositórios
- `internal/infrastructure/persistence/postgres/` — Repositórios PostgreSQL
- `internal/infrastructure/persistence/meilisearch/` — Cliente e indexer Meilisearch
- `internal/infrastructure/config/config.go` — Carregamento de configuração
- `internal/interface/http/handler/` — Handlers HTTP
- `internal/interface/http/dto/dto.go` — DTOs de request/response
- `internal/interface/http/middleware/auth.go` — Middleware JWT
- `internal/interface/http/router/router.go` — Definição de rotas
- `internal/usecase/` — Casos de uso
- `docs/swagger.yaml`, `docs/swagger.json`, `docs/docs.go` — Documentação Swagger
- `postman/OBM_API.postman_collection.json` — Collection Postman
- `postman/OBM_API_Local.postman_environment.json` — Ambiente Postman
- `scripts/seed_users.go` — Seed de usuários
- `scripts/gen_postman.go` — Gerador de collection Postman
- `scripts/convert_sql.go` — Conversor MySQL → PostgreSQL
- `migrations/postgres/001_obm_schema.sql` — Schema completo do banco
- `README.md` — Documentação para usuários
- `TESTING.md` — Guia de teste

## Como Incrementar a Versão (SemVer)

1. **Atualize este documento** aumentando a trilha conforme o padrão de versões semânticas:
- **Patch** (ex: 1.0.0 → 1.0.1): Correções de bugs, refactorings, limpezas que não entregam novos comportamentos
- **Minor** (ex: 1.0.0 → 1.1.0): Novos endpoints, filtros, funcionalidades retrocompatíveis
- **Major** (ex: 1.0.0 → 2.0.0): Breaking changes na API (mudança de formato de resposta, remoção de endpoints, etc.)

2. **Crie a Release Note** em `docs/release-notes/vX.Y.Z.md` descrevendo as mudanças

3. **Crie a tag git** com `git tag -a vX.Y.Z -m "Release vX.Y.Z"`

## Histórico de Versões

| Versão | Data | Descrição |
|--------|------|-----------|
| `v1.4.0` | 25/05/2026 | **feat: CMED Extended Joins** — Endpoints AMP+CMED e Supplier+CMED. Normalização de CNPJ. Cache Redis com graceful degradation. 23 testes. |
| `v1.3.0` | 25/05/2026 | **feat: CMED Conformidade** — Integração preços regulados ANVISA. Endpoints CMED (list, ID, registro, EAN, histórico). JOIN AMPP+CMED com cache Redis. CLI importação XLSX. Index Meilisearch obm_cmed. Redis no Docker Compose. 32 testes. |
| `v1.1.0` | 24/05/2026 | **feat: CLI de importação** — Pipeline completo de importação (ZIP/SQL/MySQL → PostgreSQL + Meilisearch). Conversor refatorado como pacote reutilizável. Validação pós-importação. Metadados. Revisão ortográfica pt-BR. |
| `v1.0.0` | 24/05/2026 | **feat: Release inicial da API OBM** — API REST completa com endpoints para VMP, AMP, VTM, VMPP, AMPP, DCB, Ingredientes, Fornecedores e Domínios. Busca global Meilisearch. Autenticação JWT. Swagger UI. Postman. Docker Compose. Testes. Documentação de usuário. |
