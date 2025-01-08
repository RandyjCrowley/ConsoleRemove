# Console.log Manager

A powerful Go-based tool for managing `console.log` statements in your JavaScript/TypeScript projects. This tool helps developers find, remove, and restore console.log statements across their codebase, making it easier to clean up debugging code before deployment while maintaining the ability to restore them when needed.

## Features

- **Find Mode**: Locate all `console.log` statements in your codebase
- **Delete Mode**: Remove `console.log` statements with automatic backup creation
- **Revert Mode**: Restore previously removed `console.log` statements from backups
- **Smart Parsing**: Handles both single-line and multi-line console.log statements
- **Comment Awareness**: Ignores console.log statements in comments
- **Extensive Skip List**: Automatically skips common directories like node_modules, build directories, and VCS folders
- **Support for Multiple File Types**: Works with .js, .jsx, .ts, .tsx, .vue, .mjs, .cjs, .html, and .md files
- **Safe Operation**: Creates backup files before making any changes

## Installation

```bash
# Clone the repository
git clone https://github.com/RandyjCrowley/ConsoleRemove
cd ConsoleRemove

# Build the binary
go build -o console-log-manager
```

## Usage

```bash
# Find all console.log statements
go run main.go <directory_path>

# Remove console.log statements (creates .bak files)
go run main.go <directory_path> delete

# Restore from backups
go run main.go <directory_path> revert
```

## Examples

```bash
# Search for console.log statements in the current directory
go run main.go .

# Remove console.log statements from a specific project
go run main.go /path/to/your/project delete

# Restore previously removed console.log statements
go run main.go /path/to/your/project revert
```

## Supported File Types

- JavaScript (.js)
- TypeScript (.ts)
- React/JSX (.jsx)
- TypeScript React (.tsx)
- Vue (.vue)
- ES Modules (.mjs)
- CommonJS Modules (.cjs)
- HTML (.html)
- Markdown (.md)

## Future Ideas

1. **Configuration File Support**
   - Custom skip patterns
   - Configurable file extensions
   - Project-specific settings

2. **Enhanced Functionality**
   - Support for other console methods (console.warn, console.error, etc.)
   - Regular expression pattern matching for more flexible searches
   - Selective removal based on patterns or content
   - Statistics reporting (number of statements found/removed)

3. **Output Options**
   - JSON output format
   - HTML report generation
   - Integration with CI/CD pipelines

4. **Code Quality Features**
   - Detection of nested console.logs
   - Identification of potential performance impacts
   - Code style preservation
   - Support for source maps

5. **Safety Improvements**
   - Dry run mode
   - Backup compression
   - Backup rotation
   - Checksum verification

6. **Development Tools**
   - Git hooks integration
   - Webpack/Rollup plugin
   - npm package distribution

7. **Performance Optimizations**
   - Parallel processing for large codebases
   - Incremental processing
   - Memory optimization for large files
   - Caching mechanism

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
