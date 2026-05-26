# CMED Conformidade — Specification

## Problem Statement

A OBM possui dados ontológicos de medicamentos (VMP/AMP/AMPP) mas não possui dados de preços regulados pelo governo brasileiro. A planilha de Conformidade CMED contém ~25.276 medicamentos com preços PF/PMC, registro sanitário, EAN, tarja e outros dados regulatórios. Hoje não é possível buscar um medicamento e ver seu preço regulado, nem cruzar dados ontológicos com dados de preço.

## Goals

- [ ] Importar dados da planilha CMED para PostgreSQL e Meilisearch
- [ ] Relacionar medicamentos CMED com AMPP da OBM via Registro Sanitário (nu_sanreg)
- [ ] Busca unificada: buscar por nome, registro, EAN, tarja — retornando dados OBM + preços
- [ ] Versionamento: manter histórico de preços por data de referência
- [ ] JOIN AMPP↔CMED com cache Redis para queries combinadas

## Out of Scope

| Feature | Reason |
|---|---|
| Import automático/agendado da planilha CMED | On-demand via CLI |
| Download da planilha do site da CMED | Acesso restrito, arquivo fornecido manualmente |
| Alterar tabelas OBM existentes (tb_ampp, tb_amp, etc.) | Separação de responsabilidades — OBM = ontologia, CMED = preço/regulação |
| Cálculos de reajuste ou comparação automática de preços | Apenas exposição dos dados brutos |
| Frontend/UI | Apenas API REST |

---

## User Stories

### P1: Importar planilha CMED para PostgreSQL ⭐ MVP

**User Story**: Como operador do sistema, quero importar a planilha CMED via CLI para que os dados de preço fiquem disponíveis na API.

**Why P1**: Sem dados importados, nenhuma outra funcionalidade funciona.

**Acceptance Criteria**:

1. WHEN executar `go run cmd/cmed_import/main.go --source planilha.xlsx --date 2025-05-08` THEN sistema SHALL importar todos os registros da planilha para `tb_cmed_conformidade`
2. WHEN `--header-row` não for informado THEN sistema SHALL usar valor default 42
3. WHEN `--header-row N` for informado THEN sistema SHALL ler headers na linha N e dados a partir de N+1
4. WHEN planilha tiver coluna REGISTRO THEN sistema SHALL mapeá-la como `nu_sanreg` (BIGINT)
5. WHEN registro já existe para mesma data de referência THEN sistema SHALL fazer UPSERT (ON CONFLICT DO UPDATE)
6. WHEN import concluir THEN sistema SHALL invalidar todas as chaves de cache `cmed:*` e `ampp_cmed:*` no Redis
7. WHEN `--skip-index=false` (default) THEN sistema SHALL reindexar Meilisearch `obm_cmed` após import
8. WHEN `--skip-index=true` THEN sistema SHALL pular reindexação do Meilisearch
9. WHEN header row não contiver colunas obrigatórias (REGISTRO, PRODUTO, SUBSTÂNCIA) THEN sistema SHALL retornar erro com lista de colunas faltantes

**Independent Test**: Importar planilha, verificar registros no PostgreSQL via `SELECT COUNT(*) FROM tb_cmed_conformidade`.

---

### P2: Buscar medicamentos CMED via API REST

**User Story**: Como consumidor da API, quero buscar medicamentos CMED por nome, registro, EAN ou filtros para obter preços regulados.

**Why P2**: Core search capability — necessário para qualquer uso dos dados importados.

**Acceptance Criteria**:

1. WHEN `GET /api/v1/cmed?nome=Dipirona` THEN sistema SHALL retornar medicamentos CMED com nome contendo "Dipirona"
2. WHEN `GET /api/v1/cmed/registro/1018003900019` THEN sistema SHALL retornar medicamento com nu_sanreg=1018003900019
3. WHEN `GET /api/v1/cmed/ean/7891106000956` THEN sistema SHALL retornar medicamento com EAN correspondente
4. WHEN `GET /api/v1/cmed?tarja=Vermelha` THEN sistema SHALL filtrar por tarja
5. WHEN `GET /api/v1/cmed?dt_referencia=2025-05-08` THEN sistema SHALL filtrar por data de referência
6. WHEN `GET /api/v1/cmed/:id` THEN sistema SHALL retornar medicamento por ID interno incluindo JSONB de preços detalhados
7. WHEN busca não encontrar resultados THEN sistema SHALL retornar `200` com lista vazia `[]`
8. WHEN registro não encontrado THEN sistema SHALL retornar `404`

**Independent Test**: Chamar `GET /api/v1/cmed?nome=ORENCIA` e verificar resposta com dados de preço.

---

### P3: Busca global no Meilisearch incluindo CMED

**User Story**: Como consumidor da API, quero usar a busca global (`GET /api/v1/search?entity=cmed`) para encontrar medicamentos CMED junto com VMPs/AMPs.

**Why P3**: Integração com search existente — melhora descobribilidade.

**Acceptance Criteria**:

1. WHEN `GET /api/v1/search?q=Dipirona&entity=cmed` THEN sistema SHALL retornar hits do index `obm_cmed`
2. WHEN `GET /api/v1/search?q=1018003900019` THEN sistema SHALL buscar no campo `nu_sanreg` do index CMED
3. WHEN entity não informado THEN sistema SHALL incluir `cmed` no default entities (`vmp`, `amp`, `supplier`, `cmed`)
4. WHEN hit do CMED retornado THEN SHALL conter `entity=cmed`, `nome` (no_produto), `codigo` (nu_sanreg), `fabricante` (no_laboratorio)

