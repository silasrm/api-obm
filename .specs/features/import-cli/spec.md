# Feature Spec: OBM Import CLI

## Summary
CLI tool that takes an OBM data file (.zip, .sql, or MySQL connection string), converts from MySQL to PostgreSQL, imports into the database, and reindexes Meilisearch.

## Background
The OBM portal publishes data as a MySQL dump inside a ZIP file (e.g., `portal-obm-20250530.zip` containing `portal-obm-20250530.sql`, ~15 MB zipped, ~140 MB uncompressed). Access is restricted — users must download manually. The current pipeline requires manually running `convert_sql.go`, then placing the file in `migrations/postgres/` for Docker init. This is fragile and undocumented.

## Requirements

### R1: Source resolution
The CLI accepts `--source` with three formats:
- **ZIP file** (local path ending in `.zip`): Extract the `.sql` file inside
- **SQL file** (local path ending in `.sql`): Use directly
- **MySQL connection string** (`mysql://user:pass@host:port/db`): Execute `mysqldump` to produce a SQL dump

### R2: MySQL to PostgreSQL conversion
Reuse existing `convert_sql.go` logic as a reusable package (`internal/importer/converter.go`):
- Same regex-based line-by-line conversion (bigint→BIGINT, tinyint→BOOLEAN, INSERT IGNORE→ON CONFLICT DO NOTHING, etc.)
- Skip tables: `roles`, `role_users`, `sessions`
- Append `CREATE TABLE users` and `CREATE TABLE obm_metadata` at end
- In `--full` mode: prepend `DROP SCHEMA public CASCADE; CREATE SCHEMA public;`

### R3: PostgreSQL import
Connect to PostgreSQL (using same env vars as API) and execute the converted SQL:
- Split SQL into statements (by `;`)
- Execute in batches within transactions
- Progress output: print table names and row counts as imported
- Timeout: 30 minutes default for large dumps

### R4: Meilisearch reindex
After successful import, trigger reindex (reuse existing `ReindexUsecase`):
- Same flow as `POST /api/v1/admin/reindex` and startup sync
- Index VMP, AMP, Supplier (no change to current scope)

### R5: Validation
Optional `--validate` flag that after import:
- Counts records per table (`SELECT COUNT(*) FROM <table>`)
- Checks FK integrity (VMP→VTM, AMP→VMP, AMP→Supplier)
- Reports orphaned references as warnings
- Prints summary

### R6: Metadata tracking
New `obm_metadata` table to track import versioning:
- `obm_version`: date extracted from filename (e.g., `20250530`)
- `import_date`: timestamp of import
- `source_file`: original filename
- `record_counts`: JSON with counts per table
- `tables_imported`: number of tables

### R7: Streaming pipeline
The conversion and import should stream without writing the full 140 MB SQL to disk (except with `--convert-only`):
- Use `io.Pipe` for converter→importer streaming
- ZIP extraction also streams (no temp files for the .sql)

### R8: CLI flags
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--source` | string | (required) | Path to .zip/.sql or mysql:// URI |
| `--output` | string | (temp) | Output path for --convert-only |
| `--convert-only` | bool | false | Only convert, don't import |
| `--reindex-only` | bool | false | Only reindex Meilisearch |
| `--skip-index` | bool | false | Skip Meilisearch reindex after import |
| `--validate` | bool | false | Run post-import validation |
| `--full` | bool | true | Drop + recreate schema before import |

### R9: Backward compatibility
`scripts/convert_sql.go` becomes a thin wrapper calling `internal/importer/converter.go` so existing workflows still work.

## Out of scope
- Scheduled/automatic imports
- Incremental/delta updates (always full replace)
- Expanding Meilisearch indexes beyond VMP/AMP/Supplier
- CSV/XML/JSON input formats
- Downloading from portal URL (access is restricted)

## Success criteria
1. `go run cmd/import/main.go --source=portal-obm-20250530.zip` imports all 67 tables, reindexes Meilisearch, and records metadata
2. `go run cmd/import/main.go --source=dump.sql --convert-only --output=out.sql` produces a valid PostgreSQL file
3. `go run cmd/import/main.go --reindex-only` reindexes without touching PostgreSQL data
4. `go run cmd/import/main.go --source=dump.zip --validate` prints validation report
5. `go run scripts/convert_sql.go -input=dump.sql -output=out.sql` still works (backward compat)
