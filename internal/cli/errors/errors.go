package errors

const (
	ErrorNoInputFile = "AppError: no pg_dump file provided.\n\nPlease specify the pg_dump file.\nUse --help or -h for help"
	ErrorManyArgs    = "AppError: too many arguments provided.\n\nPlease specify only the pg_dump file.\nUse --help or -h for help"
	ErrNoSuchFile    = "AppError: no such file or directory.\n\nPlease verify the input file and try again.\nUse --help or -h for help"
	ErrIsDirectory   = "AppError: invalid input.\n\nThe provided path is a directory. Specify a file path and try again.\nUse --help or -h for help"
)
