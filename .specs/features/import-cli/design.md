# Design: OBM Import CLI

## Architecture

```
cmd/import/main.go
       │
       ▼
internal/importer/pipeline.go  ─── orchestrates the full flow
       │
       ├── source.go     ─── resolves --source to io.Reader
       ├── converter.go  ─── MySQL→PostgreSQL (refactored from scripts/convert_sql.go)
       ├── pgimporter.go ─── executes SQL against PostgreSQL
       ├── validator.go  ─── post-import integrity checks
       ├── metadata.go   ─── obm_metadata table CRUD
       └── (reuses) internal/usecase/reindex.go ─── Meilisearch reindex
```

## Component Details

### source.go
```
SourceResolver
  Resolve(source string) (io.Reader, error)
  
  .zip  → archive/zip → find *.sql inside → open via io.ReaderAt
  .sql  → os.Open directly
  .sql.gz → os.Open + compress/gzip
  mysql:// → parse URL → exec mysqldump → stdout pipe
```

For ZIP: Go's `archive/zip` requires `io.ReaderAt` (not streaming). We extract the `.sql` entry to a temp file, then return a reader for it. This is necessary because the converter needs to scan line-by-line and the ZIP format doesn't support true streaming of compressed entries.

For `mysql://`: Parse with `net/url`, construct `mysqldump` command, pipe stdout. Requires `mysqldump` binary on PATH.

### converter.go
Extract `Converter` struct and `Convert(in io.Reader, out io.Writer) error` from `scripts/convert_sql.go`. No changes to regex logic. Add:

- `FullDrop bool` field → prepends `DROP SCHEMA public CASCADE; CREATE SCHEMA public;`
- `AppendMetadata bool` field → appends `CREATE TABLE IF NOT EXISTS obm_metadata (...)`
- Public API: `func Convert(in io.Reader, out io.Writer, opts Options) error`

### pgimporter.go
```
PGImporter
  New(pool *pgxpool.Pool)
  Import(ctx context.Context, reader io.Reader) error
  
  - Read SQL in chunks (bufio.Scanner, 32MB buffer like converter)
  - Accumulate statement text until ';'
  - Execute each statement via pool.Exec()
  - Track: current table name, rows inserted, errors
  - Print progress to stdout
  - Wrap in transaction? No — 140MB SQL in a single tx is too large.
    Use autocommit per statement instead (same as psql).
```

### validator.go
```
Validator
  New(pool *pgxpool.Pool)
  Validate(ctx context.Context) (*ValidationReport, error)
  
  - List tables from information_schema
  - SELECT COUNT(*) for each
  - FK checks as NOT IN subqueries
  - Return structured report
```

### metadata.go
```
MetadataManager
  New(pool *pgxpool.Pool)
  CreateTable(ctx context.Context) error
  Set(ctx context.Context, key, value string) error
  Get(ctx context.Context, key string) (string, error)
  
  CREATE TABLE IF NOT EXISTS obm_metadata (
    key VARCHAR(100) PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
  );
```

### pipeline.go
```
Pipeline
  New(cfg Config) *Pipeline
  Run(ctx context.Context) error
  
  Config: Source, Output, ConvertOnly, ReindexOnly, SkipIndex, Validate, Full
  
  Flow:
  1. If reindex-only → reindex and return
  2. Resolve source → io.Reader
  3. If convert-only → converter.Convert(reader, outputFile) and return
  4. Create io.Pipe: converter writes to PipeWriter, pgimporter reads from PipeReader
  5. Run converter and importer in goroutines
  6. Wait for both
  7. If validate → run validator
  8. Write metadata (obm_version, import_date, etc.)
  9. If !skip-index → reindex Meilisearch
```

## Key Design Decisions

1. **No single transaction for import** — 140 MB of INSERTs in one tx would blow up WAL. Use autocommit per statement (same behavior as `psql -f file.sql`).

2. **Pipe-based streaming** — `io.Pipe` connects converter output directly to importer input without disk. Both run in goroutines.

3. **ZIP requires temp file** — Go's `archive/zip` doesn't stream entries. We extract to a temp file and clean up after. For .sql files, we read directly.

4. **Reuse existing reindex** — The `ReindexUsecase` from `internal/usecase/reindex.go` already does everything needed. The import CLI just needs to instantiate it the same way `main.go` does.

5. **No schema migration tracking** — The import is always full-replace (`DROP SCHEMA`). No need for migration versioning like Flyway.

6. **obm_metadata table** — Simple key-value store. Created by the converter (appended to SQL) or by metadata.go if it doesn't exist. This lets even the `--convert-only` flow produce a SQL file with the metadata table.
