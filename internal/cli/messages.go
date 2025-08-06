package cli

const HelpContent = `pgump-mapper is a tool for mapping and exploring pg_dump content.

Usage:

	pgdump-mapper [Options] <file-path>

Example:

	pgdump-mapper example.db.sql.txt --html

Options:
	--help, -h	Give a bit of help about the command line arguments and options.
	--json		Export as JSON.
	--yaml		Export as YAML.
	--html		Export as HTML. (default)
	--sqlite    Export as SQLite.
	--cli       Export as CLI table.
	--table     Filter by table (valid only with --cli).
	--columns   Filter by Columns (valid only with --cli).`
