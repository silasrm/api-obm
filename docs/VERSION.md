# Centralizador de Versao - API OBM

> **ATENCAO**: Este e o arquivo de referencia e auditoria da versao atual da API. Todas as builds, tags no git ou release notes devem referenciar a versao travada aqui.

## Versao Atual

```
Versao: 1.0.0
Release Date: 2026-05-24
Timezone: America/Sao_Paulo
Status: Released
```

### v1.0.0 (2026-05-24)

**Status:** Released

**Mudancas:**

- feat: API REST da Ontologia Brasileira de Medicamentos (OBM) — modelo dm+d do NHS adaptado para o Brasil
- feat: Endpoints de consulta para VMP, AMP, VTM, VMPP, AMPP, DCB, Ingredientes, Fornecedores e Dominios
- feat: Busca global via Meilisearch com filtros por nome, codigo, fabricante, descricao e status ativo
- feat: Endpoints de detalhe para VMP (com VTM, dominios, ingredientes, rotas, formas, classes ATC) e AMP (com VMP, fornecedor, ingredientes, rotas)
- feat: 32 tipos de dominio (forma farmaceutica, via, classe ATC, categoria de controle, etc.)
- feat: Autenticacao JWT (HS256) com login/geracao de token
- feat: Paginacao cursor-based em todos os endpoints de listagem
- feat: Health check publico (PostgreSQL + Meilisearch)
- feat: Reindexacao administrativa do Meilisearch
- feat: Swagger UI integrado (`/swagger/index.html`)
- feat: Collection Postman com ambiente local
- feat: Infraestrutura Docker Compose (PostgreSQL 16, Meilisearch v1.8, API)
- feat: Seed de usuarios iniciais (admin/admin123, viewer/viewer123)
- feat: Suite de testes unitarios e de integracao
- fix: Correcao da descricao de "Observatorio de Medicamentos" para "Ontologia Brasileira de Medicamentos" em todos os arquivos
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
- `internal/domain/entity/entities.go` — Entidades de dominio
- `internal/domain/repository/interfaces.go` — Contratos dos repositorios
- `internal/infrastructure/persistence/postgres/` — Repositorios PostgreSQL
- `internal/infrastructure/persistence/meilisearch/` — Cliente e indexer Meilisearch
- `internal/infrastructure/config/config.go` — Carregamento de configuracao
- `internal/interface/http/handler/` — Handlers HTTP
- `internal/interface/http/dto/dto.go` — DTOs de request/response
- `internal/interface/http/middleware/auth.go` — Middleware JWT
- `internal/interface/http/router/router.go` — Definicao de rotas
- `internal/usecase/` — Casos de uso
- `docs/swagger.yaml`, `docs/swagger.json`, `docs/docs.go` — Documentacao Swagger
- `postman/OBM_API.postman_collection.json` — Collection Postman
- `postman/OBM_API_Local.postman_environment.json` — Ambiente Postman
- `scripts/seed_users.go` — Seed de usuarios
- `scripts/gen_postman.go` — Gerador de collection Postman
- `scripts/convert_sql.go` — Conversor MySQL → PostgreSQL
- `migrations/postgres/001_obm_schema.sql` — Schema completo do banco
- `README.md` — Documentacao para usuarios
- `TESTING.md` — Guia de teste

## Como Incrementar a Versao (SemVer)

1. **Atualize este documento** aumentando a trilha conforme o padrao de versoes semanticas:
   - **Patch** (ex: 1.0.0 → 1.0.1): Correcoes de bugs, refactorings, limpezas que nao entregam novos comportamentos
   - **Minor** (ex: 1.0.0 → 1.1.0): Novos endpoints, filtros, funcionalidades retrocompativeis
   - **Major** (ex: 1.0.0 → 2.0.0): Breaking changes na API (mudanca de formato de resposta, remocao de endpoints, etc.)

2. **Crie a Release Note** em `docs/release-notes/vX.Y.Z.md` descrevendo as mudancas

3. **Crie a tag git** com `git tag -a vX.Y.Z -m "Release vX.Y.Z"`

## Historico de Versoes

| Versao | Data | Descricao |
|--------|------|-----------|
| `v1.0.0` | 24/05/2026 | **feat: Release inicial da API OBM** — API REST completa com endpoints para VMP, AMP, VTM, VMPP, AMPP, DCB, Ingredientes, Fornecedores e Dominios. Busca global Meilisearch. Autenticacao JWT. Swagger UI. Postman. Docker Compose. Testes. Documentacao de usuario. |
