package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/silasrm/api-obm/internal/importer"
	"github.com/silasrm/api-obm/internal/infrastructure/config"
)

func main() {
	sourceFlag := flag.String("source", "", "Caminho para .zip/.sql ou string mysql://user:pass@host:port/db")
	outputFlag := flag.String("output", "", "Caminho de saida para --convert-only (padrao: migrations/postgres/001_obm_schema.sql)")
	convertOnlyFlag := flag.Bool("convert-only", false, "Somente converter, sem importar")
	reindexOnlyFlag := flag.Bool("reindex-only", false, "Somente reindexar Meilisearch")
	skipIndexFlag := flag.Bool("skip-index", false, "Pular reindexacao Meilisearch apos importacao")
	validateFlag := flag.Bool("validate", false, "Executar validacao pos-importacao")
	fullFlag := flag.Bool("full", true, "Drop + recreate schema antes de importar (padrao)")
	flag.Parse()

	godotenv.Load()
	cfg := config.Load()

	if !*reindexOnlyFlag && *sourceFlag == "" {
		fmt.Fprintln(os.Stderr, "Erro: --source e obrigatorio (exceto com --reindex-only)")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Uso:")
		fmt.Fprintln(os.Stderr, "  go run cmd/import/main.go --source=portal-obm-20250530.zip")
		fmt.Fprintln(os.Stderr, "  go run cmd/import/main.go --source=dump.sql --convert-only")
		fmt.Fprintln(os.Stderr, "  go run cmd/import/main.go --source=mysql://user:pass@host/db")
		fmt.Fprintln(os.Stderr, "  go run cmd/import/main.go --reindex-only")
		os.Exit(1)
	}

	pipelineCfg := importer.PipelineConfig{
		Source:      *sourceFlag,
		Output:      *outputFlag,
		ConvertOnly: *convertOnlyFlag,
		ReindexOnly: *reindexOnlyFlag,
		SkipIndex:   *skipIndexFlag,
		Validate:    *validateFlag,
		Full:        *fullFlag,
		PGConfig:    cfg.PostgreSQL,
		MeiliConfig: cfg.Meilisearch,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	pipeline := importer.NewPipeline(pipelineCfg)
	if err := pipeline.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		os.Exit(1)
	}
}