**Independent Test**: Chamar `GET /api/v1/search?q=ORENCIA&entity=cmed` e verificar hit retornado.

---

### P4: JOIN AMPP↔CMED com cache Redis

**User Story**: Como consumidor da API, quero obter dados completos de um AMPP (ontologia + preços CMED) em uma única chamada.

**Why P4**: Caso de uso principal — ver medicamento OBM com preços governo.

**Acceptance Criteria**:

1. WHEN `GET /api/v1/ampp/:id/cmed` THEN sistema SHALL retornar AMPP + dados CMED vinculados via nu_sanreg
2. WHEN AMPP não possui nu_sanreg THEN sistema SHALL retornar AMPP sem dados CMED (cmed = null)
3. WHEN AMPP possui nu_sanreg mas sem CMED correspondente THEN sistema SHALL retornar AMPP com cmed = null
4. WHEN resposta está em cache Redis THEN sistema SHALL retornar do cache sem query PostgreSQL
5. WHEN resposta não está em cache THEN sistema SHALL buscar no PostgreSQL, cachear com TTL 24h, e retornar
6. WHEN import CMED executa THEN sistema SHALL invalidar chaves `ampp_cmed:*` no Redis

**Independent Test**: Chamar `GET /api/v1/ampp/{id}/cmed` para AMPP com nu_sanreg populado e verificar dados de preço na resposta.

---

### P5: Histórico de preços por versão

**User Story**: Como consumidor da API, quero ver a evolução de preços de um medicamento ao longo do tempo.

**Why P5**: Value-add do versionamento — comparar preços entre versões da planilha.

**Acceptance Criteria**:

1. WHEN `GET /api/v1/cmed/:id/historico` THEN sistema SHALL retornar todas as versões do medicamento ordenadas por dt_referencia DESC
2. WHEN medicamento tem apenas 1 versão THEN sistema SHALL retornar array com 1 elemento
3. WHEN medicamento não existe THEN sistema SHALL retornar 404

**Independent Test**: Importar 2 versões da planilha, chamar historico e verificar array com 2 elementos.

---

## Edge Cases

- WHEN planilha tem REGISTRO vazio ou "-" THEN sistema SHALL registrar nu_sanreg como NULL
- WHEN planilha tem EAN como "-" THEN sistema SHALL registrar como NULL
- WHEN valor de preço está vazio THEN sistema SHALL registrar como NULL (NUMERIC)
- WHEN nu_sanreg na planilha não existe na tabela AMPP THEN registro CMED é criado normalmente (sem vínculo OBM, ainda assim buscável)
- WHEN AMPP tem nu_sanreg mas nenhuma versão CMED corresponde THEN endpoint `/ampp/:id/cmed` retorna cmed=null
- WHEN múltiplos AMPPs compartilham o mesmo nu_sanreg THEN JOIN retorna o AMPP específico + dados CMED (1:1 via nu_sanreg)
- WHEN Redis não está disponível THEN sistema SHALL logar warning e fazer query direta no PostgreSQL (degradação graceful)
- WHEN header row informado não existe na planilha THEN sistema SHALL retornar erro

---

## Requirement Traceability

| Requirement ID | Story | Phase | Status |
|---|---|---|---|
| CMED-01 | P1: Import CLI — parse XLSX | Design | Pending |
| CMED-02 | P1: Import CLI --header-row param | Design | Pending |
| CMED-03 | P1: Import CLI --date param | Design | Pending |
| CMED-04 | P1: Import CLI --skip-index param | Design | Pending |
| CMED-05 | P1: UPSERT on conflict | Design | Pending |
| CMED-06 | P1: Cache invalidation on import | Design | Pending |
| CMED-07 | P1: Header validation | Design | Pending |
| CMED-08 | P2: List CMED with filters | Design | Pending |
| CMED-09 | P2: Get CMED by registro | Design | Pending |
| CMED-10 | P2: Get CMED by EAN | Design | Pending |
| CMED-11 | P2: Get CMED by ID | Design | Pending |
| CMED-12 | P3: Meilisearch index obm_cmed | Design | Pending |
| CMED-13 | P3: Search entity cmed | Design | Pending |
| CMED-14 | P4: JOIN AMPP+CMED endpoint | Design | Pending |
| CMED-15 | P4: Redis cache with TTL | Design | Pending |
| CMED-16 | P4: Graceful degradation without Redis | Design | Pending |
| CMED-17 | P5: Price history endpoint | Design | Pending |
| CMED-18 | Infra: Redis in docker-compose | Design | Pending |
| CMED-19 | Infra: Redis config in config.go | Design | Pending |
| CMED-20 | Infra: go-redis dependency | Design | Pending |
| CMED-21 | Infra: excelize dependency | Design | Pending |
| CMED-22 | DB: Migration tb_cmed_conformidade | Design | Pending |
| CMED-23 | DB: Entity CMEDConformidade | Design | Pending |

**Coverage**: 23 total, 0 mapped to tasks, 23 unmapped

---

## Success Criteria

- [ ] Importar planilha CMED completa (~25K registros) em < 2 minutos
- [ ] Busca por REGISTRO retorna resultado em < 100ms
- [ ] JOIN AMPP+CMED com cache em < 50ms, sem cache em < 200ms
- [ ] Zero alterações nas tabelas OBM existentes
- [ ] Redis graceful degradation — API funciona sem Redis (sem cache, com query direta)
