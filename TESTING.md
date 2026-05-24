# Guia de Teste - API OBM

## Pré-requisitos

- Docker e Docker Compose
- Go 1.25+
- curl ou Postman/Insomnia
- jq (opcional, para formatar JSON)

## 1. Subir infraestrutura

```bash
docker compose up postgres meilisearch -d
```

Aguardar os healthchecks passarem (~30s). O PostgreSQL estará disponível na porta **5433** e o Meilisearch na porta **7701** do host. O schema SQL será carregado automaticamente pelo volume `./migrations/postgres:/docker-entrypoint-initdb.d`.

> **Atenção**: o arquivo `001_obm_schema.sql` tem ~1.1 milhão de linhas. O carregamento inicial pode levar vários minutos. Monitore com:
> ```bash
> docker compose logs -f postgres
> ```

Verificar se os serviços estão saudáveis:

```bash
docker compose ps
```

Ambos devem mostrar `(healthy)` na coluna STATUS.

## 2. Configurar o .env

Copie o exemplo e ajuste se necessário:

```bash
cp .env.example .env
```

O `.env` já vem com as portas corretas (`PG_PORT=5433`, `MEILI_URL=http://localhost:7701`).

## 3. Criar usuarios iniciais

```bash
go run scripts/seed_users.go
```

Cria `admin/admin123` e `viewer/viewer123`.

## 4. Rodar a API

```bash
GIN_MODE=debug go run ./cmd/api/
```

Ou em modo release (padrão):

```bash
go run ./cmd/api/
```

## 5. Obter token JWT

```bash
curl -s -X POST http://localhost:8094/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

Resposta esperada:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 86400
}
```

Salve o token:

```bash
TOKEN="<cole_o_token_aqui>"
```

## 6. Testar endpoints

### Health Check (sem auth)

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

### Busca Global

```bash
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/search?q=paracetamol&limit=5" | jq
```

Filtros disponíveis: `filter[nome]`, `filter[codigo]`, `filter[fabricante]`, `filter[descricao]`, `filter[ativo]`

```bash
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/search?q=dipirona&entity=vmp,amp&limit=5" | jq
```

### VMP - Virtual Medicinal Product

```bash
# Listar VMPs
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmp?limit=5" | jq

# Filtrar por nome
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmp?nome=paracetamol&limit=5" | jq

# VMP por ID
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmp/1" | jq

# VMP detalhado (com VTM, domínios, ingredientes, rotas, etc.)
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmp/1/detail" | jq
```

### AMP - Actual Medicinal Product

```bash
# Listar AMPs
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/amp?limit=5" | jq

# Filtrar por fabricante
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/amp?fabricante=medley&limit=5" | jq

# AMP por ID
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/amp/1" | jq

# AMP detalhado (com VMP, fornecedor, domínios, ingredientes, etc.)
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/amp/1/detail" | jq
```

### VTM - Virtual Therapeutic Moiety

```bash
# Listar VTMs
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vtm?limit=5" | jq

# VTM por ID
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vtm/1" | jq
```

### VMPP - Virtual Medicinal Product Pack

```bash
# Listar VMPPs
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmpp?limit=5" | jq

# VMPP por ID
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmpp/1" | jq
```

### AMPP - Actual Medicinal Product Pack

```bash
# Listar AMPPs
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/ampp?limit=5" | jq

# AMPP por ID
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/ampp/1" | jq
```

### Fornecedores (Suppliers)

```bash
# Listar fornecedores
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/suppliers?limit=5" | jq

# Fornecedor por ID
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/suppliers/1" | jq
```

### DCB - Denominação Comum Brasileira

```bash
# Listar DCBs
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/dcb?limit=5" | jq

# DCB por ID
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/dcb/1" | jq
```

### Ingredientes (Ingredient Substances)

```bash
# Listar ingredientes
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/ingredients?limit=5" | jq

# Ingrediente por ID
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/ingredients/1" | jq
```

### Domínios

Tipos disponíveis: `form`, `route`, `flavour`, `legal_category`, `licensing_authority`, `availability_restriction`, `med_class_br`, `anvs_class_br`, `atc_class_br`, `control_drug_category`, `df_indicator`, `discontinued_ind`, `pres_status`, `virtual_product_non_avail`, `basis_of_name`, `basis_of_strnth`, `brimunologico`, `catmat_br`, `country`, `healthcare_prof_br`, `indicacao_farmpop_br`, `ingredient_source_br`, `lic_auth_change_reason`, `local_aplicacao`, `monitoring_reason_br`, `name_change_reason`, `ont_form_route`, `package`, `phpid`, `preserv_cond_br`, `rename_comp_br`, `unit_of_measure`

