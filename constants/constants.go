package constants

const (
	HelpContent = "pgump-mapper is a tool for mapping and exploring pg_dump content.\n\nUsage:\n\n\tpgdump-mapper <file-path>\n\nOptions:\n\n\t--help, -h\tGive a bit of help about the command line arguments and options."

	ErrorNoInputFile = "AppError: no pg_dump file provided.\n\nPlease specify the pg_dump file.\nUse --help or -h for help"
	ErrorManyArgs    = "AppError: too many arguments provided.\n\nPlease specify only the pg_dump file.\nUse --help or -h for help"
	ErrNoSuchFile    = "AppError: no such file or directory.\n\nPlease verify the input file and try again.\nUse --help or -h for help"
	ErrIsDirectory   = "AppError: invalid input.\n\nThe provided path is a directory. Specify a file path and try again.\nUse --help or -h for help"
)
