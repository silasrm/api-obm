# OBM API — Ontologia Brasileira de Medicamentos

API REST para consulta de dados da [Ontologia Brasileira de Medicamentos (OBM)](https://portal-obm.saude.gov.br/), seguindo o modelo **dm+d** (Dictionary of Medicine and Devices) do NHS adaptado para o Brasil.

---

## Sumário

- [Sobre a OBM](#sobre-a-obm)
- [Conceitos da Ontologia](#conceitos-da-ontologia)
- [Autenticação](#autenticação)
- [Busca Global](#busca-global)
- [Referência de Endpoints](#referência-de-endpoints)
- [VMP — Virtual Medicinal Product](#vmp--virtual-medicinal-product)
- [AMP — Actual Medicinal Product](#amp--actual-medicinal-product)
- [VTM — Virtual Therapeutic Moiety](#vtm--virtual-therapeutic-moiety)
- [VMPP — Virtual Medicinal Product Pack](#vmpp--virtual-medicinal-product-pack)
- [AMPP — Actual Medicinal Product Pack](#ampp--actual-medicinal-product-pack)
- [DCB — Denominação Comum Brasileira](#dcb--denominação-comum-brasileira)
- [Ingredientes](#ingredientes)
- [Fornecedores](#fornecedores)
- [Domínios](#domínios)
- [Admin](#admin)
- [Paginação](#paginação)
- [Códigos de Erro](#códigos-de-erro)
- [CLI de Importação](#cli-de-importação)
- [Instalação Local](#instalação-local)
- [Exemplos Práticos de Uso](#exemplos-práticos-de-uso)

---

## Sobre a OBM

A **Ontologia Brasileira de Medicamentos (OBM)** é um padrão nacional de base de medicamentos para utilização em sistemas de prescrição e dispensação eletrônicas, instituída pela [Portaria GM/MS No 6.093, de 16 de dezembro de 2024](https://www.in.gov.br/en/web/dou/-/portaria-gm/ms-n-6.093-de-16-de-dezembro-de-2024-602264704).

Seus objetivos principais são:

- **Integrar e padronizar** dados de diferentes sistemas de informações em saúde
- **Normalizar registros** de prescrições e dispensações
- **Promover a interoperabilidade** por meio da Rede Nacional de Dados em Saúde (RNDS)
- **Potencializar a segurança do paciente** por meio de identificação unívoca e inequívoca de medicamentos
- **Seguir práticas internacionais** para descrição e categorização de medicamentos

A estrutura da OBM está baseada no modelo **dm+d** (Dictionary of Medicine and Devices) do **NHS** (National Health Service) do Reino Unido. Trata-se de dado público, atualizado, acessível, processável por máquina, em formato não proprietário, livre de licenças e com rastreabilidade das modificações via versionamento.

**Portal oficial:** [https://portal-obm.saude.gov.br/](https://portal-obm.saude.gov.br/)

---

## Conceitos da Ontologia

A OBM organiza os medicamentos em uma hierarquia inspirada no dm+d, com cinco níveis principais de abstração:

```
VTM (Virtual Therapeutic Moiety)
└── VMP (Virtual Medicinal Product)
    ├── VMPP (Virtual Medicinal Product Pack)
    └── AMP (Actual Medicinal Product)
        └── AMPP (Actual Medicinal Product Pack)
```

| Conceito | Sigla | O que representa | Exemplo |
|----------|-------|------------------|---------|
| **Virtual Therapeutic Moiety** | VTM | Princípio ativo genérico, sem forma farmacêutica | Paracetamol |
| **Virtual Medicinal Product** | VMP | Princípio ativo + forma farmacêutica + dose | Paracetamol 500mg comprimido |
| **Virtual Medicinal Product Pack** | VMPP | Apresentação virtual do VMP (quantidade) | Paracetamol 500mg comprimido — 20 comprimidos |
| **Actual Medicinal Product** | AMP | Produto comercial de um fabricante | Paracetamol 500mg comprimido — Medley |
| **Actual Medicinal Product Pack** | AMPP | Embalagem comercial do AMP (com código EAN) | Paracetamol 500mg comprimido — Medley — caixa 20 |

### Entidades complementares

| Entidade | Descrição |
|----------|-----------|
| **DCB** (Denominação Comum Brasileira) | Denominação oficial de substâncias ativas conforme ANVISA |
| **Ingredient Substance** | Substância ativa que compõe um medicamento (com código CAS e DCB) |
| **Supplier** (Fornecedor) | Fabricante ou detentor do registro sanitário do AMP |
| **Domain** | Tabelas de domínio/classificação (forma farmacêutica, via, classe ATC, etc.) |

### Relacionamentos principais

- Um **VTM** possui vários **VMPs** (diferentes formas/doses do mesmo princípio ativo)
- Um **VMP** pertence a um **VTM**
- Um **VMP** possui vários **AMPs** (produtos de diferentes fabricantes)
- Um **AMP** pertence a um **VMP** e a um **Supplier**
- Um **VMP** possui vários **VMPPs** (apresentações)
- Um **AMP** possui vários **AMPPs** (embalagens comerciais com EAN)
- **VMPs** e **AMPs** possuem **Ingredientes** com concentração
- **VMPs** e **AMPs** estão ligados a **Domínios** (forma farmacêutica, via, classe ATC, etc.)

---

## Autenticação

Todos os endpoints de dados (sob `/api/v1/`) exigem autenticação via **JWT Bearer Token**. Apenas `/auth/login` e `/health` são públicos.

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

### 2. Usar o token nas requisições

Inclua o header `Authorization: Bearer <token>` em toda requisição protegida:

```bash
TOKEN="eyJhbGciOiJIUzI1NiIs..."

curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmp?limit=5"
```

### 3. Expiração

O token é válido por **24 horas** (configurável via `JWT_EXPIRATION_HOURS`). Após expirar, faça login novamente.

---

## Busca Global

O endpoint de busca utiliza o [Meilisearch](https://www.meilisearch.com/) para busca full-text em VMPs, AMPs e Fornecedores.

```
GET /api/v1/search
```

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| `q` | string | Sim | Termo de busca |
| `entity` | string | Não | Entidades: `vmp`, `amp`, `supplier`. Separadas por vírgula. Padrão: todas |
| `limit` | int | Não | Limite de resultados (padrão: 20, máx: 100) |
| `cursor` | string | Não | Cursor de paginação |
| `filter[nome]` | string | Não | Filtro por nome |
| `filter[codigo]` | string | Não | Filtro por código |
| `filter[fabricante]` | string | Não | Filtro por fabricante |
| `filter[descricao]` | string | Não | Filtro por descrição |
| `filter[ativo]` | string | Não | Filtro por status ativo |

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

## Referência de Endpoints

### VMP — Virtual Medicinal Product

| Método | Rota | Descrição |
|--------|------|-----------|
| `GET` | `/api/v1/vmp` | Listar VMPs com filtros e paginação |
| `GET` | `/api/v1/vmp/:id` | Obter VMP por ID |
| `GET` | `/api/v1/vmp/:id/detail` | Obter VMP com detalhes completos |

**Parâmetros de listagem (`GET /api/v1/vmp`):**

| Parâmetro | Tipo | Descrição |
|-----------|------|-----------|
| `limit` | int | Limite por página (padrão: 20, máx: 100) |
| `cursor` | string | Cursor de paginação |
| `nome` | string | Filtro por nome (busca parcial, case-insensitive) |
| `codigo` | string | Filtro por código NU_VPID (busca exata) |
| `ativo` | boolean | Filtro por status ativo (padrão: `true`) |

**Exemplos:**

```bash
# Listar VMPs
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmp?limit=5"

# Filtrar por nome
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmp?nome=paracetamol&limit=5"

# Buscar por código
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
  "anvs_class": { "no_descr": "Isento de prescrição" },
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

| Método | Rota | Descrição |
|--------|------|-----------|
| `GET` | `/api/v1/amp` | Listar AMPs com filtros e paginação |
| `GET` | `/api/v1/amp/:id` | Obter AMP por ID |
| `GET` | `/api/v1/amp/:id/detail` | Obter AMP com detalhes completos |

**Parâmetros de listagem (`GET /api/v1/amp`):**

| Parâmetro | Tipo | Descrição |
|-----------|------|-----------|
| `limit` | int | Limite por página (padrão: 20, máx: 100) |
| `cursor` | string | Cursor de paginação |
| `nome` | string | Filtro por nome (busca parcial) |
| `codigo` | string | Filtro por código NU_APID (busca exata) |
| `fabricante` | string | Filtro por nome do fabricante (busca parcial) |
| `ativo` | boolean | Filtro por status ativo (padrão: `true`) |

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
  "med_class": { "no_descr": "Isento de prescrição" },
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

| Método | Rota | Descrição |
|--------|------|-----------|
| `GET` | `/api/v1/vtm` | Listar VTMs com filtros e paginação |
| `GET` | `/api/v1/vtm/:id` | Obter VTM por ID |

**Parâmetros de listagem:**

| Parâmetro | Tipo | Descrição |
|-----------|------|-----------|
| `limit` | int | Limite por página (padrão: 20, máx: 100) |
| `cursor` | string | Cursor de paginação |
| `nome` | string | Filtro por nome (busca parcial) |
| `ativo` | boolean | Filtro por status ativo (padrão: `true`) |

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

| Método | Rota | Descrição |
|--------|------|-----------|
| `GET` | `/api/v1/vmpp` | Listar VMPPs com filtros e paginação |
| `GET` | `/api/v1/vmpp/:id` | Obter VMPP por ID |

**Parâmetros de listagem:** `limit`, `cursor`, `nome` (busca parcial), `ativo`

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

| Método | Rota | Descrição |
|--------|------|-----------|
| `GET` | `/api/v1/ampp` | Listar AMPPs com filtros e paginação |
| `GET` | `/api/v1/ampp/:id` | Obter AMPP por ID |

**Parâmetros de listagem:** `limit`, `cursor`, `nome` (busca parcial), `ativo`

```bash
# Listar AMPPs
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/ampp?limit=5"

# AMPP por ID
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/ampp/1"
```

**Resposta AMPP (inclui códigos EAN):**

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

### DCB — Denominação Comum Brasileira

| Método | Rota | Descrição |
|--------|------|-----------|
| `GET` | `/api/v1/dcb` | Listar DCBs com filtros e paginação |
| `GET` | `/api/v1/dcb/:id` | Obter DCB por ID |

**Parâmetros de listagem:** `limit`, `cursor`, `nome` (busca parcial em DS_DCB), `ativo`

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

| Método | Rota | Descrição |
|--------|------|-----------|
| `GET` | `/api/v1/ingredients` | Listar Ingredient Substances com filtros e paginação |
| `GET` | `/api/v1/ingredients/:id` | Obter ingrediente por ID |

**Parâmetros de listagem:** `limit`, `cursor`, `nome` (busca parcial), `ativo`

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

| Método | Rota | Descrição |
|--------|------|-----------|
| `GET` | `/api/v1/suppliers` | Listar fornecedores com filtros e paginação |
| `GET` | `/api/v1/suppliers/:id` | Obter fornecedor por ID |

**Parâmetros de listagem:**

| Parâmetro | Tipo | Descrição |
|-----------|------|-----------|
| `limit` | int | Limite por página (padrão: 20, máx: 100) |
| `cursor` | string | Cursor de paginação |
| `nome` | string | Filtro por nome (busca parcial) |
| `codigo` | string | Filtro por código NU_CD (busca exata) |
| `ativo` | boolean | Filtro por status ativo (padrão: `true`) |

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

### Domínios

Domínios são tabelas de classificação e referência que categorizam os medicamentos (forma farmacêutica, via de administração, classe ATC, categoria de controle, etc.).

| Método | Rota | Descrição |
|--------|------|-----------|
| `GET` | `/api/v1/domains/:domain` | Listar domínios por tipo |
| `GET` | `/api/v1/domains/:domain/:id` | Obter domínio por tipo e ID |

**Parâmetros de listagem:** `limit`, `cursor`, `nome` (busca parcial em NO_DESCR e NO_DESCR_PT_BR), `ativo`

**Tipos de domínio disponíveis:**

| Slug | Tabela | Descrição |
|------|--------|-----------|
| `form` | `td_form` | Forma farmacêutica |
| `route` | `td_route` | Via de administração |
| `flavour` | `td_flavour` | Sabor |
| `legal-category` | `td_legal_category` | Categoria legal |
| `licensing-authority` | `td_licensing_authority` | Autoridade de licenciamento |
| `avail-restriction` | `td_availability_restriction` | Restrição de disponibilidade |
| `med-class` | `td_med_class_br` | Classe medicamentosa BR |
| `anvs-class` | `td_anvs_class_br` | Classe ANVISA |
| `atc-class` | `td_atc_class_br` | Classe ATC BR |
| `control-drug` | `td_control_drug_category` | Categoria de controle de drogas |
| `df-indicator` | `td_df_indicator` | Indicador DF |
| `discontinued-ind` | `td_discontinued_ind` | Indicador de descontinuação |
| `pres-status` | `td_virtual_product_pres_status` | Status de prescrição |
| `non-avail` | `td_virtual_product_non_avail` | Motivo de indisponibilidade |
| `basis-of-name` | `td_basis_of_name` | Base do nome |
| `basis-of-strnth` | `td_basis_of_strnth` | Base de concentração |
| `brimunologico` | `td_brimunologico` | Imunológico BR |
| `catmat` | `td_catmat_br` | CATMAT BR |
| `country` | `td_country` | País |
| `unit-of-measure` | `td_unit_of_measure` | Unidade de medida |
| `package` | `td_package` | Tipo de embalagem |
| `ont-form-route` | `td_ont_form_route` | Forma/Rota OBM |
| `preserv-cond` | `td_preserv_cond_br` | Condição de preservação BR |
| `rename-comp` | `td_rename_comp_br` | RENAME Complementar BR |
| `ingredient-source` | `td_ingredient_source_br` | Fonte do ingrediente BR |
| `healthcare-prof` | `td_healthcare_prof_br` | Profissional de saúde BR |
| `indicacao-farmpop` | `td_indicacao_farmpop_br` | Indicação Farmácia Popular BR |
| `monitoring-reason` | `td_monitoring_reason_br` | Motivo de monitoramento BR |
| `name-change-reason` | `td_name_change_reason` | Motivo de mudança de nome |
| `lic-auth-change-reason` | `td_lic_auth_change_reason` | Motivo de mudança de licença |
| `phpid` | `td_phpid` | PHPID |
| `local-aplicacao` | `td_local_aplicacao` | Local de aplicação |

```bash
# Listar formas farmacêuticas
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/domains/form?limit=5"

# Listar vias de administração
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/domains/route?limit=5"

# Listar classes ATC
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/domains/atc-class?limit=5"

# Domínio específico por ID
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/domains/form/1"
```

---

### Admin

| Método | Rota | Descrição |
|--------|------|-----------|
| `GET` | `/health` | Health check (sem autenticação) |
| `POST` | `/api/v1/admin/reindex` | Reindexar Meilisearch |

**Health check (público):**

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

> Status `"degraded"` indica que PostgreSQL ou Meilisearch está indisponível.

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

## Paginação

Todos os endpoints de listagem utilizam **paginação baseada em cursor**. Não há numeração de páginas — o cursor codifica o offset da próxima página.

### Como funciona

1. Faça a primeira requisição sem cursor
2. A resposta contém o campo `cursor` com o valor para a próxima página
3. Use esse valor no parâmetro `cursor` da próxima requisição
4. Quando `cursor` for vazio (`""`), não há mais páginas

### Parâmetros

| Parâmetro | Descrição | Padrão | Máximo |
|-----------|-----------|--------|--------|
| `limit` | Registros por página | 20 | 100 |
| `cursor` | Cursor da próxima página | — | — |

### Exemplo

```bash
# Primeira página
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
# Próxima página (usar o cursor retornado)
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

## Códigos de Erro

| HTTP | Mensagem | Condição |
|------|----------|----------|
| `400` | `invalid request body` | JSON malformado no login |
| `400` | `invalid query parameters` | Parâmetros de busca inválidos |
| `400` | `invalid id` | Parâmetro `:id` não numérico |
| `401` | `missing authorization header` | Header `Authorization` ausente |
| `401` | `invalid token` | Token JWT inválido ou expirado |
| `401` | `invalid token claims` | Claims do JWT não extraíveis |
| `401` | `invalid credentials` | Usuário/senha incorretos ou usuário inativo |
| `404` | `not found` | Registro não encontrado |
| `500` | (mensagem variável) | Erro interno (banco, Meilisearch, etc.) |

### Formato da resposta de erro

```json
{
  "error": "descrição do erro",
  "code": 401
}
```

---

## CLI de Importação

O CLI de importação (`cmd/import/`) permite carregar dados da OBM no banco PostgreSQL e reindexar o Meilisearch em um único comando. Ele converte automaticamente dumps MySQL para o formato PostgreSQL, importa os dados e atualiza os índices de busca.

O dump OBM (~137 MB descomprimido) não está versionado no git por exceder o limite de arquivo do GitHub (100 MB). Por isso, a importação é feita sob demanda a partir de uma fonte local.

### Fontes aceitas

| Formato | Exemplo | Descrição |
|---------|---------|-----------|
| ZIP | `portal-obm-20250530.zip` | Extrai o arquivo `.sql` interno automaticamente |
| SQL | `portal-obm-20250530.sql` | Dump MySQL utilizado diretamente |
| SQL.GZ | `dump.sql.gz` | Dump compactado com gzip |
| MySQL | `mysql://user:pass@host:3306/db` | Executa `mysqldump` automaticamente (necessário ter o binário no PATH) |

### Flags

| Flag | Tipo | Padrão | Descrição |
|------|------|--------|-----------|
| `--source` | string | (obrigatório) | Caminho para `.zip`/`.sql`/`.sql.gz` ou URI `mysql://` |
| `--output` | string | `migrations/postgres/001_obm_schema.sql` | Caminho de saída para `--convert-only` |
| `--convert-only` | bool | `false` | Somente converter MySQL→PostgreSQL, sem importar |
| `--reindex-only` | bool | `false` | Somente reindexar o Meilisearch |
| `--skip-index` | bool | `false` | Pular a reindexação do Meilisearch após a importação |
| `--validate` | bool | `false` | Executar validação pós-importação (contagem de registros e integridade referencial) |
| `--full` | bool | `true` | Remover e recriar o schema antes de importar (modo padrão) |

### Fluxo da importação

1. **Resolução da fonte** — o CLI identifica o tipo de entrada (ZIP, SQL, gzip ou MySQL) e prepara um leitor de dados
2. **Conversão** — o dump MySQL é convertido para PostgreSQL linha a linha (streaming via `io.Pipe`, sem arquivo temporário em disco)
3. **Importação** — os comandos SQL convertidos são executados diretamente no PostgreSQL
4. **Validação** (opcional, `--validate`) — contagem de registros por tabela e verificação de integridade referencial (VMP→VTM, AMP→VMP, AMP→Fornecedor)
5. **Metadados** — a tabela `obm_metadata` registra a versão dos dados, data da importação, arquivo de origem e contagem de registros
6. **Reindexação** — o Meilisearch é reindexado automaticamente com os dados de VMP, AMP e Fornecedores

### Exemplos

```bash
# Importação completa (ZIP → PostgreSQL + Meilisearch)
go run cmd/import/main.go --source=portal-obm-20250530.zip

# Importação com validação pós-importação
go run cmd/import/main.go --source=portal-obm-20250530.zip --validate

# Apenas converter MySQL→PostgreSQL (sem importar)
go run cmd/import/main.go --source=dump.sql --convert-only --output=migrations/postgres/001_obm_schema.sql

# Apenas reindexar o Meilisearch
go run cmd/import/main.go --reindex-only

# Importar a partir de uma conexão MySQL direta (necessário ter mysqldump no PATH)
go run cmd/import/main.go --source=mysql://user:pass@host:3306/dbportalobm

# Importar sem reindexar o Meilisearch
go run cmd/import/main.go --source=dump.sql --skip-index
```

### Método alternativo: SQL no diretório de migrations

Se preferir que o Docker Compose carregue o banco na primeira inicialização, converta o dump separadamente:

```bash
go run scripts/convert_sql.go -input <caminho_do_dump_mysql> -output migrations/postgres/001_obm_schema.sql
```

O PostgreSQL do Docker Compose carrega automaticamente qualquer arquivo `.sql` colocado em `migrations/postgres/` na primeira inicialização (via `docker-entrypoint-initdb.d`).

> **Atenção:** o arquivo `migrations/postgres/001_obm_schema.sql` possui ~1,1 milhão de linhas. O carregamento inicial pode levar vários minutos.

---

## Instalação Local

### Pré-requisitos

- [Docker](https://www.docker.com/) e Docker Compose
- [Go 1.25+](https://go.dev/dl/)
- `curl` ou Postman/Insomnia
- `jq` (opcional, para formatar JSON)

### 1. Subir infraestrutura

```bash
docker compose up postgres meilisearch -d
```

Aguarde os healthchecks passarem (~30 s). O PostgreSQL estará na porta **5433** e o Meilisearch na porta **7701**.

Verifique se os serviços estão saudáveis:

```bash
docker compose ps
```

Ambos devem mostrar `(healthy)` na coluna STATUS.

### 2. Configurar o .env

```bash
cp .env.example .env
```

O `.env` já vem com as portas corretas. Variáveis disponíveis:

| Variável | Padrão | Descrição |
|----------|--------|-----------|
| `PG_HOST` | `localhost` | Host do PostgreSQL |
| `PG_PORT` | `5433` | Porta do PostgreSQL (5432 dentro do Docker) |
| `PG_USER` | `obm` | Usuário do PostgreSQL |
| `PG_PASSWORD` | `obm123` | Senha do PostgreSQL |
| `PG_DATABASE` | `dbportalobm` | Nome do banco |
| `PG_SSLMODE` | `disable` | SSL mode do PostgreSQL |
| `MEILI_URL` | `http://localhost:7701` | URL do Meilisearch (7700 dentro do Docker) |
| `MEILI_API_KEY` | `obm-meili-master-key` | Chave API do Meilisearch |
| `MEILI_INDEX_PREFIX` | `obm_` | Prefixo dos índices (produz `obm_vmp`, `obm_amp`, `obm_supplier`) |
| `JWT_SECRET` | `obm-secret-key-change-in-prod` | Chave de assinatura JWT (altere em produção!) |
| `JWT_EXPIRATION_HOURS` | `24` | Expiração do token em horas |
| `SERVER_PORT` | `8094` | Porta do servidor API |
| `GIN_MODE` | `release` | Modo do Gin (`debug` ou `release`) |
| `SYNC_ON_STARTUP` | `true` | Reindexar Meilisearch ao iniciar |

### 3. Importar dados

Com a infraestrutura rodando, importe o dump OBM com o CLI de importação (veja a seção [CLI de Importação](#cli-de-importação) para detalhes completos):

```bash
go run cmd/import/main.go --source=portal-obm-20250530.zip
```

Ou, se preferir carregar via Docker Compose, coloque o SQL convertido em `migrations/postgres/` e reinicie o container do PostgreSQL.

### 4. Criar usuários iniciais

```bash
go run scripts/seed_users.go
```

Cria os usuários padrão:

| Usuário | Senha | Ativo |
|---------|-------|-------|
| `admin` | `admin123` | sim |
| `viewer` | `viewer123` | sim |

### 5. Executar a API

```bash
go run ./cmd/api/
```

Ou em modo debug:

```bash
GIN_MODE=debug go run ./cmd/api/
```

O servidor iniciará na porta **8094**. Se `SYNC_ON_STARTUP=true`, a reindexação do Meilisearch ocorrerá automaticamente.

### 6. Verificar o funcionamento

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

### 7. Swagger UI

Acesse no navegador:

```
http://localhost:8094/swagger/index.html
```

Permite testar todos os endpoints interativamente com autenticação Bearer.

### 8. Postman

Uma collection Postman está disponível em `postman/OBM_API.postman_collection.json`.

1. Abra o Postman
2. Clique em **Import**
3. Selecione o arquivo `postman/OBM_API.postman_collection.json`
4. Importe também o ambiente `postman/OBM_API_Local.postman_environment.json`
5. Selecione o ambiente **OBM API - Local** no canto superior direito
6. Faça login via o request **Auth > Login**, copie o token e cole na variável `token` do ambiente

### 9. Executar via Docker (build completo)

```bash
docker compose up --build -d
docker compose logs -f api
```

### 10. Parar os serviços

```bash
docker compose down
```

Para remover os dados (volumes):

```bash
docker compose down -v
```

---

## Exemplos Práticos de Uso

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

### Exemplo 3: Buscar por código de registro

```bash
# Buscar AMP por código NU_APID
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/amp?codigo=9876543" | jq

# Buscar VMP por código NU_VPID
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmp?codigo=3456789" | jq
```

### Exemplo 4: Consultar formas farmacêuticas e vias de administração

```bash
# Listar formas farmacêuticas disponíveis
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/domains/form?limit=20" | jq

# Listar vias de administração
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/domains/route?limit=20" | jq

# Listar classes ATC
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/domains/atc-class?limit=20" | jq

# Buscar domínio por nome
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/domains/form?nome=comprimido&limit=5" | jq
```

### Exemplo 5: Percorrer todas as páginas de resultados

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

  # Obter próximo cursor
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
# VMP detalhado mostra ingredientes com concentração
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

| Método | Rota | Auth | Descrição |
|--------|------|------|-----------|
| `POST` | `/auth/login` | Não | Login, retorna JWT |
| `GET` | `/health` | Não | Health check (PostgreSQL + Meilisearch) |
| `GET` | `/swagger/*` | Não | Swagger UI |
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
| `GET` | `/api/v1/domains/:domain` | Sim | Listar domínios por tipo |
| `GET` | `/api/v1/domains/:domain/:id` | Sim | Domínio por tipo e ID |
| `POST` | `/api/v1/admin/reindex` | Sim | Reindexar Meilisearch |
