package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type PostmanCollection struct {
	Info     PostmanInfo      `json:"info"`
	Item     []PostmanItem    `json:"item"`
	Variable []PostmanVar     `json:"variable"`
	Auth     *PostmanAuth     `json:"auth,omitempty"`
}

type PostmanInfo struct {
	Name        string `json:"name"`
	Schema      string `json:"schema"`
	Description string `json:"description"`
}

type PostmanVar struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
}

type PostmanAuth struct {
	Type   string            `json:"type"`
	Bearer []PostmanAuthItem `json:"bearer,omitempty"`
}

type PostmanAuthItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type PostmanItem struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Item        []PostmanItem     `json:"item,omitempty"`
	Request      *PostmanRequest  `json:"request,omitempty"`
	Response     []interface{}    `json:"response"`
}

type PostmanRequest struct {
	Method  string           `json:"method"`
	Header  []PostmanHeader  `json:"header"`
	Body    *PostmanBody     `json:"body,omitempty"`
	URL     PostmanURL       `json:"url"`
}

type PostmanHeader struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Type    string `json:"type"`
}

type PostmanBody struct {
	Mode    string `json:"mode"`
	Raw     string `json:"raw"`
	Options struct {
		Raw struct {
			Language string `json:"language"`
		} `json:"raw"`
	} `json:"options"`
}

type PostmanURL struct {
	Raw   string   `json:"raw"`
	Host  []string `json:"host"`
	Path  []string `json:"path"`
	Query []PostmanQuery `json:"query,omitempty"`
}

type PostmanQuery struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Disabled bool  `json:"disabled"`
	Description string `json:"description,omitempty"`
}

