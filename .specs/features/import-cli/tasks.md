# Tasks: OBM Import CLI

## Task 1: Refactor converter into package — DONE
- **What**: Extract `Converter` from `scripts/convert_sql.go` into `internal/importer/converter.go` with public API `Convert(in io.Reader, out io.Writer, opts Options) error`. Add `FullDrop` and `AppendMetadata` options. Add `obm_metadata` CREATE TABLE to the append block alongside `users`.
- **Where**: `internal/importer/converter.go` (new), `scripts/convert_sql.go` (modified)
- **Depends on**: —
- **Status**: DONE — converter.go (474 lines), scripts/convert_sql.go is thin wrapper (47 lines)

## Task 2: Source resolver — DONE
- **What**: Create `internal/importer/source.go` with `Resolve(source string) (io.ReadCloser, error)`. Handles `.zip` (extract .sql to temp file), `.sql` (open directly), `.sql.gz` (open + gzip), `mysql://` (exec mysqldump). For ZIP, return reader from temp file; caller is responsible for cleanup via ReadCloser.
- **Where**: `internal/importer/source.go` (new)
- **Depends on**: —
- **Status**: DONE — source.go (202 lines), includes ExtractVersion()

## Task 3: PostgreSQL importer — DONE
- **What**: Create `internal/importer/pgimporter.go` with `PGImporter` struct and `Import(ctx context.Context, reader io.Reader) error`. Reads SQL statements, executes each via `pgxpool.Pool.Exec()`. Tracks and prints progress (table name, row count). Uses 32MB buffer scanner.
- **Where**: `internal/importer/pgimporter.go` (new)
- **Depends on**: —
- **Status**: DONE — pgimporter.go (93 lines), ImportStats with TablesCreated/RowsInserted/Duration/Errors

## Task 4: Validator — DONE
- **What**: Create `internal/importer/validator.go` with `Validator` struct and `Validate(ctx context.Context) (*ValidationReport, error)`. Queries `information_schema.tables` for table list, counts rows per table, checks FK integrity for VMP→VTM, AMP→VMP, AMP→Supplier. Returns structured report with counts, warnings, errors.
- **Where**: `internal/importer/validator.go` (new)
- **Depends on**: —
- **Status**: DONE — validator.go (121 lines), FK checks + row counts + summary

## Task 5: Metadata manager — DONE
- **What**: Create `internal/importer/metadata.go` with `MetadataManager`. Creates `obm_metadata` table if not exists. Provides `Set(key, value)` and `Get(key)` methods. After import, stores `obm_version`, `import_date`, `source_file`, `record_counts`, `tables_imported`.
- **Where**: `internal/importer/metadata.go` (new)
- **Depends on**: —
- **Status**: DONE — metadata.go (76 lines), EnsureTable/Set/Get/RecordImport

## Task 6: Pipeline orchestrator — DONE
- **What**: Create `internal/importer/pipeline.go` with `Pipeline` struct and `Run(ctx context.Context) error`. Orchestrates the full flow: resolve source → convert → import → validate → metadata → reindex. Uses `io.Pipe` for streaming between converter and importer (both in goroutines). Handles `--convert-only` (write to file), `--reindex-only` (skip to reindex), `--skip-index`, `--validate` flags.
- **Where**: `internal/importer/pipeline.go` (new)
- **Depends on**: T1, T2, T3, T4, T5
- **Status**: DONE — pipeline.go (180 lines), streaming via io.Pipe + reindex integration

## Task 7: CLI entrypoint — DONE
- **What**: Create `cmd/import/main.go` with flag parsing, `.env` loading (godotenv), config loading, pool creation, and pipeline invocation. Flags: `--source`, `--output`, `--convert-only`, `--reindex-only`, `--skip-index`, `--validate`, `--full`. Reads same env vars as the API.
- **Where**: `cmd/import/main.go` (new)
- **Depends on**: T6
- **Status**: DONE — main.go (59 lines), 30-min timeout, same env vars as API

## Task 8: Update backward-compat wrapper — DONE
- **What**: Modify `scripts/convert_sql.go` to import and call `internal/importer/converter.Convert()` instead of having its own logic. Becomes a thin CLI wrapper (~20 lines).
- **Where**: `scripts/convert_sql.go` (modified)
- **Depends on**: T1
- **Status**: DONE — already a thin wrapper calling `importer.Convert()`

## Task 9: Update README — DONE
- **What**: Update "Obter o dump do banco de dados" section in README.md with import CLI usage examples. Add section on the import command with all flags and examples for ZIP, SQL, and MySQL sources.
- **Where**: `README.md` (modified)
- **Depends on**: T7
- **Status**: DONE — step 1 of "Instalação Local" now documents the import CLI with flags table, source types, and examples

## Remaining: E2E test with real dump
- Test with `/Users/silasrm/Downloads/portal-obm-20250530.zip` or `portal-obm-20250530.sql`
- Requires running PostgreSQL + Meilisearch (Docker Compose)
