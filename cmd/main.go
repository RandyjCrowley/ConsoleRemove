package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	chunkSize = 1024 * 1024 // 1MB chunks for reading
)

var skipPatterns = []string{
	// Version Control
	".git",
	".svn",
	".hg",
	".bzr",
	"_darcs",

	// Node.js
	"node_modules",
	"bower_components",

	// Python
	"venv",
	"env",
	".env",
	"__pycache__",
	".pytest_cache",
	"*.egg-info",
	".eggs",
	".tox",

	// Ruby
	"vendor/bundle",
	".bundle",

	// PHP
	"vendor",
	"composer/cache",

	// Java/Kotlin
	"target",
	".gradle",
	"build",
	".maven",
	".m2",

	// Rust
	"target",
	"cargo/registry",

	// Go
	"vendor",
	".go/cache",
	"go/pkg",

	// .NET
	"bin",
	"obj",
	"packages",

	// Build outputs and distributions
	"dist",
	"build",
	"out",
	"output",
	"release",
	"releases",
	"public/build",
	".next",
	".nuxt",
	".vuepress/dist",

	// IDE and Editor
	".idea",
	".vscode",
	".vs",
	"*.sublime-workspace",
	".atom",

	// Documentation
	"docs/_build",
	"site",
	".docusaurus",

	// Cache and Temp
	".cache",
	"tmp",
	".tmp",
	"temp",
	".temp",

	// Coverage and Tests
	"coverage",
	".nyc_output",
	"htmlcov",
	".coverage",

	// Logs
	"logs",
	"*.log",

	// CI/CD
	".jenkins",
	".github",
	".gitlab",
	".circleci",

	// Dependencies lockfiles (optional)
	"package-lock.json",
	"yarn.lock",
	"Gemfile.lock",
	"composer.lock",
	"poetry.lock",

	// Minified/Generated
	"*.min.js",
	"*.min.css",

	// Extensions
	"raycast",

	// Local State
	".local",
	".vim",

	// User Defined
	"Library",
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("Usage: go run main.go <directory_path> [delete|revert]")
		fmt.Println("  - Without argument: Find all console.log statements")
		fmt.Println("  - delete: Remove console.log statements and create .bak files")
		fmt.Println("  - revert: Restore files from .bak files and remove the backups")
		os.Exit(1)
	}

	rootDir := os.Args[1]
	mode := ""
	if len(os.Args) == 3 {
		mode = os.Args[2]
	}

	switch mode {
	case "delete":
		err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				// Handle permission errors and other issues gracefully
				if os.IsPermission(err) {
					fmt.Printf("Skipping (permission denied): %s\n", path)
					return filepath.SkipDir
				}
				if info == nil {
					fmt.Printf("Skipping (access error): %s\n", path)
					return filepath.SkipDir
				}
				return nil
			}
			return processFile(path, info, nil, true)
		})
		if err != nil {
			fmt.Printf("Error walking through directory: %v\n", err)
			os.Exit(1)
		}
	case "revert":
		err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				if os.IsPermission(err) {
					fmt.Printf("Skipping (permission denied): %s\n", path)
					return filepath.SkipDir
				}
				if info == nil {
					fmt.Printf("Skipping (access error): %s\n", path)
					return filepath.SkipDir
				}
				return nil
			}
			return revertFile(path, info, nil)
		})
		if err != nil {
			fmt.Printf("Error walking through directory: %v\n", err)
			os.Exit(1)
		}
	case "":
		err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				if os.IsPermission(err) {
					fmt.Printf("Skipping (permission denied): %s\n", path)
					return filepath.SkipDir
				}
				if info == nil {
					fmt.Printf("Skipping (access error): %s\n", path)
					return filepath.SkipDir
				}
				return nil
			}
			return processFile(path, info, nil, false)
		})
		if err != nil {
			fmt.Printf("Error walking through directory: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown argument: %s\n", mode)
		os.Exit(1)
	}
}

func revertFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	// Check if this is a .bak file
	if strings.HasSuffix(path, ".bak") {
		originalPath := strings.TrimSuffix(path, ".bak")

		// Read backup file
		backupContent, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading backup file %s: %v", path, err)
		}

		// Write content back to original file
		err = os.WriteFile(originalPath, backupContent, 0644)
		if err != nil {
			return fmt.Errorf("error writing original file %s: %v", originalPath, err)
		}

		// Remove backup file
		err = os.Remove(path)
		if err != nil {
			return fmt.Errorf("error removing backup file %s: %v", path, err)
		}

		fmt.Printf("Reverted file: %s (backup deleted)\n", originalPath)
	}

	return nil
}

func shouldSkipPath(path string) bool {
	normalizedPath := filepath.ToSlash(path)
	for _, pattern := range skipPatterns {
		if strings.Contains(normalizedPath, pattern) {
			return true
		}
	}
	return false
}

