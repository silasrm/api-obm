package importer

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/silasrm/api-obm/internal/infrastructure/config"
	meilisearchrepo "github.com/silasrm/api-obm/internal/infrastructure/persistence/meilisearch"
	"github.com/silasrm/api-obm/internal/infrastructure/persistence/postgres"
	"github.com/silasrm/api-obm/internal/usecase"
)

type PipelineConfig struct {
	Source       string
	Output       string
	ConvertOnly  bool
	ReindexOnly  bool
	SkipIndex    bool
	Validate     bool
	Full         bool
	PGConfig     config.PostgresConfig
	MeiliConfig  config.MeilisearchConfig
}

type Pipeline struct {
	cfg  PipelineConfig
	pool *pgxpool.Pool
}

func NewPipeline(cfg PipelineConfig) *Pipeline {
	return &Pipeline{cfg: cfg}
}

func (p *Pipeline) Run(ctx context.Context) error {
	if p.cfg.ReindexOnly {
		return p.reindex(ctx)
	}

	reader, err := Resolve(p.cfg.Source)
	if err != nil {
		return fmt.Errorf("resolve source: %w", err)
	}
	defer reader.Close()

	if p.cfg.ConvertOnly {
		return p.convertOnly(ctx, reader)
	}

	if p.cfg.Full {
		fmt.Println("Modo full: o schema será recriado")
	}

	pool, err := postgres.NewPool(ctx, p.cfg.PGConfig)
	if err != nil {
		return fmt.Errorf("connect postgres: %w", err)
	}
	defer pool.Close()
	p.pool = pool

	pr, pw := io.Pipe()

	var importStats *ImportStats
	var importErr error
	importDone := make(chan struct{})

	go func() {
		defer close(importDone)
		importer := NewPGImporter(pool)
		importStats, importErr = importer.Import(ctx, pr)
	}()

	converterOpts := Options{
		FullDrop:        p.cfg.Full,
		AppendMetadata:  true,
	}

	convertErr := Convert(reader, pw, converterOpts)
	pw.Close()

	<-importDone

	if convertErr != nil {
		return fmt.Errorf("converter: %w", convertErr)
	}
	if importErr != nil {
		return fmt.Errorf("importer: %w", importErr)
	}

	fmt.Printf("Importação concluída: %d tabelas criadas\n", importStats.TablesCreated)

	if p.cfg.Validate {
		validator := NewValidator(pool)
		report, err := validator.Validate(ctx)
		if err != nil {
			log.Printf("Validação falhou: %v", err)
		} else {
			totalRows := int64(0)
			for _, count := range report.TableCounts {
				totalRows += count
			}
			fmt.Printf("Validação: %d tabelas, %d registros totais, %d avisos de integridade referencial\n",
				len(report.TableCounts), totalRows, len(report.FKWarnings))
			for _, w := range report.FKWarnings {
				fmt.Printf("  ⚠ %s\n", w)
			}
		}
	}

	version := ExtractVersion(p.cfg.Source)
	mm := NewMetadataManager(pool)
	if err := mm.RecordImport(ctx, version, p.cfg.Source, importStats.RowsInserted); err != nil {
		log.Printf("Aviso: falhou ao gravar metadados: %v", err)
	} else {
		fmt.Printf("Metadados gravados: versão=%s\n", version)
	}

	if !p.cfg.SkipIndex {
		if err := p.reindex(ctx); err != nil {
			log.Printf("Aviso: reindexação falhou: %v", err)
		}
	}

	return nil
}

func (p *Pipeline) convertOnly(ctx context.Context, reader io.ReadCloser) error {
	outPath := p.cfg.Output
	if outPath == "" {
		outPath = "migrations/postgres/001_obm_schema.sql"
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer outFile.Close()

	opts := Options{
		FullDrop:       p.cfg.Full,
		AppendMetadata: true,
	}

	if err := Convert(reader, outFile, opts); err != nil {
		return fmt.Errorf("convert: %w", err)
	}

	fmt.Printf("Conversão concluída: %s\n", outPath)
	return nil
}

func (p *Pipeline) reindex(ctx context.Context) error {
	if p.pool == nil {
		pool, err := postgres.NewPool(ctx, p.cfg.PGConfig)
		if err != nil {
			return fmt.Errorf("connect postgres: %w", err)
		}
		defer pool.Close()
		p.pool = pool
	}

	syncRepo := postgres.NewSyncRepo(p.pool)
	meiliRepo := meilisearchrepo.NewMeilisearchRepo(p.cfg.MeiliConfig)
	reindexUC := usecase.NewReindexUsecase(syncRepo, meiliRepo)

	fmt.Println("Reindexando Meilisearch...")
	indexed, err := reindexUC.Reindex(ctx)
	if err != nil {
		return fmt.Errorf("reindex: %w", err)
	}

	fmt.Println("Reindexação concluída:")
	for entity, count := range indexed {
		fmt.Printf("  %s: %d documentos\n", entity, count)
	}
	return nil
}