```bash
# Listar domínios por tipo
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/domains/form?limit=5" | jq

curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/domains/route?limit=5" | jq

# Domínio por tipo e ID
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/domains/form/1" | jq
```

### Admin - Reindexar Meilisearch

```bash
curl -s -X POST -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/admin/reindex" | jq
```

Esperado:

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

## 7. Paginação com cursor

Todas as rotas de listagem suportam paginação baseada em cursor:

```bash
# Primeira página
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmp?limit=3" | jq

# Use o valor de "cursor" da resposta anterior para a próxima página
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmp?limit=3&cursor=<cursor_value>" | jq
```

## 8. Testes de erro

### Auth inválida (401)

```bash
curl -s -H "Authorization: Bearer invalid_token" \
  "http://localhost:8094/api/v1/vmp" | jq
```

Esperado: `{"error":"invalid token","code":401}`

### Sem Authorization header (401)

```bash
curl -s "http://localhost:8094/api/v1/vmp" | jq
```

Esperado: `{"error":"missing authorization header","code":401}`

### ID inválido (400)

```bash
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmp/abc" | jq
```

Esperado: `{"error":"invalid id","code":400}`

### Registro não encontrado (200 com null)

```bash
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8094/api/v1/vmp/999999999" | jq
```

### Login inválido (401)

```bash
curl -s -X POST http://localhost:8094/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"wrong"}' | jq
```

Esperado: `{"error":"invalid credentials","code":401}`

## 9. Swagger UI

Acesse no navegador:

```
http://localhost:8094/swagger/index.html
```

Permite testar todos os endpoints interativamente com autenticação Bearer.

## 10. Postman

Uma collection Postman está disponível em `postman/OBM_API.postman_collection.json` com todos os endpoints organizados por categoria.

### Importar no Postman

1. Abra o Postman
2. Clique em **Import**
3. Selecione o arquivo `postman/OBM_API.postman_collection.json`
4. Importe também o ambiente `postman/OBM_API_Local.postman_environment.json`
5. Selecione o ambiente **OBM API - Local** no canto superior direito

### Configurar token

1. Faça login via o request **Auth > Login**
2. Copie o token da resposta
3. Clique no ambiente **OBM API - Local** e cole o token na variável `token`
4. Todos os requests protegidos usarão automaticamente o Bearer token

### Regenerar a collection

```bash
go run scripts/gen_postman.go
```

## 10. Testar via Docker (build completo)

```bash
docker compose up --build -d
docker compose logs -f api
```

## 11. Parar os serviços

```bash
docker compose down
```

Para remover os dados (volumes):

```bash
docker compose down -v
```

## Tabela de Rotas

| Método | Rota | Auth | Descrição |
|--------|------|------|-----------|
| POST | `/auth/login` | Não | Login, retorna JWT |
| GET | `/health` | Não | Health check (PG + Meilisearch) |
| GET | `/swagger/*` | Não | Swagger UI |
| GET | `/api/v1/search` | Sim | Busca global Meilisearch |
| GET | `/api/v1/vmp` | Sim | Listar VMPs |
| GET | `/api/v1/vmp/:id` | Sim | VMP por ID |
| GET | `/api/v1/vmp/:id/detail` | Sim | VMP detalhado |
| GET | `/api/v1/amp` | Sim | Listar AMPs |
| GET | `/api/v1/amp/:id` | Sim | AMP por ID |
| GET | `/api/v1/amp/:id/detail` | Sim | AMP detalhado |
| GET | `/api/v1/vtm` | Sim | Listar VTMs |
| GET | `/api/v1/vtm/:id` | Sim | VTM por ID |
| GET | `/api/v1/vmpp` | Sim | Listar VMPPs |
| GET | `/api/v1/vmpp/:id` | Sim | VMPP por ID |
| GET | `/api/v1/ampp` | Sim | Listar AMPPs |
| GET | `/api/v1/ampp/:id` | Sim | AMPP por ID |
| GET | `/api/v1/suppliers` | Sim | Listar Fornecedores |
| GET | `/api/v1/suppliers/:id` | Sim | Fornecedor por ID |
| GET | `/api/v1/dcb` | Sim | Listar DCBs |
| GET | `/api/v1/dcb/:id` | Sim | DCB por ID |
| GET | `/api/v1/ingredients` | Sim | Listar Ingredientes |
| GET | `/api/v1/ingredients/:id` | Sim | Ingrediente por ID |
| GET | `/api/v1/domains/:domain` | Sim | Listar domínios por tipo |
| GET | `/api/v1/domains/:domain/:id` | Sim | Domínio por tipo e ID |
| POST | `/api/v1/admin/reindex` | Sim | Reindexar Meilisearch |