# OBM API — Ontologia Brasileira de Medicamentos

API REST para consulta de dados da [Ontologia Brasileira de Medicamentos (OBM)](https://portal-obm.saude.gov.br/), seguindo o modelo **dm+d** (Dictionary of Medicine and Devices) do NHS adaptado para o Brasil.

---

## Sumario

- [Sobre a OBM](#sobre-a-obm)
- [Conceitos da Ontologia](#conceitos-da-ontologia)
- [Autenticacao](#autenticacao)
- [Busca Global](#busca-global)
- [Referencia de Endpoints](#referencia-de-endpoints)
  - [VMP — Virtual Medicinal Product](#vmp--virtual-medicinal-product)
  - [AMP — Actual Medicinal Product](#amp--actual-medicinal-product)
  - [VTM — Virtual Therapeutic Moiety](#vtm--virtual-therapeutic-moiety)
  - [VMPP — Virtual Medicinal Product Pack](#vmpp--virtual-medicinal-product-pack)
  - [AMPP — Actual Medicinal Product Pack](#ampp--actual-medicinal-product-pack)
  - [DCB — Denominacao Comum Brasileira](#dcb--denominacao-comum-brasileira)
  - [Ingredientes](#ingredientes)
  - [Fornecedores](#fornecedores)
  - [Dominios](#dominios)
  - [Admin](#admin)
- [Paginacao](#paginacao)
- [Codigos de Erro](#codigos-de-erro)
- [Instalacao Local](#instalacao-local)
- [Exemplos Praticos de Uso](#exemplos-praticos-de-uso)

---

## Sobre a OBM

A **Ontologia Brasileira de Medicamentos (OBM)** e um padrao nacional de base de medicamentos para utilizacao em sistemas de prescricao e dispensacao eletronicas, instituida pela [Portaria GM/MS No 6.093, de 16 de dezembro de 2024](https://www.in.gov.br/en/web/dou/-/portaria-gm/ms-n-6.093-de-16-de-dezembro-de-2024-602264704).

Seus objetivos principais sao:

- **Integrar e padronizar** dados de diferentes sistemas de informacoes em saude
- **Normalizar registros** de prescricoes e dispensacoes
- **Promover a interoperabilidade** por meio da Rede Nacional de Dados em Saude (RNDS)
- **Potencializar a seguranca do paciente** por meio de identificacao univoca e inequivoca de medicamentos
- **Seguir praticas internacionais** para descricao e categorizacao de medicamentos

A estrutura da OBM esta baseada no modelo **dm+d** (Dictionary of Medicine and Devices) do **NHS** (National Health Service) do Reino Unido. Trata-se de dado publico, atualizado, acessivel, processavel por maquina, em formato nao proprietario, livre de licencas e com rastreabilidade das modificacoes via versionamento.

**Portal oficial:** [https://portal-obm.saude.gov.br/](https://portal-obm.saude.gov.br/)

---

## Conceitos da Ontologia

A OBM organiza os medicamentos em uma hierarquia inspirada no dm+d, com cinco niveis principais de abstracao:

```
VTM (Virtual Therapeutic Moiety)
 └── VMP (Virtual Medicinal Product)
      ├── VMPP (Virtual Medicinal Product Pack)
      └── AMP (Actual Medicinal Product)
           └── AMPP (Actual Medicinal Product Pack)
```

| Conceito | Sigla | O que representa | Exemplo |
|----------|-------|------------------|---------|
| **Virtual Therapeutic Moiety** | VTM | Principio ativo generico, sem forma farmaceutica | Paracetamol |
| **Virtual Medicinal Product** | VMP | Principio ativo + forma farmaceutica + dose | Paracetamol 500mg comprimido |
| **Virtual Medicinal Product Pack** | VMPP | Apresentacao virtual do VMP (quantidade) | Paracetamol 500mg comprimido — 20 comprimidos |
| **Actual Medicinal Product** | AMP | Produto comercial de um fabricante | Paracetamol 500mg comprimido — Medley |
| **Actual Medicinal Product Pack** | AMPP | Embalagem comercial do AMP (com codigo EAN) | Paracetamol 500mg comprimido — Medley — caixa 20 |

### Entidades complementares

| Entidade | Descricao |
|----------|-----------|
| **DCB** (Denominacao Comum Brasileira) | Denominacao oficial de substancias ativas conforme ANVISA |
| **Ingredient Substance** | Substancia ativa que compoe um medicamento (com codigo CAS e DCB) |
| **Supplier** (Fornecedor) | Fabricante ou detentor do registro sanitario do AMP |
| **Domain** | Tabelas de dominio/classificacao (forma farmaceutica, via, classe ATC, etc.) |

### Relacionamentos principais

- Um **VTM** possui varios **VMPs** (diferentes formas/doses do mesmo principio ativo)
- Um **VMP** pertence a um **VTM**
- Um **VMP** possui varios **AMPs** (produtos de diferentes fabricantes)
- Um **AMP** pertence a um **VMP** e a um **Supplier**
- Um **VMP** possui varios **VMPPs** (apresentacoes)
- Um **AMP** possui varios **AMPPs** (embalagens comerciais com EAN)
- **VMPs** e **AMPs** possuem **Ingredientes** com concentracao
- **VMPs** e **AMPs** estao ligados a **Dominios** (forma farmaceutica, via, classe ATC, etc.)

---

## Autenticacao

Todos os endpoints de dados (sob `/api/v1/`) exigem autenticacao via **JWT Bearer Token**. Apenas `/auth/login` e `/health` sao publicos.

### 1. Obter o token

```bash
curl -s -X POST http://localhost:8094/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

Resposta:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 86400
}
```

### 2. Usar o token nas requisicoes

Inclua o header `Authorization: Bearer <token>` em toda requisicao protegida:

```bash
TOKEN="eyJhbGciOiJIUzI1NiIs..."

curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmp?limit=5"
```

### 3. Expiracao

O token e valido por **24 horas** (configuravel via `JWT_EXPIRATION_HOURS`). Apos expirar, faca login novamente.

---

## Busca Global

O endpoint de busca utiliza o [Meilisearch](https://www.meilisearch.com/) para busca full-text em VMPs, AMPs e Fornecedores.

```
GET /api/v1/search
```

| Parametro | Tipo | Obrigatorio | Descricao |
|-----------|------|-------------|-----------|
| `q` | string | Sim | Termo de busca |
| `entity` | string | Nao | Entidades: `vmp`, `amp`, `supplier`. Separadas por virgula. Padrao: todas |
| `limit` | int | Nao | Limite de resultados (padrao: 20, max: 100) |
| `cursor` | string | Nao | Cursor de paginacao |
| `filter[nome]` | string | Nao | Filtro por nome |
| `filter[codigo]` | string | Nao | Filtro por codigo |
| `filter[fabricante]` | string | Nao | Filtro por fabricante |
| `filter[descricao]` | string | Nao | Filtro por descricao |
| `filter[ativo]` | string | Nao | Filtro por status ativo |

### Exemplos

**Buscar paracetamol em todas as entidades:**

```bash
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/search?q=paracetamol&limit=5"
```

**Buscar apenas em VMPs e AMPs:**

```bash
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/search?q=dipirona&entity=vmp,amp&limit=5"
```

**Buscar com filtro por fabricante:**

```bash
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/search?q=paracetamol&filter[fabricante]=medley&limit=5"
```

### Resposta

```json
{
  "query": "paracetamol",
  "entities": ["vmp", "amp", "supplier"],
  "hits": [
    {
      "entity": "vmp",
      "id": 1234,
      "nome": "Paracetamol 500mg comprimido",
      "codigo": "3456789",
      "fabricante": "",
      "descricao": "",
      "relevance_score": 9.5
    },
    {
      "entity": "amp",
      "id": 5678,
      "nome": "Paracetamol 500mg comprimido revestido",
      "codigo": "9876543",
      "fabricante": "MEDLEY",
      "descricao": "Paracetamol 500mg...",
      "relevance_score": 8.2
    }
  ],
  "cursor": "MjA=",
  "limit": 5,
  "total": 42
}
```

---

## Referencia de Endpoints

### VMP — Virtual Medicinal Product

| Metodo | Rota | Descricao |
|--------|------|-----------|
| `GET` | `/api/v1/vmp` | Listar VMPs com filtros e paginacao |
| `GET` | `/api/v1/vmp/:id` | Obter VMP por ID |
| `GET` | `/api/v1/vmp/:id/detail` | Obter VMP com detalhes completos |

**Parametros de listagem (`GET /api/v1/vmp`):**

| Parametro | Tipo | Descricao |
|-----------|------|-----------|
| `limit` | int | Limite por pagina (padrao: 20, max: 100) |
| `cursor` | string | Cursor de paginacao |
| `nome` | string | Filtro por nome (busca parcial, case-insensitive) |
| `codigo` | string | Filtro por codigo NU_VPID (busca exata) |
| `ativo` | boolean | Filtro por status ativo (padrao: `true`) |

**Exemplos:**

```bash
# Listar VMPs
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmp?limit=5"

# Filtrar por nome
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmp?nome=paracetamol&limit=5"

# Buscar por codigo
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmp?codigo=3456789"

# VMP por ID
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmp/1234"

# VMP detalhado (VTM, ingredientes, rotas, formas, classes ATC, etc.)
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmp/1234/detail"
```

**Resposta do detail:**

```json
{
  "vmp": {
    "co_seq_id": 1234,
    "nu_vpid": "3456789",
    "no_nm": "Paracetamol 500mg comprimido",
    "co_vtmid": 100,
    "st_registro_ativo": "ACTIVE",
    "nu_snomed": "..."
  },
  "vtm": {
    "co_seq_id": 100,
    "no_nm": "Paracetamol",
    "nu_vtmid": 1001
  },
  "basis_of_name": { "no_descr": "Substance" },
  "pres_status": { "no_descr": "Valid as a prescribable product" },
  "anvs_class": { "no_descr": "Isento de prescricao" },
  "forms": [
    { "no_descr": "Tablet" }
  ],
  "routes": [
    { "no_descr": "Oral" }
  ],
  "atc_classes": [
    { "no_descr": "N02BE01" }
  ],
  "ingredients": [
    {
      "co_iscd": 500,
      "qt_strnt_nmrtr_val": 500,
      "co_strnt_nmrtr_uomcd": 30,
      "ingredient_substance": {
        "co_seq_id": 500,
        "no_nm": "Paracetamol",
        "nu_isid": "SUB12345"
      }
    }
  ],
  "control_drug_infos": [],
  "catmats": [],
  "renames_br": [],
  "local_aplicacao": []
}
```

---

### AMP — Actual Medicinal Product

| Metodo | Rota | Descricao |
|--------|------|-----------|
| `GET` | `/api/v1/amp` | Listar AMPs com filtros e paginacao |
| `GET` | `/api/v1/amp/:id` | Obter AMP por ID |
| `GET` | `/api/v1/amp/:id/detail` | Obter AMP com detalhes completos |

**Parametros de listagem (`GET /api/v1/amp`):**

| Parametro | Tipo | Descricao |
|-----------|------|-----------|
| `limit` | int | Limite por pagina (padrao: 20, max: 100) |
| `cursor` | string | Cursor de paginacao |
| `nome` | string | Filtro por nome (busca parcial) |
| `codigo` | string | Filtro por codigo NU_APID (busca exata) |
| `fabricante` | string | Filtro por nome do fabricante (busca parcial) |
| `ativo` | boolean | Filtro por status ativo (padrao: `true`) |

**Exemplos:**

```bash
# Listar AMPs
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/amp?limit=5"

# Filtrar por fabricante
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/amp?fabricante=medley&limit=5"

# Filtrar por nome e fabricante
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/amp?nome=paracetamol&fabricante=medley"

# AMP por ID
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/amp/5678"

# AMP detalhado (VMP, fornecedor, ingredientes, rotas, etc.)
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/amp/5678/detail"
```

**Resposta do detail:**

```json
{
  "amp": {
    "co_seq_id": 5678,
    "nu_apid": "9876543",
    "no_nm": "Paracetamol 500mg comprimido revestido — Medley",
    "co_vpid": 1234,
    "co_suppcd": 200,
    "nu_nreg": 12345,
    "st_registro_ativo": "ACTIVE"
  },
  "vmp": {
    "co_seq_id": 1234,
    "no_nm": "Paracetamol 500mg comprimido",
    "nu_vpid": "3456789"
  },
  "supplier": {
    "co_seq_id": 200,
    "no_descr": "MEDLEY INDUSTRIA FARMACEUTICA",
    "nu_cnpj": "12.345.678/0001-90"
  },
  "lic_auth": { "no_descr": "ANVISA" },
  "med_class": { "no_descr": "Isento de prescricao" },
  "avail_restriction": { "no_descr": "None" },
  "routes": [
    { "no_descr": "Oral" }
  ],
  "preserv_conds": [],
  "ingredients": [
    {
      "co_isid": 500,
      "qt_strnth": 500,
      "ingredient_substance": {
        "co_seq_id": 500,
        "no_nm": "Paracetamol",
        "nu_isid": "SUB12345"
      }
    }
  ]
}
```

---

### VTM — Virtual Therapeutic Moiety

| Metodo | Rota | Descricao |
|--------|------|-----------|
| `GET` | `/api/v1/vtm` | Listar VTMs com filtros e paginacao |
| `GET` | `/api/v1/vtm/:id` | Obter VTM por ID |

**Parametros de listagem:**

| Parametro | Tipo | Descricao |
|-----------|------|-----------|
| `limit` | int | Limite por pagina (padrao: 20, max: 100) |
| `cursor` | string | Cursor de paginacao |
| `nome` | string | Filtro por nome (busca parcial) |
| `ativo` | boolean | Filtro por status ativo (padrao: `true`) |

**Exemplos:**

```bash
# Listar VTMs
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vtm?limit=5"

# Filtrar por nome
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vtm?nome=paracetamol"

# VTM por ID
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vtm/100"
```

---

### VMPP — Virtual Medicinal Product Pack

| Metodo | Rota | Descricao |
|--------|------|-----------|
| `GET` | `/api/v1/vmpp` | Listar VMPPs com filtros e paginacao |
| `GET` | `/api/v1/vmpp/:id` | Obter VMPP por ID |

**Parametros de listagem:** `limit`, `cursor`, `nome` (busca parcial), `ativo`

```bash
# Listar VMPPs
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmpp?limit=5"

# VMPP por ID
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmpp/1"
```

---

### AMPP — Actual Medicinal Product Pack

| Metodo | Rota | Descricao |
|--------|------|-----------|
| `GET` | `/api/v1/ampp` | Listar AMPPs com filtros e paginacao |
| `GET` | `/api/v1/ampp/:id` | Obter AMPP por ID |

**Parametros de listagem:** `limit`, `cursor`, `nome` (busca parcial), `ativo`

```bash
# Listar AMPPs
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/ampp?limit=5"

# AMPP por ID
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/ampp/1"
```

**Resposta AMPP (inclui codigos EAN):**

```json
{
  "co_seq_id": 1,
  "nu_appid": "AMP123456",
  "no_nm": "Paracetamol 500mg comprimido — 20 comprimidos",
  "nu_ean13a": "7891234567890",
  "nu_ean13b": "",
  "nu_ean13c": "",
  "st_registro_ativo": "ACTIVE"
}
```

---

### DCB — Denominacao Comum Brasileira

| Metodo | Rota | Descricao |
|--------|------|-----------|
| `GET` | `/api/v1/dcb` | Listar DCBs com filtros e paginacao |
| `GET` | `/api/v1/dcb/:id` | Obter DCB por ID |

**Parametros de listagem:** `limit`, `cursor`, `nome` (busca parcial em DS_DCB), `ativo`

```bash
# Listar DCBs
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/dcb?limit=5"

# Filtrar por nome
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/dcb?nome=paracetamol"

# DCB por ID
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/dcb/1"
```

---

### Ingredientes

| Metodo | Rota | Descricao |
|--------|------|-----------|
| `GET` | `/api/v1/ingredients` | Listar Ingredient Substances com filtros e paginacao |
| `GET` | `/api/v1/ingredients/:id` | Obter ingrediente por ID |

**Parametros de listagem:** `limit`, `cursor`, `nome` (busca parcial), `ativo`

```bash
# Listar ingredientes
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/ingredients?limit=5"

# Filtrar por nome
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/ingredients?nome=paracetamol"

# Ingrediente por ID
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/ingredients/500"
```

---

### Fornecedores

| Metodo | Rota | Descricao |
|--------|------|-----------|
| `GET` | `/api/v1/suppliers` | Listar fornecedores com filtros e paginacao |
| `GET` | `/api/v1/suppliers/:id` | Obter fornecedor por ID |

**Parametros de listagem:**

| Parametro | Tipo | Descricao |
|-----------|------|-----------|
| `limit` | int | Limite por pagina (padrao: 20, max: 100) |
| `cursor` | string | Cursor de paginacao |
| `nome` | string | Filtro por nome (busca parcial) |
| `codigo` | string | Filtro por codigo NU_CD (busca exata) |
| `ativo` | boolean | Filtro por status ativo (padrao: `true`) |

```bash
# Listar fornecedores
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/suppliers?limit=5"

# Filtrar por nome
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/suppliers?nome=medley"

# Fornecedor por ID
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/suppliers/200"
```

---

### Dominios

Dominios sao tabelas de classificacao e referencia que categorizam os medicamentos (forma farmaceutica, via de administracao, classe ATC, categoria de controle, etc.).

| Metodo | Rota | Descricao |
|--------|------|-----------|
| `GET` | `/api/v1/domains/:domain` | Listar dominios por tipo |
| `GET` | `/api/v1/domains/:domain/:id` | Obter dominio por tipo e ID |

**Parametros de listagem:** `limit`, `cursor`, `nome` (busca parcial em NO_DESCR e NO_DESCR_PT_BR), `ativo`

**Tipos de dominio disponiveis:**

| Slug | Tabela | Descricao |
|------|--------|-----------|
| `form` | `td_form` | Forma farmaceutica |
| `route` | `td_route` | Via de administracao |
| `flavour` | `td_flavour` | Sabor |
| `legal-category` | `td_legal_category` | Categoria legal |
| `licensing-authority` | `td_licensing_authority` | Autoridade de licenciamento |
| `avail-restriction` | `td_availability_restriction` | Restricao de disponibilidade |
| `med-class` | `td_med_class_br` | Classe medicamentosa BR |
| `anvs-class` | `td_anvs_class_br` | Classe ANVISA |
| `atc-class` | `td_atc_class_br` | Classe ATC BR |
| `control-drug` | `td_control_drug_category` | Categoria de controle de drogas |
| `df-indicator` | `td_df_indicator` | Indicador DF |
| `discontinued-ind` | `td_discontinued_ind` | Indicador de descontinuacao |
| `pres-status` | `td_virtual_product_pres_status` | Status de prescricao |
| `non-avail` | `td_virtual_product_non_avail` | Motivo de indisponibilidade |
| `basis-of-name` | `td_basis_of_name` | Base do nome |
| `basis-of-strnth` | `td_basis_of_strnth` | Base de concentracao |
| `brimunologico` | `td_brimunologico` | Imunologico BR |
| `catmat` | `td_catmat_br` | CATMAT BR |
| `country` | `td_country` | Pais |
| `unit-of-measure` | `td_unit_of_measure` | Unidade de medida |
| `package` | `td_package` | Tipo de embalagem |
| `ont-form-route` | `td_ont_form_route` | Forma/Rota OBM |
| `preserv-cond` | `td_preserv_cond_br` | Condicao de preservacao BR |
| `rename-comp` | `td_rename_comp_br` | RENAME Complementar BR |
| `ingredient-source` | `td_ingredient_source_br` | Fonte do ingrediente BR |
| `healthcare-prof` | `td_healthcare_prof_br` | Profissional de saude BR |
| `indicacao-farmpop` | `td_indicacao_farmpop_br` | Indicacao Farmacia Popular BR |
| `monitoring-reason` | `td_monitoring_reason_br` | Motivo de monitoramento BR |
| `name-change-reason` | `td_name_change_reason` | Motivo de mudanca de nome |
| `lic-auth-change-reason` | `td_lic_auth_change_reason` | Motivo de mudanca de licenca |
| `phpid` | `td_phpid` | PHPID |
| `local-aplicacao` | `td_local_aplicacao` | Local de aplicacao |

```bash
# Listar formas farmaceuticas
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/domains/form?limit=5"

# Listar vias de administracao
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/domains/route?limit=5"

# Listar classes ATC
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/domains/atc-class?limit=5"

# Dominio especifico por ID
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/domains/form/1"
```

---

### Admin

| Metodo | Rota | Descricao |
|--------|------|-----------|
| `GET` | `/health` | Health check (sem autenticacao) |
| `POST` | `/api/v1/admin/reindex` | Reindexar Meilisearch |

**Health check (publico):**

```bash
curl -s http://localhost:8094/health
```

```json
{
  "status": "ok",
  "postgres": "ok",
  "meilisearch": "ok"
}
```

> Status `"degraded"` indica que PostgreSQL ou Meilisearch esta indisponivel.

**Reindexar Meilisearch:**

```bash
curl -s -X POST -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/admin/reindex"
```

```json
{
  "status": "completed",
  "indexed": {
    "vmp": 12345,
    "amp": 67890,
    "supplier": 500
  }
}
```

---

## Paginacao

Todos os endpoints de listagem utilizam **paginacao baseada em cursor**. Nao ha numeracao de paginas — o cursor codifica o offset da proxima pagina.

### Como funciona

1. Faca a primeira requisicao sem cursor
2. A resposta contem o campo `cursor` com o valor para a proxima pagina
3. Use esse valor no parametro `cursor` da proxima requisicao
4. Quando `cursor` for vazio (`""`), nao ha mais paginas

### Parametros

| Parametro | Descricao | Padrao | Maximo |
|-----------|-----------|--------|--------|
| `limit` | Registros por pagina | 20 | 100 |
| `cursor` | Cursor da proxima pagina | — | — |

### Exemplo

```bash
# Primeira pagina
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmp?limit=3"
```

```json
{
  "items": [...],
  "cursor": "Mw==",
  "limit": 3,
  "total": 15000
}
```

```bash
# Proxima pagina (usar o cursor retornado)
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmp?limit=3&cursor=Mw=="
```

```json
{
  "items": [...],
  "cursor": "Ng==",
  "limit": 3,
  "total": 15000
}
```

---

## Codigos de Erro

| HTTP | Mensagem | Condicao |
|------|----------|----------|
| `400` | `invalid request body` | JSON malformado no login |
| `400` | `invalid query parameters` | Parametros de busca invalidos |
| `400` | `invalid id` | Parametro `:id` nao numerico |
| `401` | `missing authorization header` | Header `Authorization` ausente |
| `401` | `invalid token` | Token JWT invalido ou expirado |
| `401` | `invalid token claims` | Claims do JWT nao extraiveis |
| `401` | `invalid credentials` | Usuario/senha incorretos ou usuario inativo |
| `404` | `not found` | Registro nao encontrado |
| `500` | (mensagem variavel) | Erro interno (banco, Meilisearch, etc.) |

### Formato da resposta de erro

```json
{
  "error": "descricao do erro",
  "code": 401
}
```

---

## Instalacao Local

### Pre-requisitos

- [Docker](https://www.docker.com/) e Docker Compose
- [Go 1.25+](https://go.dev/dl/)
- `curl` ou Postman/Insomnia
- `jq` (opcional, para formatar JSON)

### 1. Subir infraestrutura

```bash
docker compose up postgres meilisearch -d
```

Aguarde os healthchecks passarem (~30s). O PostgreSQL estara na porta **5433** e o Meilisearch na porta **7701**.

> **Atencao:** o arquivo `migrations/postgres/001_obm_schema.sql` tem ~1.1 milhao de linhas. O carregamento inicial pode levar varios minutos. Monitore com:
> ```bash
> docker compose logs -f postgres
> ```

Verifique se os servicos estao saudaveis:

```bash
docker compose ps
```

Ambos devem mostrar `(healthy)` na coluna STATUS.

### 2. Configurar o .env

```bash
cp .env.example .env
```

O `.env` ja vem com as portas corretas. Variaveis disponiveis:

| Variavel | Padrao | Descricao |
|----------|--------|-----------|
| `PG_HOST` | `localhost` | Host do PostgreSQL |
| `PG_PORT` | `5433` | Porta do PostgreSQL (5432 dentro do Docker) |
| `PG_USER` | `obm` | Usuario do PostgreSQL |
| `PG_PASSWORD` | `obm123` | Senha do PostgreSQL |
| `PG_DATABASE` | `dbportalobm` | Nome do banco |
| `PG_SSLMODE` | `disable` | SSL mode do PostgreSQL |
| `MEILI_URL` | `http://localhost:7701` | URL do Meilisearch (7700 dentro do Docker) |
| `MEILI_API_KEY` | `obm-meili-master-key` | Chave API do Meilisearch |
| `MEILI_INDEX_PREFIX` | `obm_` | Prefixo dos indices (produz `obm_vmp`, `obm_amp`, `obm_supplier`) |
| `JWT_SECRET` | `obm-secret-key-change-in-prod` | Chave de assinatura JWT (altere em producao!) |
| `JWT_EXPIRATION_HOURS` | `24` | Expiracao do token em horas |
| `SERVER_PORT` | `8094` | Porta do servidor API |
| `GIN_MODE` | `release` | Modo do Gin (`debug` ou `release`) |
| `SYNC_ON_STARTUP` | `true` | Reindexar Meilisearch ao iniciar |

### 3. Criar usuarios iniciais

```bash
go run scripts/seed_users.go
```

Cria os usuarios padrao:

| Usuario | Senha | Ativo |
|---------|-------|-------|
| `admin` | `admin123` | sim |
| `viewer` | `viewer123` | sim |

### 4. Rodar a API

```bash
go run ./cmd/api/
```

Ou em modo debug:

```bash
GIN_MODE=debug go run ./cmd/api/
```

O servidor iniciara na porta **8094**. Se `SYNC_ON_STARTUP=true`, a reindexacao do Meilisearch ocorrera automaticamente.

### 5. Verificar se esta funcionando

```bash
curl -s http://localhost:8094/health | jq
```

Esperado:

```json
{
  "status": "ok",
  "postgres": "ok",
  "meilisearch": "ok"
}
```

### 6. Swagger UI

Acesse no navegador:

```
http://localhost:8094/swagger/index.html
```

Permite testar todos os endpoints interativamente com autenticacao Bearer.

### 7. Postman

Uma collection Postman esta disponivel em `postman/OBM_API.postman_collection.json`.

1. Abra o Postman
2. Clique em **Import**
3. Selecione o arquivo `postman/OBM_API.postman_collection.json`
4. Importe tambem o ambiente `postman/OBM_API_Local.postman_environment.json`
5. Selecione o ambiente **OBM API - Local** no canto superior direito
6. Faca login via o request **Auth > Login**, copie o token e cole na variavel `token` do ambiente

### 8. Rodar via Docker (build completo)

```bash
docker compose up --build -d
docker compose logs -f api
```

### 9. Parar os servicos

```bash
docker compose down
```

Para remover os dados (volumes):

```bash
docker compose down -v
```

---

## Exemplos Praticos de Uso

### Exemplo 1: Buscar medicamento e obter detalhes completos

```bash
# 1. Login
TOKEN=$(curl -s -X POST http://localhost:8094/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r '.token')

# 2. Buscar paracetamol
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/search?q=paracetamol&entity=vmp&limit=3" | jq

# 3. Obter detalhes do VMP encontrado
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmp/1234/detail" | jq
```

### Exemplo 2: Listar medicamentos de um fabricante

```bash
# Buscar AMPs do fabricante Medley
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/amp?fabricante=medley&limit=10" | jq

# Ou via busca global com filtro
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/search?q=&entity=amp&filter[fabricante]=medley&limit=10" | jq
```

### Exemplo 3: Buscar por codigo de registro

```bash
# Buscar AMP por codigo NU_APID
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/amp?codigo=9876543" | jq

# Buscar VMP por codigo NU_VPID
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmp?codigo=3456789" | jq
```

### Exemplo 4: Consultar formas farmaceuticas e vias de administracao

```bash
# Listar formas farmaceuticas disponiveis
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/domains/form?limit=20" | jq

# Listar vias de administracao
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/domains/route?limit=20" | jq

# Listar classes ATC
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/domains/atc-class?limit=20" | jq

# Buscar dominio por nome
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/domains/form?nome=comprimido&limit=5" | jq
```

### Exemplo 5: Percorrer todas as paginas de resultados

```bash
# Script para paginar por todos os VMPs
cursor=""
while true; do
  if [ -z "$cursor" ]; then
    response=$(curl -s -H "Authorization: Bearer $TOKEN" \
      "http://localhost:8094/api/v1/vmp?limit=50")
  else
    response=$(curl -s -H "Authorization: Bearer $TOKEN" \
      "http://localhost:8094/api/v1/vmp?limit=50&cursor=$cursor")
  fi

  # Processar os itens
  echo "$response" | jq '.items[] | {id: .co_seq_id, nome: .no_nm}'

  # Obter proximo cursor
  cursor=$(echo "$response" | jq -r '.cursor')

  # Se cursor vazio, chegamos ao fim
  if [ -z "$cursor" ] || [ "$cursor" = "null" ] || [ "$cursor" = "" ]; then
    echo "Fim dos resultados"
    break
  fi
done
```

### Exemplo 6: Obter ingredientes de um VMP

```bash
# VMP detalhado mostra ingredientes com concentracao
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmp/1234/detail" | jq '.ingredients'
```

### Exemplo 7: Consultar fornecedor de um AMP

```bash
# AMP detalhado mostra dados do fornecedor
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/amp/5678/detail" | jq '.supplier'
```

---

## Tabela Completa de Rotas

| Metodo | Rota | Auth | Descricao |
|--------|------|------|-----------|
| `POST` | `/auth/login` | Nao | Login, retorna JWT |
| `GET` | `/health` | Nao | Health check (PostgreSQL + Meilisearch) |
| `GET` | `/swagger/*` | Nao | Swagger UI |
| `GET` | `/api/v1/search` | Sim | Busca global Meilisearch |
| `GET` | `/api/v1/vmp` | Sim | Listar VMPs |
| `GET` | `/api/v1/vmp/:id` | Sim | VMP por ID |
| `GET` | `/api/v1/vmp/:id/detail` | Sim | VMP detalhado |
| `GET` | `/api/v1/amp` | Sim | Listar AMPs |
| `GET` | `/api/v1/amp/:id` | Sim | AMP por ID |
| `GET` | `/api/v1/amp/:id/detail` | Sim | AMP detalhado |
| `GET` | `/api/v1/vtm` | Sim | Listar VTMs |
| `GET` | `/api/v1/vtm/:id` | Sim | VTM por ID |
| `GET` | `/api/v1/vmpp` | Sim | Listar VMPPs |
| `GET` | `/api/v1/vmpp/:id` | Sim | VMPP por ID |
| `GET` | `/api/v1/ampp` | Sim | Listar AMPPs |
| `GET` | `/api/v1/ampp/:id` | Sim | AMPP por ID |
| `GET` | `/api/v1/suppliers` | Sim | Listar Fornecedores |
| `GET` | `/api/v1/suppliers/:id` | Sim | Fornecedor por ID |
| `GET` | `/api/v1/dcb` | Sim | Listar DCBs |
| `GET` | `/api/v1/dcb/:id` | Sim | DCB por ID |
| `GET` | `/api/v1/ingredients` | Sim | Listar Ingredientes |
| `GET` | `/api/v1/ingredients/:id` | Sim | Ingrediente por ID |
| `GET` | `/api/v1/domains/:domain` | Sim | Listar dominios por tipo |
| `GET` | `/api/v1/domains/:domain/:id` | Sim | Dominio por tipo e ID |
| `POST` | `/api/v1/admin/reindex` | Sim | Reindexar Meilisearch |