func main() {
	baseURL := "{{base_url}}"

	collection := PostmanCollection{
		Info: PostmanInfo{
			Name:        "OBM API",
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
			Description: "API da Ontologia Brasileira de Medicamentos (OBM). Gerencia dados de medicamentos seguindo o padrão dm%d adaptado para o Brasil.",
		},
		Variable: []PostmanVar{
			{Key: "base_url", Value: "http://localhost:8094", Enabled: true},
			{Key: "token", Value: "", Enabled: true},
		},
		Auth: &PostmanAuth{
			Type: "bearer",
			Bearer: []PostmanAuthItem{
				{Key: "token", Value: "{{token}}", Type: "string"},
			},
		},
	}

	authHeader := []PostmanHeader{
		{Key: "Authorization", Value: "Bearer {{token}}", Type: "string"},
		{Key: "Content-Type", Value: "application/json", Type: "string"},
	}

	collection.Item = []PostmanItem{
		makeAuthFolder(baseURL),
		makeHealthFolder(baseURL),
		makeSearchFolder(baseURL, authHeader),
		makeVMPFolder(baseURL, authHeader),
		makeAMPFolder(baseURL, authHeader),
		makeVTMFolder(baseURL, authHeader),
		makeVMPPFolder(baseURL, authHeader),
		makeAMPPFolder(baseURL, authHeader),
		makeSupplierFolder(baseURL, authHeader),
		makeDCBFolder(baseURL, authHeader),
		makeIngredientFolder(baseURL, authHeader),
		makeDomainFolder(baseURL, authHeader),
		makeAdminFolder(baseURL, authHeader),
	}

	data, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile("postman/OBM_API.postman_collection.json", data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Postman collection generated: postman/OBM_API.postman_collection.json")
}

func req(method, name, url string, headers []PostmanHeader, query []PostmanQuery, body *PostmanBody) PostmanItem {
	parts := strings.Split(strings.TrimPrefix(url, "{{base_url}}"), "/")
	path := make([]string, 0)
	for _, p := range parts {
		if p != "" {
			path = append(path, p)
		}
	}
	return PostmanItem{
		Name: name,
		Request: &PostmanRequest{
			Method: method,
			Header: headers,
			Body:   body,
			URL: PostmanURL{
				Raw:   url,
				Host:  []string{`{{base_url}}`},
				Path:  path,
				Query: query,
			},
		},
		Response: []interface{}{},
	}
}

func makeAuthFolder(baseURL string) PostmanItem {
	loginBody := &PostmanBody{Mode: "raw", Raw: `{"username":"admin","password":"admin123"}`}
	loginBody.Options.Raw.Language = "json"
	return PostmanItem{
		Name: "Auth",
		Item: []PostmanItem{
			req("POST", "Login", baseURL+"/auth/login", []PostmanHeader{
				{Key: "Content-Type", Value: "application/json", Type: "string"},
			}, nil, loginBody),
		},
	}
}

func makeHealthFolder(baseURL string) PostmanItem {
	return PostmanItem{
		Name: "Health",
		Item: []PostmanItem{
			req("GET", "Health Check", baseURL+"/health", nil, nil, nil),
		},
	}
}

func makeSearchFolder(baseURL string, h []PostmanHeader) PostmanItem {
	return PostmanItem{
		Name: "Search",
		Item: []PostmanItem{
			req("GET", "Busca Global", baseURL+"/api/v1/search", h, []PostmanQuery{
				{Key: "q", Value: "paracetamol", Disabled: false, Description: "Termo de busca"},
				{Key: "entity", Value: "vmp,amp", Disabled: true, Description: "Entidades (vmp,amp,supplier)"},
				{Key: "limit", Value: "20", Disabled: true},
				{Key: "cursor", Value: "", Disabled: true, Description: "Cursor de paginação"},
				{Key: "filter[nome]", Value: "", Disabled: true},
				{Key: "filter[codigo]", Value: "", Disabled: true},
				{Key: "filter[fabricante]", Value: "", Disabled: true},
			}, nil),
		},
	}
}

func makeVMPFolder(baseURL string, h []PostmanHeader) PostmanItem {
	return PostmanItem{
		Name: "VMP",
		Item: []PostmanItem{
			req("GET", "Listar VMPs", baseURL+"/api/v1/vmp", h, []PostmanQuery{
				{Key: "limit", Value: "20", Disabled: true},
				{Key: "cursor", Value: "", Disabled: true},
				{Key: "nome", Value: "", Disabled: true},
				{Key: "codigo", Value: "", Disabled: true},
				{Key: "ativo", Value: "true", Disabled: true},
			}, nil),
			req("GET", "VMP por ID", baseURL+"/api/v1/vmp/1", h, nil, nil),
			req("GET", "VMP Detalhado", baseURL+"/api/v1/vmp/1/detail", h, nil, nil),
		},
	}
}

func makeAMPFolder(baseURL string, h []PostmanHeader) PostmanItem {
	return PostmanItem{
		Name: "AMP",
		Item: []PostmanItem{
			req("GET", "Listar AMPs", baseURL+"/api/v1/amp", h, []PostmanQuery{
				{Key: "limit", Value: "20", Disabled: true},
				{Key: "cursor", Value: "", Disabled: true},
				{Key: "nome", Value: "", Disabled: true},
				{Key: "codigo", Value: "", Disabled: true},
				{Key: "fabricante", Value: "", Disabled: true},
				{Key: "ativo", Value: "true", Disabled: true},
			}, nil),
			req("GET", "AMP por ID", baseURL+"/api/v1/amp/1", h, nil, nil),
			req("GET", "AMP Detalhado", baseURL+"/api/v1/amp/1/detail", h, nil, nil),
		},
	}
}

func makeVTMFolder(baseURL string, h []PostmanHeader) PostmanItem {
	return PostmanItem{
		Name: "VTM",
		Item: []PostmanItem{
			req("GET", "Listar VTMs", baseURL+"/api/v1/vtm", h, []PostmanQuery{
				{Key: "limit", Value: "20", Disabled: true},
				{Key: "cursor", Value: "", Disabled: true},
				{Key: "nome", Value: "", Disabled: true},
				{Key: "ativo", Value: "true", Disabled: true},
			}, nil),
			req("GET", "VTM por ID", baseURL+"/api/v1/vtm/1", h, nil, nil),
		},
	}
}

func makeVMPPFolder(baseURL string, h []PostmanHeader) PostmanItem {
	return PostmanItem{
		Name: "VMPP",
		Item: []PostmanItem{
			req("GET", "Listar VMPPs", baseURL+"/api/v1/vmpp", h, []PostmanQuery{
				{Key: "limit", Value: "20", Disabled: true},
				{Key: "cursor", Value: "", Disabled: true},
				{Key: "nome", Value: "", Disabled: true},
				{Key: "ativo", Value: "true", Disabled: true},
			}, nil),
			req("GET", "VMPP por ID", baseURL+"/api/v1/vmpp/1", h, nil, nil),
		},
	}
}

func makeAMPPFolder(baseURL string, h []PostmanHeader) PostmanItem {
	return PostmanItem{
		Name: "AMPP",
		Item: []PostmanItem{
			req("GET", "Listar AMPPs", baseURL+"/api/v1/ampp", h, []PostmanQuery{
				{Key: "limit", Value: "20", Disabled: true},
				{Key: "cursor", Value: "", Disabled: true},
				{Key: "nome", Value: "", Disabled: true},
				{Key: "ativo", Value: "true", Disabled: true},
			}, nil),
			req("GET", "AMPP por ID", baseURL+"/api/v1/ampp/1", h, nil, nil),
		},
	}
}

func makeSupplierFolder(baseURL string, h []PostmanHeader) PostmanItem {
	return PostmanItem{
		Name: "Suppliers",
		Item: []PostmanItem{
			req("GET", "Listar Fornecedores", baseURL+"/api/v1/suppliers", h, []PostmanQuery{
				{Key: "limit", Value: "20", Disabled: true},
				{Key: "cursor", Value: "", Disabled: true},
				{Key: "nome", Value: "", Disabled: true},
				{Key: "codigo", Value: "", Disabled: true},
				{Key: "ativo", Value: "true", Disabled: true},
			}, nil),
			req("GET", "Fornecedor por ID", baseURL+"/api/v1/suppliers/1", h, nil, nil),
		},
	}
}

func makeDCBFolder(baseURL string, h []PostmanHeader) PostmanItem {
	return PostmanItem{
		Name: "DCB",
		Item: []PostmanItem{
			req("GET", "Listar DCBs", baseURL+"/api/v1/dcb", h, []PostmanQuery{
				{Key: "limit", Value: "20", Disabled: true},
				{Key: "cursor", Value: "", Disabled: true},
				{Key: "nome", Value: "", Disabled: true},
				{Key: "ativo", Value: "true", Disabled: true},
			}, nil),
			req("GET", "DCB por ID", baseURL+"/api/v1/dcb/1", h, nil, nil),
		},
	}
}

func makeIngredientFolder(baseURL string, h []PostmanHeader) PostmanItem {
	return PostmanItem{
		Name: "Ingredients",
		Item: []PostmanItem{
			req("GET", "Listar Ingredientes", baseURL+"/api/v1/ingredients", h, []PostmanQuery{
				{Key: "limit", Value: "20", Disabled: true},
				{Key: "cursor", Value: "", Disabled: true},
				{Key: "nome", Value: "", Disabled: true},
				{Key: "ativo", Value: "true", Disabled: true},
			}, nil),
			req("GET", "Ingrediente por ID", baseURL+"/api/v1/ingredients/1", h, nil, nil),
		},
	}
}

func makeDomainFolder(baseURL string, h []PostmanHeader) PostmanItem {
	domainTypes := []struct {
		name string
		key  string
	}{
		{"Formas Farmacêuticas", "form"},
		{"Vias de Administração", "route"},
		{"Sabores", "flavour"},
		{"Categorias Legais", "legal_category"},
		{"Autoridades Licenciadas", "licensing_authority"},
		{"Restrições de Disponibilidade", "availability_restriction"},
		{"Classificação ANVS", "med_class_br"},
		{"Classificação ATC", "atc_class_br"},
		{"Categorias Controladas", "control_drug_category"},
		{"Indicadores DF", "df_indicator"},
		{"Países", "country"},
		{"Unidades de Medida", "unit_of_measure"},
		{"Preservação", "preserv_cond_br"},
		{"Rename", "rename_comp_br"},
		{"Brimunologico", "brimunologico"},
		{"Catmat", "catmat_br"},
		{"Monitoramento", "monitoring_reason_br"},
		{"Local Aplicação", "local_aplicacao"},
	}

	items := []PostmanItem{
		req("GET", "Listar Domínios", baseURL+"/api/v1/domains/form", h, []PostmanQuery{
			{Key: "limit", Value: "20", Disabled: true},
			{Key: "cursor", Value: "", Disabled: true},
			{Key: "nome", Value: "", Disabled: true},
		}, nil),
		req("GET", "Domínio por ID", baseURL+"/api/v1/domains/form/1", h, nil, nil),
	}

	for _, d := range domainTypes {
		items = append(items, req("GET", d.name, baseURL+"/api/v1/domains/"+d.key, h, []PostmanQuery{
			{Key: "limit", Value: "20", Disabled: true},
			{Key: "nome", Value: "", Disabled: true},
		}, nil))
	}

	return PostmanItem{
		Name: "Domains",
		Item: items,
	}
}

func makeAdminFolder(baseURL string, h []PostmanHeader) PostmanItem {
	return PostmanItem{
		Name: "Admin",
		Item: []PostmanItem{
			req("POST", "Reindexar Meilisearch", baseURL+"/api/v1/admin/reindex", h, nil, nil),
		},
	}
}
