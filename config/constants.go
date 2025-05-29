package constants

const (
	HelpMsg = "pgump-mapper is a tool for mapping and exploring pg_dump content.\n\nUsage:\n\n\tpgdump-mapper <file-path>\n\nOptions:\n\n\t--help, -h\tGive a bit of help about the command line arguments and options.\n\n"

	ErrorNoInputFile    = "no pg_dump file provided"
	ErrorNoInputFileMsg = "Error: no pg_dump file provided.\n\nPlease specify the pg_dump file.\nUse --help or -h for help.\n"
	ErrorManyArgs       = "too many arguments provided"
	ErrorManyArgsMsg    = "Error: too many arguments provided.\n\nPlease specify only the pg_dump file.\nUse --help or -h for help.\n"
	ErrNoSuchFile       = "no such file or directory"
	ErrNoSuchFileMsg    = "Error: no such file or directory.\n\nPlease verify the input file and try again.\nUse --help or -h for help.\n"
	ErrIsDirectory      = "invalid input"
	ErrIsDirectoryMsg   = "Error: invalid input.\n\nThe provided path is a directory. Specify a file path and try again.\nUse --help or -h for help.\n"
)
