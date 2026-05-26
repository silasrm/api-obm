package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/infrastructure/config"
	"github.com/silasrm/api-obm/internal/infrastructure/persistence/meilisearch"
	redisrepo "github.com/silasrm/api-obm/internal/infrastructure/persistence/redis"
	"github.com/silasrm/api-obm/internal/infrastructure/persistence/postgres"
)

func main() {
	source := flag.String("source", "", "Caminho do arquivo XLSX da CMED")
	dateStr := flag.String("date", "", "Data de referência (YYYY-MM-DD)")
	headerRow := flag.Int("header-row", 42, "Linha do cabeçalho na planilha (1-based)")
	skipIndex := flag.Bool("skip-index", false, "Pular reindexação do Meilisearch")
	flag.Parse()

	if *source == "" || *dateStr == "" {
		fmt.Fprintln(os.Stderr, "Uso: cmed_import --source planilha.xlsx --date 2025-05-08 [--header-row 42] [--skip-index]")
		os.Exit(1)
	}

	dtReferencia, err := time.Parse("2006-01-02", *dateStr)
	if err != nil {
		log.Fatalf("Data inválida %q: %v", *dateStr, err)
	}

	cfg := config.Load()

	ctx := context.Background()

	pool, err := postgres.NewPool(ctx, cfg.PostgreSQL)
	if err != nil {
		log.Fatalf("PostgreSQL connection failed: %v", err)
	}
	defer pool.Close()

	cmedRepo := postgres.NewCMEDRepo(pool)
	cacheRepo := redisrepo.NewCacheRepo(cfg.Redis)
	meiliRepo := meilisearch.NewMeilisearchRepo(cfg.Meilisearch)

	log.Printf("Abrindo planilha: %s", *source)
	f, err := excelize.OpenFile(*source)
	if err != nil {
		log.Fatalf("Erro ao abrir planilha: %v", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Fatalf("Erro ao ler planilha: %v", err)
	}

	if len(rows) < *headerRow {
		log.Fatalf("Planilha tem %d linhas, mas header-row=%d", len(rows), *headerRow)
	}

	headerMap := make(map[string]int)
	for col, cell := range rows[*headerRow-1] {
		name := strings.TrimSpace(cell)
		if name != "" {
			headerMap[strings.ToUpper(name)] = col
		}
	}

	requiredCols := []string{"REGISTRO", "PRODUTO", "SUBSTÂNCIA"}
	var missing []string
	for _, col := range requiredCols {
		if _, ok := headerMap[col]; !ok {
			missing = append(missing, col)
		}
	}
	if len(missing) > 0 {
		log.Fatalf("Colunas obrigatórias não encontradas no header (linha %d): %v\nColunas encontradas: %v", *headerRow, missing, mapKeys(headerMap))
	}

	dataStartRow := *headerRow
	var records []entity.CMEDConformidade
	var skipped int

	for rowIdx := dataStartRow; rowIdx < len(rows); rowIdx++ {
		row := rows[rowIdx]
		if len(row) == 0 {
			continue
		}

		reg := getCell(row, headerMap, "REGISTRO")
		if reg == "" || reg == "-" {
			skipped++
			continue
		}

		record := entity.CMEDConformidade{
			DTReferencia:    dtReferencia.Format("2006-01-02"),
			STRegistroAtivo: "ACTIVE",
		}

		if v, err := strconv.ParseInt(strings.TrimSpace(reg), 10, 64); err == nil {
			record.NUSanReg = &v
		}

		record.NUGgrem = parseStringPtr(getCell(row, headerMap, "CÓDIGO GGREM"))
		record.DSSubstancia = parseStringPtr(getCell(row, headerMap, "SUBSTÂNCIA"))
		record.NUCnpj = parseStringPtr(getCell(row, headerMap, "CNPJ"))
		record.NOLaboratorio = parseStringPtr(getCell(row, headerMap, "LABORATÓRIO"))
		record.NOProduto = parseStringPtr(getCell(row, headerMap, "PRODUTO"))
		record.DSApresentacao = parseStringPtr(getCell(row, headerMap, "APRESENTAÇÃO"))
		record.DSClasseTerapeutica = parseStringPtr(getCell(row, headerMap, "CLASSE TERAPÊUTICA"))
		record.TPProduto = parseStringPtr(getCell(row, headerMap, "TIPO DE PRODUTO (STATUS DO PRODUTO)"))
		record.TPRegimePreco = parseStringPtr(getCell(row, headerMap, "REGIME DE PREÇO"))

		record.NUEAN1 = parseCleanStringPtr(getCell(row, headerMap, "EAN 1"))
		record.NUEAN2 = parseCleanStringPtr(getCell(row, headerMap, "EAN 2"))
		record.NUEAN3 = parseCleanStringPtr(getCell(row, headerMap, "EAN 3"))

		record.VRPFSemImpostos = parseFloatPtr(getCell(row, headerMap, "PF Sem Impostos"))
		record.VRPF0 = parseFloatPtr(getCell(row, headerMap, "PF 0%"))
		record.VRPF12 = parseFloatPtr(getCell(row, headerMap, "PF 12 %"))
		record.VRPF17 = parseFloatPtr(getCell(row, headerMap, "PF 17 %"))
		record.VRPF18 = parseFloatPtr(getCell(row, headerMap, "PF 18 %"))
		record.VRPF20 = parseFloatPtr(getCell(row, headerMap, "PF 20 %"))
		record.VRPMCSemImpostos = parseFloatPtr(getCell(row, headerMap, "PMC Sem Impostos"))
		record.VRPMC0 = parseFloatPtr(getCell(row, headerMap, "PMC 0 %"))
		record.VRPMC12 = parseFloatPtr(getCell(row, headerMap, "PMC 12 %"))
		record.VRPMC17 = parseFloatPtr(getCell(row, headerMap, "PMC 17 %"))
		record.VRPMC18 = parseFloatPtr(getCell(row, headerMap, "PMC 18 %"))
		record.VRPMC20 = parseFloatPtr(getCell(row, headerMap, "PMC 20 %"))

		record.JSPrecosPF = buildPriceJSON(row, headerMap, "PF")
		record.JSPrecosPMC = buildPriceJSON(row, headerMap, "PMC")

		record.STRestricaoHospitalar = parseCleanStringPtr(getCell(row, headerMap, "RESTRIÇÃO HOSPITALAR"))
		record.STCap = parseCleanStringPtr(getCell(row, headerMap, "CAP"))
		record.STConfaz87 = parseCleanStringPtr(getCell(row, headerMap, "CONFAZ 87"))
		record.STIcms0 = parseCleanStringPtr(getCell(row, headerMap, "ICMS 0%"))
		record.DSAnaliseRecural = parseStringPtr(getCell(row, headerMap, "ANÁLISE RECURSAL"))
		record.DSListaPisCofins = parseStringPtr(getCell(row, headerMap, "LISTA DE CONCESSÃO DE CRÉDITO TRIBUTÁRIO (PIS/COFINS)"))
		record.STComercializacao = parseCleanStringPtr(getCell(row, headerMap, "COMERCIALIZAÇÃO 2025"))
		record.DSTarja = parseStringPtr(getCell(row, headerMap, "TARJA"))
		record.DSDestinacaoComercial = parseCleanStringPtr(getCell(row, headerMap, "DESTINAÇÃO COMERCIAL"))

		records = append(records, record)
	}

	log.Printf("Registros parseados: %d, ignorados: %d", len(records), skipped)

	count, err := cmedRepo.UpsertBatch(ctx, records)
	if err != nil {
		log.Fatalf("Erro no upsert: %v", err)
	}
	log.Printf("Registros importados/atualizados: %d", count)

	if cacheRepo != nil {
		log.Println("Invalidando cache Redis...")
		if err := cacheRepo.DeleteByPattern(ctx, "cmed:*"); err != nil {
			log.Printf("Warning: erro ao invalidar cache cmed:*: %v", err)
		}
		if err := cacheRepo.DeleteByPattern(ctx, "ampp_cmed:*"); err != nil {
			log.Printf("Warning: erro ao invalidar cache ampp_cmed:*: %v", err)
		}
	}

	if !*skipIndex {
		log.Println("Reindexando Meilisearch...")
		if err := meiliRepo.ConfigureIndexes(ctx); err != nil {
			log.Printf("Warning: erro ao configurar indexes: %v", err)
		}
		syncRepo := postgres.NewSyncRepo(pool)
		docs, err := syncRepo.GetAllCMED(ctx)
		if err != nil {
			log.Printf("Warning: erro ao buscar CMED para indexação: %v", err)
		} else {
			if err := meiliRepo.IndexCMEDs(ctx, docs); err != nil {
				log.Printf("Warning: erro ao indexar CMED: %v", err)
			} else {
				log.Printf("Indexados %d documentos CMED no Meilisearch", len(docs))
			}
		}
	}

	log.Println("Importação concluída com sucesso!")
}

func getCell(row []string, headerMap map[string]int, colName string) string {
	col, ok := headerMap[strings.ToUpper(colName)]
	if !ok || col >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[col])
}

func parseStringPtr(s string) *string {
	if s == "" || s == "-" {
		return nil
	}
	return &s
}

func parseCleanStringPtr(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" || s == "-" || s == "- (*)" {
		return nil
	}
	return &s
}

func parseFloatPtr(s string) *float64 {
	s = strings.TrimSpace(s)
	if s == "" || s == "-" {
		return nil
	}
	s = strings.ReplaceAll(s, ",", ".")
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &v
}

func mapKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func buildPriceJSON(row []string, headerMap map[string]int, prefix string) json.RawMessage {
	prices := make(map[string]interface{})

	for colName, colIdx := range headerMap {
		if strings.HasPrefix(colName, prefix+" ") {
			val := ""
			if colIdx < len(row) {
				val = strings.TrimSpace(row[colIdx])
			}
			if val != "" && val != "-" {
				val = strings.ReplaceAll(val, ",", ".")
				if f, err := strconv.ParseFloat(val, 64); err == nil {
					prices[colName] = f
				}
			}
		}
	}

	if len(prices) == 0 {
		return nil
	}
	data, _ := json.Marshal(prices)
	return data
}
