# CMED Extended Joins Specification

## Problem Statement

Atualmente a CMED Conformidade só se relaciona com AMPP via `nu_sanreg`. Entidades como AMP e Supplier possuem chaves exatas (`nu_nreg` e `nu_cnpj`) que permitem JOIN direto com CMED, mas esses relacionamentos não estão expostos na API. Isso obriga o consumidor a fazer múltiplas chamadas para obter preço regulado a partir de um AMP ou listar todos os produtos CMED de um fornecedor.

## Goals

- [ ] Excluir endpoint AMP+CMED via `AMP.nu_nreg = CMED.nu_sanreg` (registro único + VMP pai)
- [ ] Excluir endpoint Supplier+CMED via `Supplier.nu_cnpj = CMED.nu_cnpj` (lista de produtos)
- [ ] Normalização de CNPJ para garantir JOIN confiável
- [ ] Cache Redis com graceful degradation (mesmo padrão do AMPP+CMED)
- [ ] Testes unitários e de handler para os novos endpoints
- [ ] Documentação atualizada (README, VERSION, release notes)

## Out of Scope

| Feature | Reason |
| --- | --- |
| Joins fuzzy (VTM, DCB, Ingredient) | Requer normalização de texto e matching aproximado — escopo futuro |
| AMPP+CMED via EAN | 9 combinações cruzadas — escopo futuro |
| Paginação no endpoint Supplier+CMED | Retorno limitado por dt_referencia + st_registro_ativo |
| Novas migrações SQL | Tabela CMED já existe, nenhuma alteração de schema necessária |

---

## User Stories

### P1: AMP+CMED Join ⭐ MVP

**User Story**: Como consumidor da API, quero consultar o preço CMED diretamente a partir de um AMP, para obter dados ontológicos e preço regulado em uma única chamada.

**Why P1**: AMP é a entidade central do modelo dm+d. `AMP.nu_nreg` e `CMED.nu_sanreg` representam o mesmo Registro Sanitário ANVISA — é um JOIN exato e confiável.

**Acceptance Criteria**:

1. WHEN `GET /api/v1/amp/:id/cmed` com AMP válido e `nu_nreg` preenchido THEN system SHALL retornar AMP + VMP pai + CMED (registro único ou mais recente)
2. WHEN `GET /api/v1/amp/:id/cmed` com AMP válido e `nu_nreg` vazio THEN system SHALL retornar AMP + VMP pai com `cmed: null`
3. WHEN `GET /api/v1/amp/:id/cmed` com AMP inexistente THEN system SHALL retornar 404
4. WHEN `GET /api/v1/amp/:id/cmed?dt_referencia=2026-05-01` THEN system SHALL retornar CMED da data específica
5. WHEN CMED lookup falhar THEN system SHALL retornar AMP + VMP com `cmed: null` (graceful degradation)

**Independent Test**: `curl /api/v1/amp/1/cmed` retorna AMP + VMP + preço CMED em uma resposta.

---

### P2: Supplier+CMED Join

**User Story**: Como consumidor da API, quero listar todos os produtos CMED de um fornecedor, para ver quais medicamentos de um laboratório têm preço regulado.

**Why P2**: `Supplier.nu_cnpj` e `CMED.nu_cnpj` são o mesmo CNPJ — JOIN exato. Um fornecedor tem N produtos CMED.

**Acceptance Criteria**:

1. WHEN `GET /api/v1/suppliers/:id/cmed` com Supplier válido e CNPJ preenchido THEN system SHALL retornar Supplier + lista de CMED (filtrado por dt_referencia ou mais recente)
2. WHEN `GET /api/v1/suppliers/:id/cmed` com Supplier válido e CNPJ vazio THEN system SHALL retornar Supplier com `cmed: []`
3. WHEN `GET /api/v1/suppliers/:id/cmed` com Supplier inexistente THEN system SHALL retornar 404
4. WHEN `GET /api/v1/suppliers/:id/cmed?dt_referencia=2026-05-01` THEN system SHALL retornar apenas CMED da data específica
5. WHEN CMED lookup falhar THEN system SHALL retornar Supplier com `cmed: []` (graceful degradation)
6. WHEN CNPJ tem formatação diferente (pontuação) THEN system SHALL normalizar removendo não-dígitos antes do JOIN

**Independent Test**: `curl /api/v1/suppliers/1/cmed` retorna Supplier + lista de produtos CMED.

---

## Edge Cases

- WHEN `AMP.nu_nreg` é NULL ou zero THEN system SHALL retornar `cmed: null` sem erro
- WHEN `Supplier.nu_cnpj` é NULL ou string vazia THEN system SHALL retornar `cmed: []` sem erro
- WHEN CNPJ em Supplier tem formato "00.000.000/0000-00" e CMED tem "00000000000000" THEN JOIN deve funcionar via normalização
- WHEN Redis está indisponível THEN cache é ignorado e query vai direto ao PostgreSQL
- WHEN `dt_referencia` fornecida não existe THEN CMED retorna null/[] (sem erro)

---

## Requirement Traceability

| ID | Story | Phase | Status |
| --- | --- | --- | --- |
| AMP-01 | P1: AMP+CMED | Execute | Verified |
| AMP-02 | P1: AMP+CMED | Execute | Verified |
| AMP-03 | P1: AMP+CMED | Execute | Verified |
| AMP-04 | P1: AMP+CMED | Execute | Verified |
| AMP-05 | P1: AMP+CMED | Execute | Verified |
| SUP-01 | P2: Supplier+CMED | Execute | Verified |
| SUP-02 | P2: Supplier+CMED | Execute | Verified |
| SUP-03 | P2: Supplier+CMED | Execute | Verified |
| SUP-04 | P2: Supplier+CMED | Execute | Verified |
| SUP-05 | P2: Supplier+CMED | Execute | Verified |
| SUP-06 | P2: Supplier+CMED | Execute | Verified |
| CNPJ-01 | Cross-cutting | Execute | Verified |
| CACHE-01 | Cross-cutting | Execute | Verified |
| DOC-01 | Cross-cutting | Execute | Verified |

**Coverage**: 14 total, 14 mapped to tasks, 0 unmapped

---

## Success Criteria

- [x] Endpoint `GET /api/v1/amp/:id/cmed` retorna AMP + VMP + CMED
- [x] Endpoint `GET /api/v1/suppliers/:id/cmed` retorna Supplier + lista CMED
- [x] Normalização CNPJ garante JOIN independente de formatação
- [x] Cache Redis funciona com graceful degradation
- [x] Testes unitários passam (amp_cmed, supplier_cmed, handler)
- [x] Build e testes existentes passam (0 regressões)
- [x] README, VERSION.md e release notes atualizados
