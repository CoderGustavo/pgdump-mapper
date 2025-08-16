# pgdump-mapper

## Overview

pgdump-mapper is a Go application designed to read PostgreSQL dump files and export the data in various formats. It provides a command-line interface for easy configuration and usage, making it a versatile tool for database management and data migration.

## Features

- **Read PostgreSQL Dump Files**: Efficiently parses PostgreSQL dump files, handling table definitions, primary keys, and foreign keys.
- **Export Options**: Supports exporting data in multiple formats:
  - **JSON**
  - **JSON Pretty**
  - **YAML**
  - **HTML**
  - **SQLite** 

## Project Structure

```
pgdump-mapper
├── internal
│   ├── cli/            # Command-line interface related code
│   ├── file/           # Main logic for reading and exporting data
│   └── models/         # Data models used throughout the project
├── README.md           # Documentation for the project
└── main.go             # Entry point for the application
```

## Usage Instructions

1. **Clone the Repository**:
   ```
   git clone https://github.com/hedibertosilva/pgdump-mapper.git
   cd pgdump-mapper
   ```

2. **Build the Application**:
   ```
   go build -o pgdump-mapper main.go
   ```

3. **Install it**:
   ```
   cp pgdump-mapper /home/$USER/.local/bin
   ```

4. **Run the Application**:
   ```
   pgdump-mapper <pgdump-file> --<help|json|yaml|html|sqlite> [filters: --table --columns]
   ```
## Contributing Code

Feel free to contribute. Contact me via hed.cavalcante@gmail.com