# Centralizador de Versão — API OBM

> **ATENÇÃO**: Este é o arquivo de referência e auditoria da versão atual da API. Todas as builds, tags no git ou release notes devem referenciar a versão travada aqui.

## Versão Atual

```
Versão: 1.1.0
Release Date: 2026-05-24
Timezone: America/Sao_Paulo
Status: Released
```

### v1.1.0 (2026-05-24)

**Status:** Released

**Mudanças:**

- feat: CLI de importação (`cmd/import/`) com pipeline completo — resolução de fonte (ZIP, SQL, gzip, MySQL), conversão MySQL→PostgreSQL, importação via streaming, validação pós-importação, rastreamento de metadados e reindexação do Meilisearch
- feat: Conversor MySQL→PostgreSQL refatorado como pacote reutilizável (`internal/importer/converter.go`)
- feat: Resolvedor de fonte com suporte a ZIP, SQL, SQL.GZ e URI `mysql://` (`internal/importer/source.go`)
- feat: Importador PostgreSQL com streaming via `io.Pipe` e progresso em tempo real (`internal/importer/pgimporter.go`)
- feat: Validador pós-importação com contagem de registros e verificação de integridade referencial (`internal/importer/validator.go`)
- feat: Gerenciador de metadados da importação — tabela `obm_metadata` (`internal/importer/metadata.go`)
- feat: Orquestrador de pipeline com suporte a `--convert-only`, `--reindex-only`, `--skip-index`, `--validate`, `--full` (`internal/importer/pipeline.go`)
- feat: Documentação do projeto (.specs/project/ — PROJECT.md, ROADMAP.md, STATE.md)
- fix: Revisão ortográfica completa em português brasileiro em todas as mensagens do importador e documentação
- fix: Binário `import` adicionado ao `.gitignore`
- refactor: `scripts/convert_sql.go` agora é um invólucro fino chamando `internal/importer/converter.go`

**Arquivos novos:**

- `cmd/import/main.go` — Entry point do CLI de importação
- `internal/importer/converter.go` — Conversor MySQL→PostgreSQL (pacote reutilizável)
- `internal/importer/source.go` — Resolvedor de fonte (ZIP/SQL/gzip/MySQL)
- `internal/importer/pgimporter.go` — Importador PostgreSQL com streaming
- `internal/importer/validator.go` — Validador pós-importação
- `internal/importer/metadata.go` — Gerenciador de metadados (tabela `obm_metadata`)
- `internal/importer/pipeline.go` — Orquestrador do pipeline de importação
- `.specs/project/PROJECT.md` — Visão, objetivos, stack e arquitetura do projeto
- `.specs/project/ROADMAP.md` — Roadmap de funcionalidades e milestones
- `.specs/project/STATE.md` — Memória persistente (decisões, bloqueios, lições)
- `docs/release-notes/v1.1.0.md` — Release notes desta versão

**Arquivos modificados:**

- `scripts/convert_sql.go` — Refatorado para invólucro fino do pacote `internal/importer`
- `README.md` — Seção dedicada "CLI de Importação" + revisão ortográfica pt-BR
- `.gitignore` — Adicionado binário `import`

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
| `v1.1.0` | 24/05/2026 | **feat: CLI de importação** — Pipeline completo de importação (ZIP/SQL/MySQL → PostgreSQL + Meilisearch). Conversor refatorado como pacote reutilizável. Validação pós-importação. Metadados. Revisão ortográfica pt-BR. |
| `v1.0.0` | 24/05/2026 | **feat: Release inicial da API OBM** — API REST completa com endpoints para VMP, AMP, VTM, VMPP, AMPP, DCB, Ingredientes, Fornecedores e Domínios. Busca global Meilisearch. Autenticação JWT. Swagger UI. Postman. Docker Compose. Testes. Documentação de usuário. |
