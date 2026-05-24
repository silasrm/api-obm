# Tasks: OBM Import CLI

## Task 1: Refactor converter into package
- **What**: Extract `Converter` from `scripts/convert_sql.go` into `internal/importer/converter.go` with public API `Convert(in io.Reader, out io.Writer, opts Options) error`. Add `FullDrop` and `AppendMetadata` options. Add `obm_metadata` CREATE TABLE to the append block alongside `users`.
- **Where**: `internal/importer/converter.go` (new), `scripts/convert_sql.go` (modified)
- **Depends on**: —
- **Reuses**: All regex logic from `scripts/convert_sql.go`
- **Done when**: `go run scripts/convert_sql.go -input test.sql -output test_out.sql` still works; `converter.Convert()` is callable from other packages
- **Tests**: Unit test with small MySQL snippet verifying output is valid PostgreSQL
- **Gate**: `go build ./scripts/...` and `go vet ./internal/importer/...`

## Task 2: Source resolver
- **What**: Create `internal/importer/source.go` with `Resolve(source string) (io.ReadCloser, error)`. Handles `.zip` (extract .sql to temp file), `.sql` (open directly), `.sql.gz` (open + gzip), `mysql://` (exec mysqldump). For ZIP, return reader from temp file; caller is responsible for cleanup via ReadCloser.
- **Where**: `internal/importer/source.go` (new)
- **Depends on**: —
- **Reuses**: —
- **Done when**: Can open a ZIP and return a reader for the SQL inside; can open .sql directly
- **Tests**: Unit test with a small test ZIP file
- **Gate**: `go vet ./internal/importer/...`

## Task 3: PostgreSQL importer
- **What**: Create `internal/importer/pgimporter.go` with `PGImporter` struct and `Import(ctx context.Context, reader io.Reader) error`. Reads SQL statements, executes each via `pgxpool.Pool.Exec()`. Tracks and prints progress (table name, row count). Uses 32MB buffer scanner.
- **Where**: `internal/importer/pgimporter.go` (new)
- **Depends on**: —
- **Reuses**: `pgxpool` connection from existing config pattern
- **Done when**: Can import a small PostgreSQL SQL file into a test database
- **Tests**: Integration test with a tiny SQL file
- **Gate**: `go vet ./internal/importer/...`

## Task 4: Validator
- **What**: Create `internal/importer/validator.go` with `Validator` struct and `Validate(ctx context.Context) (*ValidationReport, error)`. Queries `information_schema.tables` for table list, counts rows per table, checks FK integrity for VMP→VTM, AMP→VMP, AMP→Supplier. Returns structured report with counts, warnings, errors.
- **Where**: `internal/importer/validator.go` (new)
- **Depends on**: —
- **Reuses**: `pgxpool`
- **Done when**: Prints a validation report after import
- **Tests**: Integration test with test data
- **Gate**: `go vet ./internal/importer/...`

## Task 5: Metadata manager
- **What**: Create `internal/importer/metadata.go` with `MetadataManager`. Creates `obm_metadata` table if not exists. Provides `Set(key, value)` and `Get(key)` methods. After import, stores `obm_version`, `import_date`, `source_file`, `record_counts`, `tables_imported`.
- **Where**: `internal/importer/metadata.go` (new)
- **Depends on**: —
- **Reuses**: `pgxpool`
- **Done when**: Can create table and write/read metadata values
- **Tests**: Unit test with mock pool or integration test
- **Gate**: `go vet ./internal/importer/...`

## Task 6: Pipeline orchestrator
- **What**: Create `internal/importer/pipeline.go` with `Pipeline` struct and `Run(ctx context.Context) error`. Orchestrates the full flow: resolve source → convert → import → validate → metadata → reindex. Uses `io.Pipe` for streaming between converter and importer (both in goroutines). Handles `--convert-only` (write to file), `--reindex-only` (skip to reindex), `--skip-index`, `--validate` flags.
- **Where**: `internal/importer/pipeline.go` (new)
- **Depends on**: T1, T2, T3, T4, T5
- **Reuses**: `internal/usecase/reindex.go` (ReindexUsecase), `internal/infrastructure/config/config.go` (Load), `internal/infrastructure/persistence/postgres/postgres.go` (NewPool)
- **Done when**: Full pipeline runs end-to-end from source to reindex
- **Tests**: Integration test (small dump)
- **Gate**: `go vet ./internal/importer/...`

## Task 7: CLI entrypoint
- **What**: Create `cmd/import/main.go` with flag parsing, `.env` loading (godotenv), config loading, pool creation, and pipeline invocation. Flags: `--source`, `--output`, `--convert-only`, `--reindex-only`, `--skip-index`, `--validate`, `--full`. Reads same env vars as the API.
- **Where**: `cmd/import/main.go` (new)
- **Depends on**: T6
- **Reuses**: `internal/infrastructure/config/config.go`, `godotenv`
- **Done when**: `go run cmd/import/main.go --source=test.sql --convert-only` runs without error
- **Tests**: Manual test
- **Gate**: `go build ./cmd/import/...`

## Task 8: Update backward-compat wrapper
- **What**: Modify `scripts/convert_sql.go` to import and call `internal/importer/converter.Convert()` instead of having its own logic. Becomes a thin CLI wrapper (~20 lines).
- **Where**: `scripts/convert_sql.go` (modified)
- **Depends on**: T1
- **Reuses**: `internal/importer/converter.go`
- **Done when**: `go run scripts/convert_sql.go -input test.sql -output out.sql` works as before
- **Tests**: Manual test
- **Gate**: `go build ./scripts/...`

## Task 9: Update README
- **What**: Update "Obter o dump do banco de dados" section in README.md with import CLI usage examples. Add section on the import command with all flags and examples for ZIP, SQL, and MySQL sources.
- **Where**: `README.md` (modified)
- **Depends on**: T7
- **Reuses**: —
- **Done when**: README documents the import CLI
- **Tests**: —
- **Gate**: —
