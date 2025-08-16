package cli

const HelpContent = `pgump-mapper is a tool for mapping and exploring pg_dump content.

Usage:

	pgdump-mapper <file> [Options] [Filter Options]

Examples:

	pgdump-mapper example.db.sql.txt --json --table=tableA
	pgdump-mapper example.db.sql.txt --json-pretty --table=tableA
	pgdump-mapper example.db.sql.txt --json --schema=public --table=tableB --cache
	pgdump-mapper example.db.sql.txt --cli --table=tableA --columns=createad_at,id

Options:
	--help, -h	Help
	--json		Export as JSON
	--json-pretty	Export as JSON Pretty
	--yaml		Export as YAML
	--html		Export as HTML (default)
	--sqlite	Export as SQLite
	--cli		Export as CLI table	
	--cache		Read and save cache to /tmp/pgdump-mapper directory
Filter Options:
	--schema	Filter by Schema (valid for json, default: public)
	--table		Filter by Table (valid for json, cli and yaml). Only one at a time.
	--columns	Filter by Columns (valid for cli). Multiples separated by commas.`
