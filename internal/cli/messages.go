package cli

const HelpContent = `pgump-mapper is a tool for mapping and exploring pg_dump content.

Usage:

	pgdump-mapper [Options] <file-path>

Example:

	pgdump-mapper example.db.sql.txt --json --table=tableA

Options:
	--help, -h		Give a bit of help about the command line arguments and options
	--json			Export as JSON
	--json-pretty		Export as JSON Pretty
	--yaml			Export as YAML
	--html			Export as HTML (default)
	--sqlite		Export as SQLite
	--cli			Export as CLI table
	--schema		Filter by Schema (default schema: public)
	--table			Filter by Table (valid for json, cli and yaml)
	--columns		Filter by Columns (valid for json, cli and yaml)`