func processFile(path string, info os.FileInfo, err error, shouldDelete bool) error {
	if err != nil {
		return err
	}

	if shouldSkipPath(path) {
		if info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}

	if info.IsDir() {
		return nil
	}

	if !isRelevantFile(path) {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file %s: %v", path, err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", path, err)
	}

	if shouldDelete {
		newContent, changed := removeConsoleLog(content, path)
		if changed {
			// Create backup file
			backupPath := path + ".bak"
			err = os.WriteFile(backupPath, content, 0644)
			if err != nil {
				return fmt.Errorf("error creating backup file %s: %v", backupPath, err)
			}

			// Write new content
			err = os.WriteFile(path, newContent, 0644)
			if err != nil {
				return fmt.Errorf("error writing file %s: %v", path, err)
			}
			fmt.Printf("Updated file: %s (backup created at %s)\n", path, backupPath)
		}
	} else {
		return searchInContent(content, path)
	}

	return nil
}

func removeConsoleLog(content []byte, path string) ([]byte, bool) {
	lines := bytes.Split(content, []byte{'\n'})
	var newLines [][]byte
	inComment := false
	inConsoleLog := false
	parenCount := 0
	changed := false

	for i := 0; i < len(lines); i++ {
		line := string(lines[i])
		shouldKeepLine := true

		// Handle comments
		if !inComment && strings.Contains(line, "/*") {
			inComment = true
		}
		if inComment && strings.Contains(line, "*/") {
			inComment = false
		}
		if !inComment && strings.HasPrefix(strings.TrimSpace(line), "//") {
			newLines = append(newLines, lines[i])
			continue
		}
		if inComment {
			newLines = append(newLines, lines[i])
			continue
		}

		if inConsoleLog {
			shouldKeepLine = false
			changed = true
			for _, ch := range line {
				if ch == '(' {
					parenCount++
				} else if ch == ')' {
					parenCount--
					if parenCount == 0 {
						inConsoleLog = false
						break
					}
				}
			}
		} else if strings.Contains(line, "console.log") {
			startIdx := strings.Index(line, "console.log")
			prefix := line[:startIdx]

			// Check if this console.log is in a comment
			if strings.Contains(prefix, "//") {
				newLines = append(newLines, lines[i])
				continue
			}

			rest := line[startIdx:]
			openParenIndex := strings.Index(rest, "(")

			if openParenIndex != -1 {
				// Count parentheses in the rest of the line
				parenCount = 0
				for _, ch := range rest[openParenIndex:] {
					if ch == '(' {
						parenCount++
					} else if ch == ')' {
						parenCount--
						if parenCount == 0 {
							break
						}
					}
				}

				if parenCount == 0 {
					// Single-line console.log
					shouldKeepLine = false
					changed = true
				} else {
					// Multi-line console.log starts here
					inConsoleLog = true
					shouldKeepLine = false
					changed = true
				}
			}
		}

		if shouldKeepLine {
			newLines = append(newLines, lines[i])
		}
	}

	return bytes.Join(newLines, []byte{'\n'}), changed
}

func isRelevantFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))

	// Skip minified files
	if strings.Contains(path, "-min") || strings.Contains(path, ".min") {
		return false
	}

	relevantExts := map[string]bool{
		".js":   true,
		".jsx":  true,
		".ts":   true,
		".tsx":  true,
		".vue":  true,
		".mjs":  true,
		".cjs":  true,
		".html": true,
		".md":   true,
	}
	return relevantExts[ext]
}

func searchInContent(content []byte, path string) error {
	lines := bytes.Split(content, []byte{'\n'})
	inComment := false
	inConsoleLog := false
	consoleLogStart := 0
	parenCount := 0
	var consoleLogLines []string

	for i := 0; i < len(lines); i++ {
		line := string(lines[i])
		trimmedLine := strings.TrimSpace(line)

		// Handle comments
		if !inComment {
			if strings.HasPrefix(trimmedLine, "//") {
				continue
			}
			if strings.Contains(line, "/*") {
				inComment = true
				continue
			}
		}
		if inComment {
			if strings.Contains(line, "*/") {
				inComment = false
			}
			continue
		}

		// If we're in a console.log statement, continue collecting lines
		if inConsoleLog {
			consoleLogLines = append(consoleLogLines, line)
			for _, ch := range line {
				if ch == '(' {
					parenCount++
				} else if ch == ')' {
					parenCount--
					if parenCount == 0 {
						fullStatement := strings.Join(consoleLogLines, "\n")
						if len(fullStatement) > 500 {
							fullStatement = fullStatement[:497] + "..."
						}
						fmt.Printf("File: %s\nLines %d-%d: %s\n\n",
							path, consoleLogStart, i+1, fullStatement)
						inConsoleLog = false
						break
					}
				}
			}
			continue
		}

		// Look for new console.log statements
		if strings.Contains(line, "console.log") {
			startIdx := strings.Index(line, "console.log")
			prefix := line[:startIdx]

			if strings.Contains(prefix, "//") {
				continue
			}

			consoleLogStart = i + 1
			consoleLogLines = []string{line}

			rest := line[startIdx:]
			openParenIndex := strings.Index(rest, "(")
			if openParenIndex == -1 {
				inConsoleLog = true
				parenCount = 0
				continue
			}

			inConsoleLog = true
			parenCount = 0
			for _, ch := range rest[openParenIndex:] {
				if ch == '(' {
					parenCount++
				} else if ch == ')' {
					parenCount--
					if parenCount == 0 {
						fmt.Printf("File: %s\nLine %d: %s\n\n",
							path, i+1, strings.TrimSpace(line))
						inConsoleLog = false
						break
					}
				}
			}
		}
	}

	return nil
}

